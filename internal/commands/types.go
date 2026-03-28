package commands

import (
	"time"
)

type BotAdapter interface {
	GetPlayerID() string
	GetUUID() string
	GetPosition() (x, y, z float64)
	GetRotation() (yaw, pitch float32)
	SendChat(msg string) error
	SendCommand(cmd string) error
	SendPrivateMessage(target, msg string) error
	SetYawPitch(yaw, pitch float32) error
	LookAt(x, y, z float64) error
	GetNearbyPlayers() []PlayerInfo
	GetPlayerByName(name string) (PlayerInfo, bool)
	DistanceTo(x, y, z float64) float64
	IsOnline() bool
	// 实体交互
	SetHeldSlot(slot int16) error        // 切换快捷栏槽位 (0-8)
	InteractEntity(entityID int32) error // 右键点击实体
}

type Message struct {
	Type       string
	PlainText  string
	RawJSON    string
	Sender     string
	SenderUUID string
	IsPrivate  bool
	Timestamp  time.Time
}

type PlayerInfo struct {
	Name     string
	UUID     string
	EntityID int32
	Position struct{ X, Y, Z float64 }
}

type ChatContext struct {
	Bot     BotAdapter
	Message Message
	Sender  string
	Args    []string
}

type CommandResult struct {
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

func (s StateType) String() string {
	switch s {
	case StateIdle:
		return "idle"
	case StatePreparing:
		return "preparing"
	case StateExecuting:
		return "executing"
	case StateCooldown:
		return "cooldown"
	case StateFailed:
		return "failed"
	default:
		return "unknown"
	}
}

type Command interface {
	Name() string
	Description() string
	Usage() string
	Execute(ctx *ChatContext) *CommandResult
	Tick(ctx *ChatContext) *CommandResult
	Cleanup()
	Stop()
	State() StateType
	Target() string
}
