package web

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	authsession "gmcc/internal/auth/session"
	"gmcc/internal/cluster"
	"gmcc/internal/logx"
	"gmcc/internal/resource"
	"gmcc/internal/webtypes"
)

// handleGetStatus 获取集群状态
func (s *Server) handleGetStatus(c *gin.Context) {
	if s.clusterManager == nil {
		c.JSON(503, gin.H{"error": "cluster manager not initialized"})
		return
	}

	status := s.clusterManager.GetClusterStatus()
	c.JSON(200, status)
}

// handleGetAccounts 获取账号列表
func (s *Server) handleGetAccounts(c *gin.Context) {
	if s.resourceManager == nil {
		c.JSON(503, gin.H{"error": "resource manager not initialized"})
		return
	}

	records, err := s.resourceManager.ListAccounts()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to load accounts"})
		return
	}

	accounts := make([]webtypes.AccountView, 0, len(records))
	for _, record := range records {
		accounts = append(accounts, s.buildAccountView(record))
	}

	c.JSON(200, gin.H{"accounts": accounts})
}

// handleGetAccount 获取单个账号详情
func (s *Server) handleGetAccount(c *gin.Context) {
	id := c.Param("id")

	if s.resourceManager == nil {
		c.JSON(503, gin.H{"error": "resource manager not initialized"})
		return
	}

	record, err := s.resourceManager.GetAccount(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "account not found"})
		return
	}

	c.JSON(200, s.buildAccountView(record))
}

// handleGetInstances 获取实例列表
func (s *Server) handleGetInstances(c *gin.Context) {
	if s.clusterManager == nil {
		c.JSON(503, gin.H{"error": "cluster manager not initialized"})
		return
	}

	instances := s.clusterManager.ListInstances()
	views := make([]webtypes.InstanceView, 0, len(instances))
	for _, inst := range instances {
		views = append(views, s.buildInstanceView(inst))
	}

	c.JSON(200, gin.H{"instances": views})
}

// handleGetInstance 获取单个实例详情
func (s *Server) handleGetInstance(c *gin.Context) {
	id := c.Param("id")

	if s.clusterManager == nil {
		c.JSON(503, gin.H{"error": "cluster manager not initialized"})
		return
	}

	info, err := s.clusterManager.GetInstanceInfo(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "instance not found"})
		return
	}

	c.JSON(200, s.buildInstanceView(info))
}

// handleStartInstance 启动实例
func (s *Server) handleStartInstance(c *gin.Context) {
	id := c.Param("id")

	if s.clusterManager == nil {
		c.JSON(503, webtypes.OperationResponse{Success: false, Error: "cluster manager not initialized"})
		return
	}

	if err := s.clusterManager.StartInstance(id); err != nil {
		s.logOperation(c, "instance_start", id, "", false, err.Error())
		c.JSON(400, webtypes.OperationResponse{Success: false, Error: err.Error()})
		return
	}

	s.logOperation(c, "instance_start", id, "", true, "")
	c.JSON(200, webtypes.OperationResponse{Success: true, Message: "Instance started"})
}

// handleStopInstance 停止实例
func (s *Server) handleStopInstance(c *gin.Context) {
	id := c.Param("id")

	if s.clusterManager == nil {
		c.JSON(503, webtypes.OperationResponse{Success: false, Error: "cluster manager not initialized"})
		return
	}

	if err := s.clusterManager.StopInstance(id); err != nil {
		s.logOperation(c, "instance_stop", id, "", false, err.Error())
		c.JSON(400, webtypes.OperationResponse{Success: false, Error: err.Error()})
		return
	}

	s.logOperation(c, "instance_stop", id, "", true, "")
	c.JSON(200, webtypes.OperationResponse{Success: true, Message: "Instance stopped"})
}

