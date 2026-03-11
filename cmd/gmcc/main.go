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
	fmt.Println(" ═══ gmcc - Minecraft 控制台客户端 ═══ ")
	fmt.Printf(" 版本: %s\n", Version)
	fmt.Println(" ─────────────────────────────────────────────")
	fmt.Println()

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Printf("[错误] 配置加载失败: %v\n", err)
		os.Exit(1)
	}

	if err := logx.Init(cfg.Log.LogDir, cfg.Log.EnableFile, cfg.Log.MaxSize, cfg.Log.Debug); err != nil {
		fmt.Printf("[错误] 日志初始化失败: %v\n", err)
		os.Exit(1)
	}
	defer logx.Close()

	fmt.Printf("[✓] 配置加载成功\n")
	fmt.Printf("  玩家: %s\n", cfg.Account.PlayerID)
	fmt.Printf("  服务器: %s\n", cfg.Server.Address)
	fmt.Println()

	term := terminal.New()
	if err := term.Start(); err != nil {
		fmt.Printf("[\033[31m错误\033[0m] 终端初始化失败: %v\n", err)
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
				term.Printf("[\033[31m发送失败\033[0m] %v", err)
			}
		}
	})

	term.SetCommandHook(func(cmd string) {
		if client.IsReady() {
			if err := client.SendCommand(cmd); err != nil {
				term.Printf("[\033[31m命令失败\033[0m] %v", err)
			}
		}
	})

	term.Printf("\n正在连接 %s ...", cfg.Server.Address)

	errCh := make(chan error, 1)
	go func() {
		errCh <- client.Run(ctx)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			term.Printf("[\033[31m客户端退出\033[0m] %v", err)
		} else {
			term.PrintLine("[\033[33m提示\033[0m] 客户端已断开连接")
		}
	case <-ctx.Done():
		term.PrintLine("[\033[33m提示\033[0m] 正在断开连接...")
	}

	term.PrintLine("[\033[90m提示\033[0m] 按 Ctrl+C 退出")
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
