# 实体跟踪系统设计文档

**日期**: 2026-03-26  
**版本**: 1.0  
**协议**: Minecraft Java 1.21.11 (protocol 774)

## 概述

本设计文档描述 gmcc 客户端如何实现周围玩家坐标跟踪功能。系统支持实时获取周围玩家的位置和移动信息，同时预留接口供未来扩展支持其他实体类型。

## 设计目标

1. **核心功能**: 实时跟踪周围玩家的坐标位置
2. **预留扩展**: 为未来支持其他实体类型（怪物、掉落物等）预留接口
3. **性能优化**: 通过去重机制避免高频回调风暴
4. **简单易用**: 提供清晰的查询API和事件回调

## 架构设计

### 分层架构

```
┌─────────────────────────────────────┐
│         User Application            │
│   - GetNearbyPlayers()              │
│   - OnPlayerEnter/Leave/Move         │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│    player.NearbyTracker             │
│   - 筛选玩家实体                     │
│   - 关联 player_info               │
│   - 提供玩家专用API                  │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│    entity.Tracker                   │
│   - 管理所有实体生命周期              │
│   - 位置更新与去重                   │
│   - 协议包处理                       │
└─────────────────────────────────────┘
```

### 核心数据结构

#### Entity (internal/entity/entity.go)

```go
type Entity struct {
    ID       int32
    Type     string        // "minecraft:player" 或其他实体类型
    UUID     [16]byte      // 可选，不是所有实体都有
    Position Position
    Velocity Vector3
    OnGround bool
    LastUpdate time.Time
}

type Position struct {
    X, Y, Z float64
}

type Vector3 struct {
    X, Y, Z float64
}
```

#### NearbyPlayer (internal/player/nearby.go)

```go
type NearbyPlayer struct {
    *entity.Entity
    Username string  // 来自 player_info
    Latency  int32   // 延迟(ms)
}
```

## 协议实现

### 新增协议常量

```go
// internal/mcclient/protocol/v774.go
const (
    PlayClientAddEntity       int32 = 0x01  // 实体生成
    PlayClientTeleportEntity  int32 = 0x48  // 实体传送
    PlayClientMoveEntityPos   int32 = 0x09  // 位置增量更新
    PlayClientRemoveEntities  int32 = 0x4B  // 实体移除
)
```

### 包处理器

| 包ID | 处理器 | 功能 |
|------|--------|------|
| 0x01 | handleAddEntity | 解析实体生成，关联UUID，标记玩家类型 |
| 0x48 | handleTeleportEntity | 处理完整位置同步，触发移动回调 |
| 0x09 | handleMoveEntityPos | 处理位置增量更新（delta * 4096） |
| 0x4B | handleRemoveEntities | 移除实体，触发离开回调 |
| 0x42 | handlePlayerInfoUpdate | 维护 player_info，供实体生成时关联 |

### 玩家识别机制

通过 entity_type 判断实体类型。玩家实体的 type 为 `minecraft:player`（注册表ID）。当收到 `add_entity` 包时：

1. 解析 entity_type 字段
2. 如果 type 是 "minecraft:player"，在 player_info 表中查找匹配的 UUID
3. 如果找到，创建 `NearbyPlayer`，关联 username
4. 如果未找到，作为普通实体存储（可能是尚未收到 player_info）

### 去重机制

为避免高频移动回调，实现100ms窗口期去重：

```go
type Tracker struct {
    mu            sync.RWMutex
    entities      map[int32]*Entity
    byUUID        map[[16]byte]*Entity
    pendingUpdates map[int32]*pendingUpdate
    callbacks     Callbacks
}

type pendingUpdate struct {
    entityID int32
    newPos   Position
    timer    *time.Timer
}

func (t *Tracker) updatePosition(id int32, newPos Position) {
    t.mu.Lock()
    defer t.mu.Unlock()
    
    // 检查是否有待处理的同一实体更新
    if pending, ok := t.pendingUpdates[id]; ok {
        pending.newPos = newPos
        return
    }
    
    // 延迟100ms后执行回调
    timer := time.AfterFunc(100*time.Millisecond, func() {
        t.mu.Lock()
        e, exists := t.entities[id]
        if !exists {
            // 实体已被移除，取消回调
            delete(t.pendingUpdates, id)
            t.mu.Unlock()
            return
        }
        oldPos := e.Position
        e.Position = newPos
        delete(t.pendingUpdates, id)
        
        // 复制数据用于回调
        entityCopy := *e
        t.mu.Unlock()
        
        if t.callbacks.OnMove != nil {
            t.callbacks.OnMove(&entityCopy, oldPos)
        }
    })
    
    t.pendingUpdates[id] = &pendingUpdate{id, newPos, timer}
}
```

## 用户API

### 查询接口

```go
// GetNearbyPlayers 获取周围所有玩家
func (t *NearbyTracker) GetNearbyPlayers() []*NearbyPlayer

// GetNearbyPlayer 通过UUID查找特定玩家
func (t *NearbyTracker) GetNearbyPlayer(uuid [16]byte) (*NearbyPlayer, bool)

// PlayersWithinDistance 获取指定距离内的玩家
func (t *NearbyTracker) PlayersWithinDistance(
    center entity.Position, 
    distance float64
) []*NearbyPlayer
```

### 回调接口

```go
type PlayerCallbacks struct {
    OnPlayerEnter func(p *NearbyPlayer)  // 玩家进入视距
    OnPlayerLeave func(p *NearbyPlayer)  // 玩家离开视距
    OnPlayerMove  func(p *NearbyPlayer, oldPos entity.Position)
}

func (t *NearbyTracker) SetCallbacks(callbacks PlayerCallbacks)
```

