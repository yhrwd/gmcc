package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gmcc/internal/config"
	"gmcc/internal/logx"
	"gmcc/internal/tui"
)

var Version = "dev"

func main() {
	configPath := "config.yaml"
	if v := os.Getenv("GMCC_CONFIG"); v != "" {
		configPath = v
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[错误] 配置加载失败: %v\n", err)
		os.Exit(1)
	}

	cfg.Log.Debug = false

	if err := logx.Init(cfg.Log.LogDir, cfg.Log.EnableFile, cfg.Log.MaxSize, cfg.Log.Debug); err != nil {
		fmt.Fprintf(os.Stderr, "[错误] 日志初始化失败: %v\n", err)
		os.Exit(1)
	}
	defer logx.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ui := tui.New(cfg)
	if err := ui.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "[错误] %v\n", err)
		os.Exit(1)
	}
}
