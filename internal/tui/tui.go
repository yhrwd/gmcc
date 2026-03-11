package tui

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

type Color struct {
	IsRGB   bool
	R, G, B uint8
	Index   uint8
}

var (
	ColorBlack   = Color{Index: 0}
	ColorRed     = Color{Index: 1}
	ColorGreen   = Color{Index: 2}
	ColorYellow  = Color{Index: 3}
	ColorBlue    = Color{Index: 4}
	ColorMagenta = Color{Index: 5}
	ColorCyan    = Color{Index: 6}
	ColorWhite   = Color{Index: 7}
	ColorDefault = Color{Index: 255}
)

func RGBColor(r, g, b uint8) Color { return Color{IsRGB: true, R: r, G: g, B: b} }
func IndexColor(i uint8) Color     { return Color{Index: i} }

func (c Color) Encode() string {
	if c.IsRGB {
		return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.R, c.G, c.B)
	}
	if c.Index < 8 {
		return fmt.Sprintf("\x1b[%dm", 30+c.Index)
	}
	if c.Index == 255 {
		return "\x1b[39m"
	}
	return fmt.Sprintf("\x1b[38;5;%dm", c.Index)
}

func (c Color) EncodeBg() string {
	if c.IsRGB {
		return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", c.R, c.G, c.B)
	}
	if c.Index < 8 {
		return fmt.Sprintf("\x1b[%dm", 40+c.Index)
	}
	if c.Index == 255 {
		return "\x1b[49m"
	}
	return fmt.Sprintf("\x1b[48;5;%dm", c.Index)
}

type Style struct {
	Fg       Color
	Bg       Color
	Bold     bool
	Dim      bool
	Italic   bool
	Undeline bool
	Blink    bool
	Reverse  bool
}

var (
	StyleDefault = Style{Fg: ColorDefault, Bg: ColorDefault}
	StyleBold    = Style{Fg: ColorDefault, Bg: ColorDefault, Bold: true}
	StyleDim     = Style{Fg: ColorDefault, Bg: ColorDefault, Dim: true}
)

func (s Style) Encode() string {
	var codes []string
	codes = append(codes, s.Fg.Encode())
	codes = append(codes, s.Bg.EncodeBg())
	if s.Bold {
		codes = append(codes, "\x1b[1m")
	}
	if s.Dim {
		codes = append(codes, "\x1b[2m")
	}
	if s.Italic {
		codes = append(codes, "\x1b[3m")
	}
	if s.Undeline {
		codes = append(codes, "\x1b[4m")
	}
	if s.Blink {
		codes = append(codes, "\x1b[5m")
	}
	if s.Reverse {
		codes = append(codes, "\x1b[7m")
	}
	return strings.Join(codes, "")
}

func (s Style) Reset() string { return "\x1b[0m" }

type EventType int

const (
	EventKey EventType = iota
	EventResize
	EventMouse
	EventQuit
)

type KeyCode int

const (
	KeyUnknown KeyCode = iota
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyHome
	KeyEnd
	KeyPgUp
	KeyPgDn
	KeyDelete
	KeyInsert
	KeyBackspace
	KeyEnter
	KeyTab
	KeyEscape
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
)

type KeyEvent struct {
	Type  EventType
	Rune  rune
	Key   KeyCode
	Ctrl  bool
	Alt   bool
	Shift bool
}

type EventHandler func(e interface{}) bool

type Screen struct {
	mu       sync.RWMutex
	width    int
	height   int
	buffer   [][]Cell
	oldBuf   [][]Cell
	style    Style
	altBuf   bool
	rawMode  bool
	oldState *term.State
}

type Cell struct {
	Char  rune
	Style Style
}

func NewScreen() (*Screen, error) {
	s := &Screen{}
	if err := s.initSize(); err != nil {
		return nil, err
	}
	s.buffer = s.newBuffer()
	s.oldBuf = s.newBuffer()
	s.style = StyleDefault
	return s, nil
}

