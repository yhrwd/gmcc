# GMCC Manager 技术规格文档

## 1. 项目概述

### 1.1 目标
为 gmcc 增加多账号集群管理能力，支持同时挂机 10 个以内 Minecraft 账号，实现：
- 统一的 HTTP API 管理接口
- 聊天消息去重记录（仅保留主账号公共聊天，保留所有私聊）
- WebSocket 实时状态推送
- 主账号自动切换机制
- 自动断线重连
- 日志轮转防止文件过大

### 1.2 技术约束
- **账号数量**：≤10 个
- **服务器**：固定单服务器
- **内存预算**：1GB 以内
- **日志策略**：按大小切分，单文件 ≤100MB
- **指令支持**：仅 `!ride`
- **部署方式**：服务器长期运行

---

## 2. 系统架构

### 2.1 总体架构

```
┌────────────────────────────────────────────────────────────────┐
│                        HTTP Client                             │
│                    (curl / Dashboard / Scripts)                  │
└───────────────────────────┬────────────────────────────────────┘
                            │ HTTP / WebSocket
                            ▼
┌────────────────────────────────────────────────────────────────┐
│                      API Server (Gin)                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐      │
│  │ REST API     │  │ WebSocket    │  │ Auth Middleware  │      │
│  │ /api/*       │  │ /ws          │  │ Bearer Token     │      │
│  └──────────────┘  └──────────────┘  └──────────────────┘      │
└───────────────────────────┬────────────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────────────┐
│                      Bot Manager                                 │
│  ┌────────────────────────────────────────────────────────┐    │
│  │  Bot Pool (max 10)                                     │    │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐     ┌─────────┐   │    │
│  │  │ Bot #1  │ │ Bot #2  │ │ Bot #3  │ ... │ Bot #N  │   │    │
│  │  │(Primary)│ │         │ │         │     │         │   │    │
│  │  └────┬────┘ └────┬────┘ └────┬────┘     └────┬────┘   │    │
│  └───────┼───────────┼───────────┼───────────────┼────────┘    │
│          │           │           │               │              │
│          └───────────┴───────────┴───────────────┘              │
│                      │                                         │
│  ┌───────────────────┼──────────────────────────────────┐     │
│  │ Reconnect Manager │ Command Router │ Chat Deduplicator │     │
│  └───────────────────┴──────────────────────────────────┘     │
└───────────────────────────┬────────────────────────────────────┘
                            │
                            ▼
                   ┌─────────────────┐
                   │  Target Server  │
                   │ mc.hypixel.net  │
                   └─────────────────┘
```

### 2.2 目录结构

```
gmcc/
├── cmd/
│   ├── gmcc/                    # 原有单账号客户端
│   └── manager/                 # 新增集群管理器
│       └── main.go
├── internal/
│   ├── manager/                 # 核心管理模块
│   │   ├── manager.go           # BotManager 主结构
│   │   ├── bot.go               # Bot 实例包装
│   │   ├── pool.go              # Bot 池管理
│   │   ├── reconnect.go         # 重连管理器
│   │   ├── chat.go              # 聊天去重逻辑
│   │   └── command.go           # 指令路由
│   ├── api/
│   │   ├── server.go            # HTTP 服务器
│   │   ├── handlers.go          # REST handlers
│   │   ├── websocket.go         # WebSocket 处理
│   │   └── middleware.go        # 中间件
│   └── logrotate/               # 日志轮转模块
│       └── rotator.go
└── docs/
    └── manager_spec.md          # 本文件
```

---

## 3. 核心模块设计

### 3.1 BotManager

**职责**：管理所有 Bot 实例的生命周期

```go
type BotManager struct {
    // Bot 存储
    bots       map[string]*Bot        // ID -> Bot
    byPlayerID map[string]string      // player_id -> bot_id
    botsMu     sync.RWMutex
    
    // 主账号管理
    primaryBotID string               // 当前主账号
    primaryMu    sync.RWMutex
    
    // 组件
    reconnector *ReconnectManager    // 重连管理器
    chatLogger  *ChatLogger          // 聊天日志管理器
    cmdRouter   *CommandRouter        // 指令路由器
    wsHub       *WebSocketHub         // WebSocket 广播中心
    
    // 配置
    config *ManagerConfig
    
    // 上下文
    ctx    context.Context
    cancel context.CancelFunc
}
```

