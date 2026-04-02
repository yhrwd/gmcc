package web

import (
	"time"
)

// AccountView 账号展示模型（公开）
type AccountView struct {
	ID             string    `json:"id"`
	PlayerID       string    `json:"player_id"`
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

// ClusterStatus 集群状态
type ClusterStatus struct {
	Status           string        `json:"cluster_status"`
	TotalInstances   int           `json:"total_instances"`
	RunningInstances int           `json:"running_instances"`
	Uptime           time.Duration `json:"uptime"`
}

// AuthVerifyRequest 密码验证请求
type AuthVerifyRequest struct {
	Password string `json:"password" binding:"required"`
	Action   string `json:"action,omitempty"`
	Target   string `json:"target,omitempty"`
}

// AuthVerifyResponse 密码验证响应
type AuthVerifyResponse struct {
	Success    bool      `json:"success"`
	Token      string    `json:"token,omitempty"`
	PasswordID string    `json:"password_id,omitempty"`
	ExpiresAt  time.Time `json:"expires_at,omitempty"`
	Error      string    `json:"error,omitempty"`
}

// OperationRequest 操作请求（受保护API）
type OperationRequest struct {
	Password string `json:"password" binding:"required"`
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
	PasswordID       string    `json:"password_id" yaml:"password_id"`
	Action           string    `json:"action" yaml:"action"`
	TargetInstanceID string    `json:"target_instance_id,omitempty" yaml:"target_instance_id,omitempty"`
	TargetAccountID  string    `json:"target_account_id,omitempty" yaml:"target_account_id,omitempty"`
	Details          string    `json:"details,omitempty" yaml:"details,omitempty"`
	Success          bool      `json:"success" yaml:"success"`
	ErrorMsg         string    `json:"error_msg,omitempty" yaml:"error_msg,omitempty"`
	ClientIP         string    `json:"client_ip" yaml:"client_ip"`
	UserAgent        string    `json:"user_agent,omitempty" yaml:"user_agent,omitempty"`
}

// EncryptedToken 加密后的Token结构
type EncryptedToken struct {
	Version    int       `json:"version"`
	Algorithm  string    `json:"algorithm"`
	KDF        string    `json:"kdf"`
	ScryptN    int       `json:"scrypt_n"`
	ScryptR    int       `json:"scrypt_r"`
	ScryptP    int       `json:"scrypt_p"`
	Salt       []byte    `json:"salt"`
	Nonce      []byte    `json:"nonce"`
	Ciphertext []byte    `json:"ciphertext"`
	PlayerID   string    `json:"player_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// MicrosoftAuthInitResponse Microsoft认证初始化响应
type MicrosoftAuthInitResponse struct {
	Success                 bool   `json:"success"`
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

// MicrosoftAuthPollRequest Microsoft认证轮询请求
type MicrosoftAuthPollRequest struct {
	Password   string `json:"password" binding:"required"`
	DeviceCode string `json:"device_code" binding:"required"`
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
	Password        string `json:"password" binding:"required"`
	ID              string `json:"id" binding:"required"`
	PlayerID        string `json:"player_id" binding:"required"`
	ServerAddress   string `json:"server_address" binding:"required"`
	UseOfficialAuth bool   `json:"use_official_auth"`
}
