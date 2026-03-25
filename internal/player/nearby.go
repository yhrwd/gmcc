package player

import (
	"sync"

	"gmcc/internal/entity"
)

// NearbyPlayer 表示周围的玩家实体
type NearbyPlayer struct {
	*entity.Entity
	Username string // 来自 player_info
}

// PlayerInfoLookup 函数类型，用于查找玩家信息
type PlayerInfoLookup func(uuid [16]byte) (username string, found bool)

// NearbyTracker 跟踪周围玩家
type NearbyTracker struct {
	mu            sync.RWMutex
	entityTracker *entity.Tracker
	lookupPlayer  PlayerInfoLookup
	players       map[int32]*NearbyPlayer
	callbacks     PlayerCallbacks
}

// PlayerCallbacks 定义玩家事件回调
type PlayerCallbacks struct {
	OnPlayerEnter func(p *NearbyPlayer)
	OnPlayerLeave func(p *NearbyPlayer)
	OnPlayerMove  func(p *NearbyPlayer, oldPos entity.Position)
}

// NewNearbyTracker 创建新的附近玩家跟踪器
func NewNearbyTracker(tracker *entity.Tracker, lookup PlayerInfoLookup) *NearbyTracker {
	nt := &NearbyTracker{
		entityTracker: tracker,
		lookupPlayer:  lookup,
		players:       make(map[int32]*NearbyPlayer),
	}

	// 设置实体回调
	tracker.SetCallbacks(entity.Callbacks{
		OnSpawn:  nt.handleEntitySpawn,
		OnMove:   nt.handleEntityMove,
		OnRemove: nt.handleEntityRemove,
	})

	return nt
}

// SetCallbacks 设置玩家回调
func (nt *NearbyTracker) SetCallbacks(callbacks PlayerCallbacks) {
	nt.mu.Lock()
	defer nt.mu.Unlock()
	nt.callbacks = callbacks
}

// handleEntitySpawn 处理实体生成
func (nt *NearbyTracker) handleEntitySpawn(e *entity.Entity) {
	if !e.IsPlayer() {
		return
	}

	nt.mu.Lock()
	defer nt.mu.Unlock()

	// 查找玩家信息
	username, found := nt.lookupPlayer(e.UUID)
	if !found {
		// 未找到玩家信息，仍然添加但用户名为空
		username = ""
	}

	player := &NearbyPlayer{
		Entity:   e,
		Username: username,
	}

	nt.players[e.ID] = player

	if nt.callbacks.OnPlayerEnter != nil {
		go nt.callbacks.OnPlayerEnter(player)
	}
}

// handleEntityMove 处理实体移动
func (nt *NearbyTracker) handleEntityMove(e *entity.Entity, oldPos entity.Position) {
	if !e.IsPlayer() {
		return
	}

	nt.mu.RLock()
	player, exists := nt.players[e.ID]
	callbacks := nt.callbacks
	nt.mu.RUnlock()

	if !exists {
		return
	}

	if callbacks.OnPlayerMove != nil {
		go callbacks.OnPlayerMove(player, oldPos)
	}
}

// handleEntityRemove 处理实体移除
func (nt *NearbyTracker) handleEntityRemove(e *entity.Entity) {
	if !e.IsPlayer() {
		return
	}

	nt.mu.Lock()
	player, exists := nt.players[e.ID]
	if !exists {
		nt.mu.Unlock()
		return
	}

	delete(nt.players, e.ID)
	callbacks := nt.callbacks
	nt.mu.Unlock()

	if callbacks.OnPlayerLeave != nil {
		go callbacks.OnPlayerLeave(player)
	}
}

// GetNearbyPlayers 获取周围所有玩家
func (nt *NearbyTracker) GetNearbyPlayers() []*NearbyPlayer {
	nt.mu.RLock()
	defer nt.mu.RUnlock()

	result := make([]*NearbyPlayer, 0, len(nt.players))
	for _, p := range nt.players {
		result = append(result, p)
	}
	return result
}

// GetNearbyPlayer 通过UUID查找特定玩家
func (nt *NearbyTracker) GetNearbyPlayer(uuid [16]byte) (*NearbyPlayer, bool) {
	nt.mu.RLock()
	defer nt.mu.RUnlock()

	for _, p := range nt.players {
		if p.UUID == uuid {
			return p, true
		}
	}
	return nil, false
}

// PlayersWithinDistance 获取指定距离内的玩家
func (nt *NearbyTracker) PlayersWithinDistance(center entity.Position, distance float64) []*NearbyPlayer {
	nt.mu.RLock()
	defer nt.mu.RUnlock()

	distanceSq := distance * distance
	result := make([]*NearbyPlayer, 0)

	for _, p := range nt.players {
		if p.DistanceTo(center) <= distanceSq {
			result = append(result, p)
		}
	}

	return result
}

// Count 返回周围玩家数量
func (nt *NearbyTracker) Count() int {
	nt.mu.RLock()
	defer nt.mu.RUnlock()
	return len(nt.players)
}

// Clear 清空所有玩家数据
func (nt *NearbyTracker) Clear() {
	nt.mu.Lock()
	defer nt.mu.Unlock()
	nt.players = make(map[int32]*NearbyPlayer)
}
