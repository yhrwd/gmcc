package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/term"

	"gmcc/internal/config"
	"gmcc/internal/logx"
	"gmcc/internal/mcclient"
)

var Version = "dev"

type App struct {
	client     *mcclient.Client
	cfg        *config.Config
	mu         sync.RWMutex
	logs       []string
	maxLogs    int
	termWidth  int
	termHeight int
	oldState   *term.State
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

	// 禁用控制台debug输出，只写文件
	cfg.Log.Debug = false

	if err := logx.Init(cfg.Log.LogDir, cfg.Log.EnableFile, cfg.Log.MaxSize, cfg.Log.Debug); err != nil {
		fmt.Fprintf(os.Stderr, "[错误] 日志初始化失败: %v\n", err)
		os.Exit(1)
	}
	defer logx.Close()

	runTUI(cfg)
}

func runTUI(cfg *config.Config) {
	app := &App{
		maxLogs: 200,
		logs:    make([]string, 0, 200),
		cfg:     cfg,
	}

	// 获取终端尺寸
	app.termWidth, app.termHeight, _ = term.GetSize(int(os.Stdout.Fd()))

	// 设置原始模式
	var err error
	app.oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "[错误] 无法设置终端模式: %v\n", err)
		os.Exit(1)
	}
	defer term.Restore(int(os.Stdin.Fd()), app.oldState)

	// 启用备用屏幕
	fmt.Print("\x1b[?1049h")
	defer fmt.Print("\x1b[?1049l")

	app.client = mcclient.New(cfg)
	app.client.SetChatHandler(func(msg mcclient.ChatMessage) {
		var text string
		if msg.RawJSON != "" {
			// 原始JSON只写日志文件
			logx.Debugf("chat raw: %s", msg.RawJSON)
			comp, err := mcclient.ParseTextComponent(msg.RawJSON)
			if err == nil {
				text = comp.ToPlain()
			} else {
				text = msg.PlainText
			}
		} else {
			text = msg.PlainText
		}
		if text != "" {
			app.addLog(text)
		}
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.client.Run(ctx)
	}()

	// 显示初始信息
	app.addLog(fmt.Sprintf("gmcc v%s", Version))
	app.addLog(fmt.Sprintf("玩家: %s", cfg.Account.PlayerID))
	app.addLog(fmt.Sprintf("服务器: %s", cfg.Server.Address))
	app.addLog("正在连接...")
	app.addLog("")

	// 输入缓冲
	input := ""
	history := make([]string, 0, 100)
	histIdx := -1

	// 渲染初始界面
	app.render(input)

	// 输入循环
	inputCh := make(chan byte, 256)
	go func() {
		buf := make([]byte, 1)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil || n == 0 {
				close(inputCh)
				return
			}
			select {
			case inputCh <- buf[0]:
			default:
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			app.addLog("[提示] 断开连接")
			time.Sleep(500 * time.Millisecond)
			return
		case err := <-errCh:
			if err != nil {
				app.addLog(fmt.Sprintf("[错误] %v", err))
			}
			time.Sleep(500 * time.Millisecond)
			return
		case b, ok := <-inputCh:
			if !ok {
				return
			}

			// 处理输入
			switch b {
			case 3: // Ctrl+C
				app.addLog("[提示] 退出")
				time.Sleep(500 * time.Millisecond)
				return

			case 13, 10: // Enter
				if input == "" {
					continue
				}

				// 保存历史
				if len(history) == 0 || history[len(history)-1] != input {
					history = append(history, input)
					if len(history) > 100 {
						history = history[1:]
					}
				}
				histIdx = -1

				// 处理输入
				app.processInput(input)
				input = ""
				app.render(input)

			case 127, 8: // Backspace
				if len(input) > 0 {
					input = input[:len(input)-1]
					app.render(input)
				}

			case 27: // Escape sequence
				// 读取后续字符
				seq := make([]byte, 0, 4)
				seqDone := false
				timeout := time.NewTimer(50 * time.Millisecond)
				for !seqDone {
					select {
					case b2, ok2 := <-inputCh:
						if !ok2 {
							timeout.Stop()
							continue
						}
						seq = append(seq, b2)
						if len(seq) >= 2 && seq[0] == '[' {
							// 处理方向键
							switch seq[1] {
							case 'A': // Up
								if len(history) > 0 {
									if histIdx == -1 {
										histIdx = len(history) - 1
									} else if histIdx > 0 {
										histIdx--
									}
									input = history[histIdx]
									app.render(input)
								}
								seqDone = true
							case 'B': // Down
								if histIdx != -1 {
									if histIdx < len(history)-1 {
										histIdx++
										input = history[histIdx]
									} else {
										histIdx = -1
										input = ""
									}
									app.render(input)
								}
								seqDone = true
							case 'C', 'D': // Left, Right
								seqDone = true
							}
						}
						if len(seq) >= 3 {
							seqDone = true
						}
					case <-timeout.C:
						seqDone = true
					}
				}

			default:
				// 普通字符
				if b >= 32 && b < 127 {
					input += string(b)
					app.render(input)
				}
			}
		}
	}
}

