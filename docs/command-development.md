# 命令开发指南

## 快速开始

### 简单命令（推荐入门）

使用 `RegisterSimpleCommand` 一行注册：

```go
import "gmcc/internal/commands"

// 创建路由器
router := commands.NewRouter(bot, "!")

// 查看坐标
router.RegisterSimpleCommand("pos", func(ctx *commands.ChatContext) *commands.CommandResult {
    x, y, z := ctx.Bot.GetPosition()
    return &commands.CommandResult{
        Success: true,
        Message: fmt.Sprintf("坐标: %.1f, %.1f, %.1f", x, y, z),
    }
})

// 查看状态
router.RegisterSimpleCommand("status", func(ctx *commands.ChatContext) *commands.CommandResult {
    name := ctx.Bot.GetPlayerID()
    online := ctx.Bot.IsOnline()
    status := "离线"
    if online {
        status = "在线"
    }
    return &commands.CommandResult{
        Success: true,
        Message: fmt.Sprintf("%s 当前%s", name, status),
    }
})

// 执行命令
router.RegisterSimpleCommand("home", func(ctx *commands.ChatContext) *commands.CommandResult {
    if err := ctx.Bot.SendCommand("home"); err != nil {
        return &commands.CommandResult{
            Success: false,
            Message: fmt.Sprintf("执行失败: %v", err),
        }
    }
    return &commands.CommandResult{
        Success: true,
        Message: "已执行 /home",
    }
})

// 带参数的命令
router.RegisterSimpleCommand("echo", func(ctx *commands.ChatContext) *commands.CommandResult {
    if len(ctx.Args) == 0 {
        return &commands.CommandResult{
            Success: false,
            Message: "用法: !echo <文本>",
        }
    }
    text := strings.Join(ctx.Args, " ")
    return &commands.CommandResult{
        Success: true,
        Message: text,
    }
})

// 查看附近玩家
router.RegisterSimpleCommand("nearby", func(ctx *commands.ChatContext) *commands.CommandResult {
    players := ctx.Bot.GetNearbyPlayers()
    if len(players) == 0 {
        return &commands.CommandResult{
            Success: true,
            Message: "附近没有玩家",
        }
    }
    
    var names []string
    for _, p := range players {
        names = append(names, p.Name)
    }
    return &commands.CommandResult{
        Success: true,
        Message: fmt.Sprintf("附近玩家: %s", strings.Join(names, ", ")),
    }
})

// 计算距离
router.RegisterSimpleCommand("distance", func(ctx *commands.ChatContext) *commands.CommandResult {
    if len(ctx.Args) < 3 {
        return &commands.CommandResult{
            Success: false,
            Message: "用法: !distance <x> <y> <z>",
        }
    }
    
    x, _ := strconv.ParseFloat(ctx.Args[0], 64)
    y, _ := strconv.ParseFloat(ctx.Args[1], 64)
    z, _ := strconv.ParseFloat(ctx.Args[2], 64)
    
    dist := ctx.Bot.DistanceTo(x, y, z)
    return &commands.CommandResult{
        Success: true,
        Message: fmt.Sprintf("距离: %.1f 格", dist),
    }
})
```

---

## 命令类型

| 类型 | 适用场景 | 实现复杂度 | 示例 |
|------|---------|-----------|------|
| **简单命令** | 一次性操作，立即返回结果 | 低 | `!pos`, `!status`, `!calc` |
| **复杂命令** | 需要持续执行、状态管理 | 高 | `!ride`, `!follow`, `!attack` |

---

## 复杂命令开发

### 什么时候需要复杂命令？

- **持续执行**：需要每帧更新（如追踪玩家）
- **状态管理**：多个状态流转
- **等待条件**：等待玩家靠近、等待超时等
- **中断处理**：中途中断、失败恢复

### 状态机设计

```
┌──────────────┐  Execute()   ┌──────────────┐
│   StateIdle  │──────────────│StatePreparing│
└──────────────┘              └──────┬───────┘
       ^                              │ Tick()
       │                              │ 找到玩家
       │                              ▼
       │                         ┌──────────────┐
       │                         │StateExecuting│
       │                         └──────┬───────┘
       │                              │ Tick()
       │                              │ 条件满足
       │                              ▼
       │                         ┌──────────────┐
       │                   Tick()│ StateCooldown│
       │                         └──────┬───────┘
       │                              │ 超时
       └──────────────────────────────┘
```

### 状态说明

