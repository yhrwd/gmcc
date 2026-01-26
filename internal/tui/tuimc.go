package tui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/shlex"

	"gmcc/pkg/fileprocess"
	"gmcc/pkg/logger"
)

// ----------------------------------------------------------------------------
// 消息类型
// ----------------------------------------------------------------------------
type MsgUpdateImage string
type MsgUpdateUserID string
type MsgUpdateStatus struct {
	Addr    string
	Latency string
}
type MsgLog string
type msgCommandResult struct{ output string }

// ----------------------------------------------------------------------------
// 样式
// ----------------------------------------------------------------------------
var (
	colorBlue, colorGreen, colorYellow, colorWhite = lipgloss.Color("#BBDEFB"), lipgloss.Color("#C8E6C9"), lipgloss.Color("#FFF9C4"), lipgloss.Color("#FFFFFF")

	separatorStyle = lipgloss.NewStyle().Foreground(colorWhite)
	labelStyle     = lipgloss.NewStyle().Foreground(colorBlue).Bold(true)
	valueStyle     = lipgloss.NewStyle().Foreground(colorWhite)

	inputFieldStyle = lipgloss.NewStyle().Foreground(colorYellow).Background(lipgloss.Color("#263238")).Height(1)
	promptStyle     = lipgloss.NewStyle().Foreground(colorGreen)
)

func horizontalRule(width int) string {
	return separatorStyle.Render(strings.Repeat("─", width))
}

// ----------------------------------------------------------------------------
// Model
// ----------------------------------------------------------------------------
const maxLogs = 5000

type model struct {
	input      textinput.Model
	quitting   bool
	imagePath  string
	imageCache string

	logs []string

	vp viewport.Model // 使用 viewport 管理日志显示

	userID  string
	addr    string
	latency string

	width, height int

	history    []string
	historyIdx int
}

func NewModel() *model {
	ti := textinput.New()
	ti.Placeholder = "输入命令..."
	ti.Focus()
	ti.CharLimit = 256
	ti.TextStyle = inputFieldStyle
	ti.PlaceholderStyle = valueStyle

	return &model{
		input:      ti,
		imageCache: "No Image",
		logs:       []string{"TUI 启动成功～"},
		userID:     "Unknown",
		addr:       "Disconnected",
		latency:    "0ms",
		history:    []string{},
		historyIdx: -1,
	}
}

func (m *model) Init() tea.Cmd { return textinput.Blink }

// ----------------------------------------------------------------------------
// 日志操作
// ----------------------------------------------------------------------------
func (m *model) addLog(line string) {
	m.logs = append(m.logs, line)
	if len(m.logs) > maxLogs {
		m.logs = m.logs[len(m.logs)-maxLogs:]
	}
	m.updateViewport()
}

func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	var result strings.Builder
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		if i > 0 {
			result.WriteString("\n")
		}

		// 处理 ANSI 转义序列的长度计算
		visibleLen := lipgloss.Width(line)
		if visibleLen <= width {
			result.WriteString(line)
			continue
		}

		// 逐字符切割（简单实现）
		currentWidth := 0
		for _, r := range line {
			charWidth := lipgloss.Width(string(r))
			if currentWidth+charWidth > width {
				result.WriteString("\n")
				currentWidth = 0
			}
			result.WriteRune(r)
			currentWidth += charWidth
		}
	}

	return result.String()
}

func (m *model) updateViewport() {
	if m.vp.Width <= 0 || m.vp.Height <= 0 {
		return
	}

	wrappedLines := make([]string, len(m.logs))
	for i, line := range m.logs {
		wrappedLines[i] = wrapText(line, m.vp.Width)
	}

	m.vp.SetContent(strings.Join(wrappedLines, "\n"))
	m.vp.GotoBottom()
}

// ----------------------------------------------------------------------------
// 命令系统
// ----------------------------------------------------------------------------
type Command struct {
	Name        string
	Description string
	Sub         map[string]Command
	Run         func(m *model, args []string) tea.Cmd
}