func (s *Screen) initSize() error {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w = 80
		h = 24
	}
	s.width = w
	s.height = h
	return nil
}

func (s *Screen) newBuffer() [][]Cell {
	buf := make([][]Cell, s.height)
	for i := range buf {
		buf[i] = make([]Cell, s.width)
		for j := range buf[i] {
			buf[i][j] = Cell{Char: ' ', Style: StyleDefault}
		}
	}
	return buf
}

func (s *Screen) Size() (int, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.width, s.height
}

func (s *Screen) EnableRawMode() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.rawMode {
		return nil
	}
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	s.oldState = oldState
	s.rawMode = true
	return nil
}

func (s *Screen) DisableRawMode() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.rawMode || s.oldState == nil {
		return nil
	}
	err := term.Restore(int(os.Stdin.Fd()), s.oldState)
	s.rawMode = false
	s.oldState = nil
	return err
}

func (s *Screen) EnterAltScreen() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.altBuf {
		fmt.Print("\x1b[?1049h")
		s.altBuf = true
	}
}

func (s *Screen) ExitAltScreen() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.altBuf {
		fmt.Print("\x1b[?1049l")
		s.altBuf = false
	}
}

func (s *Screen) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buffer = s.newBuffer()
}

func (s *Screen) SetCell(x, y int, char rune, style Style) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if x >= 0 && x < s.width && y >= 0 && y < s.height {
		s.buffer[y][x] = Cell{Char: char, Style: style}
	}
}

func (s *Screen) SetText(x, y int, text string, style Style) {
	s.mu.Lock()
	defer s.mu.Unlock()
	runes := []rune(text)
	for i, r := range runes {
		px := x + i
		if px >= 0 && px < s.width && y >= 0 && y < s.height {
			if r == '\t' {
				r = ' '
			}
			s.buffer[y][px] = Cell{Char: r, Style: style}
		}
	}
}

func (s *Screen) HideCursor()        { fmt.Print("\x1b[?25l") }
func (s *Screen) ShowCursor()        { fmt.Print("\x1b[?25h") }
func (s *Screen) SetCursor(x, y int) { fmt.Printf("\x1b[%d;%dH", y+1, x+1) }

func (s *Screen) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	var sb strings.Builder
	sb.WriteString("\x1b[?25l")

	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			cell := s.buffer[y][x]
			if y < len(s.oldBuf) && x < len(s.oldBuf[y]) && s.oldBuf[y][x] == cell {
				continue
			}
			sb.WriteString(fmt.Sprintf("\x1b[%d;%dH", y+1, x+1))
			sb.WriteString(cell.Style.Encode())
			sb.WriteRune(cell.Char)
		}
	}

	sb.WriteString(s.style.Encode())
	sb.WriteString("\x1b[?25h")
	fmt.Print(sb.String())

	s.oldBuf = s.buffer
	s.buffer = s.newBuffer()
}

func (s *Screen) Resize() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.initSize(); err != nil {
		return err
	}
	s.buffer = s.newBuffer()
	s.oldBuf = s.newBuffer()
	return nil
}

type EventLoop struct {
	mu       sync.RWMutex
	running  bool
	quit     chan struct{}
	events   chan interface{}
	handlers map[EventType][]EventHandler
	screen   *Screen
}

func NewEventLoop(screen *Screen) *EventLoop {
	return &EventLoop{
		quit:     make(chan struct{}),
		events:   make(chan interface{}, 256),
		handlers: make(map[EventType][]EventHandler),
		screen:   screen,
	}
}

func (el *EventLoop) On(et EventType, h EventHandler) {
	el.mu.Lock()
	defer el.mu.Unlock()
	el.handlers[et] = append(el.handlers[et], h)
}

func (el *EventLoop) Start() {
	el.mu.Lock()
	el.running = true
	el.mu.Unlock()
	go el.readInput()
	go el.process()
}

func (el *EventLoop) Stop() {
	el.mu.Lock()
	defer el.mu.Unlock()
	if !el.running {
		return
	}
	el.running = false
	close(el.quit)
}

