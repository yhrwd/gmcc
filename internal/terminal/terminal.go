package terminal

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Terminal struct {
	mu           sync.Mutex
	app          *tview.Application
	pages        *tview.Pages
	chatView     *tview.TextView
	inputField   *tview.InputField
	cmdHook      func(string)
	msgHook      func(string)
	history      []string
	historyIndex int
	currentInput string
	stopped      bool
	connected    bool
}

func New() *Terminal {
	return &Terminal{
		history:      []string{},
		historyIndex: -1,
	}
}

func (t *Terminal) Start() error {
	t.app = tview.NewApplication()
	t.app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		return false
	})

	t.chatView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	t.chatView.SetBorderPadding(1, 1, 1, 0)

	t.inputField = tview.NewInputField().
		SetPlaceholder("输入聊天消息或命令 /help...").
		SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetLabel("> ").
		SetLabelColor(tview.Styles.PrimaryTextColor)
	t.inputField.SetBorderPadding(0, 0, 1, 1)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(t.chatView, 0, 1, false).
		AddItem(t.inputField, 1, 0, true)

	t.pages = tview.NewPages().
		AddPage("main", flex, true, true)

	t.inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			t.submitInput()
		} else if key == tcell.KeyUp {
			t.historyUp()
		} else if key == tcell.KeyDown {
			t.historyDown()
		} else if key == tcell.KeyTab {
			t.completeCommand()
		}
	})

	t.app.SetRoot(t.pages, true)
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			t.app.Stop()
			return nil
		}
		return event
	})

	go func() {
		if err := t.app.Run(); err != nil {
		}
	}()

	return nil
}

func (t *Terminal) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.stopped = true
	if t.app != nil {
		t.app.Stop()
	}
}

func (t *Terminal) SetCommandHook(hook func(string)) {
	t.cmdHook = hook
}

func (t *Terminal) SetMessageHook(hook func(string)) {
	t.msgHook = hook
}

func (t *Terminal) submitInput() {
	input := t.inputField.GetText()
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	t.addToHistory(input)

	if strings.HasPrefix(input, "/") {
		if t.cmdHook != nil {
			t.cmdHook(strings.TrimPrefix(input, "/"))
		}
	} else {
		if t.msgHook != nil {
			t.msgHook(input)
		}
	}

	t.inputField.SetText("")
	t.historyIndex = -1
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

func (t *Terminal) historyUp() {
	if len(t.history) == 0 {
		return
	}

	if t.historyIndex == -1 {
		t.currentInput = t.inputField.GetText()
		t.historyIndex = len(t.history) - 1
	} else if t.historyIndex > 0 {
		t.historyIndex--
	}

	t.inputField.SetText(t.history[t.historyIndex])
}

func (t *Terminal) historyDown() {
	if t.historyIndex == -1 {
		return
	}

	if t.historyIndex < len(t.history)-1 {
		t.historyIndex++
		t.inputField.SetText(t.history[t.historyIndex])
	} else {
		t.historyIndex = -1
		t.inputField.SetText(t.currentInput)
	}
}

func (t *Terminal) completeCommand() {
	input := t.inputField.GetText()
	if !strings.HasPrefix(input, "/") {
		return
	}

	cmds := []string{
		"/help", "/quit", "/clear", "/tps", "/money",
		"/balance", "/pay", "/msg", "/tell", "/r",
		"/afk", "/near", "/spawn", "/warp", "/bal",
	}

	var matches []string
	for _, cmd := range cmds {
		if strings.HasPrefix(cmd, input) {
			matches = append(matches, cmd)
		}
	}

	if len(matches) == 1 {
		t.inputField.SetText(matches[0] + " ")
	} else if len(matches) > 1 {
		t.AppendMessage(strings.Join(matches, "  "))
	}
}

func (t *Terminal) AppendMessage(msg string) {
	t.app.QueueUpdateDraw(func() {
		fmt.Fprintln(t.chatView, msg)
		t.chatView.ScrollToEnd()
	})
}

func (t *Terminal) PrintLine(line string) {
	t.AppendMessage(line)
}

func (t *Terminal) Printf(format string, args ...any) {
	t.AppendMessage(fmt.Sprintf(format, args...))
}

func (t *Terminal) SetConnected(connected bool) {
	t.connected = connected
	t.app.QueueUpdateDraw(func() {
		if connected {
			t.pages.ShowPage("main")
		}
	})
}

func (t *Terminal) SetTitle(title string) {
	t.app.QueueUpdateDraw(func() {
		t.pages.SetTitle(title)
	})
}

func (t *Terminal) ShowConnect() {
	t.app.QueueUpdateDraw(func() {
		t.pages.ShowPage("connect")
	})
}

func (t *Terminal) ShowMain() {
	t.app.QueueUpdateDraw(func() {
		t.pages.ShowPage("main")
	})
}

func (t *Terminal) SetStatusBar(text string) {
	t.app.QueueUpdateDraw(func() {
	})
}

func IsTerminal() bool {
	return true
}

func ReadPassword(prompt string) (string, error) {
	return "", nil
}
