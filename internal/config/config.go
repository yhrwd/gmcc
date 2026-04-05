package config

import (
	"fmt"
	"os"
	"time"

	authsession "gmcc/internal/auth/session"

	yaml "gopkg.in/yaml.v3"
)

// Config 统一配置结构
type Config struct {
	// 认证配置
	Auth AuthConfig `yaml:"auth"`

	// 集群配置
	Cluster ClusterConfig `yaml:"cluster"`

	// Web面板配置
	Web WebConfig `yaml:"web"`

	// 日志配置
	Log LogConfig `yaml:"log"`

	// 向后兼容配置
	Server ServerConfig `yaml:"server"`

	// 集群运行时注入配置（不参与序列化）
	ClusterRuntime ClusterRuntimeConfig `yaml:"-"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	Vault AuthVaultConfig `yaml:"vault"`
}

// AuthVaultConfig 认证凭据保险库配置
type AuthVaultConfig struct {
	Path    string `yaml:"path"`
	KeyEnv  string `yaml:"key_env"`
	ScryptN int    `yaml:"scrypt_n"`
	ScryptR int    `yaml:"scrypt_r"`
	ScryptP int    `yaml:"scrypt_p"`
	SaltLen int    `yaml:"salt_len"`
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
	ID            string `yaml:"id"`
	ServerAddress string `yaml:"server_address"`
	Enabled       bool   `yaml:"enabled"`
}

// WebConfig Web面板配置
type WebConfig struct {
	Bind string        `yaml:"bind"` // 监听地址
	Auth WebAuthConfig `yaml:"auth"`
	CORS CORSConfig    `yaml:"cors"`
}

// WebAuthConfig Web认证配置
type WebAuthConfig struct {
	AuditLogRetentionDays int `yaml:"audit_log_retention_days"`
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
		Auth: AuthConfig{
			Vault: AuthVaultConfig{
				Path:    ".authvault",
				KeyEnv:  "GMCC_AUTH_VAULT_KEY",
				ScryptN: 1 << 20,
				ScryptR: 8,
				ScryptP: 1,
				SaltLen: 32,
			},
		},
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
				AuditLogRetentionDays: 30,
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
	return []byte(fmt.Sprintf(`# GMCC 配置文件
# 配置文件路径可通过环境变量 GMCC_CONFIG 指定，默认 config.yaml

auth:
  # 认证保险库配置
  vault:
    # 存储路径
    path: "%s"
    # 从环境变量读取主密钥
    key_env: "%s"
    # scrypt 参数
    scrypt_n: %d
    scrypt_r: %d
    scrypt_p: %d
    salt_len: %d

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
    # 审计日志保留天数
    audit_log_retention_days: %d

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
		cfg.Auth.Vault.Path,
		cfg.Auth.Vault.KeyEnv,
		cfg.Auth.Vault.ScryptN,
		cfg.Auth.Vault.ScryptR,
		cfg.Auth.Vault.ScryptP,
		cfg.Auth.Vault.SaltLen,
		cfg.Cluster.Global.MaxInstances,
		cfg.Cluster.Global.ReconnectPolicy.Enabled,
		cfg.Cluster.Global.ReconnectPolicy.MaxRetries,
		cfg.Cluster.Global.ReconnectPolicy.BaseDelay,
		cfg.Cluster.Global.ReconnectPolicy.MaxDelay,
		cfg.Cluster.Global.ReconnectPolicy.Multiplier,
		cfg.Cluster.Accounts,
		cfg.Web.Bind,
		cfg.Web.Auth.AuditLogRetentionDays,
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
	if c.Auth.Vault.Path == "" {
		c.Auth.Vault.Path = ".authvault"
	}
	if c.Auth.Vault.KeyEnv == "" {
		c.Auth.Vault.KeyEnv = "GMCC_AUTH_VAULT_KEY"
	}
	if c.Auth.Vault.ScryptN == 0 {
		c.Auth.Vault.ScryptN = 1 << 20
	}
	if c.Auth.Vault.ScryptR == 0 {
		c.Auth.Vault.ScryptR = 8
	}
	if c.Auth.Vault.ScryptP == 0 {
		c.Auth.Vault.ScryptP = 1
	}
	if c.Auth.Vault.SaltLen == 0 {
		c.Auth.Vault.SaltLen = 32
	}

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
	if c.Web.Auth.AuditLogRetentionDays == 0 {
		c.Web.Auth.AuditLogRetentionDays = 30
	}

	// 日志默认值
	if c.Log.LogDir == "" {
		c.Log.LogDir = "logs"
	}
	if c.Log.MaxSize == 0 {
		c.Log.MaxSize = 512
	}
}
