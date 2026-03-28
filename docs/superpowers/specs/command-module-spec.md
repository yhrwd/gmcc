# GMCC 指令模块规格文档

## 1. 项目概述

### 1.1 目标
为 gmcc 构建一个**模块化、可热更新**的指令系统，支持集群化管理。该系统完全独立于核心客户端，可以作为单独模块使用。

### 1.2 核心指令
- **!ride**：骑乘指令（等待模式，机器人原地等待目标靠近后执行）

### 1.3 设计原则
- **低耦合**：指令模块完全独立，通过接口与Bot通信
- **可扩展**：支持动态加载新指令模块
- **可配置**：支持YAML配置白名单、超时时间等
- **轻量级**：内存占用小，适合10个Bot并发

---

## 2. 系统架构

### 2.1 整体架构

```
┌────────────────────────────────────────────────────────────────┐
│                        Bot (mcclient.Client)                    │
│                                                                   │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │                    CommandRouter (指令路由器)               │  │
│  │                                                              │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐ │  │
│  │  │   Parser    │→ │   Auth      │→ │  State Machine    │ │  │
│  │  │  私聊解析    │  │  UUID鉴权   │  │  状态管理        │  │
│  │  │  指令提取    │  │  白名单检查  │  │  IDLE→PREPARING │  │
│  │  └──────────────┘  └──────────────┘  │       ↓          │  │
│  │                                      │  EXECUTING       │  │
│  │                                      │       ↓          │  │
│  │                                      │  COOLDOWN→IDLE  │  │
│  │                                      └──────────────────┘ │  │
│  │                                                              │  │
│  │  ┌──────────────────────────────────────────────────────┐   │  │
│  │  │                   !ride Module                       │   │  │
│  │  │  1. 视角锁定目标    2. 等待目标靠近                  │   │  │
│  │  │  3. 目标≤4格时执行  4. 发送反馈                      │   │  │
│  │  └──────────────────────────────────────────────────────┘   │  │
│  │                                                              │  │
│  └────────────────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────┘
```

### 2.2 目录结构

```
internal/
├── commands/                    # 指令系统根包
│   ├── commands.go             # 主入口，导出核心组件
│   ├── types.go                # 类型和接口定义
│   ├── router.go               # 指令路由器实现
│   ├── auth.go                 # 鉴权模块
│   ├── state.go                # 状态机
│   ├── config.go               # 配置结构体
│   │
│   ├── modules/                # 指令模块
│   │   ├── ride/
│   │   │   ├── ride.go        # !ride 指令实现
│   │   │   ├── look.go        # 视角锁定逻辑
│   │   │   └── config.go      # 模块配置
│   │   └── module.go           # 模块接口
│   │
│   ├── adapter/
│   │   └── bot_adapter.go     # Bot适配器实现
│   │
│   └── utils/
│       ├── chat.go             # 聊天消息解析工具
│       └── math.go             # 数学计算工具
│
├── mcclient/
│   └── client.go               # 现有客户端（保持不变）
```

---

## 3. 核心接口定义

### 3.1 Bot适配器接口

```go
// BotAdapter 是Bot实例的抽象接口
type BotAdapter interface {
    // 基本信息
    GetPlayerID() string              // 获取玩家名
    GetUUID() string                  // 获取UUID
    GetPosition() (x, y, z float64)   // 获取当前位置
    GetRotation() (yaw, pitch float32) // 获取当前视角
    
    // 聊天操作
    SendChat(msg string) error        // 发送聊天消息
    SendCommand(cmd string) error      // 发送命令（如 /ride）
    SendPrivateMessage(target, msg string) error  // 发送私聊
    
    // 视角操作
    SetYawPitch(yaw, pitch float32)   // 设置视角
    LookAt(x, y, z float64) error     // 看向指定坐标
    
    // 查询
    GetNearbyPlayers() []PlayerInfo   // 获取附近玩家列表
    GetPlayerByName(name string) (PlayerInfo, bool)  // 按名称查找玩家
    DistanceTo(x, y, z float64) float64  // 计算到目标距离
    
    // 状态
    IsOnline() bool                   // 是否在线
}
```

### 3.2 消息结构

```go
// Message 聊天消息结构
type Message struct {
    Type        string    // 消息类型：player_chat, system_chat
    PlainText   string    // 纯文本内容
    RawJSON     string    // 原始JSON
    Sender      string    // 发送者玩家名
    SenderUUID  string    // 发送者UUID
    IsPrivate   bool      // 是否私聊
    Timestamp   time.Time // 接收时间
}

// PlayerInfo 玩家信息
type PlayerInfo struct {
    Name      string
    UUID      string
    Position  struct{ X, Y, Z float64 }
    EntityID  int32
}
```

### 3.3 指令上下文