// handleRestartInstance 重启实例
func (s *Server) handleRestartInstance(c *gin.Context) {
	id := c.Param("id")

	if s.clusterManager == nil {
		c.JSON(503, webtypes.OperationResponse{Success: false, Error: "cluster manager not initialized"})
		return
	}

	if err := s.clusterManager.RestartInstance(id); err != nil {
		s.logOperation(c, "instance_restart", id, "", false, err.Error())
		c.JSON(400, webtypes.OperationResponse{Success: false, Error: err.Error()})
		return
	}

	s.logOperation(c, "instance_restart", id, "", true, "")
	c.JSON(200, webtypes.OperationResponse{Success: true, Message: "Instance restarted"})
}

// handleDeleteInstance 删除实例
func (s *Server) handleDeleteInstance(c *gin.Context) {
	id := c.Param("id")

	if s.clusterManager == nil {
		c.JSON(503, webtypes.OperationResponse{Success: false, Error: "cluster manager not initialized"})
		return
	}

	if err := s.clusterManager.DeleteInstance(id); err != nil {
		s.logOperation(c, "instance_delete", id, "", false, err.Error())
		c.JSON(400, webtypes.OperationResponse{Success: false, Error: err.Error()})
		return
	}

	s.logOperation(c, "instance_delete", id, "", true, "")
	c.JSON(200, webtypes.OperationResponse{Success: true, Message: "Instance deleted"})
}

// handleCreateAccount 创建账号
func (s *Server) handleCreateAccount(c *gin.Context) {
	var req webtypes.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, webtypes.OperationResponse{Success: false, Error: "invalid request"})
		return
	}

	if s.resourceManager == nil {
		c.JSON(503, webtypes.OperationResponse{Success: false, Error: "resource manager not initialized"})
		return
	}

	if _, err := s.resourceManager.CreateAccount(resource.CreateAccountInput{
		AccountID: req.ID,
		Enabled:   true,
		Label:     req.Label,
		Note:      req.Note,
	}); err != nil {
		s.logOperation(c, "account_create", "", req.ID, false, err.Error())
		c.JSON(400, webtypes.OperationResponse{Success: false, Error: err.Error()})
		return
	}

	s.logOperation(c, "account_create", "", req.ID, true, "")
	c.JSON(200, webtypes.OperationResponse{Success: true, Message: "Account created"})
}

// handleCreateInstance 创建实例
func (s *Server) handleCreateInstance(c *gin.Context) {
	var req webtypes.CreateInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, webtypes.OperationResponse{Success: false, Error: "invalid request"})
		return
	}

	if s.clusterManager == nil {
		c.JSON(503, webtypes.OperationResponse{Success: false, Error: "cluster manager not initialized"})
		return
	}

	account := cluster.AccountEntry{
		ID:            req.AccountID,
		ServerAddress: req.ServerAddress,
		Enabled:       true,
	}
	if req.Enabled != nil {
		account.Enabled = *req.Enabled
	}
	if req.AutoStart && !account.Enabled {
		s.logOperation(c, "instance_create", req.ID, req.AccountID, false, "disabled instance cannot auto start")
		c.JSON(400, webtypes.OperationResponse{Success: false, Error: "disabled instance cannot auto start"})
		return
	}

	if err := s.clusterManager.CreateInstance(req.ID, account); err != nil {
		s.logOperation(c, "instance_create", req.ID, req.AccountID, false, err.Error())
		c.JSON(400, webtypes.OperationResponse{Success: false, Error: err.Error()})
		return
	}

	if req.AutoStart {
		if err := s.clusterManager.StartInstance(req.ID); err != nil {
			s.logOperation(c, "instance_create", req.ID, req.AccountID, false, err.Error())
			c.JSON(400, webtypes.OperationResponse{Success: false, Error: err.Error()})
			return
		}
	}

	s.logOperation(c, "instance_create", req.ID, req.AccountID, true, "")
	c.JSON(200, webtypes.OperationResponse{Success: true, Message: "Instance created"})
}

