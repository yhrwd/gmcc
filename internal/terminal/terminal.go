package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"golang.org/x/term"
)

type Terminal struct {
	mu           sync.Mutex
	width        int
	height       int
	inputLine    string
	cursorPos    int
	history      []string
	historyIndex int
	completions  []string
	compIndex    int
	showComps    bool
	messageHook  func(string)
	cmdHook      func(string)
	oldState     *term.State
}

func New() *Terminal {
	return &Terminal{
		history: make([]string, 0, 100),
	}
}

func (t *Terminal) Start() error {
	var err error
	t.oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	t.updateSize()
	go t.readLoop()

	return nil
}

func (t *Terminal) Stop() {
	if t.oldState != nil {
		term.Restore(int(os.Stdin.Fd()), t.oldState)
	}
}

func (t *Terminal) updateSize() {
	w, h, err := term.GetSize(int(os.Stdin.Fd()))
	if err == nil {
		t.width = w
		t.height = h
	}
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
	t.printLine(fmt.Sprintf(format, args...))
}

func (t *Terminal) PrintLine(line string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.printLine(line)
}

func (t *Terminal) printLine(line string) {
	t.clearInputLine()
	fmt.Print("\r\033[K")
	fmt.Println(line)
	t.drawInputLine()
}

func (t *Terminal) clearInputLine() {
	fmt.Print("\r\033[K")
}

func (t *Terminal) drawInputLine() {
	prompt := "> "
	input := t.inputLine
	display := prompt + input

	fmt.Print("\r\033[K")
	fmt.Print(display)

	if t.cursorPos < len(input) {
		fmt.Printf("\033[%dD", len(input)-t.cursorPos)
	}
}

func (t *Terminal) readLoop() {
	buf := make([]byte, 128)

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}

		t.handleInput(buf[:n])
	}
}

func (t *Terminal) handleInput(data []byte) {
	t.mu.Lock()
	defer t.mu.Unlock()

	i := 0
	for i < len(data) {
		b := data[i]

		switch b {
		case 13, 10:
			t.executeLine()
		case 127, 8:
			t.backspace()
		case 27:
			if i+2 < len(data) && data[i+1] == '[' {
				end := i + 2
				for end < len(data) && data[end] >= '0' && data[end] <= '9' {
					end++
				}
				if end < len(data) {
					end++
				}
				t.handleEscapeSeq(data[i:end])
				i = end
				continue
			}
		case 1:
			t.cursorPos = 0
			t.drawInputLine()
		case 5:
			t.cursorPos = len(t.inputLine)
			t.drawInputLine()
		case 11:
			t.inputLine = t.inputLine[:t.cursorPos]
			t.drawInputLine()
		case 21:
			t.inputLine = t.inputLine[t.cursorPos:]
			t.cursorPos = 0
			t.drawInputLine()
		case 9:
			t.tabComplete()
		case 4:
			if len(t.inputLine) == 0 {
				fmt.Print("\r\n")
				t.Stop()
				os.Exit(0)
			}
		default:
			if b >= 32 && b < 127 {
				t.insertChar(rune(b))
			}
		}
		i++
	}
}

func (t *Terminal) handleEscapeSeq(seq []byte) {
	if len(seq) < 2 {
		return
	}

	if seq[0] == '[' {
		switch seq[1] {
		case 'A':
			if t.historyIndex > 0 {
				t.historyIndex--
				t.inputLine = t.history[t.historyIndex]
				t.cursorPos = len(t.inputLine)
				t.drawInputLine()
			}
		case 'B':
			if t.historyIndex < len(t.history)-1 {
				t.historyIndex++
				t.inputLine = t.history[t.historyIndex]
			} else {
				t.historyIndex = len(t.history)
				t.inputLine = ""
			}
			t.cursorPos = len(t.inputLine)
			t.drawInputLine()
		case 'C':
			if t.cursorPos < len(t.inputLine) {
				t.cursorPos++
				t.drawInputLine()
			}
		case 'D':
			if t.cursorPos > 0 {
				t.cursorPos--
				t.drawInputLine()
			}
		case 'H':
			t.cursorPos = 0
			t.drawInputLine()
		case 'F':
			t.cursorPos = len(t.inputLine)
			t.drawInputLine()
		case '3':
			if len(seq) > 2 && seq[2] == '~' {
				if t.cursorPos < len(t.inputLine) {
					t.inputLine = t.inputLine[:t.cursorPos] + t.inputLine[t.cursorPos+1:]
					t.drawInputLine()
				}
			}
		}
	}
}

func (t *Terminal) insertChar(c rune) {
	if t.cursorPos == len(t.inputLine) {
		t.inputLine += string(c)
	} else {
		t.inputLine = t.inputLine[:t.cursorPos] + string(c) + t.inputLine[t.cursorPos:]
	}
	t.cursorPos++
	t.drawInputLine()
}

func (t *Terminal) backspace() {
	if t.cursorPos > 0 {
		t.inputLine = t.inputLine[:t.cursorPos-1] + t.inputLine[t.cursorPos:]
		t.cursorPos--
		t.drawInputLine()
	}
}

func (t *Terminal) executeLine() {
	line := strings.TrimSpace(t.inputLine)
	if line == "" {
		t.inputLine = ""
		t.cursorPos = 0
		fmt.Print("\r\n")
		t.drawInputLine()
		return
	}

	t.history = append(t.history, line)
	t.historyIndex = len(t.history)

	t.inputLine = ""
	t.cursorPos = 0
	fmt.Print("\r\n")

	if strings.HasPrefix(line, "/") {
		if t.cmdHook != nil {
			t.cmdHook(strings.TrimPrefix(line, "/"))
		}
	} else {
		if t.messageHook != nil {
			t.messageHook(line)
		}
	}

	t.drawInputLine()
}

func (t *Terminal) tabComplete() {
	if t.completions == nil || len(t.completions) == 0 {
		return
	}

	t.compIndex = (t.compIndex + 1) % len(t.completions)
	comp := t.completions[t.compIndex]

	parts := strings.Split(t.inputLine, " ")
	if len(parts) > 0 {
		parts[len(parts)-1] = comp
		t.inputLine = strings.Join(parts, " ")
	} else {
		t.inputLine = comp
	}
	t.cursorPos = len(t.inputLine)
	t.drawInputLine()
}

func (t *Terminal) SetCompletions(comps []string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.completions = comps
	t.compIndex = -1
}

func (t *Terminal) GetInputLine() string {
	return t.inputLine
}

func ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	defer fmt.Println()
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	return string(b), err
}

func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func NewReader() *bufio.Reader {
	return bufio.NewReader(os.Stdin)
}