**核心方法**：

```go
// Bot 生命周期
func (m *BotManager) CreateBot(cfg *config.Config) (*Bot, error)
func (m *BotManager) RemoveBot(id string) error
func (m *BotManager) StartBot(id string) error
func (m *BotManager) StopBot(id string) error

// 主账号管理
func (m *BotManager) SetPrimaryBot(id string) error      // 手动设置
func (m *BotManager) AutoSelectPrimary() string          // 自动选择第一个在线的
func (m *BotManager) OnBotOffline(id string)             // 主账号下线回调

// 广播
func (m *BotManager) BroadcastToAll(msg string) error     // 发送消息到所有 Bot
func (m *BotManager) SendCommand(id, cmd string) error  // 发送命令到指定 Bot
```

### 3.2 Bot 实例

```go
type Bot struct {
    ID       string
    PlayerID string
    Config   config.Config
    
    // 状态
    Status      BotStatus
    Client      *mcclient.Client
    CancelFunc  context.CancelFunc
    
    // 统计
    OnlineSince      time.Time
    MessagesReceived int64
    MessagesSent     int64
    CommandsExecuted int64
    LastError        string
    
    // 重连计数
    reconnectCount int
    
    mu sync.RWMutex
}

type BotStatus string
const (
    StatusOffline     BotStatus = "offline"
    StatusConnecting  BotStatus = "connecting"
    StatusAuthenticating BotStatus = "authenticating"
    StatusOnline      BotStatus = "online"
    StatusError       BotStatus = "error"
    StatusReconnecting BotStatus = "reconnecting"
)
```

### 3.3 主账号切换机制

**策略**：
1. **初始化**：第一个成功上线的 Bot 自动成为主账号
2. **主账号主动下线**：立即在所有在线 Bot 中选择新的主账号（按创建顺序）
3. **主账号异常断开**：触发重连机制，重连期间临时选择次主账号
4. **主账号重连成功**：恢复原主账号身份

```go
type PrimaryManager struct {
    currentID    string
    pendingID    string    // 重连中的临时主账号
    mu           sync.RWMutex
}

func (pm *PrimaryManager) OnBotOnline(botID string) {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    
    // 如果当前无主账号，或者这是原主账号重连成功
    if pm.currentID == "" {
        pm.currentID = botID
        logx.Infof("主账号已设置: %s", botID)
    } else if botID == pm.pendingID {
        // 原主账号重连成功，恢复身份
        pm.currentID = botID
        pm.pendingID = ""
        logx.Infof("主账号 %s 恢复上线", botID)
    }
}

func (pm *PrimaryManager) OnBotOffline(botID string, hasReconnect bool) string {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    
    if pm.currentID != botID {
        return "" // 不是主账号下线，无需处理
    }
    
    if hasReconnect {
        // 主账号还会重连，设置临时主账号
        pm.pendingID = pm.currentID
        pm.currentID = "" // 等待新主账号
        return "need_new_primary"
    } else {
        // 主账号彻底下线
        pm.currentID = ""
        pm.pendingID = ""
        return "need_new_primary"
    }
}
```

---

## 4. 聊天去重与日志系统

### 4.1 去重策略