func (a *App) addLog(text string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.logs = append(a.logs, text)
	if len(a.logs) > a.maxLogs {
		a.logs = a.logs[1:]
	}
}

func (a *App) render(input string) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// 清屏
	fmt.Print("\x1b[2J\x1b[H")

	// 获取终端尺寸
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	if w == 0 {
		w = 80
	}
	if h == 0 {
		h = 24
	}

	// 标题
	title := fmt.Sprintf(" gmcc v%s ", Version)
	if len(title) < w {
		title = title + strings.Repeat(" ", w-len(title))
	}
	fmt.Printf("\x1b[1;36m%s\x1b[0m\r\n", title)

	// 信息栏
	player := a.client.Player
	if player != nil {
		hp, _, food, _ := player.GetHealth()
		x, y, z := player.GetPosition()
		info := fmt.Sprintf(" HP:%.0f Food:%d Pos:%.0f,%.0f,%.0f Mode:%s", hp, food, x, y, z, player.GameMode.String())
		if len(info) > w {
			info = info[:w]
		}
		fmt.Printf("\x1b[1;32m%s\x1b[0m\r\n", info)
	} else {
		info := fmt.Sprintf(" %s ", a.cfg.Server.Address)
		if len(info) < w {
			info = info + strings.Repeat(" ", w-len(info))
		}
		fmt.Printf("\x1b[1;33m%s\x1b[0m\r\n", info)
	}

	// 分隔线
	line := strings.Repeat("─", w)
	fmt.Printf("%s\r\n", line)

	// 消息区
	msgHeight := h - 6 // 预留: 标题1 + 信息1 + 分隔线1 + 输入提示1 + 输入1 = 5
	if msgHeight < 1 {
		msgHeight = 1
	}

	start := 0
	if len(a.logs) > msgHeight {
		start = len(a.logs) - msgHeight
	}

	for i, log := range a.logs[start:] {
		if i >= msgHeight {
			break
		}
		// 截断超长行
		runes := []rune(log)
		if len(runes) > w {
			runes = runes[:w]
		}
		fmt.Printf("%s\r\n", string(runes))
	}

	// 填充空行
	for i := len(a.logs[start:]); i < msgHeight; i++ {
		fmt.Println()
	}

	// 分隔线
	fmt.Printf("%s\r\n", line)

	// 输入区
	fmt.Print("\x1b[1;37m> \x1b[0m")

	// 输入内容
	inputRunes := []rune(input)
	if len(inputRunes) > w-3 {
		inputRunes = inputRunes[:w-3]
	}
	fmt.Print(string(inputRunes))
}

func (a *App) processInput(line string) {
	if strings.HasPrefix(line, "/") {
		cmd := strings.TrimPrefix(line, "/")
		if a.client.IsReady() {
			if err := a.client.SendCommand(cmd); err != nil {
				a.addLog(fmt.Sprintf("\x1b[31m[错误] %v\x1b[0m", err))
			}
		} else {
			a.addLog("\x1b[33m[提示] 尚未连接\x1b[0m")
		}
	} else {
		if a.client.IsReady() {
			if err := a.client.SendMessage(line); err != nil {
				a.addLog(fmt.Sprintf("\x1b[31m[错误] %v\x1b[0m", err))
			}
		} else {
			a.addLog("\x1b[33m[提示] 尚未连接\x1b[0m")
		}
	}
}
