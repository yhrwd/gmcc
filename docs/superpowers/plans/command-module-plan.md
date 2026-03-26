# GMCC 指令模块实现计划

## 项目概述

为 gmcc 构建模块化、可热更新的指令系统，首期仅实现 `!ride` 指令。

---

只要任务没完成，或超时之前，就保持待机状态，

## 实现阶段

### Phase 1: 核心框架

#### 1.1 类型定义 (`internal/commands/types.go`)

**目标**：定义核心接口和类型

**任务**：
- [ ] `BotAdapter` 接口定义
- [ ] `Message` 结构体
- [ ] `Context` 结构体
- [ ] `Result` 结构体
- [ ] `StateType` 枚举
- [ ] `Command` 接口

**产出**：
```go
// internal/commands/types.go
package commands

type BotAdapter interface {
    GetPlayerID() string
    GetUUID() string
    GetPosition() (x, y, z float64)
    SendChat(msg string) error
    SendCommand(cmd string) error
    SendPrivateMessage(target, msg string) error
    SetYawPitch(yaw, pitch float32) error
    LookAt(x, y, z float64) error
    GetNearbyPlayers() []PlayerInfo
    GetPlayerByName(name string) (PlayerInfo, bool)
    DistanceTo(x, y, z float64) float64
    IsOnline() bool
}

type Message struct {
    Type        string
    PlainText   string
    RawJSON     string
    Sender      string
    SenderUUID  string
    IsPrivate   bool
    Timestamp   time.Time
}

type PlayerInfo struct {
    Name     string
    UUID     string
    Position struct{ X, Y, Z float64 }
    EntityID int32
}

type Context struct {
    Bot     BotAdapter
    Message Message
    Sender  string
    Args    []string
}

type Result struct {
    Success   bool
    Message   string
    NextState StateType
    Cooldown  time.Duration
    Error     error
}

type StateType int
const (
    StateIdle StateType = iota
    StatePreparing
    StateExecuting
    StateCooldown
    StateFailed
)

type Command interface {
    Name() string
    Description() string
    Usage() string
    Init(bot BotAdapter, cfg *ModuleConfig) error
    Execute(ctx *Context) *Result
    Tick(ctx *Context) *Result
    Cleanup()
    State() StateType
    Target() string
}
```

#### 1.2 Bot适配器 (`internal/commands/adapter/bot_adapter.go`)

**目标**：实现BotAdapter接口，封装mcclient.Client

**任务**：
- [ ] `ClientAdapter` 结构体
- [ ] 基本信息获取方法
- [ ] 聊天发送方法
- [ ] 视角操作方法（发送Player Rotation包）
- [ ] 玩家查询方法

**关键实现**：
```go
func (c *ClientAdapter) SetYawPitch(yaw, pitch float32) error {
    // Minecraft协议：Player Rotation 包 (0x1B / 0x16)
    // 格式：Yaw (float32), Pitch (float32), OnGround (bool)
    
    payload := make([]byte, 0, 9)
    buf := bytes.NewBuffer(payload)
    
    binary.Write(buf, binary.BigEndian, yaw)
    binary.Write(buf, binary.BigEndian, pitch)
    binary.Write(buf, binary.BigEndian, true)
    
    return c.client.SendPacket(protocol.PlayClientPlayerRotation, buf.Bytes())
}

func (c *ClientAdapter) GetPlayerByName(name string) (PlayerInfo, bool) {
    // 从NearbyPlayers获取玩家信息
    // 需要结合Entity Tracker获取位置
    players := c.client.NearbyPlayers.GetAll()
    for _, p := range players {
        if strings.EqualFold(p.Name, name) {
            return PlayerInfo{
                Name:     p.Name,
                UUID:     packet.FormatUUID(p.UUID),
                EntityID: p.EntityID,
                // Position 需要从entity tracker获取
            }, true
        }
    }
    return PlayerInfo{}, false
}
```

#### 1.3 指令路由器 (`internal/commands/router.go`)

**目标**：实现消息接收和指令分发

**任务**：
- [ ] `Router` 结构体
- [ ] 消息处理流程
- [ ] 指令解析逻辑
- [ ] 注册/注销指令方法
- [ ] 状态查询方法

**关键实现**：
```go
func (r *Router) HandleMessage(msg Message) {
    // 1. 检查在线状态
    if !r.bot.IsOnline() {
        return
    }
    
    // 2. 检查私聊
    if !msg.IsPrivate {
        return
    }
    
    // 3. 检查前缀
    if !strings.HasPrefix(msg.PlainText, r.prefix) {
        return
    }
    
    // 4. 解析指令
    cmdName, args := r.parseCommand(msg.PlainText)
    
    // 5. 鉴权
    if !r.auth.Check(&msg) {
        r.bot.SendPrivateMessage(msg.Sender, "你没有权限使用此机器人")
        return
    }
    
    // 6. 获取指令
    cmd, ok := r.commands[cmdName]
    if !ok {
        r.bot.SendPrivateMessage(msg.Sender, fmt.Sprintf("未知指令: %s", cmdName))
        return
    }
    
    // 7. 执行
    ctx := &Context{
        Bot:    r.bot,
        Message: msg,
        Sender: msg.Sender,
        Args:   args,
    }
    
    result := cmd.Execute(ctx)
    
    // 8. 反馈
    if result.Message != "" {
        r.bot.SendPrivateMessage(msg.Sender, result.Message)
    }
}
```