```go
// Context 指令执行上下文
type Context struct {
    Bot     BotAdapter              // Bot实例
    Message Message                 // 触发指令的消息
    Sender  string                 // 发送者
    Args    []string               // 指令参数
    State   *StateContext          // 状态机上下文
}

// StateContext 状态机上下文（用于维护指令执行过程中的状态）
type StateContext struct {
    CommandName string             // 当前指令名
    Target     string              // 目标玩家
    StartTime  time.Time           // 开始时间
    Metadata   map[string]any      // 指令特定数据
}
```

### 3.4 指令结果

```go
// Result 指令执行结果
type Result struct {
    Success     bool      // 是否成功
    Message     string    // 结果消息（用于私聊反馈）
    NextState   StateType // 下一状态
    Cooldown    time.Duration // 冷却时间
    Error       error     // 错误信息
}

// StateType 状态类型
type StateType int

const (
    StateIdle       StateType = iota  // 空闲，等待指令
    StatePreparing                     // 准备中，鉴权成功
    StateExecuting                     // 执行中，等待目标靠近
    StateCooldown                      // 冷却中
    StateFailed                        // 失败
)
```

### 3.5 指令接口

```go
// Command 指令接口
type Command interface {
    // 元数据
    Name() string                    // 指令名（如 "ride"）
    Description() string             // 描述
    Usage() string                   // 用法示例
    
    // 初始化
    Init(bot BotAdapter, cfg *ModuleConfig) error
    
    // 执行
    Execute(ctx *Context) *Result
    
    // 状态更新（每帧调用）
    Tick(ctx *Context) *Result
    
    // 清理
    Cleanup()
    
    // 状态查询
    State() StateType
    Target() string
}
```

---

## 4. 消息解析器

### 4.1 设计目的

不同服务器的私聊消息格式可能不同，需要将消息解析逻辑独立为可插拔组件，方便手动适配。

### 4.2 Parser接口

```go
// MessageParser 消息解析器接口
// 用户可通过实现此接口适配不同服务器的私信格式
type MessageParser interface {
    // Parse 解析原始聊天消息，返回标准化的Message
    // 返回 nil 表示这不是一条有效的指令消息
    Parse(rawChat RawChat) *Message
}

// RawChat 原始聊天消息（来自mcclient）
type RawChat struct {
    Type        string    // player_chat, system_chat
    PlainText   string    // 纯文本内容
    RawJSON     string    // 原始JSON
    SenderName  string    // 发送者名称（可能为空）
    SenderUUID  [16]byte  // 发送者UUID
    Timestamp   time.Time
}
```

### 4.3 默认实现

```go
// DefaultParser 默认消息解析器
// 适用于大多数标准Minecraft服务器
// 支持 /tell 和 /msg 两种私聊格式
type DefaultParser struct {
    prefix string
}

func NewDefaultParser(prefix string) *DefaultParser {
    return &DefaultParser{prefix: prefix}
}

func (p *DefaultParser) Parse(raw RawChat) *Message {
    // 1. 判断是否为私聊
    isPrivate, sender, content := p.detectPrivateChat(raw)
    if !isPrivate {
        return nil
    }
    
    // 2. 移除私聊前缀后检查指令前缀
    if !strings.HasPrefix(content, p.prefix) {
        return nil
    }
    
    return &Message{
        Type:       raw.Type,
        PlainText:  content,
        RawJSON:    raw.RawJSON,
        Sender:     sender,
        SenderUUID: packet.FormatUUID(raw.SenderUUID),
        IsPrivate:  true,
        Timestamp:  raw.Timestamp,
    }
}

// detectPrivateChat 检测私聊格式
// 返回: (是否私聊, 发送者, 私聊内容)
func (p *DefaultParser) detectPrivateChat(raw RawChat) (bool, string, string) {
    text := raw.PlainText
    
    // 格式1: "玩家私聊消息" (某些服务器)
    // 格式2: "[玩家 -> 你] 消息"
    // 格式3: "玩家 whispers: 消息"
    
    // 尝试匹配常见私聊格式
    patterns := []struct {
        regex   *regexp.Regexp
        extract func([]string) (string, string)
    }{
        // /tell 格式: "[Player -> 你] message"
        {
            regexp.MustCompile(`^\[([^\]]+) -> 你\] (.+)$`),
            func(m []string) (string, string) { return m[1], m[2] },
        },
        // /msg 格式: "玩家 whispers to you: message"
        {
            regexp.MustCompile(`^([^\s]+) whispers (?:to you:|to you from [^:]+:)? (.+)$`),
            func(m []string) (string, string) { return m[1], m[2] },
        },
    }
    
    for _, p := range patterns {
        if m := p.regex.FindStringSubmatch(text); m != nil {
            sender, content := p.extract(m)
            return true, sender, content
        }
    }
    
    // 如果有明确发送者（非系统消息），检查是否包含私聊标记
    if raw.SenderName != "" && raw.Type != "system_chat" {
        // 某些服务器私聊会直接显示发送者名称
        // 用户可根据服务器特性自定义此逻辑
        return false, "", ""
    }
    
    return false, "", ""
}
```