```go
type ChatLogger struct {
    primaryBotID   string
    primaryMu      sync.RWMutex
    
    // 日志文件
    publicLog      *RotatingLogger    // 公共聊天日志
    privateLogs    map[string]*RotatingLogger // Bot私聊日志
    logsMu         sync.RWMutex
    
    // 去重（防止网络延迟导致的重复）
    recentMessages *lru.Cache         // 最近100条消息缓存
}

func (cl *ChatLogger) ShouldLog(botID string, msg ChatMessage) (bool, string) {
    // 检查是否重复
    msgKey := fmt.Sprintf("%s:%s:%d", msg.Type, msg.PlainText, msg.Timestamp.Unix())
    if cl.recentMessages.Contains(msgKey) {
        return false, ""
    }
    cl.recentMessages.Add(msgKey, true)
    
    // 私聊：每个 Bot 单独记录
    if msg.IsPrivate {
        return true, cl.getPrivateLogPath(botID)
    }
    
    // 公共聊天：只有主账号记录
    cl.primaryMu.RLock()
    isPrimary := (botID == cl.primaryBotID)
    cl.primaryMu.RUnlock()
    
    if isPrimary {
        return true, cl.publicLog.Path()
    }
    
    return false, ""
}
```

### 4.2 日志轮转策略

**触发条件**：
- 单个文件达到 100MB 时自动切分
- 保留最近 10 个切分文件（共 1GB）
- 启动时检查并清理过期文件

**文件命名**：
- `public.log` - 当前公共日志
- `public.log.1` - 上一个切分文件
- `public.log.2.gz` - 更早的文件（自动压缩）
- `private/Bot001.log` - Bot001 私聊日志

```go
type RotatingLogger struct {
    path        string
    maxSize     int64        // 100MB
    maxBackups  int          // 10
    compress    bool         // true
    
    currentSize int64
    file        *os.File
    mu          sync.Mutex
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    // 检查是否需要切分
    if rl.currentSize+int64(len(p)) > rl.maxSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err = rl.file.Write(p)
    rl.currentSize += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    // 1. 关闭当前文件
    rl.file.Close()
    
    // 2. 重命名文件 (log -> log.1)
    backup := rl.path + ".1"
    os.Rename(rl.path, backup)
    
    // 3. 压缩旧文件
    go rl.compressOldBackups()
    
    // 4. 创建新文件
    f, err := os.OpenFile(rl.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        return err
    }
    rl.file = f
    rl.currentSize = 0
    
    return nil
}

func (rl *RotatingLogger) compressOldBackups() {
    // 异步压缩 log.2 -> log.2.gz
    // 删除超过 maxBackups 的旧文件
}
```

---

## 5. 断线重连机制

### 5.1 重连策略

```go
type ReconnectManager struct {
    manager     *BotManager
    intervals   []time.Duration  // [10s, 30s, 60s, 300s, 600s]
    maxAttempts int              // 5
    
    active      map[string]*ReconnectTask
    mu          sync.RWMutex
}

type ReconnectTask struct {
    BotID       string
    Attempts    int
    NextRetry   time.Time
    Cancelled   bool
}
```

**重连流程**：

```
Bot 异常断开
    │
    ▼
检查是否主动停止
    │
    ├── 是 → 不重连，标记 offline
    │
    └── 否 → 触发重连
                │
                ▼
        延迟 10s 后尝试
                │
                ▼
        创建新 Client 实例
                │
                ▼
        尝试连接 + 认证
                │
        ┌───────┴───────┐
        ▼               ▼
    成功            失败
        │               │
        ▼               ▼
    上线成功      增加重试计数
    恢复状态      延迟时间翻倍
                  (max 10分钟)
                        │
                        ▼
                  超过最大次数?
                        │
                ┌───────┴───────┐
                ▼               ▼
            是              否
                │               │
                ▼               ▼
        标记永久离线      继续重试
```

### 5.2 重连与主账号切换

```go
func (m *BotManager) handleBotDisconnect(botID string, err error) {
    bot := m.bots[botID]
    
    // 检查是否是主动停止
    if bot.Status == StatusOffline {
        return // 不重连
    }
    
    // 标记为 reconnecting
    bot.Status = StatusReconnecting
    
    // 处理主账号逻辑
    action := m.primaryMgr.OnBotOffline(botID, true)
    if action == "need_new_primary" {
        // 主账号正在重连，临时选择新主账号
        newPrimary := m.selectNextPrimary(botID)
        if newPrimary != "" {
            m.primaryMgr.SetTemporaryPrimary(newPrimary)
            logx.Infof("主账号 %s 重连中，临时主账号: %s", botID, newPrimary)
        }
    }
    
    // 启动重连任务
    m.reconnector.StartReconnect(botID)
}

func (m *BotManager) handleBotReconnectSuccess(botID string) {
    // 如果是原主账号重连成功
    if m.primaryMgr.IsPendingPrimary(botID) {
        m.primaryMgr.RestorePrimary(botID)
        logx.Infof("主账号 %s 重连成功，恢复主账号身份", botID)
    }
    
    // 广播状态更新
    m.broadcastStatus(botID, StatusOnline)
}
```

