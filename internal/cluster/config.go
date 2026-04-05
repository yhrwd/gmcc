package cluster

import "time"

// ClusterConfig 多账号集群配置
type ClusterConfig struct {
	Global   GlobalConfig   `yaml:"global"`
	Accounts []AccountEntry `yaml:"accounts"`
	Log      LogConfig      `yaml:"log"`
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

// LogConfig 日志配置
type LogConfig struct {
	LogDir     string `yaml:"log_dir"`
	MaxSize    int64  `yaml:"max_size"`
	Debug      bool   `yaml:"debug"`
	EnableFile bool   `yaml:"enable_file"`
}

// ClusterStatus 集群状态响应类型（用于Web API）
type ClusterStatus struct {
	Status           string        `json:"cluster_status"`
	TotalInstances   int           `json:"total_instances"`
	RunningInstances int           `json:"running_instances"`
	Uptime           time.Duration `json:"uptime"`
}

// DefaultClusterConfig 返回默认集群配置
func DefaultClusterConfig() ClusterConfig {
	return ClusterConfig{
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
		Log: LogConfig{
			LogDir:     "logs",
			MaxSize:    512,
			Debug:      false,
			EnableFile: true,
		},
	}
}
