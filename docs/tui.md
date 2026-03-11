# TUI 框架文档

## 概述

gmcc 实现了一个轻量级的 Terminal User Interface (TUI) 框架，不依赖任何外部 TUI 库，完全基于 Go 标准库和 `golang.org/x/term` 构建。

## 架构设计

```
┌────────────────────────────────────────────────────────────┐
│                         TUI 架构                           │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    │
│  │   Screen    │───>│   Layout    │<───│   Theme     │    │
│  │  (屏幕管理) │    │  (布局管理) │    │  (主题配置) │    │
│  └─────────────┘    └─────────────┘    └─────────────┘    │
│         │                 │                   │            │
│         ▼                 ▼                   │            │
│  ┌─────────────┐    ┌─────────────┐          │            │
│  │   Renderer  │    │  Component  │<─────────┘            │
│  │  (渲染器)   │    │  (组件系统) │                       │
│  └─────────────┘    └─────────────┘                       │
│         │                 │                                │
│         │    ┌────────────┴────────────┐                  │
│         │    ▼            ▼            ▼                  │
│         │ ┌──────┐   ┌──────┐   ┌──────┐                 │
│         │ │  Box  │   │  Text │   │ List │  ...          │
│         │ └──────┘   └──────┘   └──────┘                 │
│         │                                                 │
│         ▼                                                 │
│  ┌─────────────┐    ┌─────────────┐                      │
│  │ Event Loop  │<───│   Input     │                      │
│  │ (事件循环)  │    │  (输入处理) │                      │
│  └─────────────┘    └─────────────┘                      │
│         │                                                 │
│         ▼                                                 │
│  ┌─────────────┐                                         │
│  │   Engine    │                                         │
│  │  (主引擎)   │                                         │
│  └─────────────┘                                         │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

## 核心组件

### 1. Screen（屏幕管理）

```go
type Screen struct {
    width      int
    height     int
    buffer     [][]Cell         // 双缓冲
    oldBuffer  [][]Cell
    cursorX    int
    cursorY    int
    cursorVis  bool
    rawMode    bool
    oldState   *term.State
}

type Cell struct {
    Char  rune
    Style Style
}
```

**功能**：
- 管理终端尺寸
- 双缓冲渲染
- 光标控制
- 原始模式切换

### 2. Layout（布局管理）

```go
type Layout struct {
    root       Component
    components map[string]Component
    focusOrder []string
    focusIndex int
}
```

**布局类型**：
- **BorderLayout**: 上下左右中五区域
- **FlowLayout**: 水平/垂直流式布局
- **GridLayout**: 网格布局
- **StackLayout**: 堆叠布局

### 3. Component（组件系统）

```go
type Component interface {
    Render(w, h int) []string
    HandleEvent(e Event) bool
    SetStyle(s Style)
    GetStyle() Style
    SetSize(w, h int)
    GetSize() (int, int)
}

type Focusable interface {
    Component
    Focus()
    Blur()
    IsFocused() bool
}
```

**内置组件**：
- `Box`: 基础容器
- `Text`: 文本显示
- `Input`: 输入框
- `List`: 列表选择
- `Table`: 表格显示
- `ProgressBar`: 进度条
- `StatusBar`: 状态栏

### 4. Renderer（渲染器）

```go
type Renderer struct {
    screen    *Screen
    styleMap  map[string]string
}

func (r *Renderer) Render(layout *Layout) {
    // 1. 清屏
    r.screen.Clear()
    
    // 2. 渲染所有组件
    lines := layout.Render(r.screen.width, r.screen.height)
    
    // 3. 双缓冲比较，只更新变化的部分
    r.screen.Update(lines)
    
    // 4. 刷新显示
    r.screen.Flush()
}
```

### 5. Event Loop（事件循环）

```go
type EventLoop struct {
    events    chan Event
    handlers  map[EventType][]EventHandler
    quit      chan struct{}
}

type Event struct {
    Type EventType
    Data interface{}
}

type EventType int

const (
    EventKey EventType = iota
    EventResize
    EventMouse
    EventCustom
)

type EventHandler func(Event) bool
```

### 6. Theme（主题）

```go
type Theme struct {
    Name       string
    Colors     ColorScheme
    Styles     StyleScheme
    BorderType BorderType
}

type ColorScheme struct {
    Background   Color
    Foreground   Color
    Primary      Color
    Secondary    Color
    Accent       Color
    Error        Color
    Warning      Color
    Success      Color
    Border       Color
    BorderFocus  Color
}

