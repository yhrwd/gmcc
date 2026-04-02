# GMCC 集群化管理方案 (Cluster Manager Specification)

**文档版本**: v2.0  
**创建日期**: 2026-04-02  
**状态**: Draft

---

## 1. 设计目标

实现 GMCC 的多账号集群化运行能力，支持同时管理多个 Minecraft 客户端实例：

- **多账号并发挂机**：支持多个账号同时连接到不同或相同的服务器
- **统一管理接口**：提供 Go API 管理所有实例（供 Web 面板调用）
- **资源隔离**：各账号实例独立运行，互不干扰
- **状态监控**：实时监控各账号的连接状态和游戏状态
- **故障恢复**：自动重连和异常处理机制

---

## 2. 架构设计

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                     Web HTTP Server                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                 Cluster Manager                      │   │
│  │  ┌─────────────────────────────────────────────┐  │   │
│  │  │              Instance Manager                  │  │   │
│  │  │   ┌─────────┐ ┌─────────┐ ┌─────────┐       │  │   │
│  │  │   │Instance │ │Instance │ │Instance │       │  │   │
│  │  │   │   1     │ │   2     │ │   N     │       │  │   │
│  │  │   └────┬────┘ └────┬────┘ └────┬────┘       │  │   │
│  │  │        │           │           │            │  │   │
│  │  │   ┌────┴───────────┴───────────┴────┐       │  │   │
│  │  │   │     Goroutine per Instance      │       │  │   │
│  │  │   └─────────────────────────────────┘       │  │   │
│  │  └─────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│                      Data Storage                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │cluster.yaml  │  │  sessions/   │  │   logs/      │     │
│  │(multi-account)│  │ (encrypted   │  │ (audit log)  │     │
│  └──────────────┘  │   tokens)    │  └──────────────┘     │
│                    └──────────────┘                         │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 核心组件

#### 2.2.1 ClusterManager

**职责**：
- 初始化和配置加载
- 实例生命周期管理（启动、停止）
- 全局状态聚合
- 信号处理和优雅退出

**接口**：
```go
type ClusterManager interface {
    // 生命周期
    Start(ctx context.Context) error
    Stop() error
    
    // 实例管理
    StartInstance(instanceID string) error
    StopInstance(instanceID string) error
    GetInstance(instanceID string) (*Instance, error)
    
    // 查询
    ListInstances() []InstanceInfo
    GetClusterStatus() ClusterStatus
}
```

#### 2.2.2 InstanceManager

**职责**：
- 维护实例注册表
- 实例创建和销毁
- 并发控制

#### 2.2.3 Instance（实例）

**职责**：
- 封装单个 Minecraft 客户端
- 独立的事件循环
- 状态上报
- 自动重连

---

## 3. 数据模型

### 3.1 ClusterConfig（集群配置）

```go
// ClusterConfig 多账号集群配置
type ClusterConfig struct {
    // 全局配置
    Global GlobalConfig `yaml:"global"`
    
    // 账号列表
    Accounts []AccountEntry `yaml:"accounts"`
    
    // 日志配置
    Log LogConfig `yaml:"log"`
}

type GlobalConfig struct {
    // 最大并发实例数
    MaxInstances int `yaml:"max_instances"`
    
    // 全局重连策略
    ReconnectPolicy ReconnectPolicy `yaml:"reconnect_policy"`
}

type ReconnectPolicy struct {
    Enabled    bool          `yaml:"enabled"`
    MaxRetries int           `yaml:"max_retries"`
    BaseDelay  time.Duration `yaml:"base_delay"`
    MaxDelay   time.Duration `yaml:"max_delay"`
    Multiplier float64       `yaml:"multiplier"`
}

type AccountEntry struct {
    // 实例ID（唯一标识）
    ID string `yaml:"id"`
    
    // 账号配置
    PlayerID        string `yaml:"player_id"`
    UseOfficialAuth bool   `yaml:"use_official_auth"`
    ServerAddress   string `yaml:"server_address"`
    
    // 是否启用
    Enabled bool `yaml:"enabled"`
}
```

### 3.2 Instance（实例）

```go
type Instance struct {
    // 元数据
    ID       string         `json:"id"`
    Account  AccountEntry   `json:"account"`
    Status   InstanceStatus `json:"status"`
    
    // 运行时（不导出）
    client  *mcclient.Client
    runner  *headless.Runner
    cancel  context.CancelFunc
    errChan chan error
    
    // 状态
    mu              sync.RWMutex
    startTime       time.Time
    lastActive      time.Time
    reconnectCount  int
}

type InstanceStatus string

const (
    StatusPending     InstanceStatus = "pending"
    StatusStarting    InstanceStatus = "starting"
    StatusRunning     InstanceStatus = "running"
    StatusReconnecting InstanceStatus = "reconnecting"
    StatusStopped     InstanceStatus = "stopped"
    StatusError       InstanceStatus = "error"
)

type InstanceInfo struct {
    ID              string         `json:"id"`
    PlayerID        string         `json:"player_id"`
    ServerAddress   string         `json:"server_address"`
    Status          InstanceStatus `json:"status"`
    OnlineDuration  time.Duration  `json:"online_duration"`
    LastActive      time.Time      `json:"last_active"`
    ReconnectCount  int            `json:"reconnect_count"`
    Error           string         `json:"error,omitempty"`
}
```

