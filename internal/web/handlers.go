package web

import (
	"time"

	"github.com/gin-gonic/gin"

	"gmcc/internal/cluster"
	"gmcc/internal/logx"
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
	if s.clusterManager == nil {
		c.JSON(503, gin.H{"error": "cluster manager not initialized"})
		return
	}

	instances := s.clusterManager.ListInstances()
	accounts := make([]AccountView, 0, len(instances))

	for _, inst := range instances {
		view := AccountView{
			ID:             inst.ID,
			PlayerID:       inst.PlayerID,
			ServerAddress:  inst.ServerAddress,
			Status:         string(inst.Status),
			OnlineDuration: inst.OnlineDuration.String(),
			LastSeen:       inst.LastActive,
			HasToken:       s.tokenVault.Exists(inst.PlayerID),
		}

		// 如果实例正在运行，添加详细信息
		if inst.Status == "running" {
			view.Health = inst.Health
			view.Food = inst.Food
			if inst.Position != nil {
				view.Position = &Position{
					X: inst.Position.X,
					Y: inst.Position.Y,
					Z: inst.Position.Z,
				}
			}
		}

		accounts = append(accounts, view)
	}

	c.JSON(200, gin.H{"accounts": accounts})
}

// handleGetAccount 获取单个账号详情
func (s *Server) handleGetAccount(c *gin.Context) {
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

	view := AccountView{
		ID:             info.ID,
		PlayerID:       info.PlayerID,
		ServerAddress:  info.ServerAddress,
		Status:         string(info.Status),
		OnlineDuration: info.OnlineDuration.String(),
		LastSeen:       info.LastActive,
		HasToken:       s.tokenVault.Exists(info.PlayerID),
		Health:         info.Health,
		Food:           info.Food,
	}

	if info.Position != nil {
		view.Position = &Position{
			X: info.Position.X,
			Y: info.Position.Y,
			Z: info.Position.Z,
		}
	}

	c.JSON(200, view)
}

// handleAuthVerify 验证密码
func (s *Server) handleAuthVerify(c *gin.Context) {
	var req AuthVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, AuthVerifyResponse{
			Success: false,
			Error:   "invalid request",
		})
		return
	}

	// 验证密码
	passwordID, err := s.authManager.VerifyPassword(req.Password)
	if err != nil {
		// 记录失败的验证
		s.logOperation(c, "auth_verify", "", "", false, err.Error())
		c.JSON(401, AuthVerifyResponse{
			Success: false,
			Error:   "invalid password",
		})
		return
	}

	// 生成Token
	token, expiresAt, err := s.authManager.GenerateToken(passwordID)
	if err != nil {
		c.JSON(500, AuthVerifyResponse{
			Success: false,
			Error:   "failed to generate token",
		})
		return
	}

	// 记录成功的验证
	s.logOperation(c, "auth_verify", req.Target, "", true, "")

	c.JSON(200, AuthVerifyResponse{
		Success:    true,
		Token:      token,
		PasswordID: passwordID,
		ExpiresAt:  expiresAt,
	})
}

// handleStartInstance 启动实例
func (s *Server) handleStartInstance(c *gin.Context) {
	id := c.Param("id")

	if s.clusterManager == nil {
		c.JSON(503, OperationResponse{Success: false, Error: "cluster manager not initialized"})
		return
	}

	if err := s.clusterManager.StartInstance(id); err != nil {
		s.logOperation(c, "instance_start", id, "", false, err.Error())
		c.JSON(400, OperationResponse{Success: false, Error: err.Error()})
		return
	}

	s.logOperation(c, "instance_start", id, "", true, "")
	c.JSON(200, OperationResponse{Success: true, Message: "Instance started"})
}

// handleStopInstance 停止实例
func (s *Server) handleStopInstance(c *gin.Context) {
	id := c.Param("id")

	if s.clusterManager == nil {
		c.JSON(503, OperationResponse{Success: false, Error: "cluster manager not initialized"})
		return
	}

	if err := s.clusterManager.StopInstance(id); err != nil {
		s.logOperation(c, "instance_stop", id, "", false, err.Error())
		c.JSON(400, OperationResponse{Success: false, Error: err.Error()})
		return
	}

	s.logOperation(c, "instance_stop", id, "", true, "")
	c.JSON(200, OperationResponse{Success: true, Message: "Instance stopped"})
}