| 状态 | 说明 | 进入条件 | 退出条件 |
|------|------|---------|---------|
| StateIdle | 空闲 | 默认状态 | 收到有效指令 |
| StatePreparing | 准备中 | 鉴权成功 | 找到目标玩家 |
| StateExecuting | 执行中 | 目标在视野内 | 完成/失败/超时 |
| StateCooldown | 冷却中 | 指令执行完成 | 冷却时间结束 |
| StateFailed | 失败 | 任意阶段出错 | 返回 IDLE |

### 实现模板

```go
package mycmd

import (
    "fmt"
    "sync"
    "time"
    "gmcc/internal/commands"
)

// MyCommand 复杂命令实现
type MyCommand struct {
    mu     sync.RWMutex
    bot    commands.BotAdapter
    config *Config
    
    // 状态
    state     commands.StateType
    target    string
    startTime time.Time
}

// Config 命令配置
type Config struct {
    Timeout time.Duration
}

// NewMyCommand 创建命令实例
func NewMyCommand(cfg *Config) *MyCommand {
    if cfg == nil {
        cfg = &Config{Timeout: 30 * time.Second}
    }
    return &MyCommand{
        config: cfg,
        state:  commands.StateIdle,
    }
}

// ========================================
// 实现 Command 接口
// ========================================

func (c *MyCommand) Name() string        { return "mycommand" }
func (c *MyCommand) Description() string { return "命令描述" }
func (c *MyCommand) Usage() string       { return "mycommand [参数]" }

func (c *MyCommand) Init(bot commands.BotAdapter, _ *commands.ModuleConfig) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.bot = bot
    return nil
}

// Execute 命令触发时调用
func (c *MyCommand) Execute(ctx *commands.ChatContext) *commands.CommandResult {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    // 1. 检查当前状态
    if c.state != commands.StateIdle {
        return &commands.CommandResult{
            Success: false,
            Message: "命令执行中，请稍候",
        }
    }
    
    // 2. 解析参数
    if len(ctx.Args) == 0 {
        return &commands.CommandResult{
            Success: false,
            Message: fmt.Sprintf("用法: %s", c.Usage()),
        }
    }
    c.target = ctx.Args[0]
    
    // 3. 查找目标
    player, ok := c.bot.GetPlayerByName(c.target)
    if !ok {
        c.state = commands.StatePreparing
        c.startTime = time.Now()
        return &commands.CommandResult{
            Success:   true,
            Message:   fmt.Sprintf("正在查找玩家 %s...", c.target),
            NextState: commands.StatePreparing,
        }
    }
    
    // 4. 开始执行
    c.state = commands.StateExecuting
    c.startTime = time.Now()
    
    return &commands.CommandResult{
        Success:   true,
        Message:   fmt.Sprintf("开始执行，目标: %s", c.target),
        NextState: commands.StateExecuting,
    }
}

// Tick 每帧调用（约50ms/20TPS）
func (c *MyCommand) Tick(ctx *commands.ChatContext) *commands.CommandResult {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    switch c.state {
    case commands.StateIdle:
        return nil
        
    case commands.StatePreparing:
        return c.tickPreparing(ctx)
        
    case commands.StateExecuting:
        return c.tickExecuting(ctx)
        
    case commands.StateCooldown:
        return c.tickCooldown()
        
    case commands.StateFailed:
        c.state = commands.StateIdle
        c.target = ""
        return nil
    }
    
    return nil
}

func (c *MyCommand) tickPreparing(ctx *commands.ChatContext) *commands.CommandResult {
    player, ok := c.bot.GetPlayerByName(c.target)
    if !ok {
        if time.Since(c.startTime) > c.config.Timeout {
            c.state = commands.StateFailed
            return &commands.CommandResult{
                Success:   false,
                Message:   fmt.Sprintf("找不到玩家 %s", c.target),
                NextState: commands.StateFailed,
            }
        }
        return nil // 继续等待
    }
    
    c.state = commands.StateExecuting
    return &commands.CommandResult{
        Success:   true,
        Message:   fmt.Sprintf("找到 %s，开始执行", c.target),
        NextState: commands.StateExecuting,
    }
}

func (c *MyCommand) tickExecuting(ctx *commands.ChatContext) *commands.CommandResult {
    // 检查超时
    if time.Since(c.startTime) > c.config.Timeout {
        c.state = commands.StateFailed
        return &commands.CommandResult{
            Success:   false,
            Message:   "执行超时",
            NextState: commands.StateFailed,
        }
    }
    
    // 检查目标是否在线
    if _, ok := c.bot.GetPlayerByName(c.target); !ok {
        c.state = commands.StateFailed
        return &commands.CommandResult{
            Success:   false,
            Message:   fmt.Sprintf("玩家 %s 已离线", c.target),
            NextState: commands.StateFailed,
        }
    }
    
    // 执行业务逻辑（每帧调用）
    // ...
    
    // 检查是否完成
    if c.isComplete() {
        c.state = commands.StateCooldown
        c.startTime = time.Now()
        return &commands.CommandResult{
            Success:   true,
            Message:   "执行成功",
            NextState: commands.StateCooldown,
        }
    }
    
    return nil // 继续执行
}

func (c *MyCommand) tickCooldown() *commands.CommandResult {
    if time.Since(c.startTime) > 5*time.Second {
        c.state = commands.StateIdle
        c.target = ""
    }
    return nil
}

func (c *MyCommand) isComplete() bool {
    // 判断执行完成的条件
    return true
}

func (c *MyCommand) Cleanup() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.state = commands.StateIdle
    c.target = ""
}

func (c *MyCommand) State() commands.StateType {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.state
}

func (c *MyCommand) Target() string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.target
}
```

