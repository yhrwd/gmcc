package tui

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
	"gmcc/internal/mcclient"
)

type TUI struct {
	cfg      *config.Config
	client   *mcclient.Client
	mu       sync.RWMutex
	logs     []string
	maxLogs  int
	oldState *term.State
	running  bool
	quit     chan struct{}
	redrawCh chan struct{}
	inputBuf string
	history  []string
	histIdx  int
}

func New(cfg *config.Config) *TUI {
	return &TUI{
		cfg:      cfg,
		maxLogs:  500,
		logs:     make([]string, 0, 500),
		quit:     make(chan struct{}),
		redrawCh: make(chan struct{}, 1),
		history:  make([]string, 0, 100),
		histIdx:  -1,
	}
}

func (t *TUI) Run(ctx context.Context) error {
	var err error
	t.oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("设置终端原始模式失败: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), t.oldState)

	fmt.Print("\x1b[?1049h")
	defer fmt.Print("\x1b[?1049l")

	fmt.Print("\x1b[?25l")
	defer fmt.Print("\x1b[?25h")

	t.running = true

	t.client = mcclient.New(t.cfg)
	t.client.SetChatHandler(func(msg mcclient.ChatMessage) {
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
		if text != "" {
			t.addLog(text)
		}
	})

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- t.client.Run(ctx)
	}()

	t.addLog(fmt.Sprintf("\x1b[1;36mgmcc\x1b[0m v%s", "dev"))
	t.addLog(fmt.Sprintf("\x1b[33m玩家: %s\x1b[0m", t.cfg.Account.PlayerID))
	t.addLog(fmt.Sprintf("\x1b[33m服务器: %s\x1b[0m", t.cfg.Server.Address))
	t.addLog("\x1b[90m正在连接...\x1b[0m")
	t.addLog("")

	inputCh := make(chan byte, 256)
	go t.readInput(inputCh)

	t.render()

	for {
		select {
		case <-ctx.Done():
			t.addLog("\x1b[33m[提示] 断开连接\x1b[0m")
			time.Sleep(500 * time.Millisecond)
			return nil
		case err := <-errCh:
			if err != nil {
				t.addLog(fmt.Sprintf("\x1b[31m[错误] %v\x1b[0m", err))
				time.Sleep(2 * time.Second)
			}
			return err
		case b, ok := <-inputCh:
			if !ok {
				return nil
			}
			if t.handleInput(b, inputCh) {
				return nil
			}
		case <-t.redrawCh:
			t.render()
		}
	}
}

func (t *TUI) readInput(ch chan byte) {
	buf := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			close(ch)
			return
		}
		select {
		case ch <- buf[0]:
		default:
		}
	}
}

func (t *TUI) handleInput(b byte, inputCh chan byte) bool {
	switch b {
	case 3:
		t.addLog("\x1b[33m[提示] 退出\x1b[0m")
		t.render()
		time.Sleep(300 * time.Millisecond)
		return true
	case 13, 10:
		if t.inputBuf == "" {
			return false
		}
		t.processInput(t.inputBuf)
		if len(t.history) == 0 || t.history[len(t.history)-1] != t.inputBuf {
			t.history = append(t.history, t.inputBuf)
			if len(t.history) > 100 {
				t.history = t.history[1:]
			}
		}
		t.histIdx = -1
		t.inputBuf = ""
		t.requestRedraw()
	case 127, 8:
		if len(t.inputBuf) > 0 {
			t.inputBuf = t.inputBuf[:len(t.inputBuf)-1]
			t.requestRedraw()
		}
	case 27:
		seq := make([]byte, 0, 4)
		timeout := time.NewTimer(50 * time.Millisecond)
		defer timeout.Stop()
	seqLoop:
		for {
			select {
			case b2, ok := <-inputCh:
				if !ok {
					break seqLoop
				}
				seq = append(seq, b2)
				if len(seq) >= 2 && seq[0] == '[' {
					switch seq[1] {
					case 'A':
						if len(t.history) > 0 {
							if t.histIdx == -1 {
								t.histIdx = len(t.history) - 1
							} else if t.histIdx > 0 {
								t.histIdx--
							}
							t.inputBuf = t.history[t.histIdx]
							t.requestRedraw()
						}
						break seqLoop
					case 'B':
						if t.histIdx != -1 {
							if t.histIdx < len(t.history)-1 {
								t.histIdx++
								t.inputBuf = t.history[t.histIdx]
							} else {
								t.histIdx = -1
								t.inputBuf = ""
							}
							t.requestRedraw()
						}
						break seqLoop
					case 'C', 'D':
						break seqLoop
					}
				}
				if len(seq) >= 3 {
					break seqLoop
				}
			case <-timeout.C:
				break seqLoop
			}
		}
	default:
		if b >= 32 && b < 127 {
			t.inputBuf += string(b)
			t.requestRedraw()
		}
	}
	return false
}

func (t *TUI) processInput(line string) {
	if strings.HasPrefix(line, "/") {
		cmd := strings.TrimPrefix(line, "/")
		if t.client.IsReady() {
			if err := t.client.SendCommand(cmd); err != nil {
				t.addLog(fmt.Sprintf("\x1b[31m[错误] %v\x1b[0m", err))
			}
		} else {
			t.addLog("\x1b[33m[提示] 尚未连接\x1b[0m")
		}
	} else {
		if t.client.IsReady() {
			if err := t.client.SendMessage(line); err != nil {
				t.addLog(fmt.Sprintf("\x1b[31m[错误] %v\x1b[0m", err))
			}
		} else {
			t.addLog("\x1b[33m[提示] 尚未连接\x1b[0m")
		}
	}
}

func (t *TUI) addLog(text string) {
	t.mu.Lock()
	t.logs = append(t.logs, text)
	if len(t.logs) > t.maxLogs {
		t.logs = t.logs[1:]
	}
	t.mu.Unlock()
	t.requestRedraw()
}

func (t *TUI) requestRedraw() {
	select {
	case t.redrawCh <- struct{}{}:
	default:
	}
}

func (t *TUI) render() {
	t.mu.RLock()
	defer t.mu.RUnlock()

	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	if w == 0 {
		w = 80
	}
	if h == 0 {
		h = 24
	}

	var sb strings.Builder
	sb.WriteString("\x1b[H")
	sb.WriteString("\x1b[2J")
	sb.WriteString("\x1b[H")

	msgHeight := h - 2
	if msgHeight < 1 {
		msgHeight = 1
	}

	start := 0
	if len(t.logs) > msgHeight {
		start = len(t.logs) - msgHeight
	}

	for i, log := range t.logs[start:] {
		if i >= msgHeight {
			break
		}
		line := t.truncate(log, w)
		sb.WriteString(line)
		sb.WriteString("\r\n")
	}

	for i := len(t.logs[start:]); i < msgHeight; i++ {
		sb.WriteString("\r\n")
	}

	sb.WriteString("\x1b[1;37m")
	sb.WriteString(strings.Repeat("─", w))
	sb.WriteString("\r\n")
	sb.WriteString("\x1b[1;36m")
	sb.WriteString("> ")
	sb.WriteString("\x1b[0m")

	inputRunes := []rune(t.inputBuf)
	if len(inputRunes) > w-3 {
		inputRunes = inputRunes[:w-3]
	}
	sb.WriteString(string(inputRunes))

	sb.WriteString("\x1b[?25h")

	fmt.Print(sb.String())
}

func (t *TUI) truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) > maxLen {
		return string(runes[:maxLen])
	}
	return s
}