func (el *EventLoop) Emit(e interface{}) {
	select {
	case el.events <- e:
	default:
	}
}

func (el *EventLoop) readInput() {
	buf := make([]byte, 1024)
	for {
		select {
		case <-el.quit:
			return
		default:
		}
		n, err := os.Stdin.Read(buf)
		if err != nil {
			continue
		}
		events := el.parseInput(buf[:n])
		for _, e := range events {
			el.Emit(e)
		}
	}
}

func (el *EventLoop) parseInput(b []byte) []interface{} {
	var events []interface{}
	i := 0
	for i < len(b) {
		if b[i] == 0x1b {
			key, consumed := el.parseEscape(b[i:])
			if key != KeyUnknown {
				events = append(events, KeyEvent{Type: EventKey, Key: key})
			} else if consumed >= 2 && len(b[i:]) > 1 && b[i+1] != '[' {
				events = append(events, KeyEvent{Type: EventKey, Alt: true, Rune: rune(b[i+1])})
			}
			i += consumed
			if consumed == 0 {
				i++
			}
			continue
		}
		if b[i] >= 32 && b[i] < 127 {
			events = append(events, KeyEvent{Type: EventKey, Rune: rune(b[i])})
		} else if b[i] == 13 {
			events = append(events, KeyEvent{Type: EventKey, Key: KeyEnter})
		} else if b[i] == 9 {
			events = append(events, KeyEvent{Type: EventKey, Key: KeyTab})
		} else if b[i] == 127 || b[i] == 8 {
			events = append(events, KeyEvent{Type: EventKey, Key: KeyBackspace})
		} else if b[i] == 3 {
			events = append(events, KeyEvent{Type: EventKey, Ctrl: true, Rune: 'c'})
		}
		i++
	}
	return events
}

func (el *EventLoop) parseEscape(b []byte) (KeyCode, int) {
	if len(b) < 2 || b[1] != '[' {
		return KeyUnknown, 0
	}
	if len(b) < 3 {
		return KeyUnknown, 2
	}
	switch b[2] {
	case 'A':
		return KeyUp, 3
	case 'B':
		return KeyDown, 3
	case 'C':
		return KeyRight, 3
	case 'D':
		return KeyLeft, 3
	case 'H':
		return KeyHome, 3
	case 'F':
		return KeyEnd, 3
	case 'P':
		return KeyF1, 3
	case 'Q':
		return KeyF2, 3
	case 'R':
		return KeyF3, 3
	case 'S':
		return KeyF4, 3
	case '5':
		if len(b) >= 4 && b[3] == '~' {
			return KeyPgUp, 4
		}
	case '6':
		if len(b) >= 4 && b[3] == '~' {
			return KeyPgDn, 4
		}
	case '3':
		if len(b) >= 4 && b[3] == '~' {
			return KeyDelete, 4
		}
	case '2':
		if len(b) >= 4 && b[3] == '~' {
			return KeyInsert, 4
		}
		if len(b) >= 5 && b[4] == '~' {
			switch b[3] {
			case '0':
				return KeyF9, 5
			case '1':
				return KeyF10, 5
			case '3':
				return KeyF11, 5
			case '4':
				return KeyF12, 5
			}
		}
	case '1':
		if len(b) >= 4 {
			if b[3] == '~' {
				return KeyHome, 4
			}
			if len(b) >= 5 && b[4] == '~' {
				switch b[3] {
				case '5':
					return KeyF5, 5
				case '7':
					return KeyF6, 5
				case '8':
					return KeyF7, 5
				case '9':
					return KeyF8, 5
				}
			}
		}
	}
	return KeyUnknown, 2
}

func (el *EventLoop) process() {
	for {
		select {
		case <-el.quit:
			return
		case e := <-el.events:
			el.dispatch(e)
		}
	}
}

func (el *EventLoop) dispatch(e interface{}) {
	var et EventType
	switch e.(type) {
	case KeyEvent:
		et = EventKey
	}
	el.mu.RLock()
	handlers := el.handlers[et]
	el.mu.RUnlock()
	for _, h := range handlers {
		if !h(e) {
			break
		}
	}
}