// handleDeleteAccount 删除账号
func (s *Server) handleDeleteAccount(c *gin.Context) {
	id := c.Param("id")

	if s.resourceManager == nil {
		c.JSON(503, webtypes.OperationResponse{Success: false, Error: "resource manager not initialized"})
		return
	}

	if err := s.resourceManager.DeleteAccount(id); err != nil {
		s.logOperation(c, "account_delete", "", id, false, err.Error())
		c.JSON(400, webtypes.OperationResponse{Success: false, Error: err.Error()})
		return
	}

	s.logOperation(c, "account_delete", "", id, true, "")
	c.JSON(200, webtypes.OperationResponse{Success: true, Message: "Account deleted"})
}

// handleMicrosoftAuthInit 初始化Microsoft认证
func (s *Server) handleMicrosoftAuthInit(c *gin.Context) {
	if s.runtimeAuth == nil {
		c.JSON(503, webtypes.MicrosoftAuthInitResponse{Success: false})
		return
	}

	var req webtypes.MicrosoftAuthInitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(401, webtypes.MicrosoftAuthInitResponse{Success: false})
		return
	}

	info, err := s.runtimeAuth.BeginDeviceLogin(context.Background(), strings.TrimSpace(req.AccountID))
	if err != nil {
		s.logOperation(c, "microsoft_auth_init", "", strings.TrimSpace(req.AccountID), false, err.Error())
		c.JSON(500, webtypes.MicrosoftAuthInitResponse{Success: false})
		return
	}

	expiresIn := 0
	if !info.ExpiresAt.IsZero() {
		expiresIn = int(time.Until(info.ExpiresAt).Seconds())
		if expiresIn < 0 {
			expiresIn = 0
		}
	}
	interval := int(info.PollInterval.Seconds())
	if interval <= 0 {
		interval = 1
	}

	verificationURIComplete := info.VerificationURI
	if verificationURIComplete != "" && info.UserCode != "" && !strings.Contains(verificationURIComplete, info.UserCode) {
		verificationURIComplete = verificationURIComplete + "?code=" + info.UserCode
	}

	s.logOperation(c, "microsoft_auth_init", "", strings.TrimSpace(req.AccountID), true, "")
	c.JSON(200, webtypes.MicrosoftAuthInitResponse{
		Success:                 true,
		UserCode:                info.UserCode,
		VerificationURI:         info.VerificationURI,
		VerificationURIComplete: verificationURIComplete,
		ExpiresIn:               expiresIn,
		Interval:                interval,
		AccountID:               strings.TrimSpace(req.AccountID),
	})
}

// handleMicrosoftAuthPoll 轮询Microsoft认证
func (s *Server) handleMicrosoftAuthPoll(c *gin.Context) {
	if s.runtimeAuth == nil {
		c.JSON(503, webtypes.MicrosoftAuthPollResponse{Success: false, Status: "error", Message: "runtime auth not initialized"})
		return
	}

	var req webtypes.MicrosoftAuthPollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(401, webtypes.MicrosoftAuthPollResponse{Success: false, Status: "error", Message: "account_id required"})
		return
	}

	accountID := strings.TrimSpace(req.AccountID)
	status, authSession, statusErr := s.runtimeAuth.GetDeviceLoginStatus(accountID)

	resp := webtypes.MicrosoftAuthPollResponse{
		Success:   statusErr == nil || status == authsession.DeviceLoginStatusPending || status == authsession.DeviceLoginStatusSucceeded,
		Status:    string(status),
		Message:   microsoftPollMessage(status, statusErr),
		AccountID: accountID,
	}

	if authSession != nil {
		resp.MinecraftProfile = &webtypes.MinecraftProfile{
			ID:   authSession.ProfileID,
			Name: authSession.ProfileName,
		}
	}

	if status == authsession.DeviceLoginStatusSucceeded {
		s.logOperation(c, "microsoft_auth_poll", "", accountID, true, "")
		c.JSON(200, resp)
		return
	}

	if statusErr != nil && status != authsession.DeviceLoginStatusPending {
		s.logOperation(c, "microsoft_auth_poll", "", accountID, false, statusErr.Error())
		c.JSON(200, resp)
		return
	}

	c.JSON(200, resp)
}