---

## 命令注册

### 简单命令

```go
// 单个注册
router.RegisterSimpleCommand("pos", handler)

// 批量注册
simpleCommands := map[string]commands.SimpleHandler{
    "pos":    posHandler,
    "status": statusHandler,
    "home":   homeHandler,
}
for name, handler := range simpleCommands {
    router.RegisterSimpleCommand(name, handler)
}
```

### 复杂命令

```go
// 创建命令实例
rideCmd := ride.NewRideCommand(ride.DefaultConfig())
followCmd := follow.NewFollowCommand(follow.DefaultConfig())

// 注册到路由器
router.RegisterCommand(rideCmd)
router.RegisterCommand(followCmd)
```

---

## 接口参考

### BotAdapter

```go
type BotAdapter interface {
    // 基本信息
    GetPlayerID() string                    // 获取玩家名
    GetUUID() string                         // 获取 UUID
    GetPosition() (x, y, z float64)          // 获取位置
    GetRotation() (yaw, pitch float32)       // 获取朝向
    
    // 消息发送
    SendChat(msg string) error               // 发送聊天消息
    SendCommand(cmd string) error            // 执行命令（如 /home）
    SendPrivateMessage(target, msg string) error // 发送私信
    
    // 视角操作
    SetYawPitch(yaw, pitch float32) error   // 设置朝向
    LookAt(x, y, z float64) error            // 看向指定坐标
    
    // 玩家查询
    GetNearbyPlayers() []PlayerInfo          // 获取附近玩家
    GetPlayerByName(name string) (PlayerInfo, bool) // 按名查找
    
    // 工具方法
    DistanceTo(x, y, z float64) float64      // 计算距离
    IsOnline() bool                            // 在线状态
}

type PlayerInfo struct {
    Name     string
    UUID     string
    EntityID int32
    Position struct{ X, Y, Z float64 }
}
```

### ChatContext

```go
type ChatContext struct {
    Bot     BotAdapter    // Bot 实例
    Message Message       // 原始消息
    Sender  string        // 发送者名字
    Args    []string      // 命令参数（不含命令名）
}

type Message struct {
    Type       string    // 消息类型
    PlainText  string    // 纯文本内容
    RawJSON    string    // 原始 JSON
    Sender     string    // 发送者
    SenderUUID string    // 发送者 UUID
    IsPrivate  bool      // 是否私聊
    Timestamp  time.Time // 接收时间
}
```

### CommandResult

```go
type CommandResult struct {
    Success   bool          // 是否成功
    Message   string        // 反馈消息（私聊发送给用户）
    NextState StateType     // 下一个状态（用于复杂命令）
    Cooldown  time.Duration // 冷却时间
    Error     error         // 错误信息
}
```

### StateType

```go
type StateType int

const (
    StateIdle      StateType = iota // 空闲
    StatePreparing                   // 准备中
    StateExecuting                   // 执行中
    StateCooldown                    // 冷却中
    StateFailed                      // 失败
)
```

---

## 完整示例

### 示例1: 传送命令

