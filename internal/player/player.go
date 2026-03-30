package player

import (
	"fmt"
	"sync"
	"time"
)

type GameMode int

type SlotData struct {
	ID    int32
	Count int32
}

type ContainerState struct {
	WindowID   int32
	WindowType int32
	StateID    int32
	Open       bool
}

func (s *SlotData) IDToString() string {
	if s == nil || s.ID == 0 {
		return ""
	}
	return fmt.Sprintf("minecraft:%d", s.ID)
}

const (
	GameModeSurvival GameMode = iota
	GameModeCreative
	GameModeAdventure
	GameModeSpectator
)

type Player struct {
	mu sync.RWMutex

	EntityID  int32
	UUID      [16]byte
	Name      string
	GameMode  GameMode
	Dimension string

	X, Y, Z    float64
	Yaw, Pitch float32
	OnGround   bool

	Health       float32
	MaxHealth    float32
	Food         int32
	Saturation   float32
	Air          int32
	EntityHealth float32
	Level        int32
	Experience   float32
	TotalExp     int32

	Invulnerable bool
	Flying       bool
	CanFly       bool
	InstantBreak bool
	FlyingSpeed  float32
	FieldOfView  float32

	Inventory     *Inventory
	HeldSlot      int8
	OpenContainer *ContainerState

	JoinTime   time.Time
	LastUpdate time.Time

	OnHealthChange    func(health, maxHealth float32, food int32)
	OnPositionChange  func(x, y, z float64)
	OnInventoryChange func(slot int8, item *Item)
	OnGameModeChange  func(mode GameMode)
}

func NewPlayer() *Player {
	return &Player{
		Inventory:    NewInventory(),
		Health:       20,
		MaxHealth:    20,
		Food:         20,
		Saturation:   5,
		Air:          300,
		EntityHealth: 20,
		JoinTime:     time.Now(),
		LastUpdate:   time.Now(),
		FlyingSpeed:  0.05,
		FieldOfView:  0.1,
		OnGround:     true,
	}
}

func (p *Player) SetEntityID(id int32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.EntityID = id
}

func (p *Player) SetUUID(uuid [16]byte) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.UUID = uuid
}

func (p *Player) SetName(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Name = name
}

func (p *Player) SetGameMode(mode GameMode) {
	p.mu.Lock()
	oldMode := p.GameMode
	p.GameMode = mode
	cb := p.OnGameModeChange
	p.mu.Unlock()

	if oldMode != mode && cb != nil {
		go cb(mode)
	}
}

func (p *Player) SetDimension(dim string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Dimension = dim
}

func (p *Player) UpdatePosition(x, y, z float64, yaw, pitch float32, relative int8) {
	p.mu.Lock()
	if relative&0x01 != 0 {
		x += p.X
	}
	if relative&0x02 != 0 {
		y += p.Y
	}
	if relative&0x04 != 0 {
		z += p.Z
	}
	if relative&0x08 != 0 {
		yaw += p.Yaw
	}
	if relative&0x10 != 0 {
		pitch += p.Pitch
	}
	p.X, p.Y, p.Z = x, y, z
	p.Yaw, p.Pitch = yaw, pitch
	p.LastUpdate = time.Now()
	cb := p.OnPositionChange
	p.mu.Unlock()

	if cb != nil {
		go cb(x, y, z)
	}
}

func (p *Player) GetPosition() (float64, float64, float64) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.X, p.Y, p.Z
}

func (p *Player) GetRotation() (float32, float32) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Yaw, p.Pitch
}

func (p *Player) GetMovementState() (float64, float64, float64, float32, float32, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.X, p.Y, p.Z, p.Yaw, p.Pitch, p.OnGround
}

func (p *Player) UpdateHealth(health, maxHealth float32, food int32, saturation float32) {
	p.mu.Lock()
	oldHealth := p.Health
	p.Health = health
	if maxHealth > 0 {
		p.MaxHealth = maxHealth
	}
	p.Food = food
	p.Saturation = saturation
	p.LastUpdate = time.Now()
	cb := p.OnHealthChange
	p.mu.Unlock()

	if oldHealth != health && cb != nil {
		go cb(health, maxHealth, food)
	}
}

func (p *Player) GetHealth() (float32, float32, int32, float32) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Health, p.MaxHealth, p.Food, p.Saturation
}

func (p *Player) UpdateAir(air int32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Air = air
}

func (p *Player) GetAir() int32 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Air
}

func (p *Player) UpdateEntityHealth(health float32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.EntityHealth = health
}

func (p *Player) GetEntityHealth() float32 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.EntityHealth
}

func (p *Player) UpdateExperience(level int32, experience, totalExp float32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Level = level
	p.Experience = experience
	p.TotalExp = int32(totalExp)
	p.LastUpdate = time.Now()
}