type Border struct {
	TopLeft, TopRight, BottomLeft, BottomRight string
	Horizontal, Vertical                       string
}

var (
	BorderUnicode = Border{"┌", "┐", "└", "┘", "─", "│"}
	BorderASCII   = Border{"+", "+", "+", "+", "-", "|"}
	BorderDouble  = Border{"╔", "╗", "╚", "╝", "═", "║"}
)

type Component interface {
	Render(w, h int) []string
	HandleEvent(e interface{}) bool
}

type Box struct {
	Title   string
	Content []string
	Border  Border
	focused bool
}

func NewBox() *Box                            { return &Box{Border: BorderUnicode} }
func (b *Box) Focus()                         { b.focused = true }
func (b *Box) Blur()                          { b.focused = false }
func (b *Box) IsFocused() bool                { return b.focused }
func (b *Box) HandleEvent(e interface{}) bool { return false }

func (b *Box) Render(w, h int) []string {
	lines := make([]string, h)
	for i := range lines {
		lines[i] = strings.Repeat(" ", w)
	}
	if h < 2 {
		return lines
	}
	top := b.Border.TopLeft + strings.Repeat(b.Border.Horizontal, w-2) + b.Border.TopRight
	if b.Title != "" && w > 4 {
		titleStr := " " + b.Title + " "
		if len(titleStr) > w-4 {
			titleStr = titleStr[:w-4]
		}
		pos := (w - len(titleStr)) / 2
		if pos > 0 {
			top = b.Border.TopLeft + strings.Repeat(b.Border.Horizontal, pos-1) + titleStr + strings.Repeat(b.Border.Horizontal, w-pos-len(titleStr)-1) + b.Border.TopRight
		}
	}
	lines[0] = top
	for y := 1; y < h-1; y++ {
		if y-1 < len(b.Content) {
			line := b.Content[y-1]
			runes := []rune(line)
			if len(runes) > w-2 {
				runes = runes[:w-2]
			}
			padded := string(runes) + strings.Repeat(" ", w-2-len(runes))
			lines[y] = b.Border.Vertical + padded + b.Border.Vertical
		} else {
			lines[y] = b.Border.Vertical + strings.Repeat(" ", w-2) + b.Border.Vertical
		}
	}
	lines[h-1] = b.Border.BottomLeft + strings.Repeat(b.Border.Horizontal, w-2) + b.Border.BottomRight
	return lines
}

type InputWidget struct {
	Value       string
	CursorPos   int
	Placeholder string
	Password    bool
	focused     bool
	history     []string
	histIndex   int
	OnSubmit    func(string)
}

func NewInputWidget() *InputWidget     { return &InputWidget{} }
func (i *InputWidget) Focus()          { i.focused = true }
func (i *InputWidget) Blur()           { i.focused = false }
func (i *InputWidget) IsFocused() bool { return i.focused }

func (i *InputWidget) Render(w, h int) []string {
	lines := make([]string, h)
	for j := range lines {
		lines[j] = strings.Repeat(" ", w)
	}
	display := i.Value
	if i.Password {
		display = strings.Repeat("*", len(i.Value))
	}
	if len(display) == 0 && len(i.Placeholder) > 0 && !i.focused {
		display = i.Placeholder
	}
	runes := []rune(display)
	if len(runes) > w {
		if i.CursorPos > w {
			runes = runes[i.CursorPos-w:]
		} else {
			runes = runes[:w]
		}
	}
	if h > 0 {
		lines[0] = string(runes) + strings.Repeat(" ", w-len(runes))
	}
	return lines
}