---

## 6. API 设计

### 6.1 REST API

**认证方式**：`Authorization: Bearer <token>`

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/bots` | 获取所有 Bot 列表 |
| GET | `/api/bots/:id` | 获取单个 Bot 详情 |
| POST | `/api/bots` | 创建新 Bot |
| DELETE | `/api/bots/:id` | 删除 Bot |
| POST | `/api/bots/:id/start` | 启动 Bot |
| POST | `/api/bots/:id/stop` | 停止 Bot |
| POST | `/api/bots/:id/restart` | 重启 Bot |
| POST | `/api/bots/:id/command` | 发送命令 |
| POST | `/api/batch/command` | 批量发送命令 |
| GET | `/api/logs/chat` | 获取聊天日志（SSE 流） |
| GET | `/api/logs/chat/download` | 下载日志文件 |

**请求/响应示例**：

```bash
# 创建 Bot
curl -X POST http://localhost:8080/api/bots \
  -H "Authorization: Bearer sk-abc123" \
  -H "Content-Type: application/json" \
  -d '{
    "player_id": "Bot001",
    "use_official_auth": true,
    "auto_start": true
  }'

# 响应
{
  "id": "bot-uuid-001",
  "player_id": "Bot001",
  "status": "authenticating",
  "auth_url": "https://login.microsoftonline.com/..."  // 如需要认证
}

# 发送命令
curl -X POST http://localhost:8080/api/bots/bot-uuid-001/command \
  -H "Authorization: Bearer sk-abc123" \
  -d '{"command": "msg Player Hello"}'
```

### 6.2 WebSocket 实时推送

**连接**：`ws://localhost:8080/ws?token=sk-abc123`

**消息类型**：

```json
// Bot 状态变更
{
  "type": "bot_status",
  "data": {
    "id": "bot-uuid-001",
    "player_id": "Bot001",
    "status": "online",
    "online_since": "2024-01-15T10:30:00Z",
    "is_primary": true
  },
  "timestamp": "2024-01-15T10:30:05Z"
}

// 聊天消息（仅主账号的公共聊天 + 所有私聊）
{
  "type": "chat_message",
  "data": {
    "bot_id": "bot-uuid-001",
    "bot_name": "Bot001",
    "type": "player_chat",
    "content": "Hello everyone",
    "sender": "Player123",
    "is_private": false
  },
  "timestamp": "2024-01-15T10:31:00Z"
}

// 主账号切换
{
  "type": "primary_changed",
  "data": {
    "old_primary": "bot-uuid-001",
    "new_primary": "bot-uuid-002",
    "reason": "original_primary_offline"
  },
  "timestamp": "2024-01-15T10:35:00Z"
}

// 系统日志
{
  "type": "system_log",
  "data": {
    "level": "info",
    "message": "Bot002 connected successfully"
  },
  "timestamp": "2024-01-15T10:30:01Z"
}
```

---

## 7. 指令系统

### 7.1 `!ride` 指令完整实现

**指令格式**：
- `!ride` - 骑乘发送指令的玩家
- `!ride <玩家名>` - 骑乘指定玩家

**实现机制**：

