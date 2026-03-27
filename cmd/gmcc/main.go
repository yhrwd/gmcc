package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gmcc/internal/config"
	"gmcc/internal/headless"
	"gmcc/internal/logx"
	"gmcc/internal/tui"
)

var Version = "dev"

func main() {
	configPath := "config.yaml"
	if v := os.Getenv("GMCC_CONFIG"); v != "" {
		configPath = v
	}

	// 默认启用配置自动更新
	// 可以通过环境变量 GMCC_DISABLE_AUTO_UPDATE=true 禁用
	disableAutoUpdate := os.Getenv("GMCC_DISABLE_AUTO_UPDATE") == "true"
	autoUpdate := !disableAutoUpdate

	cfg, err := config.LoadWithAutoUpdate(configPath, autoUpdate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[错误] 配置加载失败: %v\n", err)
		os.Exit(1)
	}

	if err := logx.Init(cfg.Log.LogDir, cfg.Log.EnableFile, cfg.Log.MaxSizeInBytes(), cfg.Log.Debug); err != nil {
		fmt.Fprintf(os.Stderr, "[错误] 日志初始化失败: %v\n", err)
		os.Exit(1)
	}
	defer logx.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 根据配置选择运行模式
	if cfg.Runtime.Headless {
		// 无界面模式
		runner := headless.New(cfg)
		if err := runner.Run(ctx); err != nil {
			logx.Errorf("运行错误: %v", err)
			os.Exit(1)
		}
	} else {
		// TUI 模式
		ui := tui.New(cfg)
		if err := ui.Run(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "[错误] %v\n", err)
			os.Exit(1)
		}
	}
}
