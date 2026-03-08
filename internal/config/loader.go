package config

import (
	"fmt"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// Load 读取并校验 YAML 配置。
// 如果文件不存在，会先写入默认配置并返回错误。
func Load(path string) (*Config, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("配置文件路径为空")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if writeErr := WriteDefault(path); writeErr != nil {
				return nil, fmt.Errorf("配置文件不存在，且默认配置写入失败: %w", writeErr)
			}
			return nil, fmt.Errorf("配置文件不存在，已生成默认配置: %s", path)
		}
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// WriteDefault 将默认配置写入指定路径。
func WriteDefault(path string) error {
	data, err := yaml.Marshal(Default())
	if err != nil {
		return fmt.Errorf("序列化默认配置失败: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("写入默认配置失败: %w", err)
	}
	return nil
}

// Validate 检查配置是否有效（不包含关键空值）
func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("配置为空")
	}

	var invalid []string
	if strings.TrimSpace(c.Account.PlayerID) == "" {
		invalid = append(invalid, "account.player_id")
	}
	if strings.TrimSpace(c.Server.Address) == "" {
		invalid = append(invalid, "server.address")
	}
	if strings.TrimSpace(c.Log.LogDir) == "" {
		invalid = append(invalid, "log.log_dir")
	}
	if c.Log.MaxSize <= 0 {
		invalid = append(invalid, "log.max_size")
	}
	if c.Actions.DelayMs < 0 {
		invalid = append(invalid, "actions.delay_ms")
	}

	if len(invalid) > 0 {
		return fmt.Errorf("以下配置项无效: %s", strings.Join(invalid, ", "))
	}
	return nil
}