// handleRestartInstance 重启实例
func (s *Server) handleRestartInstance(c *gin.Context) {
	id := c.Param("id")

	if s.clusterManager == nil {
		c.JSON(503, OperationResponse{Success: false, Error: "cluster manager not initialized"})
		return
	}

	if err := s.clusterManager.RestartInstance(id); err != nil {
		s.logOperation(c, "instance_restart", id, "", false, err.Error())
		c.JSON(400, OperationResponse{Success: false, Error: err.Error()})
		return
	}

	s.logOperation(c, "instance_restart", id, "", true, "")
	c.JSON(200, OperationResponse{Success: true, Message: "Instance restarted"})
}

// handleDeleteInstance 删除实例
func (s *Server) handleDeleteInstance(c *gin.Context) {
	id := c.Param("id")

	if s.clusterManager == nil {
		c.JSON(503, OperationResponse{Success: false, Error: "cluster manager not initialized"})
		return
	}

	if err := s.clusterManager.DeleteInstance(id); err != nil {
		s.logOperation(c, "instance_delete", id, "", false, err.Error())
		c.JSON(400, OperationResponse{Success: false, Error: err.Error()})
		return
	}

	s.logOperation(c, "instance_delete", id, "", true, "")
	c.JSON(200, OperationResponse{Success: true, Message: "Instance deleted"})
}

// handleCreateAccount 创建账号
func (s *Server) handleCreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, OperationResponse{Success: false, Error: "invalid request"})
		return
	}

	if s.clusterManager == nil {
		c.JSON(503, OperationResponse{Success: false, Error: "cluster manager not initialized"})
		return
	}

	// 检查实例是否已存在
	if s.clusterManager.InstanceExists(req.ID) {
		c.JSON(400, OperationResponse{Success: false, Error: "instance already exists"})
		return
	}

	// 创建账号条目
	account := cluster.AccountEntry{
		ID:              req.ID,
		PlayerID:        req.PlayerID,
		UseOfficialAuth: req.UseOfficialAuth,
		ServerAddress:   req.ServerAddress,
		Enabled:         true,
	}

	// 创建并启动实例
	if err := s.clusterManager.CreateInstance(req.ID, account); err != nil {
		s.logOperation(c, "account_create", "", req.ID, false, err.Error())
		c.JSON(400, OperationResponse{Success: false, Error: err.Error()})
		return
	}

	// 启动实例
	if err := s.clusterManager.StartInstance(req.ID); err != nil {
		logx.Warnf("自动启动实例失败 %s: %v", req.ID, err)
	}

	// 保存配置
	if err := s.clusterManager.SaveConfig(); err != nil {
		logx.Warnf("保存配置失败: %v", err)
	}

	s.logOperation(c, "account_create", "", req.ID, true, "")
	c.JSON(200, OperationResponse{Success: true, Message: "Account created"})
}

// handleDeleteAccount 删除账号
func (s *Server) handleDeleteAccount(c *gin.Context) {
	id := c.Param("id")

	if s.clusterManager == nil {
		c.JSON(503, OperationResponse{Success: false, Error: "cluster manager not initialized"})
		return
	}

	// 先删除实例
	if err := s.clusterManager.DeleteInstance(id); err != nil {
		s.logOperation(c, "account_delete", "", id, false, err.Error())
		c.JSON(400, OperationResponse{Success: false, Error: err.Error()})
		return
	}

	// 从集群配置中移除账号
	if err := s.clusterManager.RemoveAccount(id); err != nil {
		s.logOperation(c, "account_delete", "", id, false, err.Error())
		c.JSON(400, OperationResponse{Success: false, Error: err.Error()})
		return
	}

	// 保存配置
	if err := s.clusterManager.SaveConfig(); err != nil {
		logx.Warnf("保存配置失败: %v", err)
	}

	s.logOperation(c, "account_delete", "", id, true, "")
	c.JSON(200, OperationResponse{Success: true, Message: "Account deleted"})
}