### 4.4 自定义解析器示例

```go
// CustomParser 用户自定义解析器
// 根据服务器实际私聊格式手动修改此文件
// 示例: 某服务器格式 "[发送者 ➥ 接收者] 消息"
type CustomParser struct {
    prefix      string
    botName     string              // Bot玩家名，用于识别接收者
    jsonPattern  *regexp.Regexp    // JSON格式提取正则
}

func NewCustomParser(prefix string, botName string) *CustomParser {
    // 格式: [发送者 ➥ 接收者] 消息内容
    // 例如: [YHRWD ➥ YHRWD] Hello
    jsonPattern := regexp.MustCompile(`^\[([^\s➥]+)\s*➥\s*([^\]]+)\]\s*(.+)$`)
    
    return &CustomParser{
        prefix:     prefix,
        botName:    botName,
        jsonPattern: jsonPattern,
    }
}

func (p *CustomParser) Parse(raw RawChat) *Message {
    // 1. 解析JSON获取纯文本（如果PlainText不准确）
    text := p.extractPlainText(raw)
    
    // 2. 匹配私聊格式: [发送者 ➥ 接收者] 消息
    matches := p.jsonPattern.FindStringSubmatch(text)
    if matches == nil {
        return nil
    }
    
    sender := matches[1]
    receiver := matches[2]
    content := matches[3]
    
    // 3. 检查接收者是否为Bot（确认是发给Bot的私聊）
    if receiver != p.botName && !strings.EqualFold(receiver, p.botName) {
        return nil  // 不是发给Bot的消息
    }
    
    // 4. 检查指令前缀
    if !strings.HasPrefix(content, p.prefix) {
        return nil
    }
    
    // 5. 尝试从JSON中获取发送者UUID
    senderUUID := p.extractSenderUUID(raw)
    
    return &Message{
        Type:       raw.Type,
        PlainText:  content,
        RawJSON:    raw.RawJSON,
        Sender:     sender,
        SenderUUID: senderUUID,
        IsPrivate:  true,
        Timestamp:  raw.Timestamp,
    }
}

// extractPlainText 从JSON extra数组提取纯文本
func (p *CustomParser) extractPlainText(raw RawChat) string {
    if raw.PlainText != "" {
        return raw.PlainText
    }
    
    // 解析JSON提取text字段
    // JSON格式: {"extra":[{"text":"["},{"text":"YHRWD "},...],"text":""}
    type jsonExtra struct {
        Extra []struct {
            Text string `json:"text"`
        } `json:"extra"`
        Text string `json:"text"`
    }
    
    var jsonMsg jsonExtra
    if err := json.Unmarshal([]byte(raw.RawJSON), &jsonMsg); err != nil {
        return ""
    }
    
    var sb strings.Builder
    for _, e := range jsonMsg.Extra {
        sb.WriteString(e.Text)
    }
    sb.WriteString(jsonMsg.Text)
    return sb.String()
}

// extractSenderUUID 尝试从JSON中提取发送者UUID
// 某些服务器会在hover_event或特定字段中包含UUID
func (p *CustomParser) extractSenderUUID(raw RawChat) string {
    // 如果RawChat已包含UUID，直接返回
    if raw.SenderUUID != [16]byte{} {
        return packet.FormatUUID(raw.SenderUUID)
    }
    
    // 尝试从hover_event中提取
    // JSON格式可能包含: {"hover_event":{"action":"show_text","value":"..."}}
    // 需要根据具体服务器格式解析
    
   return ""  // 无法获取UUID，返回空
}
```

### 4.6 配置示例

```yaml
commands:
  prefix: "!"
  
  # Bot玩家名（用于私聊识别）
  bot_name: "Bot001"
  
  # 解析器类型: default, custom
  parser: custom
```

### 4.7 文件位置

```
internal/commands/
├── parser.go              # Parser接口定义
├── parser_default.go      # 默认解析器实现
└── parser_custom.go       # 用户自定义解析器（可手动修改）
```

---

## 5. 指令路由器

### 5.1 Router结构

```go
// Router 指令路由器
type Router struct {
    bot          BotAdapter
    commands     map[string]Command
    parser       MessageParser           // 消息解析器（可替换）
    prefix       string                 // 指令前缀，默认 "!"
    auth         *AuthManager           // 鉴权管理器
    config       *RouterConfig          // 路由器配置
    
    // 状态
    currentCmd   Command                // 当前执行的指令
    cmdMu         sync.RWMutex
}
```

### 5.2 消息处理流程