func (p *Player) GetExperience() (int32, float32, int32) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Level, p.Experience, p.TotalExp
}

func (p *Player) UpdateAbilities(flags int8, flyingSpeed, fov float32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Invulnerable = flags&0x01 != 0
	p.Flying = flags&0x02 != 0
	p.CanFly = flags&0x04 != 0
	p.InstantBreak = flags&0x08 != 0
	p.FlyingSpeed = flyingSpeed
	p.FieldOfView = fov
	p.LastUpdate = time.Now()
}

func (p *Player) GetAbilities() (bool, bool, bool, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Invulnerable, p.Flying, p.CanFly, p.InstantBreak
}

func (p *Player) SetHeldSlot(slot int8) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.HeldSlot = slot
}

func (p *Player) GetHeldSlot() int8 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.HeldSlot
}

func (p *Player) SetOpenContainer(container *ContainerState) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.OpenContainer = container
}

func (p *Player) GetOpenContainer() *ContainerState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.OpenContainer
}

func (p *Player) UpdateContainerStateID(stateID int32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.OpenContainer != nil {
		p.OpenContainer.StateID = stateID
	}
}

func (p *Player) GetHeldItem() *Item {
	p.mu.RLock()
	slot := p.HeldSlot
	p.mu.RUnlock()
	return p.Inventory.GetSlot(slot)
}

func (p *Player) UpdateInventorySlot(windowID int8, stateID int32, slot int8, item *Item) {
	if windowID != 0 {
		return
	}
	p.Inventory.SetSlot(slot, item)
	p.mu.RLock()
	cb := p.OnInventoryChange
	p.mu.RUnlock()
	if cb != nil {
		go cb(slot, item)
	}
}

func (p *Player) ClearInventory() {
	p.Inventory.Clear()
}

func (p *Player) UpdateInventory(windowID int32, items []*SlotData, carriedItem *SlotData) {
	if windowID != 0 {
		return
	}
	p.mu.Lock()
	p.Inventory.Clear()
	for i, item := range items {
		if item != nil && item.Count > 0 {
			p.Inventory.SetSlot(int8(i), &Item{
				ID:    item.IDToString(),
				Count: item.Count,
			})
		}
	}
	if carriedItem != nil && carriedItem.Count > 0 {
		p.Inventory.SetSlot(-1, &Item{
			ID:    carriedItem.IDToString(),
			Count: carriedItem.Count,
		})
	}
	p.mu.Unlock()
}

func (p *Player) UpdateSlot(windowID int32, slot int32, item *SlotData) {
	if windowID != 0 {
		return
	}
	if item != nil && item.Count > 0 {
		p.Inventory.SetSlot(int8(slot), &Item{
			ID:    item.IDToString(),
			Count: item.Count,
		})
	} else {
		p.Inventory.SetSlot(int8(slot), nil)
	}
}

func (p *Player) GetDuration() time.Duration {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return time.Since(p.JoinTime)
}

func (p *Player) GetInfo() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return map[string]interface{}{
		"name":          p.Name,
		"uuid":          formatUUID(p.UUID),
		"entity_id":     p.EntityID,
		"gamemode":      p.GameMode.String(),
		"dimension":     p.Dimension,
		"position":      []float64{p.X, p.Y, p.Z},
		"rotation":      []float32{p.Yaw, p.Pitch},
		"health":        p.Health,
		"max_health":    p.MaxHealth,
		"food":          p.Food,
		"saturation":    p.Saturation,
		"air":           p.Air,
		"entity_health": p.EntityHealth,
		"level":         p.Level,
		"experience":    p.Experience,
		"held_slot":     p.HeldSlot,
		"flying":        p.Flying,
		"can_fly":       p.CanFly,
		"join_time":     p.JoinTime.Format("2006-01-02 15:04:05"),
		"duration":      time.Since(p.JoinTime).String(),
	}
}

func (g GameMode) String() string {
	switch g {
	case GameModeSurvival:
		return "survival"
	case GameModeCreative:
		return "creative"
	case GameModeAdventure:
		return "adventure"
	case GameModeSpectator:
		return "spectator"
	default:
		return "unknown"
	}
}

func formatUUID(uuid [16]byte) string {
	hex := make([]byte, 32)
	const hexChars = "0123456789abcdef"
	for i := 0; i < 16; i++ {
		hex[i*2] = hexChars[uuid[i]>>4]
		hex[i*2+1] = hexChars[uuid[i]&0x0f]
	}
	return string(hex[0:8]) + "-" + string(hex[8:12]) + "-" + string(hex[12:16]) + "-" + string(hex[16:20]) + "-" + string(hex[20:32])
}