func (i *InputWidget) HandleEvent(e interface{}) bool {
	if !i.focused {
		return false
	}
	ke, ok := e.(KeyEvent)
	if !ok {
		return false
	}
	switch ke.Key {
	case KeyLeft:
		if i.CursorPos > 0 {
			i.CursorPos--
		}
	case KeyRight:
		if i.CursorPos < len(i.Value) {
			i.CursorPos++
		}
	case KeyHome:
		i.CursorPos = 0
	case KeyEnd:
		i.CursorPos = len(i.Value)
	case KeyBackspace:
		if i.CursorPos > 0 {
			i.Value = i.Value[:i.CursorPos-1] + i.Value[i.CursorPos:]
			i.CursorPos--
		}
	case KeyDelete:
		if i.CursorPos < len(i.Value) {
			i.Value = i.Value[:i.CursorPos] + i.Value[i.CursorPos+1:]
		}
	case KeyEnter:
		if i.OnSubmit != nil && len(i.Value) > 0 {
			i.history = append(i.history, i.Value)
			if len(i.history) > 100 {
				i.history = i.history[1:]
			}
			i.histIndex = -1
			i.OnSubmit(i.Value)
		}
	case KeyUp:
		if len(i.history) > 0 {
			if i.histIndex == -1 {
				i.histIndex = len(i.history) - 1
			} else if i.histIndex > 0 {
				i.histIndex--
			}
			if i.histIndex >= 0 && i.histIndex < len(i.history) {
				i.Value = i.history[i.histIndex]
				i.CursorPos = len(i.Value)
			}
		}
	case KeyDown:
		if i.histIndex != -1 && len(i.history) > 0 {
			if i.histIndex < len(i.history)-1 {
				i.histIndex++
				i.Value = i.history[i.histIndex]
			} else {
				i.histIndex = -1
				i.Value = ""
			}
			i.CursorPos = len(i.Value)
		}
	default:
		if ke.Rune != 0 && ke.Rune >= 32 {
			i.Value = i.Value[:i.CursorPos] + string(ke.Rune) + i.Value[i.CursorPos:]
			i.CursorPos++
		}
	}
	return true
}

type Engine struct {
	mu      sync.RWMutex
	screen  *Screen
	events  *EventLoop
	layout  Component
	running bool
	quit    chan struct{}
	fps     int
}

func NewEngine() (*Engine, error) {
	screen, err := NewScreen()
	if err != nil {
		return nil, err
	}
	events := NewEventLoop(screen)
	return &Engine{screen: screen, events: events, quit: make(chan struct{}), fps: 30}, nil
}

func (e *Engine) SetLayout(layout Component) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.layout = layout
}

func (e *Engine) On(et EventType, h EventHandler) { e.events.On(et, h) }

func (e *Engine) Start(ctx context.Context) error {
	e.mu.Lock()
	e.running = true
	e.mu.Unlock()
	e.screen.EnterAltScreen()
	defer e.screen.ExitAltScreen()
	e.screen.EnableRawMode()
	defer e.screen.DisableRawMode()
	e.screen.HideCursor()
	defer e.screen.ShowCursor()
	e.events.Start()
	defer e.events.Stop()
	e.events.On(EventKey, func(ev interface{}) bool {
		ke := ev.(KeyEvent)
		if ke.Ctrl && ke.Rune == 'c' {
			e.Stop()
			return true
		}
		return false
	})
	ticker := time.NewTicker(time.Second / time.Duration(e.fps))
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-e.quit:
			return nil
		case <-ticker.C:
			e.render()
		}
	}
}

func (e *Engine) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.running {
		return
	}
	e.running = false
	close(e.quit)
	e.events.Stop()
}

func (e *Engine) render() {
	e.mu.RLock()
	layout := e.layout
	e.mu.RUnlock()
	if layout == nil {
		return
	}
	w, h := e.screen.Size()
	lines := layout.Render(w, h)
	e.screen.Clear()
	for y, line := range lines {
		if y >= h {
			break
		}
		e.screen.SetText(0, y, line, StyleDefault)
	}
	e.screen.Flush()
}

func (e *Engine) Refresh() { e.render() }

type ComponentFunc func(w, h int) []string

func (f ComponentFunc) Render(w, h int) []string       { return f(w, h) }
func (f ComponentFunc) HandleEvent(e interface{}) bool { return false }
