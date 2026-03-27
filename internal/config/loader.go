package config

import (
	"fmt"
	"os"
	"strings"

	"gmcc/internal/logx"
	yaml "gopkg.in/yaml.v3"
)

// Load 读取并校验 YAML 配置。
// 如果文件不存在，会先写入默认配置并返回错误。
// 保持向后兼容，内部调用 LoadWithAutoUpdate。
func Load(path string) (*Config, error) {
	return LoadWithAutoUpdate(path, false)
}

// LoadWithAutoUpdate 读取并校验 YAML 配置，可选择启用自动更新。
// isAutoUpdateEnabled: 是否启用配置文件自动更新功能
func LoadWithAutoUpdate(path string, isAutoUpdateEnabled bool) (*Config, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("配置文件路径为空")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if writeErr := WriteDefault(path); writeErr != nil {
				return nil, fmt.Errorf("配置文件不存在，且默认配置写入失败: %w", writeErr)
			}
			// 启用自动更新时返回默认配置，否则返回错误
			if isAutoUpdateEnabled {
				cfg := Default()
				if err := cfg.Validate(); err != nil {
					return nil, err
				}
				return &cfg, nil
			}
			return nil, fmt.Errorf("配置文件不存在，已生成默认配置: %s", path)
		}
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		// 尝试备份损坏的配置文件
		if backupErr := createDamagedBackup(path); backupErr != nil {
			logx.Errorf("备份损坏的配置文件失败: %v", backupErr)
		}
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// 自动更新配置
	if isAutoUpdateEnabled {
		updatedCfg, err := performAutoUpdate(&cfg, path)
		if err != nil {
			logx.Errorf("配置自动更新失败: %v", err)
			// 更新失败时返回原配置，不影响程序启动
			return &cfg, nil
		}
		return updatedCfg, nil
	}

	return &cfg, nil
}

// createDamagedBackup 创建损坏配置文件的备份
func createDamagedBackup(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // 文件不存在，无需备份
	}

	backupPath := path + ".damaged." + fmt.Sprintf("%d", os.Getpid())

	return os.Rename(path, backupPath)
}

// performAutoUpdate 执行配置自动更新
func performAutoUpdate(cfg *Config, path string) (*Config, error) {
	merger := &ConfigMerger{}

	// 检查是否需要更新
	needsUpdate, err := merger.needsUpdate(cfg)
	if err != nil {
		return nil, fmt.Errorf("检查配置更新需求失败: %w", err)
	}

	if !needsUpdate {
		return cfg, nil // 无需更新
	}

	// 执行合并
	updatedCfg, changes, err := merger.MergeWithDefault(cfg)
	if err != nil {
		return nil, fmt.Errorf("合并配置失败: %w", err)
	}

	// 序列化更新的配置
	data, err := yaml.Marshal(updatedCfg)
	if err != nil {
		return nil, fmt.Errorf("序列化更新配置失败: %w", err)
	}

	// 原子更新文件
	if err := atomicUpdate(path, data); err != nil {
		return nil, fmt.Errorf("写入更新配置失败: %w", err)
	}

	// 发送通知
	notifier := GetNotifier(cfg.Runtime.Headless)
	if err := notifier.NotifyConfigUpdate(changes); err != nil {
		logx.Errorf("发送配置更新通知失败: %v", err)
	}

	logx.Infof("配置文件已自动更新，新增 %d 个配置项", len(changes))
	return updatedCfg, nil
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
