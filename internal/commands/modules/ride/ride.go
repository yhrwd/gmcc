package ride

import (
	"fmt"
	"sync"
	"time"

	"gmcc/internal/commands"
)

type RideCommand struct {
	mu     sync.RWMutex
	bot    commands.BotAdapter
	config *Config

	state     commands.StateType
	target    string
	startTime time.Time

	currentYaw   float32
	currentPitch float32
}

func NewRideCommand(cfg *Config) *RideCommand {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	return &RideCommand{
		config: cfg,
		state:  commands.StateIdle,
	}
}

func (r *RideCommand) Name() string {
	return "ride"
}

func (r *RideCommand) Description() string {
	return "骑乘指定玩家，自动追踪视角并执行骑乘命令"
}

func (r *RideCommand) Usage() string {
	return "ride [玩家名]"
}

func (r *RideCommand) Init(bot commands.BotAdapter, _ *commands.ModuleConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.bot = bot
	return nil
}

func (r *RideCommand) Execute(ctx *commands.ChatContext) *commands.CommandResult {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != commands.StateIdle {
		return &commands.CommandResult{
			Success:   false,
			Message:   "指令执行中，请稍候",
			NextState: r.state,
		}
	}

	r.target = ctx.Sender
	if len(ctx.Args) > 0 && ctx.Args[0] != "" {
		r.target = ctx.Args[0]
	}

	player, ok := r.bot.GetPlayerByName(r.target)
	if !ok {
		r.state = commands.StatePreparing
		r.startTime = time.Now()
		return &commands.CommandResult{
			Success:   true,
			Message:   fmt.Sprintf("正在查找玩家 %s，请稍候...", r.target),
			NextState: commands.StatePreparing,
		}
	}

	dist := r.bot.DistanceTo(player.Position.X, player.Position.Y, player.Position.Z)

	if dist <= r.config.RangeLimit {
		return r.executeRide(ctx)
	}

	r.state = commands.StateExecuting
	r.startTime = time.Now()
	r.currentYaw, r.currentPitch = r.bot.GetRotation()
	r.updateLookAt(player.Position.X, player.Position.Y, player.Position.Z)

	return &commands.CommandResult{
		Success:   true,
		Message:   fmt.Sprintf("目标距离 %.1f 格，请靠近后自动骑乘...", dist),
		NextState: commands.StateExecuting,
	}
}

func (r *RideCommand) Tick(ctx *commands.ChatContext) *commands.CommandResult {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch r.state {
	case commands.StatePreparing:
		return r.tickPreparing(ctx)

	case commands.StateExecuting:
		return r.tickExecuting(ctx)

	case commands.StateCooldown:
		return r.tickCooldown()

	case commands.StateFailed:
		r.state = commands.StateIdle
		r.target = ""
		return nil
	}

	return nil
}

func (r *RideCommand) tickPreparing(ctx *commands.ChatContext) *commands.CommandResult {
	player, ok := r.bot.GetPlayerByName(r.target)
	if !ok {
		if time.Since(r.startTime) > r.config.Timeout {
			r.state = commands.StateFailed
			return &commands.CommandResult{
				Success:   false,
				Message:   fmt.Sprintf("未找到玩家 %s", r.target),
				NextState: commands.StateFailed,
			}
		}
		return nil
	}

	dist := r.bot.DistanceTo(player.Position.X, player.Position.Y, player.Position.Z)
	r.state = commands.StateExecuting
	r.startTime = time.Now()
	r.currentYaw, r.currentPitch = r.bot.GetRotation()
	r.updateLookAt(player.Position.X, player.Position.Y, player.Position.Z)

	return &commands.CommandResult{
		Success:   true,
		Message:   fmt.Sprintf("已锁定 %s，距离 %.1f 格，请靠近...", r.target, dist),
		NextState: commands.StateExecuting,
	}
}

func (r *RideCommand) tickExecuting(ctx *commands.ChatContext) *commands.CommandResult {
	player, ok := r.bot.GetPlayerByName(r.target)
	if !ok {
		r.state = commands.StateFailed
		return &commands.CommandResult{
			Success:   false,
			Message:   fmt.Sprintf("目标玩家 %s 已离线", r.target),
			NextState: commands.StateFailed,
		}
	}

	if time.Since(r.startTime) > r.config.Timeout {
		r.state = commands.StateFailed
		return &commands.CommandResult{
			Success:   false,
			Message:   fmt.Sprintf("等待超时（%v）", r.config.Timeout),
			NextState: commands.StateFailed,
		}
	}

	r.smoothLookAt(player.Position.X, player.Position.Y, player.Position.Z)

	dist := r.bot.DistanceTo(player.Position.X, player.Position.Y, player.Position.Z)
	if dist <= r.config.RangeLimit {
		return r.executeRide(ctx)
	}

	return nil
}

func (r *RideCommand) tickCooldown() *commands.CommandResult {
	if time.Since(r.startTime) > r.config.Cooldown {
		r.state = commands.StateIdle
		r.target = ""
	}
	return nil
}

func (r *RideCommand) executeRide(ctx *commands.ChatContext) *commands.CommandResult {
	// 查找目标玩家的实体ID
	player, ok := r.bot.GetPlayerByName(r.target)
	if !ok {
		r.state = commands.StateFailed
		return &commands.CommandResult{
			Success:   false,
			Message:   fmt.Sprintf("找不到玩家 %s", r.target),
			NextState: commands.StateFailed,
		}
	}

	// 先切换快捷栏到槽位 0（确保手上物品正确）
	if err := r.bot.SetHeldSlot(0); err != nil {
		r.state = commands.StateFailed
		return &commands.CommandResult{
			Success:   false,
			Message:   fmt.Sprintf("切换快捷栏失败: %v", err),
			NextState: commands.StateFailed,
		}
	}

	// 右键点击玩家实体（骑乘）
	if err := r.bot.InteractEntity(player.EntityID); err != nil {
		r.state = commands.StateFailed
		return &commands.CommandResult{
			Success:   false,
			Message:   fmt.Sprintf("交互失败: %v", err),
			NextState: commands.StateFailed,
		}
	}

	r.state = commands.StateCooldown
	r.startTime = time.Now()

	return &commands.CommandResult{
		Success:   true,
		Message:   fmt.Sprintf("已向 %s 发送骑乘请求", r.target),
		NextState: commands.StateCooldown,
	}
}

func (r *RideCommand) updateLookAt(tx, ty, tz float64) {
	bx, by, bz := r.bot.GetPosition()
	targetYaw, targetPitch := calculateLookAt(bx, by, bz, tx, ty, tz)
	r.currentYaw = targetYaw
	r.currentPitch = targetPitch
	r.bot.SetYawPitch(targetYaw, targetPitch)
}

func (r *RideCommand) smoothLookAt(tx, ty, tz float64) {
	bx, by, bz := r.bot.GetPosition()
	targetYaw, targetPitch := calculateLookAt(bx, by, bz, tx, ty, tz)

	newYaw, newPitch := smoothLook(r.currentYaw, r.currentPitch, targetYaw, targetPitch, r.config.LookSmoothing)
	r.currentYaw = newYaw
	r.currentPitch = newPitch

	r.bot.SetYawPitch(newYaw, newPitch)
}

func (r *RideCommand) Cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.state = commands.StateIdle
	r.target = ""
}

func (r *RideCommand) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.state = commands.StateIdle
	r.target = ""
	r.startTime = time.Time{}
}

func (r *RideCommand) State() commands.StateType {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.state
}

func (r *RideCommand) Target() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.target
}