var commands map[string]Command

func init() {
	logger.SetUILogFunc(func(msg string) {
		Push(MsgLog(msg))
	})

	commands = map[string]Command{
		"help":   {Name: "help", Description: "显示帮助", Run: cmdHelp},
		"clear":  {Name: "clear", Description: "清空日志", Run: cmdClear},
		"status": {Name: "status", Description: "显示状态", Run: cmdStatus},
		"ping":   {Name: "ping", Description: "测试延迟", Run: cmdPing},
		"user": {Name: "user", Description: "用户命令", Sub: map[string]Command{
			"set": {Name: "set", Description: "设置用户ID", Run: cmdUserSet},
		}},
		"exit": {Name: "exit", Description: "退出程序", Run: func(m *model, args []string) tea.Cmd { return tea.Quit }},
	}
}

func cmdHelp(m *model, args []string) tea.Cmd {
	var b strings.Builder
	b.WriteString("支持命令：\n")
	keys := make([]string, 0, len(commands))
	for k := range commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		c := commands[k]
		fmt.Fprintf(&b, "  %-12s %s\n", c.Name, c.Description)
	}
	return result(b.String())
}

func cmdClear(m *model, args []string) tea.Cmd {
	m.logs = []string{"日志已清空啦～"}
	m.updateViewport()
	return nil
}

func cmdStatus(m *model, args []string) tea.Cmd {
	return result(fmt.Sprintf("当前状态：\n  ID: %s\n  Addr: %s\n  Ping: %s", m.userID, m.addr, m.latency))
}

func cmdPing(m *model, args []string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(500 * time.Millisecond)
		return msgCommandResult{output: fmt.Sprintf("Ping %s (模拟成功)", m.addr)}
	}
}

func cmdUserSet(m *model, args []string) tea.Cmd {
	if len(args) == 0 {
		return result("用法: user set <id>")
	}
	m.userID = args[0]
	return result("UserID 已更新为 " + args[0])
}

func result(s string) tea.Cmd {
	return func() tea.Msg { return msgCommandResult{output: s} }
}

// ----------------------------------------------------------------------------
// 命令解析 & 补全 & 历史
// ----------------------------------------------------------------------------
func (m *model) handleCommand(input string) tea.Cmd {
	parts, err := shlex.Split(input)
	if err != nil || len(parts) == 0 {
		return result("命令解析失败")
	}

	cmdName := parts[0]
	cmd, ok := commands[cmdName]
	if !ok {
		return result("未知命令: " + cmdName)
	}
	args := parts[1:]
	if cmd.Sub != nil && len(args) > 0 {
		if sub, ok := cmd.Sub[args[0]]; ok {
			return sub.Run(m, args[1:])
		}
		return result("未知子命令: " + args[0])
	}
	if cmd.Run == nil {
		return result("未实现该命令")
	}
	return cmd.Run(m, args)
}

func (m *model) completeInput() {
	v := strings.TrimSpace(m.input.Value())
	if v == "" {
		return
	}
	parts := strings.Fields(v)
	last := parts[len(parts)-1]
	var candidates []string

	if len(parts) == 1 {
		for k := range commands {
			if strings.HasPrefix(k, last) {
				candidates = append(candidates, k)
			}
		}
	} else {
		root, ok := commands[parts[0]]
		if ok && root.Sub != nil {
			for k := range root.Sub {
				if strings.HasPrefix(k, last) {
					candidates = append(candidates, k)
				}
			}
		}
	}
	if len(candidates) == 1 {
		parts[len(parts)-1] = candidates[0]
		m.input.SetValue(strings.Join(parts, " ") + " ")
		m.input.CursorEnd()
	}
}

func (m *model) historyPrev() {
	if len(m.history) == 0 {
		return
	}
	if m.historyIdx < 0 {
		m.historyIdx = len(m.history) - 1
	} else if m.historyIdx > 0 {
		m.historyIdx--
	}
	m.input.SetValue(m.history[m.historyIdx])
	m.input.CursorEnd()
}

