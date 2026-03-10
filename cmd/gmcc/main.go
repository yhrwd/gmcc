package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"gmcc/internal/config"
	"gmcc/internal/logx"
	"gmcc/internal/mcclient"
)

var Version = "dev"

func main() {
	configPath := "config.yaml"
	if v := os.Getenv("GMCC_CONFIG"); v != "" {
		configPath = v
	}

	logx.Infof("gmcc version: %s", Version)
	logx.Infof("正在加载配置文件: %s", configPath)

	// 加载配置
	cfg, err := config.Load(configPath)
	if err != nil {
		logx.Errorf("配置加载失败: %v", err)
		logx.Infof("请修改配置文件后重新运行程序")
		os.Exit(1)
	}

	if err := logx.Init(cfg.Log.LogDir, cfg.Log.EnableFile, cfg.Log.MaxSize, cfg.Log.Debug); err != nil {
		logx.Errorf("日志初始化失败: %v", err)
		os.Exit(1)
	}
	defer func() {
		_ = logx.Close()
	}()

	logx.Infof("配置加载成功")
	logx.Infof("Player ID: %s", cfg.Account.PlayerID)
	logx.Infof("Use Official Auth: %t", cfg.Account.UseOfficialAuth)
	logx.Infof("Server Address: %s", cfg.Server.Address)
	logx.Infof("Log Directory: %s", cfg.Log.LogDir)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client := mcclient.New(cfg)
	if err := client.Run(ctx); err != nil {
		logx.Errorf("客户端退出: %v", err)
		os.Exit(1)
	}
}
