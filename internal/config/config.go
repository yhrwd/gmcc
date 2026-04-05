package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	authsession "gmcc/internal/auth/session"

	yaml "gopkg.in/yaml.v3"
)

// Config 统一配置结构
type Config struct {
	// 集群配置
	Cluster ClusterConfig `yaml:"cluster"`

	// Web面板配置
	Web WebConfig `yaml:"web"`

	// 日志配置
	Log LogConfig `yaml:"log"`

	// 向后兼容配置
	Account AccountConfig `yaml:"account"`
	Server  ServerConfig  `yaml:"server"`

	// 集群运行时注入配置（不参与序列化）
	ClusterRuntime ClusterRuntimeConfig `yaml:"-"`
}

// ClusterRuntimeConfig 集群运行时配置（仅进程内使用）
type ClusterRuntimeConfig struct {
	AccountID   string                   `yaml:"-"`
	AuthManager *authsession.AuthManager `yaml:"-"`
}

// ClusterConfig 集群配置（与 internal/cluster/config.go 保持一致）
type ClusterConfig struct {
	Global   GlobalConfig   `yaml:"global"`
	Accounts []AccountEntry `yaml:"accounts"`
}

// GlobalConfig 全局配置
type GlobalConfig struct {
	MaxInstances    int             `yaml:"max_instances"`
	ReconnectPolicy ReconnectPolicy `yaml:"reconnect_policy"`
}

// ReconnectPolicy 重连策略
type ReconnectPolicy struct {
	Enabled    bool          `yaml:"enabled"`
	MaxRetries int           `yaml:"max_retries"`
	BaseDelay  time.Duration `yaml:"base_delay"`
	MaxDelay   time.Duration `yaml:"max_delay"`
	Multiplier float64       `yaml:"multiplier"`
}

// AccountEntry 账号条目
type AccountEntry struct {
	ID              string `yaml:"id"`
	PlayerID        string `yaml:"player_id"`
	UseOfficialAuth bool   `yaml:"use_official_auth"`
	ServerAddress   string `yaml:"server_address"`
	Enabled         bool   `yaml:"enabled"`
}

// WebConfig Web面板配置
type WebConfig struct {
	Bind       string `yaml:"bind"` // 监听地址
	Auth       WebAuthConfig
	TokenVault TokenVaultConfig
	CORS       CORSConfig
}

// WebAuthConfig Web认证配置
type WebAuthConfig struct {
	TokenExpiry           time.Duration   `yaml:"token_expiry"`
	AuditLogRetentionDays int             `yaml:"audit_log_retention_days"`
	Passwords             []PasswordEntry `yaml:"passwords"`
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

// LogConfig 日志配置
type LogConfig struct {
	LogDir     string `yaml:"log_dir"`
	MaxSize    int64  `yaml:"max_size"`
	Debug      bool   `yaml:"debug"`
	EnableFile bool   `yaml:"enable_file"`
}

// AccountConfig 账户配置（向后兼容）
type AccountConfig struct {
	PlayerID        string `yaml:"player_id"`
	UseOfficialAuth bool   `yaml:"use_official_auth"`
}

// ServerConfig 服务器配置（向后兼容）
type ServerConfig struct {
	Address string `yaml:"address"`
}

// MaxSizeInBytes 返回转换为字节的最大日志大小
func (c *LogConfig) MaxSizeInBytes() int64 {
	return c.MaxSize * 1024
}

// Default 返回默认配置
func Default() Config {
	return Config{
		Cluster: ClusterConfig{
			Global: GlobalConfig{
				MaxInstances: 10,
				ReconnectPolicy: ReconnectPolicy{
					Enabled:    true,
					MaxRetries: 0,
					BaseDelay:  2 * time.Second,
					MaxDelay:   2 * time.Minute,
					Multiplier: 1.8,
				},
			},
			Accounts: []AccountEntry{},
		},
		Web: WebConfig{
			Bind: "0.0.0.0:8080",
			Auth: WebAuthConfig{
				TokenExpiry:           5 * time.Minute,
				AuditLogRetentionDays: 30,
				Passwords: []PasswordEntry{
					{
						ID:        "default",
						Hash:      "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
						Enabled:   true,
						CreatedAt: time.Now(),
						Note:      "默认密码: password123",
					},
				},
			},
			TokenVault: TokenVaultConfig{
				StoragePath: ".tokens",
				ScryptN:     1 << 20,
				ScryptR:     8,
				ScryptP:     1,
				SaltLen:     32,
			},
			CORS: CORSConfig{
				Enabled: true,
				Origins: []string{"http://localhost:5173", "http://localhost:3000"},
			},
		},
		Log: LogConfig{
			LogDir:     "logs",
			MaxSize:    512,
			Debug:      false,
			EnableFile: true,
		},
	}
}

// Load 加载配置
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := Default()
			if err := Save(path, cfg); err != nil {
				return nil, fmt.Errorf("创建默认配置失败: %w", err)
			}
			return &cfg, nil
		}
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 填充默认值
	cfg.fillDefaults()

	return &cfg, nil
}

// Save 保存配置
func Save(path string, cfg Config) error {
	data, err := generateConfigWithComments(cfg)
	if err != nil {
		// 如果生成注释失败，回退到普通序列化
		data, err = yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("序列化配置失败: %w", err)
		}
	}
	return os.WriteFile(path, data, 0644)
}

