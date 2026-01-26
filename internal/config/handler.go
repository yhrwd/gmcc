package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 全局配置结构体
type Config struct {
	Account AccountConfig  `yaml:"account"`
	Server  ServerSetting  `yaml:"server"`
	Log     LogConfig         `yaml:"log"`
}

// AccountConfig 账号配置
type AccountConfig struct {
	PlayerID string `yaml:"player_id"`
}

// ServerSetting 服务器配置
type ServerSetting struct {
	Address string `yaml:"address"`
}

type LogConfig struct {
	LogDir     string `yaml:"log_dir"`
	MaxSize    int64  `yaml:"max_size"`
	Debug      bool   `yaml:"debug"`
	EnableFile bool   `yaml:"enable_file"`
}

var path string

func init() {
	// 获取程序当前运行目录
	exeDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("无法获取程序运行目录: %v\n", err)
		os.Exit(1)
	}
	path = filepath.Join(exeDir, "config.yaml")
}

var defaultYAML = `account:
  # 玩家ID
  player_id: "oiiaioooiiai"
server:
  # 服务器地址
  address: "127.0.0.1"
log:
  log_dir: "logs"
  max_size: 5242880   # 5 * 1024 * 1024 = 5MB
  debug: true
  enable_file: true
`

func createConfig() error {
	return os.WriteFile(path, []byte(defaultYAML), 0644)
}

func verConfig(c *Config) error {
	if c.Account.PlayerID == "" {
		return fmt.Errorf("account.player_id 不能为空")
	}

	return nil
}

func LoadConfig() (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("配置文件不存在，将创建默认配置...")

			err := createConfig()
			if err != nil {
				return nil, fmt.Errorf("创建失败: %w", err)
			}

			fmt.Printf("默认配置已创建: %s\n", path)
			fmt.Println("请修改配置文件后重新运行程序。")

			os.Exit(0)
		}
		return nil, fmt.Errorf("打开失败: %w", err)
	}
	defer f.Close()

	var config Config
	err = yaml.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("解析失败: %w", err)
	}

	if err := verConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
