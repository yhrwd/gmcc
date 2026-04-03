package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"gmcc/internal/cluster"
	"gmcc/internal/logx"
	"gmcc/internal/web/audit"
	"gmcc/internal/web/auth"
	"gmcc/internal/web/key"
	"gmcc/internal/web/vault"
	"gmcc/internal/webtypes"
)

// Server Web服务器
type Server struct {
	config         webtypes.WebConfig
	configPath     string // Web配置文件路径
	router         *gin.Engine
	httpServer     *http.Server
	clusterManager *cluster.Manager
	authManager    *auth.Manager
	tokenVault     *vault.Vault
	auditLogger    *audit.Logger
	keyManager     *key.Manager
}

// NewServer 创建Web服务器
func NewServer(config webtypes.WebConfig, configPath string, clusterManager *cluster.Manager) (*Server, error) {
	// 创建认证管理器
	authManager, err := auth.NewManager(config.Auth)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth manager: %w", err)
	}

	// 创建密钥管理器
	keyManager := key.NewManager()

	// 创建Token Vault
	keyGetter := keyManager.GetKeyGetter(
		config.TokenVault.ScryptN,
		config.TokenVault.ScryptR,
		config.TokenVault.ScryptP,
		32,
	)
	tokenVault := vault.NewVault(config.TokenVault, keyGetter)

	// 创建审计日志管理器
	logDir := "logs/audit"
	auditLogger, err := audit.NewLogger(logDir, config.Auth.AuditLogRetentionDays)
	if err != nil {
		return nil, fmt.Errorf("failed to create audit logger: %w", err)
	}

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin路由
	router := gin.New()
	router.Use(gin.Recovery())

	// 禁用自动重定向
	router.HandleMethodNotAllowed = false
	router.RedirectFixedPath = false

	server := &Server{
		config:         config,
		configPath:     configPath,
		router:         router,
		clusterManager: clusterManager,
		authManager:    authManager,
		tokenVault:     tokenVault,
		auditLogger:    auditLogger,
		keyManager:     keyManager,
	}

	// 设置路由
	server.setupRoutes()

	return server, nil
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
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

		// 认证API
		api.POST("/auth/verify", s.handleAuthVerify)
		api.POST("/auth/microsoft/init", s.handleMicrosoftAuthInit)
		api.POST("/auth/microsoft/poll", s.handleMicrosoftAuthPoll)

		// 受保护API（需要密码验证）
		protected := api.Group("")
		protected.Use(s.passwordAuthMiddleware())
		{
			// 实例操作
			protected.POST("/instances/:id/start", s.handleStartInstance)
			protected.POST("/instances/:id/stop", s.handleStopInstance)
			protected.POST("/instances/:id/restart", s.handleRestartInstance)
			protected.DELETE("/instances/:id", s.handleDeleteInstance)

			// 账号管理
			protected.POST("/accounts", s.handleCreateAccount)
			protected.DELETE("/accounts/:id", s.handleDeleteAccount)

			// 密码管理
			protected.POST("/passwords", s.handleCreatePassword)
			protected.DELETE("/passwords/:id", s.handleDeletePassword)
		}

		// 日志API
		api.GET("/logs/operations", s.passwordAuthMiddleware(), s.handleGetOperationLogs)
	}

	// API 模式 - 只提供 API 服务，不提供前端静态文件
	s.router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// API 返回 404
		if len(path) >= 4 && path[:4] == "/api" {
			c.JSON(404, webtypes.OperationResponse{
				Success: false,
				Error:   "API endpoint not found",
			})
			return
		}

		// 静态文件路径返回 404
		if len(path) >= 8 && path[:8] == "/assets/" {
			c.String(404, "Not Found")
			return
		}

		// API 模式 - 只返回 JSON 错误
		c.JSON(404, gin.H{"error": "API endpoint not found"})
	})
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

// passwordAuthMiddleware 密码认证中间件
func (s *Server) passwordAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req webtypes.OperationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(401, webtypes.OperationResponse{
				Success: false,
				Error:   "password required",
			})
			c.Abort()
			return
		}

		passwordID, err := s.authManager.VerifyPassword(req.Password)
		if err != nil {
			c.JSON(401, webtypes.OperationResponse{
				Success: false,
				Error:   "invalid password",
			})
			c.Abort()
			return
		}

		// 将password_id存入上下文
		c.Set("password_id", passwordID)
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

// GetClusterManager 获取集群管理器
func (s *Server) GetClusterManager() *cluster.Manager {
	return s.clusterManager
}

// GetAuthManager 获取认证管理器
func (s *Server) GetAuthManager() *auth.Manager {
	return s.authManager
}

// GetTokenVault 获取Token Vault
func (s *Server) GetTokenVault() *vault.Vault {
	return s.tokenVault
}

// GetAuditLogger 获取审计日志管理器
func (s *Server) GetAuditLogger() *audit.Logger {
	return s.auditLogger
}

// GetKeyManager 获取密钥管理器
func (s *Server) GetKeyManager() *key.Manager {
	return s.keyManager
}

// logOperation 记录操作日志
func (s *Server) logOperation(c *gin.Context, action string, targetInstanceID, targetAccountID string, success bool, errMsg string) {
	passwordID, _ := c.Get("password_id")
	if passwordID == nil {
		passwordID = ""
	}

	log := &webtypes.OperationLog{
		Timestamp:        time.Now(),
		PasswordID:       passwordID.(string),
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