// generateConfigWithComments 生成带注释的配置
func generateConfigWithComments(cfg Config) ([]byte, error) {
	passwordsYAML, err := yaml.Marshal(cfg.Web.Auth.Passwords)
	if err != nil {
		return nil, fmt.Errorf("序列化密码失败: %w", err)
	}

	// 为密码 YAML 添加缩进（前缀 6 个空格）
	indentedPasswords := ""
	for i, line := range strings.Split(strings.TrimSpace(string(passwordsYAML)), "\n") {
		if i > 0 {
			indentedPasswords += "\n"
		}
		indentedPasswords += "      " + line
	}

	return []byte(fmt.Sprintf(`# GMCC 配置文件
# 配置文件路径可通过环境变量 GMCC_CONFIG 指定，默认 config.yaml

cluster:
  # 集群全局配置
  global:
    # 最大实例数量 (0 表示无限制)
    max_instances: %d
    
    # 自动重连策略
    reconnect_policy:
      # 是否启用自动重连
      enabled: %t
      # 最大重试次数 (0 表示无限重试)
      max_retries: %d
      # 初始重连延迟
      base_delay: %s
      # 最大重连延迟
      max_delay: %s
      # 退避倍数
      multiplier: %.1f
   
  # 账号列表
  accounts: %v

web:
  # Web面板监听地址
  bind: "%s"
  
  # 认证配置
  auth:
    # JWT Token 过期时间
    token_expiry: %s
    # 审计日志保留天数
    audit_log_retention_days: %d
    
    # 密码列表
    passwords:
%s
  
  # Token加密配置
  token_vault:
    # 存储路径
    storage_path: "%s"
    # scrypt 参数
    scrypt_n: %d
    scrypt_r: %d
    scrypt_p: %d
    salt_len: %d
  
  # CORS配置
  cors:
    enabled: %t
    origins: %v

# 日志配置
log:
  # 日志目录
  log_dir: "%s"
  # 单个日志文件最大大小 (KB)
  max_size: %d
  # 是否启用调试模式
  debug: %t
  # 是否启用文件日志
  enable_file: %t
`,
		cfg.Cluster.Global.MaxInstances,
		cfg.Cluster.Global.ReconnectPolicy.Enabled,
		cfg.Cluster.Global.ReconnectPolicy.MaxRetries,
		cfg.Cluster.Global.ReconnectPolicy.BaseDelay,
		cfg.Cluster.Global.ReconnectPolicy.MaxDelay,
		cfg.Cluster.Global.ReconnectPolicy.Multiplier,
		cfg.Cluster.Accounts,
		cfg.Web.Bind,
		cfg.Web.Auth.TokenExpiry,
		cfg.Web.Auth.AuditLogRetentionDays,
		indentedPasswords,
		cfg.Web.TokenVault.StoragePath,
		cfg.Web.TokenVault.ScryptN,
		cfg.Web.TokenVault.ScryptR,
		cfg.Web.TokenVault.ScryptP,
		cfg.Web.TokenVault.SaltLen,
		cfg.Web.CORS.Enabled,
		cfg.Web.CORS.Origins,
		cfg.Log.LogDir,
		cfg.Log.MaxSize,
		cfg.Log.Debug,
		cfg.Log.EnableFile,
	)), nil
}

// fillDefaults 填充默认值
func (c *Config) fillDefaults() {
	// 集群默认值
	if c.Cluster.Global.MaxInstances == 0 {
		c.Cluster.Global.MaxInstances = 10
	}
	if c.Cluster.Global.ReconnectPolicy.BaseDelay == 0 {
		c.Cluster.Global.ReconnectPolicy.BaseDelay = 2 * time.Second
	}
	if c.Cluster.Global.ReconnectPolicy.MaxDelay == 0 {
		c.Cluster.Global.ReconnectPolicy.MaxDelay = 2 * time.Minute
	}
	if c.Cluster.Global.ReconnectPolicy.Multiplier == 0 {
		c.Cluster.Global.ReconnectPolicy.Multiplier = 1.8
	}

	// Web默认值
	if c.Web.Bind == "" {
		c.Web.Bind = "0.0.0.0:8080"
	}
	if c.Web.Auth.TokenExpiry == 0 {
		c.Web.Auth.TokenExpiry = 5 * time.Minute
	}
	if c.Web.Auth.AuditLogRetentionDays == 0 {
		c.Web.Auth.AuditLogRetentionDays = 30
	}
	if len(c.Web.Auth.Passwords) == 0 {
		c.Web.Auth.Passwords = []PasswordEntry{
			{
				ID:        "default",
				Hash:      "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
				Enabled:   true,
				CreatedAt: time.Now(),
				Note:      "默认密码: password123",
			},
		}
	}
	if c.Web.TokenVault.StoragePath == "" {
		c.Web.TokenVault.StoragePath = ".tokens"
	}
	if c.Web.TokenVault.ScryptN == 0 {
		c.Web.TokenVault.ScryptN = 1 << 20
	}
	if c.Web.TokenVault.ScryptR == 0 {
		c.Web.TokenVault.ScryptR = 8
	}
	if c.Web.TokenVault.ScryptP == 0 {
		c.Web.TokenVault.ScryptP = 1
	}
	if c.Web.TokenVault.SaltLen == 0 {
		c.Web.TokenVault.SaltLen = 32
	}

	// 日志默认值
	if c.Log.LogDir == "" {
		c.Log.LogDir = "logs"
	}
	if c.Log.MaxSize == 0 {
		c.Log.MaxSize = 512
	}
}