// handleMicrosoftAuthInit 初始化Microsoft认证
func (s *Server) handleMicrosoftAuthInit(c *gin.Context) {
	// 首先验证密码
	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(401, MicrosoftAuthInitResponse{Success: false})
		return
	}

	// 验证密码
	_, err := s.authManager.VerifyPassword(req.Password)
	if err != nil {
		c.JSON(401, MicrosoftAuthInitResponse{Success: false})
		return
	}

	// TODO: 实现Microsoft OAuth Device Flow
	// 这里需要调用Microsoft OAuth API获取device_code和user_code
	// 目前返回模拟数据
	c.JSON(200, MicrosoftAuthInitResponse{
		Success:                 true,
		DeviceCode:              "mock-device-code",
		UserCode:                "ABCD-EFGH",
		VerificationURI:         "https://www.microsoft.com/link",
		VerificationURIComplete: "https://www.microsoft.com/link?code=ABCD-EFGH",
		ExpiresIn:               900,
		Interval:                5,
	})
}

// handleMicrosoftAuthPoll 轮询Microsoft认证
func (s *Server) handleMicrosoftAuthPoll(c *gin.Context) {
	// 验证密码
	var req MicrosoftAuthPollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(401, MicrosoftAuthPollResponse{Success: false, Status: "error", Message: "password required"})
		return
	}

	// 验证密码
	_, err := s.authManager.VerifyPassword(req.Password)
	if err != nil {
		c.JSON(401, MicrosoftAuthPollResponse{Success: false, Status: "error", Message: "invalid password"})
		return
	}

	// TODO: 实现Microsoft OAuth Device Flow轮询
	// 这里需要调用Microsoft OAuth API轮询token
	// 目前返回pending状态
	c.JSON(200, MicrosoftAuthPollResponse{
		Success: true,
		Status:  "pending",
		Message: "Waiting for user authorization...",
	})
}

// handleGetOperationLogs 获取操作日志
func (s *Server) handleGetOperationLogs(c *gin.Context) {
	// 解析查询参数
	startStr := c.DefaultQuery("start", time.Now().AddDate(0, 0, -7).Format(time.RFC3339))
	endStr := c.DefaultQuery("end", time.Now().Format(time.RFC3339))
	passwordID := c.Query("password_id")

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

	logs, err := s.auditLogger.Query(start, end, passwordID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to query logs"})
		return
	}

	c.JSON(200, gin.H{"logs": logs})
}

// handleCreatePassword 创建密码
func (s *Server) handleCreatePassword(c *gin.Context) {
	var req struct {
		Password string `json:"password" binding:"required"`
		ID       string `json:"id" binding:"required"`
		Note     string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, OperationResponse{Success: false, Error: "invalid request"})
		return
	}

	// 生成bcrypt hash
	hash, err := s.authManager.HashPassword(req.Password)
	if err != nil {
		c.JSON(500, OperationResponse{Success: false, Error: "failed to hash password"})
		return
	}

	// 添加到配置
	entry := webtypes.PasswordEntry{
		ID:        req.ID,
		Hash:      hash,
		Enabled:   true,
		CreatedAt: time.Now(),
		Note:      req.Note,
	}

	s.config.Auth.Passwords = append(s.config.Auth.Passwords, entry)

	// 保存配置
	if s.configPath != "" {
		if err := webtypes.SaveWebConfig(s.configPath, s.config); err != nil {
			logx.Warnf("保存Web配置失败: %v", err)
		}
	}

	s.logOperation(c, "password_create", "", req.ID, true, "")
	c.JSON(200, OperationResponse{Success: true, Message: "Password created"})
}

// handleDeletePassword 删除密码
func (s *Server) handleDeletePassword(c *gin.Context) {
	id := c.Param("id")

	// 查找并删除密码
	found := false
	for i, pwd := range s.config.Auth.Passwords {
		if pwd.ID == id {
			s.config.Auth.Passwords = append(s.config.Auth.Passwords[:i], s.config.Auth.Passwords[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		c.JSON(404, OperationResponse{Success: false, Error: "password not found"})
		return
	}

	// 保存配置
	if s.configPath != "" {
		if err := webtypes.SaveWebConfig(s.configPath, s.config); err != nil {
			logx.Warnf("保存Web配置失败: %v", err)
		}
	}

	s.logOperation(c, "password_delete", "", id, true, "")
	c.JSON(200, OperationResponse{Success: true, Message: "Password deleted"})
}
