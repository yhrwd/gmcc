package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/term"

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
	fmt.Println()
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

	client := mcclient.New(cfg)

	client.SetChatHandler(func(msg mcclient.ChatMessage) {
		if msg.RawJSON != "" {
			comp, err := mcclient.ParseTextComponent(msg.RawJSON)
			if err == nil {
				fmt.Println(comp.ToANSI())
				return
			}
		}
		fmt.Println(msg.PlainText)
	})

	fmt.Printf("正在连接 %s ...\n\n", cfg.Server.Address)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- client.Run(ctx)
	}()

	history := []string{}
	historyIndex := -1
	input := ""

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err == nil {
		defer term.Restore(int(os.Stdin.Fd()), oldState)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			buf := make([]byte, 1)
			n, err := os.Stdin.Read(buf)
			if err != nil || n == 0 {
				return
			}

			switch buf[0] {
			case 13: // Enter
				fmt.Println()
				if input != "" {
					history = append(history, input)
					if len(history) > 100 {
						history = history[1:]
					}
					historyIndex = -1

					if strings.HasPrefix(input, "/") {
						if client.IsReady() {
							if err := client.SendCommand(strings.TrimPrefix(input, "/")); err != nil {
								fmt.Printf("[命令失败] %v\n", err)
							}
						}
					} else {
						if client.IsReady() {
							if err := client.SendMessage(input); err != nil {
								fmt.Printf("[发送失败] %v\n", err)
							}
						}
					}
				}
				input = ""
				fmt.Print("\033[2K\r> ")

			case 127: // Backspace
				if len(input) > 0 {
					input = input[:len(input)-1]
					fmt.Print("\b \b")
				}

			case 27: // Escape
				seq := make([]byte, 2)
				n, _ := os.Stdin.Read(seq)
				if n < 2 {
					continue
				}
				switch seq[1] {
				case 65: // Up
					if len(history) > 0 {
						fmt.Print("\r\033[K> ")
						if historyIndex == -1 {
							historyIndex = len(history) - 1
						} else if historyIndex > 0 {
							historyIndex--
						}
						input = history[historyIndex]
						fmt.Print(input)
					}
				case 66: // Down
					if historyIndex != -1 {
						fmt.Print("\r\033[K> ")
						if historyIndex < len(history)-1 {
							historyIndex++
							input = history[historyIndex]
						} else {
							historyIndex = -1
							input = ""
						}
						fmt.Print(input)
					}
				case 9: // Tab
					if strings.HasPrefix(input, "/") {
						cmds := []string{"/help", "/quit", "/tps", "/money", "/balance", "/pay", "/msg", "/tell", "/r", "/afk", "/near", "/spawn", "/warp", "/bal"}
						var matches []string
						for _, cmd := range cmds {
							if strings.HasPrefix(cmd, input) {
								matches = append(matches, cmd)
							}
						}
						if len(matches) == 1 {
							fmt.Print("\r\033[K> ")
							input = matches[0] + " "
							fmt.Print(input)
						} else if len(matches) > 1 {
							fmt.Println()
							for _, m := range matches {
								fmt.Print(m, " ")
							}
							fmt.Println()
							fmt.Print("> " + input)
						}
					}
				}

			default:
				if buf[0] >= 32 && buf[0] < 127 {
					input += string(buf[0])
					fmt.Print(string(buf[0]))
				}
			}
		}
	}()

	select {
	case err := <-errCh:
		if err != nil {
			fmt.Printf("\n[客户端退出] %v\n", err)
		}
	case <-ctx.Done():
		fmt.Println("\n[提示] 断开连接")
	}
}