```go
router.RegisterSimpleCommand("tp", func(ctx *commands.ChatContext) *commands.CommandResult {
    if len(ctx.Args) < 3 {
        return &commands.CommandResult{
            Success: false,
            Message: "用法: !tp <x> <y> <z>",
        }
    }
    
    x, _ := strconv.ParseFloat(ctx.Args[0], 64)
    y, _ := strconv.ParseFloat(ctx.Args[1], 64)
    z, _ := strconv.ParseFloat(ctx.Args[2], 64)
    
    cmd := fmt.Sprintf("tp %.1f %.1f %.1f", x, y, z)
    if err := ctx.Bot.SendCommand(cmd); err != nil {
        return &commands.CommandResult{
            Success: false,
            Message: fmt.Sprintf("传送失败: %v", err),
        }
    }
    
    return &commands.CommandResult{
        Success: true,
        Message: fmt.Sprintf("已传送到 %.1f, %.1f, %.1f", x, y, z),
    }
})
```

### 示例2: 查看玩家信息

```go
router.RegisterSimpleCommand("info", func(ctx *commands.ChatContext) *commands.CommandResult {
    name := ctx.Bot.GetPlayerID()
    uuid := ctx.Bot.GetUUID()
    x, y, z := ctx.Bot.GetPosition()
    yaw, pitch := ctx.Bot.GetRotation()
    online := ctx.Bot.IsOnline()
    
    status := "离线"
    if online {
        status = "在线"
    }
    
    return &commands.CommandResult{
        Success: true,
        Message: fmt.Sprintf(
            "玩家: %s\nUUID: %s\n状态: %s\n位置: %.1f, %.1f, %.1f\n朝向: Yaw=%.1f Pitch=%.1f",
            name, uuid, status, x, y, z, yaw, pitch,
        ),
    }
})
```

### 示例3: 认领命令

```go
router.RegisterSimpleCommand("claim", func(ctx *commands.ChatContext) *commands.CommandResult {
    if err := ctx.Bot.SendCommand("claim"); err != nil {
        return &commands.CommandResult{
            Success: false,
            Message: "执行失败",
        }
    }
    return &commands.CommandResult{
        Success: true,
        Message: "已执行 claim",
    }
})
```

---

## 目录结构

```
internal/commands/
├── types.go              # 接口定义
├── router.go             # 路由器（含 SimpleHandler）
├── auth.go               # 鉴权管理
├── config.go             # 配置
├── parser.go             # 消息解析接口
├── parser_default.go     # 默认解析器
├── adapter/
│   └── bot_adapter.go    # Bot 适配器实现
└── modules/
    ├── module.go         # 模块入口
    └── ride/             # ride 命令示例
        ├── config.go
        ├── look.go
        ├── ride.go
        └── ride_test.go
```

---

## 测试

### 测试简单命令

```go
func TestSimpleCommand(t *testing.T) {
    bot := &mockBot{online: true}
    router := commands.NewRouter(bot, "!")
    
    router.RegisterSimpleCommand("pos", func(ctx *commands.ChatContext) *commands.CommandResult {
        x, y, z := ctx.Bot.GetPosition()
        return &commands.CommandResult{
            Success: true,
            Message: fmt.Sprintf("%.1f, %.1f, %.1f", x, y, z),
        }
    })
    
    cmd, ok := router.GetCommand("pos")
    if !ok {
        t.Fatal("command not registered")
    }
    
    ctx := &commands.ChatContext{Bot: bot}
    result := cmd.Execute(ctx)
    
    if !result.Success {
        t.Errorf("expected success, got: %s", result.Message)
    }
}
```

### Mock Bot

```go
type mockBot struct {
    online   bool
    position [3]float64
    commands []string
}

func (m *mockBot) GetPlayerID() string                          { return "MockBot" }
func (m *mockBot) GetUUID() string                               { return "mock-uuid" }
func (m *mockBot) GetPosition() (x, y, z float64)                { return m.position[0], m.position[1], m.position[2] }
func (m *mockBot) GetRotation() (yaw, pitch float32)             { return 0, 0 }
func (m *mockBot) SendChat(msg string) error                     { return nil }
func (m *mockBot) SendCommand(cmd string) error                  { m.commands = append(m.commands, cmd); return nil }
func (m *mockBot) SendPrivateMessage(target, msg string) error   { return nil }
func (m *mockBot) SetYawPitch(yaw, pitch float32) error         { return nil }
func (m *mockBot) LookAt(x, y, z float64) error                  { return nil }
func (m *mockBot) IsOnline() bool                                { return m.online }
func (m *mockBot) GetNearbyPlayers() []commands.PlayerInfo       { return nil }
func (m *mockBot) GetPlayerByName(name string) (commands.PlayerInfo, bool) {
    return commands.PlayerInfo{}, false
}
func (m *mockBot) DistanceTo(x, y, z float64) float64 {
    dx := m.position[0] - x
    dy := m.position[1] - y
    dz := m.position[2] - z
    return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
```