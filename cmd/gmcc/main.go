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

	if err := logx.Init(cfg.Log.LogDir, cfg.Log.EnableFile, cfg.Log.MaxSize, cfg.Log.Debug); err != nil {
		fmt.Fprintf(os.Stderr, "[错误] 日志初始化失败: %v\n", err)
		os.Exit(1)
	}
	defer logx.Close()

	runTUI(cfg)
}

type App struct {
	engine *tui.Engine
	client *mcclient.Client
	cfg    *config.Config

	input   *tui.InputWidget
	logs    []string
	maxLogs int
}

func runTUI(cfg *config.Config) {
	engine, err := tui.NewEngine()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[错误] TUI 初始化失败: %v\n", err)
		os.Exit(1)
	}

	app := &App{
		engine:  engine,
		cfg:     cfg,
		maxLogs: 200,
		logs:    make([]string, 0, 200),
		input:   tui.NewInputWidget(),
	}

	app.input.Focus()

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

	app.engine.On(tui.EventKey, func(e interface{}) bool {
		ke := e.(tui.KeyEvent)
		if ke.Key == tui.KeyEscape {
			app.engine.Stop()
			return true
		}
		return app.input.HandleEvent(e)
	})

	app.input.OnSubmit = func(text string) {
		if text == "" {
			return
		}
		app.addLog("> " + text)
		if strings.HasPrefix(text, "/") {
			cmd := strings.TrimPrefix(text, "/")
			if app.client.IsReady() {
				if err := app.client.SendCommand(cmd); err != nil {
					app.addLog(fmt.Sprintf("[命令失败] %v", err))
				}
			}
		} else {
			if app.client.IsReady() {
				if err := app.client.SendMessage(text); err != nil {
					app.addLog(fmt.Sprintf("[发送失败] %v", err))
				}
			}
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.client.Run(ctx)
	}()

	go func() {
		select {
		case <-ctx.Done():
			app.addLog("[提示] 断开连接")
			app.engine.Stop()
		case err := <-errCh:
			if err != nil {
				app.addLog(fmt.Sprintf("[客户端退出] %v", err))
			}
			app.engine.Stop()
		}
	}()

	app.addLog(fmt.Sprintf("版本: %s", Version))
	app.addLog(fmt.Sprintf("玩家: %s", cfg.Account.PlayerID))
	app.addLog(fmt.Sprintf("服务器: %s", cfg.Server.Address))
	app.addLog("正在连接...")
	app.addLog("")

	app.engine.SetLayout(tui.ComponentFunc(func(w, h int) []string {
		return app.render(w, h)
	}))

	if err := app.engine.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "[错误] %v\n", err)
	}
}

func (a *App) addLog(text string) {
	a.logs = append(a.logs, text)
	if len(a.logs) > a.maxLogs {
		a.logs = a.logs[1:]
	}
}

func (a *App) render(w, h int) []string {
	lines := make([]string, h)
	for i := range lines {
		lines[i] = strings.Repeat(" ", w)
	}

	if h < 6 {
		return lines
	}

	lineNum := 0
	lines[lineNum] = fmt.Sprintf(" ╔%s╗", strings.Repeat("─", w-2))
	lineNum++

	player := a.client.Player
	if player != nil {
		hp, _, food, _ := player.GetHealth()
		x, y, z := player.GetPosition()
		gm := player.GameMode.String()
		status := fmt.Sprintf(" HP: %.0f | Food: %d | Pos: %.1f,%.1f,%.1f | Mode: %s ", hp, food, x, y, z, gm)
		if len(status) > w-4 {
			status = status[:w-4]
		}
		lines[lineNum] = fmt.Sprintf(" │%s│", status+strings.Repeat(" ", w-4-len(status)))
	} else {
		title := fmt.Sprintf(" gmcc v%s - %s ", Version, a.cfg.Account.PlayerID)
		if len(title) > w-4 {
			title = title[:w-4]
		}
		lines[lineNum] = fmt.Sprintf(" │%s│", title+strings.Repeat(" ", w-4-len(title)))
	}
	lineNum++
	lines[lineNum] = fmt.Sprintf(" ╚%s╝", strings.Repeat("─", w-2))
	lineNum++

	logsStart := 0
	maxLogsDisplay := h - lineNum - 3
	if maxLogsDisplay < 1 {
		maxLogsDisplay = 1
	}
	if len(a.logs) > maxLogsDisplay {
		logsStart = len(a.logs) - maxLogsDisplay
	}

	for _, log := range a.logs[logsStart:] {
		if lineNum >= h-3 {
			break
		}
		runes := []rune(log)
		if len(runes) > w {
			runes = runes[:w]
		}
		lines[lineNum] = string(runes) + strings.Repeat(" ", w-len(runes))
		lineNum++
	}

	for lineNum < h-3 {
		lines[lineNum] = strings.Repeat(" ", w)
		lineNum++
	}

	inputLine := a.input.Value
	if len(inputLine) > w-3 {
		inputLine = inputLine[:w-3]
	}
	lines[lineNum] = "> " + inputLine + strings.Repeat(" ", w-3-len(inputLine))
	lineNum++

	for lineNum < h {
		lines[lineNum] = strings.Repeat(" ", w)
		lineNum++
	}

	footer := " Ctrl+C/Esc:退出 | Up/Dn:历史 "
	if len(footer) > w {
		footer = footer[:w]
	}
	lines[h-1] = footer + strings.Repeat(" ", w-len(footer))

	return lines
}