```go
// CommandRouter 处理流程
type CommandRouter struct {
    prefix      string        // "!"
    handlers    map[string]CommandHandler
    whitelist   []string      // 允许使用指令的玩家列表
}

func (r *CommandRouter) HandleChat(bot *Bot, msg ChatMessage) {
    // 1. 检查是否是私聊
    if !msg.IsPrivate {
        return // 只处理私聊指令
    }
    
    // 2. 检查发送者是否在白名单
    if !r.isWhitelisted(msg.Sender) {
        bot.Client.SendCommand(fmt.Sprintf("msg %s 你没有权限使用此机器人", msg.Sender))
        return
    }
    
    // 3. 解析指令
    text := strings.TrimSpace(msg.PlainText)
    if !strings.HasPrefix(text, r.prefix) {
        return // 不是指令
    }
    
    // 4. 提取命令和参数
    content := strings.TrimPrefix(text, r.prefix)
    parts := strings.Fields(content)
    if len(parts) == 0 {
        return
    }
    
    cmd := strings.ToLower(parts[0])
    args := parts[1:]
    
    // 5. 执行指令
    if cmd == "ride" {
        r.handleRide(bot, msg.Sender, args)
    } else {
        // 未知指令，回复帮助
        bot.Client.SendCommand(fmt.Sprintf("msg %s 未知指令。可用指令: !ride [玩家名]", msg.Sender))
    }
}

func (r *CommandRouter) handleRide(bot *Bot, sender string, args []string) {
    // 确定目标玩家
    target := sender
    if len(args) > 0 {
        target = args[0]
    }
    
    // 发送骑乘命令
    err := bot.Client.SendCommand(fmt.Sprintf("ride %s", target))
    
    // 发送结果反馈（通过私聊）
    if err != nil {
        bot.Client.SendCommand(fmt.Sprintf("msg %s 骑乘失败: %v", sender, err))
    } else {
        bot.Client.SendCommand(fmt.Sprintf("msg %s 正在骑乘 %s", sender, target))
    }
}
```

**私聊检测逻辑**：

```go
func (c *Client) isPrivateMessage(rawJSON string) bool {
    // 检测是否为私聊的翻译键
    // Minecraft 私聊消息通常包含: "commands.message.display.incoming"
    // 或者格式如: "Player whispers to you: message"
    
    // 方法1: 检查翻译键
    if strings.Contains(rawJSON, "commands.message.display.incoming") {
        return true
    }
    
    // 方法2: 检查纯文本特征
    // 英文: "whispers to you"
    // 中文: "悄悄地对你说"
    text := ExtractPlainTextFromChatJSON(rawJSON)
    whisperPatterns := []string{
        "whispers to you",      // 英文
        "悄悄地对你说",        // 中文
        "flüstert dir zu",     // 德语
    }
    
    for _, pattern := range whisperPatterns {
        if strings.Contains(text, pattern) {
            return true
        }
    }
    
    return false
}
```

**完整消息处理流程**：

```
收到聊天消息
    │
    ▼
解析消息类型
    │
    ├── system_chat ──→ 只记录（主账号）
    │
    └── player_chat
            │
            ▼
    检测是否为私聊?
            │
    ┌───────┴───────┐
    ▼               ▼
   是              否
    │               │
    ▼               ▼
 提取发送者      检查是否主账号
    │               │
    ▼               ▼
 检查是否指令    记录公共聊天
    │               │
    ▼               ▼
 是 ──→ 处理指令 否 ──→ 忽略
    │
    ▼
 检查白名单
    │
    ▼
 发送 /ride 命令
    │
    ▼
 发送结果反馈
 (通过私聊)
```

**状态流转**：

| 步骤 | Bot 状态 | 操作 |
|------|---------|------|
| 1 | online | 接收私聊消息 "!ride" |
| 2 | online | 发送 `/ride PlayerName` |
| 3 | online | 等待服务器响应 |
| 4 | online | 发送私聊反馈 "正在骑乘 PlayerName" |

**错误处理**：

```go
type CommandError struct {
    Code    string
    Message string
}

var (
    ErrNotWhitelisted = CommandError{"NOT_WHITELISTED", "你没有权限使用此机器人"}
    ErrTargetNotFound = CommandError{"TARGET_NOT_FOUND", "目标玩家不在线"}
    ErrAlreadyRiding  = CommandError{"ALREADY_RIDING", "已经在骑乘其他实体"}
    ErrCommandFailed  = CommandError{"COMMAND_FAILED", "命令执行失败"}
)

func (r *CommandRouter) sendError(bot *Bot, target string, err CommandError) {
    // 通过私聊发送错误信息
    bot.Client.SendCommand(fmt.Sprintf("msg %s 错误: %s", target, err.Message))
}
```

