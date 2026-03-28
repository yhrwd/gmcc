package commands

import (
	"gmcc/internal/config"
)

type ModuleConfig struct {
	BotName string
	Prefix  string
	Auth    *AuthManager
}

func NewModuleConfig(cfg *config.Config, botName string) *ModuleConfig {
	mc := &ModuleConfig{
		BotName: botName,
		Prefix:  cfg.Commands.Prefix,
	}

	auth := NewAuthManager()
	auth.SetAllowAll(cfg.Commands.AllowAll)
	auth.SetWhitelist(cfg.Commands.Whitelist)
	mc.Auth = auth

	return mc
}
