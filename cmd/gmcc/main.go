package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gmcc/internal/config"
)

func main() {
	// 获取可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("无法获取可执行文件路径: %v\n", err)
		os.Exit(1)
	}

	// 构建配置文件路径
	execDir := filepath.Dir(execPath)
	configPath := filepath.Join(execDir, "config.yaml")

	fmt.Printf("正在加载配置文件: %s\n", configPath)

	// 加载配置
	cfg := config.LoadCfg(configPath)
	if cfg == nil {
		fmt.Println("配置加载失败")
		os.Exit(1)
	}

	fmt.Println("配置加载成功:")
	fmt.Printf("  Player ID: %s\n", cfg.Account.PlayerID)
	fmt.Printf("  Server Address: %s\n", cfg.Server.Address)
	fmt.Printf("  Log Directory: %s\n", cfg.Log.LogDir)
}