type StyleScheme struct {
    Normal    Style
    Focused   Style
    Disabled  Style
    Selected  Style
    Header    Style
    Footer    Style
}
```

## ANSI 转义序列

### 光标控制

```go
const (
    CSI = "\x1b["  // Control Sequence Introducer
    
    // 光标移动
    CursorUp    = CSI + "%dA"
    CursorDown  = CSI + "%dB"
    CursorRight = CSI + "%dC"
    CursorLeft  = CSI + "%dD"
    CursorHome  = CSI + "H"
    CursorPos   = CSI + "%d;%dH"
    
    // 光标可见性
    CursorShow  = CSI + "?25h"
    CursorHide  = CSI + "?25l"
    
    // 清屏
    ClearScreen = CSI + "2J"
    ClearLine   = CSI + "2K"
    ClearToEnd  = CSI + "K"
)
```

### 颜色和样式

```go
const (
    // 前景色 (16色)
    FGBlack   = CSI + "30m"
    FGRed     = CSI + "31m"
    FGGreen   = CSI + "32m"
    FGYellow  = CSI + "33m"
    FGBlue    = CSI + "34m"
    FGMagenta = CSI + "35m"
    FGCyan    = CSI + "36m"
    FGWhite   = CSI + "37m"
    
    // 背景色
    BGBlack   = CSI + "40m"
    // ... 类似前景色
    
    // 256色
    FG256  = CSI + "38;5;%dm"
    BG256  = CSI + "48;5;%dm"
    
    // RGB真彩
    FGRGB  = CSI + "38;2;%d;%d;%dm"
    BGRGB  = CSI + "48;2;%d;%d;%dm"
    
    // 样式
    Bold      = CSI + "1m"
    Dim       = CSI + "2m"
    Italic    = CSI + "3m"
    Underline = CSI + "4m"
    Blink     = CSI + "5m"
    Reverse   = CSI + "7m"
    Reset     = CSI + "0m"
)
```

### 替代屏幕缓冲区

```go
const (
    AltScreenEnter = CSI + "?1049h"
    AltScreenLeave = CSI + "?1049l"
)
```

## 事件处理

### 键盘事件

```go
type KeyEvent struct {
    Rune     rune
    Key      KeyCode
    Mod      Modifier
    Alt      bool
    Ctrl     bool
    Shift    bool
}

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
    // ...
)
```

### 解析输入

```go
func (e *EventLoop) readInput() {
    buf := make([]byte, 1024)
    for {
        n, err := os.Stdin.Read(buf)
        if err != nil {
            continue
        }
        
        events := e.parseEscapeSequence(buf[:n])
        for _, ev := range events {
            e.events <- ev
        }
    }
}

func (e *EventLoop) parseEscapeSequence(b []byte) []Event {
    if len(b) == 1 {
        if b[0] >= 32 && b[0] < 127 {
            return []Event{{Type: EventKey, Data: KeyEvent{Rune: rune(b[0])}}}
        }
        // 特殊键处理...
    }
    
    // ESC 序列处理
    if b[0] == 0x1b {
        if len(b) == 1 {
            return []Event{{Type: EventKey, Data: KeyEvent{Key: KeyEscape}}}
        }
        // [A = Up, [B = Down, [C = Right, [D = Left
        // ...
    }
    
    return nil
}
```

## 组件实现示例

### Box 组件

```go
type Box struct {
    x, y, w, h int
    content    string
    style      Style
    border     Border
    title      string
}

func (b *Box) Render(w, h int) []string {
    lines := make([]string, h)
    
    // 渲染边框
    lines[0] = b.renderTopBorder(b.w)
    for i := 1; i < h-1; i++ {
        lines[i] = b.renderLine()
    }
    lines[h-1] = b.renderBottomBorder(b.w)
    
    return lines
}

func (b *Box) HandleEvent(e Event) bool {
    // 处理事件...
    return false
}
```

### Input 组件

```go
type Input struct {
    value      string
    cursorPos  int
    style      Style
    placeholder string
    password   bool
    focused    bool
    history    []string
    histIndex  int
}

