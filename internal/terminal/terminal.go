package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Terminal struct {
	mu           sync.Mutex
	chatLines    []string
	history      []string
	historyIndex int
	currentInput string
	cmdHook      func(string)
	msgHook      func(string)
	stopCh       chan struct{}
	stopped      bool
}

func New() *Terminal {
	return &Terminal{
		chatLines:    []string{},
		history:      []string{},
		historyIndex: -1,
		stopCh:       make(chan struct{}),
	}
}

func (t *Terminal) Start() error {
	go t.readLoop()
	return nil
}

func (t *Terminal) Stop() {
	close(t.stopCh)
	t.stopped = true
}

func (t *Terminal) SetCommandHook(hook func(string)) {
	t.cmdHook = hook
}

func (t *Terminal) SetMessageHook(hook func(string)) {
	t.msgHook = hook
}

func (t *Terminal) Printf(format string, args ...any) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.chatLines = append(t.chatLines, fmt.Sprintf(format, args...))
	if len(t.chatLines) > 200 {
		t.chatLines = t.chatLines[1:]
	}
}

func (t *Terminal) PrintLine(line string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.chatLines = append(t.chatLines, line)
	if len(t.chatLines) > 200 {
		t.chatLines = t.chatLines[1:]
	}
}

func (t *Terminal) readLoop() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println()
	fmt.Println(" ═══ gmcc - Minecraft 控制台客户端 ═══")
	fmt.Println(" ─────────────────────────────────────────────")
	fmt.Println()

	for {
		select {
		case <-t.stopCh:
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

		t.addToHistory(line)

		if strings.HasPrefix(line, "/") {
			if t.cmdHook != nil {
				t.cmdHook(strings.TrimPrefix(line, "/"))
			}
		} else {
			if t.msgHook != nil {
				t.msgHook(line)
			}
		}
	}
}

func (t *Terminal) addToHistory(cmd string) {
	if len(t.history) == 0 || t.history[len(t.history)-1] != cmd {
		t.history = append(t.history, cmd)
		if len(t.history) > 100 {
			t.history = t.history[1:]
		}
	}
	t.historyIndex = -1
}

func (t *Terminal) HistoryUp() {
	if len(t.history) == 0 {
		return
	}

	if t.historyIndex == -1 {
		t.currentInput = t.inputField()
		t.historyIndex = len(t.history) - 1
	} else if t.historyIndex > 0 {
		t.historyIndex--
	}
}

func (t *Terminal) HistoryDown() {
	if t.historyIndex == -1 {
		return
	}

	if t.historyIndex < len(t.history)-1 {
		t.historyIndex++
	} else {
		t.historyIndex = -1
	}
}

func (t *Terminal) inputField() string {
	return ""
}

func IsTerminal() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	fmt.Println()
	return strings.TrimSpace(line), nil
}
