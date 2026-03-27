package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

func TestLoadWithAutoUpdate(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()

	tests := []struct {
		name          string
		configContent string
		autoUpdate    bool
		expectError   bool
		expectChanges bool
	}{
		{
			name: "existing_config_no_update",
			configContent: `
account:
  player_id: "test_player"
  use_official_auth: true
server:
  address: "test.example.com"
actions:
  delay_ms: 2000
log:
  log_dir: "test_logs"
  max_size: 1024
  debug: false
  enable_file: true
runtime:
  headless: false
packets:
  handle_container: false
`,
			autoUpdate:    true,
			expectError:   false,
			expectChanges: false,
		},
		{
			name: "minimal_config_with_update",
			configContent: `
account:
  player_id: "test_player"
  use_official_auth: true
server:
  address: "test.example.com"
log:
  log_dir: "test_logs"
  max_size: 1024
actions:
  delay_ms: 2000
runtime:
  headless: false
packets:
  handle_container: false
`,
			autoUpdate:    true,
			expectError:   false,
			expectChanges: true, // 会添加缺失的字段如 SignCommands, DefaultSignCommands
		},
		{
			name: "auto_update_disabled",
			configContent: `
account:
  player_id: "test_player"
  use_official_auth: true
server:
  address: "test.example.com"
log:
  log_dir: "test_logs"
  max_size: 1024
actions:
  delay_ms: 2000
runtime:
  headless: false
packets:
  handle_container: false
`,
			autoUpdate:    false,
			expectError:   false,
			expectChanges: false,
		},
		{
			name:          "nonexistent_config",
			configContent: "",
			autoUpdate:    true,
			expectError:   false,
			expectChanges: false, // 应该返回默认配置
		},
		{
			name: "invalid_yaml",
			configContent: `
invalid: yaml: content: [
`,
			autoUpdate:  true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建配置文件（如果需要）
			var configPath string
			if tt.configContent != "" {
				configFile := filepath.Join(tmpDir, "config.yaml")
				err := os.WriteFile(configFile, []byte(strings.TrimSpace(tt.configContent)), 0644)
				if err != nil {
					t.Fatalf("创建配置文件失败: %v", err)
				}
				configPath = configFile
			} else {
				configPath = filepath.Join(tmpDir, "nonexistent.yaml")
			}

			// 加载配置
			cfg, err := LoadWithAutoUpdate(configPath, tt.autoUpdate)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误，但没有发生")
				}
				return
			}

			if err != nil {
				t.Fatalf("加载配置失败: %v", err)
			}

			// 验证配置不为空
			if cfg == nil {
				t.Fatalf("配置为空")
			}

			// 检查是否添加了备份文件
			if tt.autoUpdate && tt.configContent != "" {
				matches, err := filepath.Glob(configPath + ".backup.*")
				if err != nil {
					t.Fatalf("检查备份文件失败: %v", err)
				}

				if tt.expectChanges && len(matches) == 0 {
					t.Errorf("期望有更新但没有生成备份文件")
				}
			}
		})
	}
}

func TestPerformAutoUpdate(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// 先创建一个配置文件（用于更新）
	initialConfig := &Config{
		Account: AccountConfig{
			PlayerID:        "test_player",
			UseOfficialAuth: true,
		},
		Server: ServerConfig{
			Address: "test.example.com",
		},
		Actions: ActionsConfig{
			DelayMs: 1000,
		},
		Log: LogConfig{
			LogDir:     "test_logs",
			MaxSize:    512,
			Debug:      false,
			EnableFile: true,
		},
		Runtime: RuntimeConfig{
			Headless: false,
		},
		Packets: PacketConfig{
			HandleContainer: false,
		},
	}

	// 写入初始配置
	data, err := yaml.Marshal(initialConfig)
	if err != nil {
		t.Fatalf("序列化初始配置失败: %v", err)
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("写入初始配置失败: %v", err)
	}

	// 创建部分配置对象（有缺失字段）
	partialConfig := &Config{
		Account: AccountConfig{
			PlayerID:        "test_player",
			UseOfficialAuth: true,
		},
		// 注意：这里故意不包含 Actions.SignCommands 和 Actions.DefaultSignCommands
		// 这些字段应该是零值（false）
		Server: ServerConfig{
			Address: "test.example.com",
		},
		Actions: ActionsConfig{
			DelayMs: 1000,
			// SignCommands 和 DefaultSignCommands 将被默认值填补
		},
		Log: LogConfig{
			LogDir:     "test_logs",
			MaxSize:    512,
			Debug:      false,
			EnableFile: true,
		},
		Runtime: RuntimeConfig{
			Headless: false,
		},
		Packets: PacketConfig{
			HandleContainer: false,
		},
	}

	// 测试自动更新
	updatedCfg, err := performAutoUpdate(partialConfig, configPath)
	if err != nil {
		t.Fatalf("自动更新失败: %v", err)
	}

	// 验证更新
	if updatedCfg == nil {
		t.Fatalf("更新后的配置为空")
	}

	// 检查文件是否被更新
	data, err = os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("读取更新后的文件失败: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "address: test.example.com") {
		t.Errorf("配置文件应该包含原有的服务器地址")
	}

	if !strings.Contains(content, "sign_commands: true") {
		t.Errorf("配置文件应该包含默认的签名命令设置")
	}
}