#### 1.4 鉴权系统 (`internal/commands/auth.go`)

**目标**：实现白名单鉴权

**任务**：
- [ ] `AuthManager` 结构体
- [ ] 白名单检查方法
- [ ] 通配符支持
- [ ] 动态更新白名单

---

### Phase 2: !ride指令（Day 3-4）

#### 2.1 状态机 (`internal/commands/state.go`)

**目标**：封装状态流转逻辑

**任务**：
- [ ] `StateContext` 结构体
- [ ] 状态常量定义
- [ ] 状态转换方法
- [ ] 超时检测

#### 2.2 !ride指令实现 (`internal/commands/modules/ride/ride.go`)

**目标**：实现骑乘指令核心逻辑

**任务**：
- [ ] `RideCommand` 结构体
- [ ] `Init` 方法
- [ ] `Execute` 方法（IDLE → PREPARING/EXECUTING）
- [ ] `Tick` 方法（状态更新）
- [ ] `Cleanup` 方法
- [ ] 视角锁定逻辑
- [ ] 距离检测
- [ ] 执行骑乘命令

**关键代码**：
```go
func (r *RideCommand) Execute(ctx *Context) *Result {
    // 检查状态
    if r.state != StateIdle {
        return &Result{
            Success: false,
            Message: "指令执行中，请稍候",
            NextState: r.state,
        }
    }
    
    // 确定目标
    r.target = ctx.Sender
    if len(ctx.Args) > 0 {
        r.target = ctx.Args[0]
    }
    
    // 查找目标
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
    
    // 检查距离
    dist := r.bot.DistanceTo(player.Position.X, player.Position.Y, player.Position.Z)
    
    if dist <= r.config.RangeLimit {
        return r.executeRide(ctx, player)
    }
    
    // 进入执行状态
    r.state = StateExecuting
    r.startTime = time.Now()
    r.updateLookAt(player.Position)
    
    return &Result{
        Success: true,
        Message: fmt.Sprintf("目标距离 %.1f 格，请靠近后自动骑乘...", dist),
        NextState: StateExecuting,
    }
}

func (r *RideCommand) Tick(ctx *Context) *Result {
    switch r.state {
    case StatePreparing:
        player, ok := r.bot.GetPlayerByName(r.target)
        if !ok {
            if time.Since(r.startTime) > r.config.Timeout {
                r.state = StateFailed
                return &Result{
                    Success: false,
                    Message: fmt.Sprintf("未找到玩家 %s", r.target),
                    NextState: StateFailed,
                }
            }
            return nil
        }
        
        dist := r.bot.DistanceTo(player.Position.X, player.Position.Y, player.Position.Z)
        r.state = StateExecuting
        r.startTime = time.Now()
        r.updateLookAt(player.Position)
        
        return &Result{
            Success: true,
            Message: fmt.Sprintf("已锁定 %s，距离 %.1f 格，请靠近...", r.target, dist),
            NextState: StateExecuting,
        }
        
    case StateExecuting:
        player, ok := r.bot.GetPlayerByName(r.target)
        if !ok {
            r.state = StateFailed
            return &Result{
                Success: false,
                Message: fmt.Sprintf("目标玩家 %s 已离线", r.target),
                NextState: StateFailed,
            }
        }
        
        // 超时检测
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
            return r.executeRide(ctx, player)
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

#### 2.3 视角锁定 (`internal/commands/modules/ride/look.go`)

**目标**：实现视角追踪逻辑

**任务**：
- [ ] 角度计算（Yaw, Pitch）
- [ ] 平滑插值
- [ ] 发送视角数据包

**关键实现**：
```go
func (r *RideCommand) updateLookAt(pos Position) {
    botPos := r.bot.GetPosition()
    
    dx := pos.X - botPos.X
    dy := pos.Y - botPos.Y
    dz := pos.Z - botPos.Z
    
    // Yaw: 水平角度 [-180, 180]
    yaw := float32(math.Atan2(-dx, dz) * 180 / math.Pi)
    
    // Pitch: 垂直角度 [-90, 90]
    horizDist := math.Sqrt(dx*dx + dz*dz)
    pitch := float32(math.Atan2(-dy, horizDist) * 180 / math.Pi)
    
    r.yawTarget = yaw
    r.pitchTarget = pitch
    
    r.bot.SetYawPitch(yaw, pitch)
}

