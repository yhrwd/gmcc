package config

import (
	"strings"
	"testing"
)

func TestConfigMerger_MergeWithDefault(t *testing.T) {
	merger := &ConfigMerger{}

	tests := []struct {
		name          string
		current       *Config
		expectChanges int
		expectFields  map[string]any
	}{
		{
			name:          "nil_config",
			current:       nil,
			expectChanges: 0,
		},
		{
			name: "full_config",
			current: &Config{
				Account: AccountConfig{
					PlayerID:        "test_player",
					UseOfficialAuth: true,
				},
				Server: ServerConfig{
					Address: "test.example.com",
				},
				Actions: ActionsConfig{
					DelayMs: 2000,
				},
				Log: LogConfig{
					LogDir:     "test_logs",
					MaxSize:    1024,
					Debug:      false,
					EnableFile: true,
				},
				Runtime: RuntimeConfig{
					Headless: true,
				},
				Packets: PacketConfig{
					HandleContainer: false,
				},
			},
			expectChanges: 0,
		},
		{
			name: "partial_config",
			current: &Config{
				Account: AccountConfig{
					PlayerID:        "test_player",
					UseOfficialAuth: true,
				},
				Server: ServerConfig{
					Address: "test.example.com",
				},
				Actions: ActionsConfig{
					DelayMs:             2000,
					SignCommands:        true,
					DefaultSignCommands: false,
				},
				Log: LogConfig{
					LogDir:     "test_logs",
					MaxSize:    1024,
					Debug:      false,
					EnableFile: true,
				},
				Runtime: RuntimeConfig{
					Headless: false,
				},
				Packets: PacketConfig{
					HandleContainer: false,
				},
			},
			expectChanges: 0,
			expectFields: map[string]any{
				"server.address": "test.example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, changes, err := merger.MergeWithDefault(tt.current)

			if tt.current == nil {
				if err == nil {
					t.Fatalf("期望错误，但没有发生错误")
				}
				if !strings.Contains(err.Error(), "当前配置为空") {
					t.Fatalf("错误消息不匹配: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("合并失败: %v", err)
			}

			if len(changes) != tt.expectChanges {
				t.Fatalf("变更数量不匹配，期望: %d, 实际: %d", tt.expectChanges, len(changes))
			}

			// 验证特定字段的值
			for path, expectedValue := range tt.expectFields {
				actualValue := getFieldValue(result, path)
				if actualValue != expectedValue {
					t.Errorf("字段 %s 的值不匹配，期望: %v, 实际: %v", path, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestConfigMerger_needsUpdate(t *testing.T) {
	merger := &ConfigMerger{}

	tests := []struct {
		name             string
		current          *Config
		expectNeedUpdate bool
	}{
		{
			name:             "nil_config",
			current:          nil,
			expectNeedUpdate: true,
		},
		{
			name: "full_config",
			current: &Config{
				Account: AccountConfig{
					PlayerID:        "test_player",
					UseOfficialAuth: true,
				},
				Server: ServerConfig{
					Address: "test.example.com",
				},
				Actions: ActionsConfig{
					DelayMs: 2000,
				},
				Log: LogConfig{
					LogDir:     "test_logs",
					MaxSize:    1024,
					Debug:      false,
					EnableFile: true,
				},
				Runtime: RuntimeConfig{
					Headless: true,
				},
				Packets: PacketConfig{
					HandleContainer: false,
				},
			},
			expectNeedUpdate: false,
		},
		{
			name: "partial_config",
			current: &Config{
				Account: AccountConfig{
					PlayerID:        "test_player",
					UseOfficialAuth: true,
				},
			},
			expectNeedUpdate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			needsUpdate, err := merger.needsUpdate(tt.current)
			if err != nil {
				t.Fatalf("检查更新需求失败: %v", err)
			}
			if needsUpdate != tt.expectNeedUpdate {
				t.Errorf("更新需求不匹配，期望: %v, 实际: %v", tt.expectNeedUpdate, needsUpdate)
			}
		})
	}
}

// getFieldValue 辅助函数：通过路径获取字段值（简化版）
func getFieldValue(cfg *Config, path string) any {
	parts := strings.Split(path, ".")

	switch parts[0] {
	case "server":
		if len(parts) > 1 {
			switch parts[1] {
			case "address":
				return cfg.Server.Address
			}
		}
	case "account":
		if len(parts) > 1 {
			switch parts[1] {
			case "player_id":
				return cfg.Account.PlayerID
			}
		}
	}

	return nil
}