```go
// HandleRawChat 处理原始聊天消息（来自mcclient）
func (r *Router) HandleRawChat(raw RawChat) {
    // 1. 检查是否在线
    if !r.bot.IsOnline() {
        return
    }
    
    // 2. 使用解析器解析消息
    msg := r.parser.Parse(raw)
    if msg == nil {
        return  // 不是有效的指令消息
    }
    
    // 3. 检查是否有指令前缀
    if !strings.HasPrefix(msg.PlainText, r.prefix) {
        return
    }
    
    // 4. 解析指令
    cmdName, args := r.parseCommand(msg.PlainText)
    
    // 5. 鉴权检查
    if !r.auth.Check(msg) {
        r.bot.SendPrivateMessage(msg.Sender, "你没有权限使用此机器人")
        return
    }
    
    // 6. 获取指令
    cmd, ok := r.commands[cmdName]
    if !ok {
        r.bot.SendPrivateMessage(msg.Sender, fmt.Sprintf("未知指令: %s", cmdName))
        return
    }
    
    // 7. 创建上下文
    ctx := &Context{
        Bot:     r.bot,
        Message: *msg,
        Sender:  msg.Sender,
        Args:    args,
    }
    
    // 8. 执行指令
    result := cmd.Execute(ctx)
    
    // 9. 发送反馈
    if result.Message != "" {
        r.bot.SendPrivateMessage(msg.Sender, result.Message)
    }
}

// SetParser 替换消息解析器（用于自定义服务器）
func (r *Router) SetParser(parser MessageParser) {
    r.parser = parser
}
```

### 5.3 指令解析

```go
func (r *Router) parseCommand(text string) (string, []string) {
    // 去除前缀
    content := strings.TrimPrefix(text, r.prefix)
    
    // 分割参数
    parts := strings.Fields(content)
    if len(parts) == 0 {
        return "", nil
    }
    
    cmd := strings.ToLower(parts[0])
    args := parts[1:]
    
    return cmd, args
}
```

---

## 6. 鉴权系统

### 6.1 AuthManager

```go
// AuthManager 鉴权管理器
type AuthManager struct {
    whitelist   map[string]bool  // UUID -> true
    allowAll    bool             // 是否允许所有人
    mu          sync.RWMutex
}
```

### 6.2 鉴权检查

```go
// Check 检查玩家是否有权限
func (a *AuthManager) Check(msg *Message) bool {
    a.mu.RLock()
    defer a.mu.RUnlock()
    
    // 允许所有人
    if a.allowAll {
        return true
    }
    
    // 检查UUID
    if a.whitelist[msg.SenderUUID] {
        return true
    }
    
    // 检查PlayerID
    if a.whitelist[msg.Sender] {
        return true
    }
    
    return false
}
```

### 6.3 配置示例

```yaml
auth:
  # 允许所有玩家（不推荐生产环境）
  allow_all: false
  
  # 白名单（UUID 或 PlayerID）
  whitelist:
    - "PlayerName1"
    - "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
```

---

## 7. 状态机

### 7.1 状态流转图

```
┌─────────┐  ┌─────────────┐  ┌─────────────┐  ┌────────────┐  ┌─────────┐
│  IDLE   │─→│ PREPARING   │─→│ EXECUTING   │─→│ COOLDOWN   │─→│  IDLE   │
└─────────┘  └─────────────┘  └─────────────┘  └────────────┘  └─────────┘
     ↑              │               │              │                │
     │              │               │              │                │
     │              ▼               ▼              │                │
     │         ┌─────────┐    ┌─────────┐         │                │
     │         │ FAILED  │    │ FAILED  │─────────┴────────────────┘
     │         └─────────┘    └─────────┘              │
     │                                                   │
     └───────────────────────────────────────────────────┘
```

### 7.2 状态说明

| 状态 | 说明 | 进入条件 | 退出条件 |
|------|------|---------|---------|
| IDLE | 空闲 | 默认状态 | 收到有效指令 |
| PREPARING | 准备中 | 鉴权成功 | 找到目标玩家 |
| EXECUTING | 执行中 | 目标在视野内 | 骑乘成功/失败/超时 |
| COOLDOWN | 冷却中 | 指令执行完成 | 冷却时间结束 |
| FAILED | 失败 | 任意阶段出错 | 返回IDLE并发送反馈 |

---

## 8. !ride 指令实现

### 8.1 模块结构

```go
// RideCommand !ride指令实现
type RideCommand struct {
    bot            BotAdapter
    config         *RideConfig
    
    // 状态
    state          StateType
    target         string              // 目标玩家名
    startTime      time.Time           // 开始时间
    lastFeedback   time.Time           // 上次反馈时间
    yawTarget      float32             // 目标Yaw
    pitchTarget    float32             // 目标Pitch
    
    // 配置
    timeout        time.Duration       // 超时时间，默认60秒
    rangeLimit     float64             // 交互范围，默认4格
}
```

### 8.2 配置结构