func microsoftPollMessage(status authsession.DeviceLoginStatus, err error) string {
	if err != nil {
		if status == authsession.DeviceLoginStatusPending {
			return "Waiting for user authorization..."
		}
		return err.Error()
	}

	switch status {
	case authsession.DeviceLoginStatusSucceeded:
		return "Microsoft authentication succeeded"
	case authsession.DeviceLoginStatusPending:
		return "Waiting for user authorization..."
	case authsession.DeviceLoginStatusExpired:
		return "Device login expired"
	case authsession.DeviceLoginStatusCancelled:
		return "Device login cancelled"
	case authsession.DeviceLoginStatusFailed:
		return "Microsoft authentication failed"
	default:
		return "Waiting for user authorization..."
	}
}

func (s *Server) buildAccountView(record resource.AccountRecord) webtypes.AccountView {
	accountID := record.Meta.AccountID
	status := s.accountAuthStatus(accountID)

	return webtypes.AccountView{
		ID:         accountID,
		PlayerID:   s.accountProfileID(accountID),
		Enabled:    record.Meta.Enabled,
		Label:      record.Meta.Label,
		Note:       record.Meta.Note,
		AuthStatus: status,
		HasToken:   status == string(authsession.AccountAuthStatusLoggedIn) || status == string(authsession.AccountAuthStatusAuthInvalid),
	}
}

func (s *Server) buildInstanceView(info cluster.InstanceInfo) webtypes.InstanceView {
	view := webtypes.InstanceView{
		ID:             info.ID,
		AccountID:      info.AccountID,
		PlayerID:       info.PlayerID,
		ServerAddress:  info.ServerAddress,
		Status:         string(info.Status),
		OnlineDuration: info.OnlineDuration.String(),
		LastSeen:       info.LastActive,
		HasToken:       s.accountHasToken(info.AccountID),
		Health:         info.Health,
		Food:           info.Food,
	}

	if info.Position != nil {
		view.Position = &webtypes.Position{
			X: info.Position.X,
			Y: info.Position.Y,
			Z: info.Position.Z,
		}
	}

	return view
}

func (s *Server) accountProfileID(accountID string) string {
	if s.runtimeAuth == nil {
		return ""
	}

	profile, err := s.runtimeAuth.GetAccountProfile(accountID)
	if err != nil {
		if !errors.Is(err, authsession.ErrAccountNotFound) {
			logx.Warnf("查询账号档案失败 %s: %v", accountID, err)
		}
		return ""
	}
	return profile.ProfileID
}

func (s *Server) accountAuthStatus(accountID string) string {
	if s.runtimeAuth == nil {
		return string(authsession.AccountAuthStatusNotLoggedIn)
	}

	status, err := s.runtimeAuth.GetAccountAuthStatus(accountID)
	if err != nil {
		if errors.Is(err, authsession.ErrAccountNotFound) {
			return string(authsession.AccountAuthStatusNotLoggedIn)
		}
		logx.Warnf("查询账号认证状态失败 %s: %v", accountID, err)
		return string(authsession.AccountAuthStatusNotLoggedIn)
	}

	return string(status)
}

func (s *Server) accountHasToken(accountID string) bool {
	status := s.accountAuthStatus(accountID)
	return status == string(authsession.AccountAuthStatusLoggedIn) || status == string(authsession.AccountAuthStatusAuthInvalid)
}

// handleGetOperationLogs 获取操作日志
func (s *Server) handleGetOperationLogs(c *gin.Context) {
	// 解析查询参数
	startStr := c.DefaultQuery("start", time.Now().AddDate(0, 0, -7).Format(time.RFC3339))
	endStr := c.DefaultQuery("end", time.Now().Format(time.RFC3339))

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid start date"})
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid end date"})
		return
	}

	logs, err := s.auditLogger.Query(start, end)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to query logs"})
		return
	}

	c.JSON(200, gin.H{"logs": logs})
}
