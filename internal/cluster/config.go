package cluster

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gmcc/internal/logx"
	yaml "gopkg.in/yaml.v3"
)

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
	ID              string `yaml:"id"`
	PlayerID        string `yaml:"player_id"`
	UseOfficialAuth bool   `yaml:"use_official_auth"`
	ServerAddress   string `yaml:"server_address"`
	Enabled         bool   `yaml:"enabled"`
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
				MaxRetries: 5,
				BaseDelay:  5 * time.Second,
				MaxDelay:   300 * time.Second,
				Multiplier: 2.0,
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

// LoadClusterConfig 从文件加载集群配置
func LoadClusterConfig(path string) (ClusterConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 配置文件不存在，返回默认配置
			cfg := DefaultClusterConfig()
			if err := SaveClusterConfig(path, cfg); err != nil {
				logx.Warnf("创建默认集群配置文件失败: %v", err)
			}
			return cfg, nil
		}
		return ClusterConfig{}, fmt.Errorf("读取集群配置文件失败: %w", err)
	}

	var cfg ClusterConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return ClusterConfig{}, fmt.Errorf("解析集群配置文件失败: %w", err)
	}

	return cfg, nil
}

// SaveClusterConfig 保存集群配置到文件
func SaveClusterConfig(path string, cfg ClusterConfig) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建配置目录失败: %w", err)
		}
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("序列化集群配置失败: %w", err)
	}

	// 原子写入：先写入临时文件，再重命名
	tempPath := path + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("写入临时配置文件失败: %w", err)
	}

	if err := os.Rename(tempPath, path); err != nil {
		// 重命名失败，尝试直接写入
		_ = os.Remove(tempPath)
		if err := os.WriteFile(path, data, 0644); err != nil {
			return fmt.Errorf("写入配置文件失败: %w", err)
		}
	}

	logx.Debugf("集群配置已保存: %s", path)
	return nil
}