**日志记录**：

```go
// 指令执行日志（每个 Bot 单独记录）
logLine := fmt.Sprintf("[%s] Command from %s: %s, Args: %v, Result: %v",
    time.Now().Format("2006-01-02 15:04:05"),
    sender,
    "ride",
    args,
    success,
)
// 写入 logs/commands/Bot001.log
```

---

## 8. 配置

```yaml
# config.yaml
manager:
  # HTTP API 配置
  api:
    host: "0.0.0.0"
    port: 8080
    token: "sk-abc123xyz789"  # 32位以上随机字符串
  
  # 固定服务器配置
  server:
    address: "mc.hypixel.net:25565"
    reconnect:
      enabled: true
      initial_delay: 10s       # 首次重连延迟
      max_delay: 600s          # 最大延迟
      multiplier: 2.0          # 延迟翻倍系数
      max_attempts: 5        # 最大重试次数
  
  # 主账号配置
  primary:
    auto_select: true          # 自动选择第一个上线的
    switch_on_disconnect: true # 主账号下线时自动切换
  
  # 日志配置
  log:
    chat:
      dir: "logs/chat"
      max_size: 100           # MB
      max_backups: 10         # 保留份数
      compress: true          # 压缩旧日志
    bots:
      dir: "logs/bots"
      max_size: 50
      max_backups: 5
  
  # 预定义账号（最多10个）
  accounts:
    - player_id: "Bot001"
      use_official_auth: true
      auto_start: true
    - player_id: "Bot002"
      use_official_auth: true
      auto_start: true
    # ... 最多10个
```

---

## 9. 部署

### 9.1 运行方式

```bash
# 编译
go build -o gmcc-manager ./cmd/manager

# 运行
./gmcc-manager -config config.yaml

# 或使用环境变量
export GMCC_API_TOKEN="sk-your-token"
export GMCC_SERVER_ADDR="mc.hypixel.net:25565"
./gmcc-manager
```

### 9.2 Systemd 服务

```ini
# /etc/systemd/system/gmcc-manager.service
[Unit]
Description=GMCC Manager
After=network.target

[Service]
Type=simple
User=minecraft
WorkingDirectory=/opt/gmcc
ExecStart=/opt/gmcc/gmcc-manager -config /opt/gmcc/config.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

---

## 10. 测试计划

### 10.1 单元测试
- Bot 生命周期管理
- 主账号切换逻辑
- 日志轮转
- 消息去重

### 10.2 集成测试
- 10 个 Bot 同时在线
- 主账号主动下线切换
- 网络异常重连
- 日志文件切分

### 10.3 压力测试
- 持续运行 24 小时
- 模拟频繁断网
- 内存泄漏检查

---

## 11. 风险与应对

| 风险 | 可能性 | 影响 | 应对措施 |
|------|--------|------|---------|
| 微软认证频率限制 | 中 | 高 | 使用 Token 缓存，减少重新认证 |
| 服务器封 IP | 低 | 高 | 支持代理配置（二期） |
| 日志磁盘满 | 中 | 中 | 自动清理，限制总大小 |
| 内存泄漏 | 低 | 中 | 定期重启，监控内存使用 |

---

## 12. 版本规划

### v1.0 (MVP)
- [x] 基础多 Bot 管理
- [x] HTTP API
- [x] 主账号自动切换
- [x] 日志轮转
- [x] 断线重连
- [x] WebSocket 推送
- [x] `!ride` 指令

### v1.1 (可选)
- [ ] 简单 Dashboard
- [ ] 更多指令支持
- [ ] 代理支持
- [ ] 性能监控

---

文档版本：v1.0  
最后更新：2024-01-15  
作者：GMCC Dev Team
