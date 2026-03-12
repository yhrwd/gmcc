package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"gmcc/internal/config"
	"gmcc/internal/logx"
	"gmcc/internal/mcclient"
)

var Version = "dev"

type App struct {
	client  *mcclient.Client
	cfg     *config.Config
	logs    []string
	maxLogs int
	mu      sync.RWMutex
	running bool
	history []string
	histIdx int
}

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

	if err := logx.Init(cfg.Log.LogDir, cfg.Log.EnableFile, cfg.Log.MaxSize, cfg.Log.Debug); err != nil {
		fmt.Fprintf(os.Stderr, "[错误] 日志初始化失败: %v\n", err)
		os.Exit(1)
	}
	defer logx.Close()

	runApp(cfg)
}

func runApp(cfg *config.Config) {
	fmt.Println()
	fmt.Println(" ═══ gmcc - Minecraft 控制台客户端 ═══")
	fmt.Printf(" 版本: %s\n", Version)
	fmt.Println(" ─────────────────────────────────────────────")
	fmt.Println()
	fmt.Printf("[✓] 配置加载成功\n")
	fmt.Printf("  玩家: %s\n", cfg.Account.PlayerID)
	fmt.Printf("  服务器: %s\n", cfg.Server.Address)
	fmt.Println()

	app := &App{
		maxLogs: 100,
		logs:    make([]string, 0, 100),
		history: make([]string, 0, 100),
		histIdx: -1,
		cfg:     cfg,
	}

	app.client = mcclient.New(cfg)
	app.client.SetChatHandler(func(msg mcclient.ChatMessage) {
		var text string
		if msg.RawJSON != "" {
			comp, err := mcclient.ParseTextComponent(msg.RawJSON)
			if err == nil {
				text = comp.ToANSI()
			} else {
				text = msg.PlainText
			}
		} else {
			text = msg.PlainText
		}
		app.addLog(text)
	})

	fmt.Printf("正在连接 %s ...\n\n", cfg.Server.Address)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.client.Run(ctx)
	}()

	// Wait for connection
	time.Sleep(2 * time.Second)

	app.runInputLoop(ctx)

	select {
	case err := <-errCh:
		if err != nil {
			fmt.Printf("\n[客户端退出] %v\n", err)
		}
	case <-ctx.Done():
		fmt.Println("\n[提示] 断开连接")
	}
}

func (a *App) addLog(text string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.logs = append(a.logs, text)
	if len(a.logs) > a.maxLogs {
		a.logs = a.logs[1:]
	}
	fmt.Println(text)
}

func (a *App) runInputLoop(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		a.addToHistory(line)
		a.processInput(line, ctx)
	}
}

func (a *App) addToHistory(cmd string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.history) == 0 || a.history[len(a.history)-1] != cmd {
		a.history = append(a.history, cmd)
		if len(a.history) > 100 {
			a.history = a.history[1:]
		}
	}
	a.histIdx = -1
}

func (a *App) processInput(line string, ctx context.Context) {
	if strings.HasPrefix(line, "/") {
		cmd := strings.TrimPrefix(line, "/")
		if a.client.IsReady() {
			if err := a.client.SendCommand(cmd); err != nil {
				a.addLog(fmt.Sprintf("[命令失败] %v", err))
			}
		} else {
			a.addLog("[提示] 尚未连接到服务器")
		}
	} else {
		if a.client.IsReady() {
			if err := a.client.SendMessage(line); err != nil {
				a.addLog(fmt.Sprintf("[发送失败] %v", err))
			}
		} else {
			a.addLog("[提示] 尚未连接到服务器")
		}
	}
}
