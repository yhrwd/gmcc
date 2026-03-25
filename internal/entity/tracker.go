// Package entity 提供Minecraft实体跟踪功能
package entity

import (
	"sync"
	"time"
)

// Callbacks 定义实体事件回调
type Callbacks struct {
	OnSpawn  func(e *Entity)
	OnMove   func(e *Entity, oldPos Position)
	OnRemove func(e *Entity)
}

// Tracker 跟踪所有实体的状态和位置
type Tracker struct {
	mu             sync.RWMutex
	entities       map[int32]*Entity
	byUUID         map[[16]byte]*Entity
	pendingUpdates map[int32]*pendingUpdate
	callbacks      Callbacks
}

// pendingUpdate 存储待处理的实体更新
type pendingUpdate struct {
	entityID int32
	newPos   Position
	timer    *time.Timer
}

// NewTracker 创建新的实体跟踪器
func NewTracker() *Tracker {
	return &Tracker{
		entities:       make(map[int32]*Entity),
		byUUID:         make(map[[16]byte]*Entity),
		pendingUpdates: make(map[int32]*pendingUpdate),
	}
}

// SetCallbacks 设置事件回调
func (t *Tracker) SetCallbacks(callbacks Callbacks) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.callbacks = callbacks
}

// SpawnEntity 添加新实体
func (t *Tracker) SpawnEntity(id int32, entityType string, uuid [16]byte, pos Position, velocity Vector3) *Entity {
	t.mu.Lock()
	defer t.mu.Unlock()

	entity := &Entity{
		ID:         id,
		Type:       entityType,
		UUID:       uuid,
		Position:   pos,
		Velocity:   velocity,
		LastUpdate: time.Now(),
	}

	t.entities[id] = entity
	t.byUUID[uuid] = entity

	if t.callbacks.OnSpawn != nil {
		// 复制数据后回调
		entityCopy := *entity
		go t.callbacks.OnSpawn(&entityCopy)
	}

	return entity
}

// UpdatePosition 更新实体位置（完整位置）
func (t *Tracker) UpdatePosition(id int32, newPos Position) {
	t.mu.Lock()
	defer t.mu.Unlock()

	_, exists := t.entities[id]
	if !exists {
		return
	}

	// 检查是否已有待处理的更新
	if pending, ok := t.pendingUpdates[id]; ok {
		pending.newPos = newPos
		return
	}

	// 创建延迟回调
	timer := time.AfterFunc(100*time.Millisecond, func() {
		t.executePositionUpdate(id, newPos)
	})

	t.pendingUpdates[id] = &pendingUpdate{
		entityID: id,
		newPos:   newPos,
		timer:    timer,
	}
}

// UpdatePositionDelta 更新实体位置（增量）
func (t *Tracker) UpdatePositionDelta(id int32, deltaX, deltaY, deltaZ int16) {
	t.mu.Lock()
	defer t.mu.Unlock()

	e, exists := t.entities[id]
	if !exists {
		return
	}

	// Delta 值需要除以 4096 转换为实际坐标
	newPos := Position{
		X: e.Position.X + float64(deltaX)/4096.0,
		Y: e.Position.Y + float64(deltaY)/4096.0,
		Z: e.Position.Z + float64(deltaZ)/4096.0,
	}

	// 检查是否已有待处理的更新
	if pending, ok := t.pendingUpdates[id]; ok {
		pending.newPos = newPos
		return
	}

	// 创建延迟回调
	timer := time.AfterFunc(100*time.Millisecond, func() {
		t.executePositionUpdate(id, newPos)
	})

	t.pendingUpdates[id] = &pendingUpdate{
		entityID: id,
		newPos:   newPos,
		timer:    timer,
	}
}

// executePositionUpdate 执行位置更新（内部方法，需要在外部加锁）
func (t *Tracker) executePositionUpdate(id int32, newPos Position) {
	t.mu.Lock()
	e, exists := t.entities[id]
	if !exists {
		// 实体已被移除，清理待处理更新
		delete(t.pendingUpdates, id)
		t.mu.Unlock()
		return
	}

	oldPos := e.Position
	e.Position = newPos
	e.LastUpdate = time.Now()
	delete(t.pendingUpdates, id)

	// 复制数据用于回调
	entityCopy := *e
	callbacks := t.callbacks
	t.mu.Unlock()

	if callbacks.OnMove != nil {
		callbacks.OnMove(&entityCopy, oldPos)
	}
}

// RemoveEntity 移除实体
func (t *Tracker) RemoveEntity(id int32) {
	t.mu.Lock()
	defer t.mu.Unlock()

	e, exists := t.entities[id]
	if !exists {
		return
	}

	// 清理待处理的更新
	if pending, ok := t.pendingUpdates[id]; ok {
		pending.timer.Stop()
		delete(t.pendingUpdates, id)
	}

	delete(t.entities, id)
	delete(t.byUUID, e.UUID)

	if t.callbacks.OnRemove != nil {
		// 复制数据后回调
		entityCopy := *e
		go t.callbacks.OnRemove(&entityCopy)
	}
}

// RemoveEntities 批量移除实体
func (t *Tracker) RemoveEntities(ids []int32) {
	for _, id := range ids {
		t.RemoveEntity(id)
	}
}

// Get 通过ID获取实体
func (t *Tracker) Get(id int32) (*Entity, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entities[id]
	return e, ok
}

// GetByUUID 通过UUID获取实体
func (t *Tracker) GetByUUID(uuid [16]byte) (*Entity, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.byUUID[uuid]
	return e, ok
}

// All 获取所有实体
func (t *Tracker) All() []*Entity {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make([]*Entity, 0, len(t.entities))
	for _, e := range t.entities {
		result = append(result, e)
	}
	return result
}

// ByType 按类型筛选实体
func (t *Tracker) ByType(entityType string) []*Entity {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make([]*Entity, 0)
	for _, e := range t.entities {
		if e.Type == entityType {
			result = append(result, e)
		}
	}
	return result
}

// Count 返回实体数量
func (t *Tracker) Count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.entities)
}

// Stop 停止跟踪器并清理资源
func (t *Tracker) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 停止所有待处理的定时器
	for _, pending := range t.pendingUpdates {
		pending.timer.Stop()
	}
	t.pendingUpdates = make(map[int32]*pendingUpdate)
	t.entities = make(map[int32]*Entity)
	t.byUUID = make(map[[16]byte]*Entity)
}
