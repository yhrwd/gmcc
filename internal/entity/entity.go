// Package entity 提供Minecraft实体跟踪功能
package entity

import (
	"time"
)

// Entity 表示Minecraft世界中的一个实体
type Entity struct {
	ID         int32
	Type       string   // "minecraft:player" 或其他实体类型
	UUID       [16]byte // 可选，不是所有实体都有
	Position   Position
	Velocity   Vector3
	OnGround   bool
	LastUpdate time.Time
}

// Position 表示三维空间位置
type Position struct {
	X, Y, Z float64
}

// Vector3 表示三维向量
type Vector3 struct {
	X, Y, Z float64
}

// IsPlayer 检查实体是否为玩家类型
func (e *Entity) IsPlayer() bool {
	return e.Type == "minecraft:player"
}

// DistanceTo 计算到另一个位置的距离
func (p Position) DistanceTo(other Position) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	dz := p.Z - other.Z
	return float64(dx*dx + dy*dy + dz*dz)
}

// DistanceTo 计算实体到另一个位置的距离
func (e *Entity) DistanceTo(other Position) float64 {
	return e.Position.DistanceTo(other)
}
