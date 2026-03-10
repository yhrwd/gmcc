package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"gmcc/internal/config"
	"gmcc/internal/logx"
	"gmcc/internal/mcclient"
	"gmcc/internal/terminal"
)

var Version = "dev"

func main() {
	configPath := "config.yaml"
	if v := os.Getenv("GMCC_CONFIG"); v != "" {
		configPath = v
	}

	if !terminal.IsTerminal() {
		runHeadless(configPath)
		return
	}

	runInteractive(configPath)
}

func runInteractive(configPath string) {
	logx.Infof("gmcc version: %s", Version)
	logx.Infof("正在加载配置文件: %s", configPath)

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Printf("配置加载失败: %v\n", err)
		os.Exit(1)
	}

	if err := logx.Init(cfg.Log.LogDir, cfg.Log.EnableFile, cfg.Log.MaxSize, cfg.Log.Debug); err != nil {
		fmt.Printf("日志初始化失败: %v\n", err)
		os.Exit(1)
	}
	defer logx.Close()

	term := terminal.New()
	if err := term.Start(); err != nil {
		fmt.Printf("终端初始化失败: %v\n", err)
		os.Exit(1)
	}
	defer term.Stop()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

	client := mcclient.New(cfg)

	client.SetChatHandler(func(msg mcclient.ChatMessage) {
		if msg.RawJSON != "" {
			comp, err := mcclient.ParseTextComponent(msg.RawJSON)
			if err == nil {
				term.PrintLine(comp.ToANSI())
				return
			}
		}
		term.PrintLine(msg.PlainText)
	})

	term.SetMessageHook(func(msg string) {
		if client.IsReady() {
			if err := client.SendMessage(msg); err != nil {
				term.PrintLine(fmt.Sprintf("\033[31m[发送失败] %v\033[0m", err))
			}
		}
	})

	term.SetCommandHook(func(cmd string) {
		if client.IsReady() {
			if err := client.SendCommand(cmd); err != nil {
				term.PrintLine(fmt.Sprintf("\033[31m[命令失败] %v\033[0m", err))
			}
		}
	})

	term.PrintLine(fmt.Sprintf("\033[36mgmcc %s - Minecraft 控制台客户端\033[0m", Version))
	term.PrintLine(fmt.Sprintf("正在连接 %s ...", cfg.Server.Address))

	errCh := make(chan error, 1)
	go func() {
		errCh <- client.Run(ctx)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			term.PrintLine(fmt.Sprintf("\033[31m客户端退出: %v\033[0m", err))
		} else {
			term.PrintLine("客户端已断开连接")
		}
	case <-ctx.Done():
		term.PrintLine("正在断开连接...")
	}

	term.PrintLine("按 Ctrl+D 或输入 /quit 退出")
}

func runHeadless(configPath string) {
	logx.Infof("gmcc version: %s", Version)
	logx.Infof("正在加载配置文件: %s", configPath)

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
	defer logx.Close()

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

func isQuitCommand(input string) bool {
	cmd := strings.ToLower(strings.TrimSpace(input))
	return cmd == "quit" || cmd == "exit" || cmd == "/quit" || cmd == "/exit"
}