```go
// RideConfig !ride模块配置
type RideConfig struct {
    // 超时时间（默认60秒）
    Timeout time.Duration
    
    // 交互范围（默认4格）
    RangeLimit float64
    
    // 冷却时间（默认5秒）
    Cooldown time.Duration
    
    // 视角更新间隔
    LookUpdateInterval time.Duration
    
    // 视角平滑系数（0-1，越小越平滑）
    LookSmoothing float32
}
```

### 8.3 Execute 方法

```go
func (r *RideCommand) Execute(ctx *Context) *Result {
    // 1. 检查当前状态
    if r.state != StateIdle {
        return &Result{
            Success: false,
            Message: "指令执行中，请稍候",
            NextState: r.state,
        }
    }
    
    // 2. 确定目标玩家（默认为发送者）
    r.target = ctx.Sender
    if len(ctx.Args) > 0 {
        r.target = ctx.Args[0]
    }
    
    // 3. 查找目标玩家
    player, ok := r.bot.GetPlayerByName(r.target)
    if !ok {
        r.state = StatePreparing
        r.startTime = time.Now()
        
        return &Result{
            Success: true,
            Message: fmt.Sprintf("正在查找玩家 %s，请稍候...", r.target),
            NextState: StatePreparing,
        }
    }
    
    // 4. 检查距离
    dist := r.bot.DistanceTo(player.Position.X, player.Position.Y, player.Position.Z)
    
    if dist <= r.config.RangeLimit {
        // 目标已在范围内，直接执行
        return r.executeRide(ctx, player)
    }
    
    // 5. 进入执行状态
    r.state = StateExecuting
    r.startTime = time.Now()
    r.updateLookAt(player.Position)
    
    return &Result{
        Success: true,
        Message: fmt.Sprintf("目标距离 %.1f 格，请靠近后自动骑乘...", dist),
        NextState: StateExecuting,
    }
}
```

### 8.4 Tick 方法（每帧更新）

```go
func (r *RideCommand) Tick(ctx *Context) *Result {
    switch r.state {
    case StateIdle:
        return nil
    
    case StatePreparing:
        // 查找目标玩家
        player, ok := r.bot.GetPlayerByName(r.target)
        if !ok {
            // 检查超时
            if time.Since(r.startTime) > r.config.Timeout {
                r.state = StateFailed
                return &Result{
                    Success: false,
                    Message: fmt.Sprintf("未找到玩家 %s", r.target),
                    NextState: StateFailed,
                }
            }
            return nil  // 继续等待
        }
        
        // 找到目标，检查距离
        dist := r.bot.DistanceTo(player.Position.X, player.Position.Y, player.Position.Z)
        if dist <= r.config.RangeLimit {
            r.state = StateExecuting
        } else {
            r.state = StateExecuting
            r.startTime = time.Now()
            r.updateLookAt(player.Position)
            return &Result{
                Success: true,
                Message: fmt.Sprintf("已锁定 %s，距离 %.1f 格，请靠近...", r.target, dist),
                NextState: StateExecuting,
            }
        }
        
    case StateExecuting:
        // 检查目标是否存在
        player, ok := r.bot.GetPlayerByName(r.target)
        if !ok {
            r.state = StateFailed
            return &Result{
                Success: false,
                Message: fmt.Sprintf("目标玩家 %s 已离线", r.target),
                NextState: StateFailed,
            }
        }
        
        // 检查超时
        if time.Since(r.startTime) > r.config.Timeout {
            r.state = StateFailed
            return &Result{
                Success: false,
                Message: fmt.Sprintf("等待超时（%v）", r.config.Timeout),
                NextState: StateFailed,
            }
        }
        
        // 更新视角（平滑追踪）
        r.smoothLookAt(player.Position)
        
        // 检查距离
        dist := r.bot.DistanceTo(player.Position.X, player.Position.Y, player.Position.Z)
        if dist <= r.config.RangeLimit {
            // 距离足够，执行骑乘
            return r.executeRide(ctx, player)
        }
        
        // 持续反馈距离（每5秒一次）
        if time.Since(r.lastFeedback) > 5*time.Second {
            r.lastFeedback = time.Now()
            return &Result{
                Success: true,
                Message: fmt.Sprintf("已锁定 %s，距离 %.1f 格，请继续靠近...", r.target, dist),
                NextState: StateExecuting,
            }
        }
        
        return nil
        
    case StateCooldown:
        if time.Since(r.startTime) > r.config.Cooldown {
            r.state = StateIdle
            r.target = ""
        }
        return nil
        
    case StateFailed:
        r.state = StateIdle
        r.target = ""
        return nil
    }
    
    return nil
}
```

### 8.5 视角锁定