func (i *Input) HandleEvent(e Event) bool {
    if !i.focused {
        return false
    }
    
    ke, ok := e.Data.(KeyEvent)
    if !ok {
        return false
    }
    
    switch ke.Key {
    case KeyLeft:
        if i.cursorPos > 0 {
            i.cursorPos--
        }
    case KeyRight:
        if i.cursorPos < len(i.value) {
            i.cursorPos++
        }
    case KeyBackspace:
        if i.cursorPos > 0 {
            i.value = i.value[:i.cursorPos-1] + i.value[i.cursorPos:]
            i.cursorPos--
        }
    case KeyEnter:
        // 提交输入
    default:
        if ke.Rune != 0 {
            i.value = i.value[:i.cursorPos] + string(ke.Rune) + i.value[i.cursorPos:]
            i.cursorPos++
        }
    }
    
    return true
}
```

### List 组件

```go
type List struct {
    items      []ListItem
    selected   int
    offset     int
    style      Style
    focused    bool
}

type ListItem struct {
    Text    string
    Value   interface{}
    Enabled bool
}

func (l *List) HandleEvent(e Event) bool {
    if !l.focused {
        return false
    }
    
    ke, ok := e.Data.(KeyEvent)
    if !ok {
        return false
    }
    
    switch ke.Key {
    case KeyUp:
        if l.selected > 0 {
            l.selected--
            l.ensureVisible()
        }
    case KeyDown:
        if l.selected < len(l.items)-1 {
            l.selected++
            l.ensureVisible()
        }
    case KeyEnter:
        // 选择当前项
    }
    
    return true
}

func (l *List) ensureVisible() {
    if l.selected < l.offset {
        l.offset = l.selected
    }
    if l.selected >= l.offset+l.height-1 {
        l.offset = l.selected - l.height + 2
    }
}
```

## 布局系统

### BorderLayout

```
┌──────────────────────────────────────┐
│              Top (Header)             │
├──────────┬─────────────────┬─────────┤
│          │                 │         │
│  Left    │     Center      │  Right  │
│          │                 │         │
├──────────┴─────────────────┴─────────┤
│           Bottom (Footer)             │
└──────────────────────────────────────┘
```

```go
type BorderLayout struct {
    top    Component
    bottom Component
    left   Component
    right  Component
    center Component
}

func (l *BorderLayout) Render(w, h int) []string {
    lines := make([]string, h)
    
    topH, bottomH, leftW, rightW := l.calculateSizes(w, h)
    
    // 渲染顶部
    if l.top != nil {
        topLines := l.top.Render(w, topH)
        copy(lines, topLines)
    }
    
    // 渲染底部
    if l.bottom != nil {
        bottomLines := l.bottom.Render(w, bottomH)
        copy(lines[h-bottomH:], bottomLines)
    }
    
    // 渲染左右和中心
    // ...
    
    return lines
}
```

## 应用集成

### 主引擎

```go
type Engine struct {
    screen    *Screen
    layout    *Layout
    events    *EventLoop
    renderer  *Renderer
    running   bool
}

func NewEngine() (*Engine, error) {
    screen, err := NewScreen()
    if err != nil {
        return nil, err
    }
    
    return &Engine{
        screen:   screen,
        layout:   NewLayout(),
        events:   NewEventLoop(),
        renderer: NewRenderer(screen),
    }, nil
}

func (e *Engine) Run(ctx context.Context) error {
    e.screen.EnterAltScreen()
    defer e.screen.ExitAltScreen()
    
    e.screen.EnableRawMode()
    defer e.screen.DisableRawMode()
    
    e.running = true
    
    // 启动事件循环
    go e.events.Start()
    defer e.events.Stop()
    
    // 主循环
    for e.running {
        select {
        case <-ctx.Done():
            e.running = false
        case ev := <-e.events.C:
            if !e.layout.HandleEvent(ev) {
                e.handleGlobalEvent(ev)
            }
        default:
            e.render()
            time.Sleep(16 * time.Millisecond)  // ~60 FPS
        }
    }
    
    return nil
}

func (e *Engine) render() {
    e.renderer.Render(e.layout)
}
```

## 最佳实践

### 1. 性能优化

- 使用双缓冲减少闪烁
- 只更新变化的部分
- 避免频繁清屏
- 使用合适的帧率（30-60 FPS）

### 2. 兼容性

- 检测终端能力
- 提供降级方案（如无真彩色时使用256色）
- 处理终端尺寸变化

### 3. 用户体验

- 提供即时反馈
- 支持常见的快捷键
- 保持界面简洁清晰
- 合理使用颜色和样式

## 命令行应用

gmcc 提供两种界面模式：

### 简单模式（默认）

基础命令行交互，适合脚本使用。

### TUI 模式

完整终端界面，提供：
- 聊天消息窗口
- 命令输入框
- 状态栏
- 玩家信息面板

配置启用：

```yaml
tui:
  enabled: true
  theme: "dark"
```