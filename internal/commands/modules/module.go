package modules

import (
	"gmcc/internal/commands"
	"gmcc/internal/commands/modules/pos"
	"gmcc/internal/commands/modules/ride"
)

type ModuleConfig struct {
	Bot    commands.BotAdapter
	Config *commands.ModuleConfig
}

func NewRideCommand(cfg *ride.Config) *ride.RideCommand {
	return ride.NewRideCommand(cfg)
}

func NewPosCommand() *pos.PosCommand {
	return pos.NewPosCommand()
}