```go
// updateLookAt 更新视角看向目标
func (r *RideCommand) updateLookAt(pos Position) {
    botPos := r.bot.GetPosition()
    
    // 计算角度
    dx := pos.X - botPos.X
    dy := pos.Y - botPos.Y  
    dz := pos.Z - botPos.Z
    
    // Yaw（水平角度）
    yaw := float32(math.Atan2(-dx, dz) * 180 / math.Pi)
    
    // Pitch（垂直角度）
    horizDist := math.Sqrt(dx*dx + dz*dz)
    pitch := float32(math.Atan2(-dy, horizDist) * 180 / math.Pi)
    
    r.yawTarget = yaw
    r.pitchTarget = pitch
}

// smoothLookAt 平滑视角追踪
func (r *RideCommand) smoothLookAt(pos Position) {
    // 计算目标角度
    botPos := r.bot.GetPosition()
    dx := pos.X - botPos.X
    dy := pos.Y - botPos.Y
    dz := pos.Z - botPos.Z
    
    yaw := float32(math.Atan2(-dx, dz) * 180 / math.Pi)
    horizDist := math.Sqrt(dx*dx + dz*dz)
    pitch := float32(math.Atan2(-dy, horizDist) * 180 / math.Pi)
    
    // 平滑插值
    smoothing := r.config.LookSmoothing
    if smoothing <= 0 {
        smoothing = 0.1 // 默认平滑系数
    }
    
    // 获取当前视角
    currentYaw, currentPitch := r.bot.GetRotation()
    
    newYaw := currentYaw + (yaw - currentYaw) * smoothing
    newPitch := currentPitch + (pitch - currentPitch) * smoothing
    
    // 设置视角
    r.bot.SetYawPitch(newYaw, newPitch)
}
```

### 8.6 执行骑乘

```go
// executeRide 执行骑乘
func (r *RideCommand) executeRide(ctx *Context, target PlayerInfo) *Result {
    // 发送骑乘命令
    err := r.bot.SendCommand(fmt.Sprintf("ride %s", target.Name))
    if err != nil {
        r.state = StateFailed
        return &Result{
            Success: false,
            Message: fmt.Sprintf("骑乘失败: %v", err),
            NextState: StateFailed,
        }
    }
    
    // 进入冷却状态
    r.state = StateCooldown
    r.startTime = time.Now()
    
    return &Result{
        Success: true,
        Message: fmt.Sprintf("骑乘 %s 成功！", target.Name),
        NextState: StateCooldown,
        Cooldown: r.config.Cooldown,
    }
}
```

---

## 9. Bot适配器实现

### 9.1 结构

```go
// ClientAdapter Bot适配器，实现BotAdapter接口
type ClientAdapter struct {
    client *mcclient.Client  // 底层客户端
    config *config.Config    // 配置
    
    // 视角状态
    currentYaw   float32
    currentPitch float32
}
```

### 9.2 接口实现

```go
func (c *ClientAdapter) GetPlayerID() string {
    if c.client.Online != nil {
        return c.client.Online.ProfileName
    }
    return c.client.OfflineName
}

func (c *ClientAdapter) GetUUID() string {
    if c.client.Online != nil {
        return packet.FormatUUID(c.client.Online.ProfileUUID)
    }
    return ""
}

func (c *ClientAdapter) GetPosition() (x, y, z float64) {
    info := c.client.Player.GetInfo()
    pos := info["position"].([]float64)
    return pos[0], pos[1], pos[2]
}

func (c *ClientAdapter) SendPrivateMessage(target, msg string) error {
    return c.client.SendCommand(fmt.Sprintf("msg %s %s", target, msg))
}

func (c *ClientAdapter) GetNearbyPlayers() []PlayerInfo {
    players := c.client.NearbyPlayers.GetAll()
    result := make([]PlayerInfo, 0, len(players))
    
    for name, info := range players {
        result = append(result, PlayerInfo{
            Name: name,
            UUID: packet.FormatUUID(info.UUID),
            // Position 需要从entity tracker获取
        })
    }
    
    return result
}

func (c *ClientAdapter) DistanceTo(x, y, z float64) float64 {
    px, py, pz := c.GetPosition()
    dx := px - x
    dy := py - y
    dz := pz - z
    return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (c *ClientAdapter) SetYawPitch(yaw, pitch float32) {
    c.currentYaw = yaw
    c.currentPitch = pitch
    // 发送视角数据包
    c.sendPlayerRotation()
}
```

### 9.3 视角数据包

```go
func (c *ClientAdapter) sendPlayerRotation() {
    // Minecraft协议：Serverbound Move Player Rotation 包 (0x1A)
    // 需要在 protocol/v774.go 中添加常量:
    //   PlayServerMovePlayerRotation int32 = 0x1A
    // 字段：Yaw (float32), Pitch (float32), OnGround (boolean)
    
    payload := make([]byte, 0, 9)
    payload = append(payload, packet.EncodeFloat32(c.currentYaw)...)
    payload = append(payload, packet.EncodeFloat32(c.currentPitch)...)
    payload = append(payload, packet.EncodeBool(true)...) // OnGround
    
    c.client.SendPacket(protocol.PlayServerMovePlayerRotation, payload)
}
```