func (r *RideCommand) smoothLookAt(pos Position) {
    botPos := r.bot.GetPosition()
    dx := pos.X - botPos.X
    dy := pos.Y - botPos.Y
    dz := pos.Z - botPos.Z
    
    targetYaw := float32(math.Atan2(-dx, dz) * 180 / math.Pi)
    horizDist := math.Sqrt(dx*dx + dz*dz)
    targetPitch := float32(math.Atan2(-dy, horizDist) * 180 / math.Pi)
    
    // 平滑插值
    smoothing := r.config.LookSmoothing
    if smoothing <= 0 {
        smoothing = 0.1
    }
    
    newYaw := r.yawTarget + (targetYaw - r.yawTarget) * smoothing
    newPitch := r.pitchTarget + (targetPitch - r.pitchTarget) * smoothing
    
    r.yawTarget = newYaw
    r.pitchTarget = newPitch
    
    r.bot.SetYawPitch(newYaw, newPitch)
}
```

#### 2.4 配置管理 (`internal/commands/config.go`)

**目标**：配置结构定义和加载

**任务**：
- [ ] `RouterConfig` 结构体
- [ ] `ModuleConfig` 结构体
- [ ] `RideConfig` 结构体
- [ ] 默认值定义
- [ ] YAML配置解析

---

### Phase 3: 集成测试（Day 5）

#### 3.1 与mcclient.Client集成

**任务**：
- [ ] 修改`client.go`，添加指令模块支持
- [ ] 暴露`SetCommandRouter`方法
- [ ] ChatHandler集成
- [ ] Tick循环（在Play状态时调用）

**集成代码**：
```go
// internal/mcclient/client.go

func (c *Client) SetupCommands(cfg *commands.RouterConfig) error {
    adapter := adapter.NewClientAdapter(c)
    router := commands.NewRouter(adapter, cfg)
    
    // 注册指令
    rideCmd := modules.NewRideCommand()
    router.RegisterCommand(rideCmd)
    
    c.commandRouter = router
    return nil
}

func (c *Client) handleChatMessage(msg ChatMessage) {
    // 转换为commands.Message
    cmdMsg := commands.Message{
        Type:       msg.Type,
        PlainText:  msg.PlainText,
        RawJSON:    msg.RawJSON,
        SenderUUID: msg.SenderUUID,
        Timestamp:  msg.ReceivedAt,
        IsPrivate:  detectPrivateMessage(msg),
    }
    
    if c.commandRouter != nil {
        c.commandRouter.HandleMessage(cmdMsg)
    }
}
```

#### 3.2 端到端测试

**测试场景**：
- [ ] 正常骑乘流程
- [ ] 目标超出范围
- [ ] 目标中途离线
- [ ] 鉴权失败
- [ ] 超时检测
- [ ] 冷却期保护
- [ ] 多指令并发

#### 3.3 性能测试

**目标**：10个Bot并发运行

**测试指标**：
- [ ] 内存占用 < 1GB
- [ ] CPU占用 < 50%
- [ ] 视角更新延迟 < 100ms

---

### Phase 4: 文档与优化（Day 6-7）

#### 4.1 代码优化

- [ ] 错误处理统一
- [ ] 日志完善
- [ ] 单元测试覆盖 > 80%

#### 4.2 文档完善

- [ ] README更新
- [ ] API文档
- [ ] 使用示例

#### 4.3 清理

- [ ] 删除调试代码
- [ ] 优化常量命名
- [ ] 代码格式化

---

## 文件清单

```
internal/commands/
├── commands.go                 # 主入口
├── types.go                   # 类型定义
├── router.go                  # 路由器
├── auth.go                    # 鉴权
├── state.go                   # 状态机
├── config.go                  # 配置
├── adapter/
│   └── bot_adapter.go         # Bot适配器
├── modules/
│   ├── module.go              # 模块接口
│   └── ride/
│       ├── ride.go            # !ride指令
│       ├── look.go            # 视角锁定
│       └── config.go          # 模块配置
└── utils/
    ├── chat.go                # 聊天工具
    └── math.go                # 数学工具
```

---

## 里程碑

| 阶段 | 完成时间 | 交付物 |
|------|---------|--------|
| Phase 1 | Day 2 | 核心框架可运行 |
| Phase 2 | Day 4 | !ride指令完成 |
| Phase 3 | Day 5 | 集成测试通过 |
| Phase 4 | Day 7 | 文档完善，发布v1.0 |

---

## 风险与应对

| 风险 | 影响 | 应对措施 |
|------|------|---------|
| 视角数据包协议不匹配 | 高 | 查阅协议文档，添加版本兼容 |
| 玩家位置获取不准确 | 中 | 结合Entity Tracker和PlayerInfo |
| 超时检测不准确 | 低 | 使用time.Since()而非time.Now() |
| 多Bot状态干扰 | 中 | 每个Bot独立Router实例 |

---

## 后续规划

1. **v1.1**：增加更多指令（如!goto, !follow等）
2. **v1.2**：支持动态加载指令模块
3. **v2.0**：集成集群管理器
