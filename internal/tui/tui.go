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
	"gmcc/internal/mcclient/chat"
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
			comp, err := chat.ParseTextComponent(msg.RawJSON)
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

	inputCh := make(chan []byte, 256)
	go t.readInput(inputCh)

	t.render()

	for {
		select {
		case <-ctx.Done():
			t.addLog("\x1b[33m[提示] 断开连接\x1b[0m")
			t.render()
			time.Sleep(500 * time.Millisecond)
			return nil
		case err := <-errCh:
			if err != nil {
				t.addLog(fmt.Sprintf("\x1b[31m[错误] %v\x1b[0m", err))
				t.render()
				time.Sleep(2 * time.Second)
			}
			return err
		case data, ok := <-inputCh:
			if !ok {
				return nil
			}
			if t.handleInput(data, inputCh) {
				return nil
			}
		case <-t.redrawCh:
			t.render()
		}
	}
}

func (t *TUI) readInput(ch chan []byte) {
	buf := make([]byte, 4)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			close(ch)
			return
		}
		data := make([]byte, n)
		copy(data, buf[:n])
		select {
		case ch <- data:
		default:
		}
	}
}

func (t *TUI) handleInput(data []byte, inputCh chan []byte) bool {
	if len(data) == 0 {
		return false
	}

	b := data[0]
	switch b {
	case 3: // Ctrl+C
		return t.handleExit()
	case 13, 10: // Enter
		t.handleEnter()
	case 127, 8: // Backspace
		t.handleBackspace()
	case 27: // ESC
		t.handleEscapeSequence(inputCh)
	default:
		if b >= 32 {
			t.inputBuf += string(data)
			t.requestRedraw()
		}
	}
	return false
}

func (t *TUI) handleExit() bool {
	t.addLog("\x1b[33m[提示] 退出\x1b[0m")
	t.render()
	time.Sleep(300 * time.Millisecond)
	return true
}

func (t *TUI) handleEnter() {
	if t.inputBuf == "" {
		return
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
}

func (t *TUI) handleBackspace() {
	if len(t.inputBuf) > 0 {
		runes := []rune(t.inputBuf)
		if len(runes) > 0 {
			t.inputBuf = string(runes[:len(runes)-1])
		}
		t.requestRedraw()
	}
}

func (t *TUI) handleEscapeSequence(inputCh chan []byte) {
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
			seq = append(seq, b2...)
			if len(seq) >= 2 && seq[0] == '[' {
				if t.handleAnsiControlSequence(seq[1]) {
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
}

func (t *TUI) handleAnsiControlSequence(cmd byte) bool {
	switch cmd {
	case 'A': // Up
		if len(t.history) > 0 {
			if t.histIdx == -1 {
				t.histIdx = len(t.history) - 1
			} else if t.histIdx > 0 {
				t.histIdx--
			}
			t.inputBuf = t.history[t.histIdx]
			t.requestRedraw()
		}
		return true
	case 'B': // Down
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
		return true
	case 'C', 'D': // Right, Left
		return true
	}
	return false
}

func (t *TUI) processInput(line string) {
	if strings.HasPrefix(line, "/") {
		cmd := strings.TrimPrefix(line, "/")
		if t.client.IsReady() {
			// 检查是否是强制使用签名或无签名的前缀
			if strings.HasPrefix(cmd, "!") {
				// 强制使用无签名命令
				cmd = strings.TrimPrefix(cmd, "!")
				if err := t.client.SendCommandUnsigned(cmd); err != nil {
					t.addLog(fmt.Sprintf("\x1b[31m[错误] %v\x1b[0m", err))
				} else {
					t.addLog(fmt.Sprintf("\x1b[90m[命令] /%s (无签名)\x1b[0m", cmd))
				}
			} else {
				// 使用默认签名行为
				signed := t.cfg.Actions.DefaultSignCommands
				if err := t.client.SendCommand(cmd); err != nil {
					t.addLog(fmt.Sprintf("\x1b[31m[错误] %v\x1b[0m", err))
				} else {
					if signed {
						t.addLog(fmt.Sprintf("\x1b[90m[命令] /%s (签名)\x1b[0m", cmd))
					} else {
						t.addLog(fmt.Sprintf("\x1b[90m[命令] /%s (无签名)\x1b[0m", cmd))
					}
				}
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
		line := truncateANSI(log, w)
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

func visibleLen(s string) int {
	len := 0
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		len++
	}
	return len
}

func truncateANSI(s string, maxLen int) string {
	if visibleLen(s) <= maxLen {
		return s
	}
	var sb strings.Builder
	visible := 0
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			sb.WriteRune(r)
			continue
		}
		if inEscape {
			sb.WriteRune(r)
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		if visible >= maxLen {
			break
		}
		sb.WriteRune(r)
		visible++
	}
	sb.WriteString("\x1b[0m")
	return sb.String()
}