---

## 10. 配置格式

### 10.1 指令系统配置

```yaml
# config.yaml

commands:
  # 指令前缀
  prefix: "!"
  
  # 鉴权配置
  auth:
    allow_all: false
    whitelist:
      - "YourPlayerName"
      - "TrustedFriend"
  
  # 指令模块配置
  modules:
    ride:
      enabled: true
      timeout: 60s        # 等待超时
      range_limit: 4.0    # 交互范围
      cooldown: 5s         # 冷却时间
      look_smoothing: 0.1  # 视角平滑系数
```

### 10.2 默认配置常量

```go
const (
    DefaultPrefix       = "!"
    DefaultTimeout      = 60 * time.Second
    DefaultRangeLimit   = 4.0
    DefaultCooldown     = 5 * time.Second
    DefaultLookSmoothing = 0.1
)
```

---

## 11. 使用示例

### 11.1 初始化

```go
// 创建Bot适配器
bot := adapter.NewClientAdapter(client, cfg)

// 创建路由器
router := commands.NewRouter(bot, &commands.RouterConfig{
    Prefix: cfg.Commands.Prefix,
})

// 设置自定义解析器（根据服务器格式）
customParser := parser.NewCustomParser(cfg.Commands.Prefix, bot.GetPlayerName())
router.SetParser(customParser)

// 配置鉴权
router.SetAuth(&commands.AuthManager{
    AllowAll: cfg.Commands.Auth.AllowAll,
    Whitelist: cfg.Commands.Auth.Whitelist,
})

// 注册指令
rideCmd := modules.NewRideCommand()
router.RegisterCommand(rideCmd)

// 设置聊天处理器
client.SetChatHandler(func(msg mcclient.ChatMessage) {
    // 转换为RawChat
    raw := commands.RawChat{
        Type:       msg.Type,
        PlainText:  msg.PlainText,
        RawJSON:    msg.RawJSON,
        SenderName: msg.SenderName,
        SenderUUID: msg.SenderUUID,
        Timestamp:  msg.ReceivedAt,
    }
    
    // 处理消息
    router.HandleRawChat(raw)
    
    // 更新指令状态
    if router.CurrentCommand() != nil {
        router.CurrentCommand().Tick(&commands.Context{Bot: bot})
    }
})
```

### 11.2 发送指令

```bash
# 玩家发送私聊给Bot
/tell Bot001 !ride

# Bot回复
[Bot] -> 你: 正在锁定你，距离 12.3 格，请靠近...

# 玩家靠近到4格内
[Bot] -> 你: 骑乘成功！
```

---

## 12. 模块化更新

### 12.1 热更新接口

```go
// ModuleLoader 模块加载器
type ModuleLoader interface {
    Load(path string) error      // 加载模块
    Unload(name string) error    // 卸载模块
    Reload(name string) error    // 热重载
    List() []string              // 列出已加载模块
}
```

### 12.2 模块接口

```go
// Module 指令模块接口
type Module interface {
    Name() string                    // 模块名
    Version() string                 // 版本
    Init(registry *Registry) error  // 初始化
    Close() error                   // 关闭
}

// Registry 模块注册表
type Registry struct {
    commands map[string]Command
    mu       sync.RWMutex
}

func (r *Registry) Register(cmd Command) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, exists := r.commands[cmd.Name()]; exists {
        return fmt.Errorf("指令 %s 已存在", cmd.Name())
    }
    
    r.commands[cmd.Name()] = cmd
    return nil
}
```

---

## 13. 错误处理

### 13.1 错误类型

```go
// CommandError 指令错误
type CommandError struct {
    Code    string
    Message string
    Cause   error
}

func (e *CommandError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Cause)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// 预定义错误
var (
    ErrNotWhitelisted = &CommandError{"AUTH_FAILED", "你没有权限使用此机器人", nil}
    ErrTargetNotFound = &CommandError{"TARGET_NOT_FOUND", "目标玩家不在线", nil}
    ErrOutOfRange     = &CommandError{"OUT_OF_RANGE", "目标超出范围", nil}
    ErrTimeout        = &CommandError{"TIMEOUT", "操作超时", nil}
    ErrAlreadyRiding  = &CommandError{"ALREADY_RIDING", "正在骑乘中", nil}
    ErrNotIdle       = &CommandError{"NOT_IDLE", "指令执行中，请稍候", nil}
)
```

### 13.2 错误反馈

```go
func (r *RideCommand) sendError(ctx *Context, err *CommandError) {
    if err == nil {
        return
    }
    
    r.bot.SendPrivateMessage(ctx.Sender, fmt.Sprintf("错误: %s", err.Message))
    r.state = StateFailed
}
```

