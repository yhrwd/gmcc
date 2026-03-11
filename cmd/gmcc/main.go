package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

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

	runInteractive(configPath)
}

func runInteractive(configPath string) {
	app := tview.NewApplication()
	pages := tview.NewPages()

	header := tview.NewTextView().
		SetDynamicColors(true).
		SetText(fmt.Sprintf(`[cyan]╔══════════════════════════════════════╗
║     [white]gmcc[cyan] %s                  ║
║     [white]Minecraft 控制台客户端[cyan]         ║
╚══════════════════════════════════════╝[white]`, Version)).
		SetTextAlign(tview.AlignCenter)

	statusText := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]正在初始化...")

	connectInfo := tview.NewTextView().
		SetDynamicColors(true).
		SetText("")

	chatView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	chatView.SetBorderPadding(1, 1, 1, 0)

	mainInput := tview.NewInputField().
		SetPlaceholder("输入聊天消息或命令 /help...").
		SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetLabel("[cyan]> [white]").
		SetLabelColor(tcell.GetColor("cyan"))

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(chatView, 0, 1, false).
		AddItem(mainInput, 1, 0, true)

	mainPage := tview.NewPages().AddPage("main", flex, true, true)

	form := tview.NewForm().
		AddTextView("", fmt.Sprintf("版本: %s", Version), 0, 1, false, false).
		AddTextView("", fmt.Sprintf("配置: %s", configPath), 0, 1, false, false).
		AddTextView("", "", 0, 1, false, false).
		AddButton("连接服务器", func() {
			pages.ShowPage("main")
			app.SetFocus(mainInput)
		})

	connectPage := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 8, 0, false).
		AddItem(statusText, 1, 0, false).
		AddItem(connectInfo, 0, 1, false).
		AddItem(form, 0, 3, true)

	pages.AddPage("connect", connectPage, true, true).
		AddPage("main", mainPage, false, false)

	pages.ShowPage("connect")

	cfg, err := config.Load(configPath)
	if err != nil {
		statusText.SetText(fmt.Sprintf("[red]配置加载失败: %v", err))
		app.SetRoot(pages, true)
		app.Run()
		return
	}

	if err := logx.Init(cfg.Log.LogDir, cfg.Log.EnableFile, cfg.Log.MaxSize, cfg.Log.Debug); err != nil {
		statusText.SetText(fmt.Sprintf("[red]日志初始化失败: %v", err))
		app.SetRoot(pages, true)
		app.Run()
		return
	}
	defer logx.Close()

	statusText.SetText(fmt.Sprintf("[green]✓ 配置加载成功[white] | 玩家: [cyan]%s[white] | 服务器: [cyan]%s",
		cfg.Account.PlayerID, cfg.Server.Address))

	connectInfo.SetText(fmt.Sprintf(`
[white]服务器: %s
正版认证: %v
签名命令: %v
`, cfg.Server.Address, cfg.Account.UseOfficialAuth, cfg.Actions.SignCommands))

	client := mcclient.New(cfg)

	client.SetChatHandler(func(msg mcclient.ChatMessage) {
		app.QueueUpdateDraw(func() {
			if msg.RawJSON != "" {
				comp, err := mcclient.ParseTextComponent(msg.RawJSON)
				if err == nil {
					fmt.Fprintln(chatView, comp.ToANSI())
					chatView.ScrollToEnd()
					return
				}
			}
			fmt.Fprintln(chatView, msg.PlainText)
			chatView.ScrollToEnd()
		})
	})

	var cmdHistory []string
	historyIndex := -1

	mainInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			input := mainInput.GetText()
			input = strings.TrimSpace(input)
			if input == "" {
				return
			}

			cmdHistory = append(cmdHistory, input)
			if len(cmdHistory) > 100 {
				cmdHistory = cmdHistory[1:]
			}
			historyIndex = -1

			if strings.HasPrefix(input, "/") {
				if client.IsReady() {
					if err := client.SendCommand(strings.TrimPrefix(input, "/")); err != nil {
						fmt.Fprintf(chatView, "[red][命令失败] %v[white]\n", err)
						chatView.ScrollToEnd()
					}
				}
			} else {
				if client.IsReady() {
					if err := client.SendMessage(input); err != nil {
						fmt.Fprintf(chatView, "[red][发送失败] %v[white]\n", err)
						chatView.ScrollToEnd()
					}
				}
			}

			mainInput.SetText("")
		} else if key == tcell.KeyUp {
			if len(cmdHistory) == 0 {
				return
			}
			if historyIndex == -1 {
				historyIndex = len(cmdHistory) - 1
			} else if historyIndex > 0 {
				historyIndex--
			}
			mainInput.SetText(cmdHistory[historyIndex])
		} else if key == tcell.KeyDown {
			if historyIndex == -1 {
				return
			}
			if historyIndex < len(cmdHistory)-1 {
				historyIndex++
				mainInput.SetText(cmdHistory[historyIndex])
			} else {
				historyIndex = -1
				mainInput.SetText("")
			}
		} else if key == tcell.KeyTab {
			input := mainInput.GetText()
			if !strings.HasPrefix(input, "/") {
				return
			}
			cmds := []string{"/help", "/quit", "/tps", "/money", "/balance", "/pay", "/msg", "/tell", "/r", "/afk", "/near", "/spawn", "/warp", "/bal", "/suicide"}
			var matches []string
			for _, cmd := range cmds {
				if strings.HasPrefix(cmd, input) {
					matches = append(matches, cmd)
				}
			}
			if len(matches) == 1 {
				mainInput.SetText(matches[0] + " ")
			} else if len(matches) > 1 {
				fmt.Fprintln(chatView, strings.Join(matches, "  "))
				chatView.ScrollToEnd()
			}
		}
	})

	app.SetRoot(pages, true)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			app.Stop()
			return nil
		}
		return event
	})

	go func() {
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
		defer stop()

		errCh := make(chan error, 1)
		go func() {
			errCh <- client.Run(ctx)
		}()

		app.QueueUpdateDraw(func() {
			pages.ShowPage("main")
			fmt.Fprintf(chatView, "[green]正在连接 %s ...\n", cfg.Server.Address)
			chatView.ScrollToEnd()
		})

		select {
		case err := <-errCh:
			app.QueueUpdateDraw(func() {
				if err != nil {
					fmt.Fprintf(chatView, "[red]客户端退出: %v[white]\n", err)
				} else {
					fmt.Fprintln(chatView, "[yellow]客户端已断开连接[white]")
				}
				chatView.ScrollToEnd()
			})
		case <-ctx.Done():
			app.QueueUpdateDraw(func() {
				fmt.Fprintln(chatView, "[yellow]正在断开连接...[white]")
				chatView.ScrollToEnd()
			})
		}
	}()

	if err := app.Run(); err != nil {
		fmt.Printf("UI Error: %v\n", err)
	}
}