---

## 4. 关键机制

### 4.1 实例生命周期

```
创建(Pending) → 启动(Starting) → 运行(Running) → 停止(Stopped)
                     ↑                │
                     └──── 重连(Reconnecting) ←── 连接断开
```

**启动流程**：
1. 验证配置
2. 创建 mcclient.Client
3. 启动 goroutine 运行事件循环
4. 状态变为 Running

**停止流程**：
1. 调用 cancel() 发送停止信号
2. 等待 goroutine 退出
3. 清理资源
4. 状态变为 Stopped

### 4.2 自动重连

```go
type ReconnectStrategy struct {
    policy      ReconnectPolicy
    attempts    int
    lastAttempt time.Time
    backoff     time.Duration
}

func (r *ReconnectStrategy) ShouldReconnect(err error) bool {
    if isFatalError(err) {
        return false
    }
    if r.policy.MaxRetries > 0 && r.attempts >= r.policy.MaxRetries {
        return false
    }
    return true
}

func (r *ReconnectStrategy) NextBackoff() time.Duration {
    // 指数退避
    r.backoff = time.Duration(float64(r.backoff) * r.policy.Multiplier)
    if r.backoff > r.policy.MaxDelay {
        r.backoff = r.policy.MaxDelay
    }
    return r.backoff
}
```

---

## 5. 配置示例

```yaml
# cluster.yaml
cluster:
  global:
    max_instances: 10
    reconnect_policy:
      enabled: true
      max_retries: 5
      base_delay: 5s
      max_delay: 300s
      multiplier: 2.0
      
  accounts:
    - id: bot_001
      player_id: "BotPlayer001"
      use_official_auth: false
      server_address: "127.0.0.1:25565"
      enabled: true
      
    - id: bot_002
      player_id: "BotPlayer002"
      use_official_auth: true
      server_address: "127.0.0.1:25565"
      enabled: true
      
    - id: bot_003
      player_id: "BotPlayer003"
      use_official_auth: false
      server_address: "192.168.1.100:25565"
      enabled: false
      
  log:
    log_dir: "logs"
    max_size: 512
    debug: false
    enable_file: true
```

---

## 6. 实现优先级

### P0 - 核心功能（MVP）
1. **ClusterConfig 配置结构** - 支持多账号配置
2. **InstanceManager 实例管理** - 注册表和并发控制
3. **基础实例生命周期** - 启动/停止
4. **多账号并发运行** - 每个账号独立 goroutine
5. **基础状态监控** - 状态查询接口

### P1 - 增强功能
1. **自动重连机制** - 指数退避策略
2. **错误处理** - 实例错误状态

### P2 - 后续扩展
1. **实例间通信** - 账号间消息传递
2. **性能指标** - 内存/CPU监控

---

## 7. 与其他模块的关系

```
┌──────────────┐      ┌──────────────────┐      ┌──────────────┐
│   Web HTTP   │ ───▶ │  ClusterManager  │ ───▶ │   Instance   │
│   Server     │      │                  │      │  (per bot)   │
└──────────────┘      └──────────────────┘      └──────────────┘
                            │                          │
                            ▼                          ▼
                      ┌──────────┐              ┌──────────┐
                      │ cluster. │              │ mcclient.│
                      │ yaml     │              │ Client   │
                      └──────────┘              └──────────┘
                           ▲                          │
                           │                          ▼
                      ┌──────────┐              ┌──────────┐
                      │  Token   │◀─────────────│  session.│
                      │  Vault   │              │  Cache   │
                      └──────────┘              └──────────┘
```

**调用关系**：
- **Web Server** → **ClusterManager**: 启动/停止/查询实例
- **ClusterManager** → **Instance**: 管理生命周期
- **Instance** → **mcclient.Client**: 实际 Minecraft 连接
- **Instance** → **session.Cache**: 读取/写入 Token（加密存储）

---

## 8. 目录结构

```
internal/
├── cluster/
│   ├── manager.go          # ClusterManager 实现
│   ├── instance.go         # Instance 定义
│   ├── registry.go         # 实例注册表
│   ├── config.go           # ClusterConfig 定义
│   └── reconnect.go        # 重连策略
├── web/
│   ├── server.go           # HTTP Server
│   ├── handlers.go         # API handlers
│   ├── auth.go             # 认证相关
│   └── token_vault.go      # Token 加密存储
└── ...
```

---

## 9. 与 Web 面板的集成

ClusterManager 作为 Web 面板的后端服务：

```go
// Web Handler 调用 ClusterManager
func (h *Handler) StartInstance(c *gin.Context) {
    instanceID := c.Param("id")
    
    // 验证密码
    if !h.authManager.VerifyPassword(c.PostForm("password")) {
        c.JSON(401, gin.H{"error": "invalid_password"})
        return
    }
    
    // 调用 ClusterManager
    if err := h.clusterManager.StartInstance(instanceID); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    // 记录操作日志
    h.auditLogger.Log("instance_start", instanceID, c.ClientIP())
    
    c.JSON(200, gin.H{"success": true})
}
```

---

**文档版本**: v2.0  
**创建日期**: 2026-04-02  
**状态**: Draft
