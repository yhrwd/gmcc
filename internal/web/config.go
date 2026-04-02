package web

import (
	"time"
)

// WebConfig Web面板配置
type WebConfig struct {
	Bind       string           `yaml:"bind"`
	StaticPath string           `yaml:"static_path"`
	Auth       AuthConfig       `yaml:"auth"`
	TokenVault TokenVaultConfig `yaml:"token_vault"`
	CORS       CORSConfig       `yaml:"cors"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	Passwords             []PasswordEntry `yaml:"passwords"`
	TokenExpiry           time.Duration   `yaml:"token_expiry"`
	AuditLogRetentionDays int             `yaml:"audit_log_retention_days"`
}

// PasswordEntry 密码条目
type PasswordEntry struct {
	ID        string    `yaml:"id"`
	Hash      string    `yaml:"hash"`
	Enabled   bool      `yaml:"enabled"`
	CreatedAt time.Time `yaml:"created_at"`
	Note      string    `yaml:"note,omitempty"`
}

// TokenVaultConfig Token Vault配置
type TokenVaultConfig struct {
	StoragePath string `yaml:"storage_path"`
	ScryptN     int    `yaml:"scrypt_n"`
	ScryptR     int    `yaml:"scrypt_r"`
	ScryptP     int    `yaml:"scrypt_p"`
	SaltLen     int    `yaml:"salt_len"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	Enabled bool     `yaml:"enabled"`
	Origins []string `yaml:"origins"`
}

// DefaultWebConfig 返回默认Web配置
func DefaultWebConfig() WebConfig {
	return WebConfig{
		Bind:       "0.0.0.0:8080",
		StaticPath: "./web",
		Auth: AuthConfig{
			TokenExpiry:           5 * time.Minute,
			AuditLogRetentionDays: 30,
			Passwords:             []PasswordEntry{},
		},
		TokenVault: TokenVaultConfig{
			StoragePath: ".tokens",
			ScryptN:     1 << 20, // 2^20
			ScryptR:     8,
			ScryptP:     1,
			SaltLen:     32,
		},
		CORS: CORSConfig{
			Enabled: true,
			Origins: []string{"http://localhost:5173", "http://localhost:3000"},
		},
	}
}