### 使用示例

```go
// 注册回调
client.NearbyPlayers.SetCallbacks(player.PlayerCallbacks{
    OnPlayerEnter: func(p *player.NearbyPlayer) {
        log.Printf("玩家 %s 进入视距", p.Username)
    },
    OnPlayerMove: func(p *player.NearbyPlayer, oldPos entity.Position) {
        log.Printf("玩家 %s 移动", p.Username)
    },
})

// 查询周围玩家
for _, p := range client.NearbyPlayers.GetNearbyPlayers() {
    fmt.Printf("%s: %.1f %.1f %.1f\n", p.Username, p.X, p.Y, p.Z)
}
```

## Client集成

### Client结构变更

```go
type Client struct {
    // ... 现有字段 ...
    
    entityTracker *entity.Tracker        // 内部：所有实体管理
    NearbyPlayers *player.NearbyTracker  // 公开：周围玩家API
}
```

### 初始化流程

```go
// initializeTrackers 在Play阶段开始时调用
func (c *Client) initializeTrackers() {
    // 创建实体跟踪器
    c.entityTracker = entity.NewTracker()
    
    // 创建玩家跟踪器，传入entityTracker和playerInfo引用
    c.NearbyPlayers = player.NewNearbyTracker(c.entityTracker, c.getPlayerInfoByUUID)
    
    // 注册包处理器
    c.RegisterHandler(protocol.PlayClientAddEntity, c.handleAddEntity)
    c.RegisterHandler(protocol.PlayClientTeleportEntity, c.handleTeleportEntity)
    c.RegisterHandler(protocol.PlayClientMoveEntityPos, c.handleMoveEntityPos)
    c.RegisterHandler(protocol.PlayClientRemoveEntities, c.handleRemoveEntities)
}

// getPlayerInfoByUUID 提供UUID到username的查找
func (c *Client) getPlayerInfoByUUID(uuid [16]byte) (username string, found bool) {
    c.playersMu.RLock()
    defer c.playersMu.RUnlock()
    
    for name, info := range c.players {
        if info.uuid == uuid {
            return name, true
        }
    }
    return "", false
}
```

### 资源清理

```go
// Shutdown 在客户端断开时调用
func (c *Client) Shutdown() {
    if c.entityTracker != nil {
        c.entityTracker.Stop()
    }
}
```

## 并发安全

### 锁策略

1. **EntityTracker**: `sync.RWMutex` 保护实体映射表
2. **NearbyTracker**: 只读访问，不额外加锁
3. **回调执行**: 在锁外执行，避免死锁

### 数据复制

回调时复制实体数据，避免外部修改影响内部状态：

```go
func (t *Tracker) handlePositionUpdate(id int32, newPos Position) {
    t.mu.Lock()
    entityCopy := *t.entities[id]
    t.mu.Unlock()
    
    // 在锁外回调
    if t.callbacks.OnMove != nil {
        t.callbacks.OnMove(&entityCopy, oldPos)
    }
}
```

## 错误处理

| 场景 | 处理方式 |
|------|----------|
| 包解析失败 | 记录日志，跳过该包 |
| UUID未找到 | 作为普通实体存储，不标记为玩家 |
| 重复实体ID | 覆盖旧数据，以服务器状态为准 |
| 未知实体类型 | 正常存储，仅记录类型字符串 |

## 测试策略

### 单元测试

1. **Tracker测试**: 增删改查操作和回调触发
2. **包解析测试**: 模拟二进制数据验证解析正确性
3. **并发测试**: 模拟高频更新验证无数据竞争

### 集成测试

1. 连接测试服务器验证实体跟踪功能
2. 测试玩家进入/离开视距场景
3. 验证位置更新和回调去重

## 未来扩展

### 支持其他实体类型

```go
// 添加实体类型过滤器
type EntityFilter struct {
    Types []string      // 如 ["minecraft:zombie", "minecraft:creeper"]
    Range float64       // 距离范围
}

func (t *Tracker) GetFilteredEntities(filter EntityFilter) []*Entity
```

### 元数据解析

```go
// 扩展Entity结构支持元数据
type Entity struct {
    // ... 现有字段 ...
    Health float32  // 从metadata解析
    CustomName string
}
```

## 实现清单

1. 添加协议常量 (v774.go)
2. 创建 entity 包和基础结构
3. 实现 EntityTracker 核心逻辑
4. 实现包处理器 (4个)
5. 扩展 player.NearbyTracker
6. Client集成和API暴露
7. 编写单元测试
8. 集成测试验证

## 附录

### 包格式详情

#### add_entity (0x01)

```
[int32]    entity_id
[UUID]     uuid
[varint]   entity_type (registry ID)
[double]   x
[double]   y
[double]   z
[byte]     pitch / 256
[byte]     yaw / 256
[byte]     head_yaw / 256
[varint]   entity_data
[double]   velocity_x
[double]   velocity_y
[double]   velocity_z
```

#### teleport_entity (0x48)

```
[int32]    entity_id
[double]   x
[double]   y
[double]   z
[byte]     relative_flags
[double]   velocity_x
[double]   velocity_y
[double]   velocity_z
[float]    yaw
[float]    pitch
[bool]     on_ground
```

#### move_entity_pos (0x09)

```
[int32]    entity_id
[short]    delta_x (blocks * 4096)
[short]    delta_y (blocks * 4096)
[short]    delta_z (blocks * 4096)
[bool]     on_ground
```

#### remove_entities (0x4B)

```
[varint]   count
[int32...] entity_ids
```