---

## 14. 日志记录

### 14.1 日志格式

```go
// 格式：[时间] [Bot] [指令] [发送者] [目标] [状态] [消息]
// 示例：[2024-01-15 10:30:05] [Bot001] [ride] [Player1] [Player1] [SUCCESS] [骑乘成功]

type CommandLog struct {
    Time       time.Time
    BotID      string
    Command    string
    Sender     string
    Target     string
    Status     string
    Message    string
}
```

### 14.2 日志输出

```go
func (r *RideCommand) log(result *Result, ctx *Context) {
    logLine := fmt.Sprintf("[%s] [%s] [%s] [%s] [%s] [%s] [%s]",
        time.Now().Format("2006-01-02 15:04:05"),
        r.bot.GetPlayerID(),
        r.Name(),
        ctx.Sender,
        r.target,
        map[bool]string{true: "SUCCESS", false: "FAILED"}[result.Success],
        result.Message,
    )
    
    // 输出到控制台
    logx.Infof(logLine)
    
    // 可选：写入文件
    r.writeToFile(logLine)
}
```

---

## 15. 实现计划

### Phase 1: 核心框架（Day 1-2）
- [ ] 定义接口和类型（types.go）
- [ ] 添加协议常量 `PlayServerMovePlayerRotation` 到 `protocol/v774.go`
- [ ] 实现消息解析器接口（parser.go）
- [ ] 实现默认解析器（parser_default.go）
- [ ] 实现自定义解析器（parser_custom.go）
- [ ] 实现Bot适配器（adapter/bot_adapter.go）
- [ ] 实现路由器（router.go）
- [ ] 实现鉴权系统（auth.go）

### Phase 2: !ride指令（Day 3-4）
- [ ] 实现状态机（state.go）
- [ ] 实现!ride模块（modules/ride/ride.go）
- [ ] 实现视角锁定（modules/ride/look.go）
- [ ] 实现配置加载（config.go）

### Phase 3: 集成测试（Day 5）
- [ ] 与mcclient.Client集成
- [ ] 端到端测试
- [ ] 性能测试（10 Bot并发）

### Phase 4: 文档与优化（Day 6-7）
- [ ] 完善文档
- [ ] 错误处理优化
- [ ] 日志完善

---

## 16. 测试用例

### 16.1 单元测试

```go
func TestAuthManager_Check(t *testing.T) {
    auth := NewAuthManager()
    auth.SetWhitelist([]string{"Player1", "uuid-123"})
    
    tests := []struct {
        name     string
        sender   string
        uuid     string
        expected bool
    }{
        {"白名单玩家名", "Player1", "", true},
        {"白名单UUID", "", "uuid-123", true},
        {"非白名单", "Stranger", "", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            msg := &Message{Sender: tt.sender, SenderUUID: tt.uuid}
            result := auth.Check(msg)
            if result != tt.expected {
                t.Errorf("期望 %v，实际 %v", tt.expected, result)
            }
        })
    }
}
```

### 16.2 集成测试

```go
func TestRideCommand_Execute(t *testing.T) {
    // 1. 创建测试Bot
    bot := NewMockBot()
    bot.AddPlayer("Player1", 10, 70, 10) // 10格外
    
    // 2. 创建指令
    cmd := NewRideCommand(bot, &RideConfig{
        Timeout:    60 * time.Second,
        RangeLimit: 4.0,
    })
    
    // 3. 执行指令
    ctx := &Context{
        Bot:    bot,
        Sender: "Player1",
        Args:   []string{},
    }
    result := cmd.Execute(ctx)
    
    // 4. 验证
    if !result.Success {
        t.Errorf("执行失败: %s", result.Message)
    }
    if result.NextState != StateExecuting {
        t.Errorf("状态错误，期望 %v，实际 %v", StateExecuting, result.NextState)
    }
}
```

---

## 17. 附录

### 17.1 常量定义

```go
// 视角范围
const (
    MaxYaw   float32 = 180.0
    MinYaw   float32 = -180.0
    MaxPitch float32 = 90.0
    MinPitch float32 = -90.0
)

// 默认参数
const (
    DefaultPrefix       = "!"
    DefaultTimeout      = 60 * time.Second
    DefaultRangeLimit   = 4.0
    DefaultCooldown     = 5 * time.Second
    DefaultLookSmoothing = 0.1
    DefaultLookInterval = 100 * time.Millisecond
)
```

### 17.2 依赖

```go
import (
    "math"
    "strings"
    "sync"
    "time"
    "gmcc/internal/logx"
    "gmcc/internal/mcclient"
    "gmcc/internal/config"
)
```

---

---

文档版本：v1.0  
最后更新：2024-01-15  
作者：GMCC Dev Team

---

> **开发指南**：命令开发详细教程请参考 [`docs/command-development.md`](./command-development.md)
