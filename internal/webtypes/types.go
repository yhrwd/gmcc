package webtypes

import "time"

// AccountView 账号展示模型（公开）
type AccountView struct {
	ID         string `json:"id"`
	PlayerID   string `json:"player_id,omitempty"`
	Enabled    bool   `json:"enabled"`
	Label      string `json:"label,omitempty"`
	Note       string `json:"note,omitempty"`
	AuthStatus string `json:"auth_status"`
	HasToken   bool   `json:"has_token"`
}

// InstanceView 实例展示模型（公开）
type InstanceView struct {
	ID             string    `json:"id"`
	AccountID      string    `json:"account_id"`
	PlayerID       string    `json:"player_id,omitempty"`
	ServerAddress  string    `json:"server_address"`
	Status         string    `json:"status"`
	OnlineDuration string    `json:"online_duration"`
	LastSeen       time.Time `json:"last_seen"`
	HasToken       bool      `json:"has_token"`
	Health         float32   `json:"health,omitempty"`
	Food           int32     `json:"food,omitempty"`
	Position       *Position `json:"position,omitempty"`
}

// Position 位置信息
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// OperationResponse 操作响应
type OperationResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	OperationID string `json:"operation_id,omitempty"`
	Error       string `json:"error,omitempty"`
}

// OperationLog 操作日志
type OperationLog struct {
	ID               string    `json:"id" yaml:"id"`
	Timestamp        time.Time `json:"timestamp" yaml:"timestamp"`
	Action           string    `json:"action" yaml:"action"`
	TargetInstanceID string    `json:"target_instance_id,omitempty" yaml:"target_instance_id,omitempty"`
	TargetAccountID  string    `json:"target_account_id,omitempty" yaml:"target_account_id,omitempty"`
	Details          string    `json:"details,omitempty" yaml:"details,omitempty"`
	Success          bool      `json:"success" yaml:"success"`
	ErrorMsg         string    `json:"error_msg,omitempty" yaml:"error_msg,omitempty"`
	ClientIP         string    `json:"client_ip" yaml:"client_ip"`
	UserAgent        string    `json:"user_agent,omitempty" yaml:"user_agent,omitempty"`
}

// MicrosoftAuthInitResponse Microsoft认证初始化响应
type MicrosoftAuthInitRequest struct {
	AccountID string `json:"account_id" binding:"required"`
}

// MicrosoftAuthInitResponse Microsoft认证初始化响应
type MicrosoftAuthInitResponse struct {
	Success                 bool   `json:"success"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
	AccountID               string `json:"account_id,omitempty"`
}

// MicrosoftAuthPollRequest Microsoft认证轮询请求
type MicrosoftAuthPollRequest struct {
	AccountID string `json:"account_id" binding:"required"`
}

// MicrosoftAuthPollResponse Microsoft认证轮询响应
type MicrosoftAuthPollResponse struct {
	Success          bool              `json:"success"`
	Status           string            `json:"status"`
	Message          string            `json:"message"`
	MinecraftProfile *MinecraftProfile `json:"minecraft_profile,omitempty"`
	AccountID        string            `json:"account_id,omitempty"`
}

// MinecraftProfile Minecraft档案
type MinecraftProfile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateAccountRequest 创建账号请求
type CreateAccountRequest struct {
	ID    string `json:"id" binding:"required"`
	Label string `json:"label,omitempty"`
	Note  string `json:"note,omitempty"`
}

// CreateInstanceRequest 创建实例请求
type CreateInstanceRequest struct {
	ID            string `json:"id" binding:"required"`
	AccountID     string `json:"account_id" binding:"required"`
	ServerAddress string `json:"server_address" binding:"required"`
	Enabled       *bool  `json:"enabled,omitempty"`
	AutoStart     bool   `json:"auto_start"`
}

// WebConfig Web面板配置
type WebConfig struct {
	Bind string     `yaml:"bind"`
	Auth AuthConfig `yaml:"auth"`
	CORS CORSConfig `yaml:"cors"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	AuditLogRetentionDays int `yaml:"audit_log_retention_days"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	Enabled bool     `yaml:"enabled"`
	Origins []string `yaml:"origins"`
}