func (m *model) historyNext() {
	if len(m.history) == 0 || m.historyIdx == -1 {
		return
	}
	if m.historyIdx >= len(m.history)-1 {
		m.historyIdx = -1
		m.input.SetValue("")
	} else {
		m.historyIdx++
		m.input.SetValue(m.history[m.historyIdx])
	}
	m.input.CursorEnd()
}

// ----------------------------------------------------------------------------
// Update
// ----------------------------------------------------------------------------
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.input.Width = m.width - 10
		// 初始化 viewport
		if m.vp.Width != m.width || m.vp.Height != m.height-8 {
			m.vp = viewport.New(m.width, m.height-8)
			m.updateViewport()
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+d":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			str := strings.TrimSpace(m.input.Value())
			if str != "" {
				m.history = append(m.history, str)
				m.historyIdx = -1
				m.addLog("> " + str)
				cmds = append(cmds, m.handleCommand(str))
				m.input.SetValue("")
				return m, tea.Batch(cmds...)
			}
		case "tab":
			m.completeInput()
			return m, nil
		case "up":
			m.vp.ScrollUp(1)
			return m, nil
		case "down":
			m.vp.ScrollDown(1)
			return m, nil
		case "pgup":
			m.vp.PageUp()
			return m, nil
		case "pgdown":
			m.vp.PageDown()
			return m, nil
		case "ctrl+p":
			m.historyPrev()
			return m, nil
		case "ctrl+n":
			m.historyNext()
			return m, nil
		}

	case msgCommandResult:
		if msg.output != "" {
			m.addLog(msg.output)
		}
	case MsgLog:
		m.addLog(string(msg))
	case MsgUpdateImage:
		m.imagePath = string(msg)
		m.imageCache = fileprocess.RenderRegion(m.imagePath)
	case MsgUpdateUserID:
		m.userID = string(msg)
	case MsgUpdateStatus:
		m.addr, m.latency = msg.Addr, msg.Latency
	}

	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	cmds = append(cmds, inputCmd)
	return m, tea.Batch(cmds...)
}

// ----------------------------------------------------------------------------
// View
// ----------------------------------------------------------------------------
func cropLines(s string, h int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > h {
		lines = lines[:h]
	}
	return strings.Join(lines, "\n")
}

func (m *model) View() string {
	if m.quitting {
		return "退出中..."
	}

	logView := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height - 8).
		Render(m.vp.View())

	imgView := lipgloss.NewStyle().
		PaddingLeft(3).
		Width(m.width / 4).
		Render(cropLines(m.imageCache, 4))

	infoView := lipgloss.NewStyle().
		PaddingLeft(10).
		Height(4).
		Render(fmt.Sprintf(
			"%s %s\n%s %s\n%s %s",
			labelStyle.Render("ID:"), valueStyle.Render(m.userID),
			labelStyle.Render("Addr:"), valueStyle.Render(m.addr),
			labelStyle.Render("Ping:"), valueStyle.Render(m.latency),
		))

	statusRow := lipgloss.JoinHorizontal(lipgloss.Top, imgView, infoView)
	console := lipgloss.JoinHorizontal(lipgloss.Top, promptStyle.Render("> "), m.input.View())
	hint := lipgloss.NewStyle().Foreground(lipgloss.Color("#78909C")).PaddingLeft(2).Render("↑/↓ 翻页 | Ctrl+P/N 历史 | Tab 补全")
	inputArea := lipgloss.JoinVertical(lipgloss.Left, console, hint)

	return lipgloss.JoinVertical(lipgloss.Left,
		logView,
		horizontalRule(m.width),
		statusRow,
		horizontalRule(m.width),
		inputArea,
	)
}

// ----------------------------------------------------------------------------
// API
// ----------------------------------------------------------------------------
var p *tea.Program

func Start() error {
	m := NewModel()
	p = tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func Push(msg tea.Msg) {
	if p != nil {
		p.Send(msg)
	}
}
