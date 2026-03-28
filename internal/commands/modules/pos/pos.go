package pos

import (
	"fmt"

	"gmcc/internal/commands"
)

type PosCommand struct {
	bot commands.BotAdapter
}

func NewPosCommand() *PosCommand {
	return &PosCommand{}
}

func (p *PosCommand) Name() string        { return "pos" }
func (p *PosCommand) Description() string { return "查询当前坐标位置" }
func (p *PosCommand) Usage() string       { return "pos" }

func (p *PosCommand) Init(bot commands.BotAdapter, _ *commands.ModuleConfig) error {
	p.bot = bot
	return nil
}

func (p *PosCommand) Execute(ctx *commands.ChatContext) *commands.CommandResult {
	x, y, z := p.bot.GetPosition()
	yaw, pitch := p.bot.GetRotation()

	return &commands.CommandResult{
		Success: true,
		Message: fmt.Sprintf("坐标: X=%.1f Y=%.1f Z=%.1f | 朝向: Yaw=%.1f Pitch=%.1f", x, y, z, yaw, pitch),
	}
}

func (p *PosCommand) Tick(_ *commands.ChatContext) *commands.CommandResult { return nil }
func (p *PosCommand) Cleanup()                                             {}
func (p *PosCommand) Stop()                                                {}
func (p *PosCommand) State() commands.StateType                            { return commands.StateIdle }
func (p *PosCommand) Target() string                                       { return "" }
