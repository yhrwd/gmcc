package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Terminal struct {
	mu          sync.Mutex
	messageHook func(string)
	cmdHook     func(string)
	stopCh      chan struct{}
}

func New() *Terminal {
	return &Terminal{
		stopCh: make(chan struct{}),
	}
}

func (t *Terminal) Start() error {
	go t.readLoop()
	return nil
}

func (t *Terminal) Stop() {
	close(t.stopCh)
}

func (t *Terminal) SetMessageHook(hook func(string)) {
	t.messageHook = hook
}

func (t *Terminal) SetCommandHook(hook func(string)) {
	t.cmdHook = hook
}

func (t *Terminal) Printf(format string, args ...any) {
	t.mu.Lock()
	defer t.mu.Unlock()
	fmt.Printf(format+"\n", args...)
}

func (t *Terminal) PrintLine(line string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	fmt.Print("\r\033[K")
	fmt.Println(line)
	fmt.Print("> ")
}

func (t *Terminal) readLoop() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")

	for {
		select {
		case <-t.stopCh:
			return
		default:
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			fmt.Print("> ")
			continue
		}

		if strings.HasPrefix(line, "/") {
			if t.cmdHook != nil {
				t.cmdHook(strings.TrimPrefix(line, "/"))
			}
		} else {
			if t.messageHook != nil {
				t.messageHook(line)
			}
		}

		fmt.Print("> ")
	}
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
