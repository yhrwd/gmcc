package web

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	authsession "gmcc/internal/auth/session"
	"gmcc/internal/cluster"
	"gmcc/internal/logx"
	"gmcc/internal/resource"
	"gmcc/internal/state"
	"gmcc/internal/systemmetrics"
	"gmcc/internal/web/audit"
	"gmcc/internal/webtypes"
	"gmcc/internal/webui"
)

// Server Web服务器
type Server struct {
	config            webtypes.WebConfig
	configPath        string // Web配置文件路径
	router            *gin.Engine
	httpServer        *http.Server
	clusterManager    *cluster.Manager
	resourceManager   accountReader
	runtimeAuth       *authsession.AuthManager
	auditLogger       *audit.Logger
	uiAssets          webui.UIAssets
	resourceCollector systemmetrics.Collector
}

type accountReader interface {
	ListAccounts() ([]resource.AccountRecord, error)
	GetAccount(accountID string) (resource.AccountRecord, error)
	CreateAccount(in resource.CreateAccountInput) (state.AccountMeta, error)
	DeleteAccount(accountID string) error
}

// NewServer 创建Web服务器
func NewServer(config webtypes.WebConfig, configPath string, clusterManager *cluster.Manager, resourceManager accountReader, runtimeAuth *authsession.AuthManager, resourceCollector systemmetrics.Collector) (*Server, error) {
	// 创建审计日志管理器
	logDir := auditLogDir(configPath)
	auditLogger, err := audit.NewLogger(logDir, config.Auth.AuditLogRetentionDays)
	if err != nil {
		return nil, err
	}

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	server := &Server{
		config:            config,
		configPath:        configPath,
		clusterManager:    clusterManager,
		resourceManager:   resourceManager,
		runtimeAuth:       runtimeAuth,
		auditLogger:       auditLogger,
		uiAssets:          webui.NewEmbeddedAssets(),
		resourceCollector: resourceCollector,
	}

	// 设置路由
	server.setupRoutes()

	return server, nil
}

func auditLogDir(configPath string) string {
	if configPath == "" {
		return filepath.Join("logs", "audit")
	}
	return filepath.Join(filepath.Dir(configPath), "logs", "audit")
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	s.router = gin.New()
	s.router.Use(gin.Recovery())
	s.router.HandleMethodNotAllowed = false
	s.router.RedirectFixedPath = false

	// CORS中间件
	if s.config.CORS.Enabled {
		s.router.Use(s.corsMiddleware())
	}

	// API路由组
	api := s.router.Group("/api")
	{
		// 公开API（无需认证）
		api.GET("/status", s.handleGetStatus)
		api.GET("/accounts", s.handleGetAccounts)
		api.GET("/accounts/:id", s.handleGetAccount)
		api.GET("/instances", s.handleGetInstances)
		api.GET("/instances/:id", s.handleGetInstance)
		api.GET("/resources", s.handleGetResources)

		// 认证API
		api.POST("/auth/microsoft/init", s.handleMicrosoftAuthInit)
		api.POST("/auth/microsoft/poll", s.handleMicrosoftAuthPoll)

		// 写操作API（无认证）
		api.POST("/instances", s.handleCreateInstance)
		api.POST("/instances/:id/start", s.handleStartInstance)
		api.POST("/instances/:id/stop", s.handleStopInstance)
		api.POST("/instances/:id/restart", s.handleRestartInstance)
		api.DELETE("/instances/:id", s.handleDeleteInstance)
		api.POST("/accounts", s.handleCreateAccount)
		api.DELETE("/accounts/:id", s.handleDeleteAccount)

		// 日志API
		api.GET("/logs/operations", s.handleGetOperationLogs)
	}

	s.router.NoRoute(s.handleNoRoute)
}

func (s *Server) handleNoRoute(c *gin.Context) {
	requestPath := c.Request.URL.Path
	if isAPIPath(requestPath) {
		c.JSON(http.StatusNotFound, webtypes.OperationResponse{
			Success: false,
			Error:   "API endpoint not found",
		})
		return
	}

	if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
		c.String(http.StatusNotFound, "Not Found")
		return
	}

	if s.uiAssets != nil {
		file, err := s.uiAssets.LookupAsset(requestPath)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if file != nil {
			serveAsset(c, file)
			return
		}

		if s.uiAssets.IsAssetLikePath(requestPath) {
			c.String(http.StatusNotFound, "Not Found")
			return
		}

		if !s.uiAssets.HasIndex() {
			serveFrontendUnavailable(c)
			return
		}

		file, err = s.uiAssets.OpenIndex()
		if err != nil || file == nil {
			c.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		serveAsset(c, file)
		return
	}

	serveFrontendUnavailable(c)
}

func isAPIPath(requestPath string) bool {
	return requestPath == "/api" || requestPath == "/api/" || strings.HasPrefix(requestPath, "/api/")
}

func serveAsset(c *gin.Context, file *webui.AssetFile) {
	if file.ContentType != "" {
		c.Header("Content-Type", file.ContentType)
	}
	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusOK)
		return
	}
	c.Data(http.StatusOK, file.ContentType, file.Content)
}

func serveFrontendUnavailable(c *gin.Context) {
	body := "<html><body><h1>Frontend unavailable</h1><p>前端尚未构建，当前服务仅提供 API。</p></body></html>"
	c.Header("Content-Type", "text/html; charset=utf-8")
	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	c.Data(http.StatusServiceUnavailable, "text/html; charset=utf-8", []byte(body))
}

func (s *Server) HasEmbeddedUI() bool {
	return s.uiAssets != nil && s.uiAssets.HasIndex()
}

// corsMiddleware CORS中间件
func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查是否在允许的源列表中
		allowed := false
		for _, o := range s.config.CORS.Origins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Run 启动服务器
func (s *Server) Run() error {
	s.httpServer = &http.Server{
		Addr:    s.config.Bind,
		Handler: s.router,
	}

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 在goroutine中启动服务器
	go func() {
		logx.Infof("Web服务器启动: http://%s", s.config.Bind)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logx.Errorf("Web服务器错误: %v", err)
		}
	}()

	// 等待退出信号
	<-quit
	logx.Infof("Web服务器关闭中...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}

// logOperation 记录操作日志
func (s *Server) logOperation(c *gin.Context, action string, targetInstanceID, targetAccountID string, success bool, errMsg string) {
	log := &webtypes.OperationLog{
		Timestamp:        time.Now(),
		Action:           action,
		TargetInstanceID: targetInstanceID,
		TargetAccountID:  targetAccountID,
		Success:          success,
		ErrorMsg:         errMsg,
		ClientIP:         c.ClientIP(),
		UserAgent:        c.Request.UserAgent(),
	}

	if err := s.auditLogger.Log(log); err != nil {
		logx.Warnf("记录操作日志失败: %v", err)
	}
}
