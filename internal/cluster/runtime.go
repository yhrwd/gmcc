package cluster

import (
	"context"
	"errors"

	"gmcc/internal/config"
	"gmcc/internal/headless"
	"gmcc/internal/mcclient"
)

var (
	ErrRunnerNetworkDisconnect = errors.New("runner network disconnect")
	ErrRunnerAuthFailed        = errors.New("runner auth failed")
	ErrRunnerStartupTimeout    = errors.New("runner startup timeout")
)

type runner interface {
	Run(ctx context.Context) error
	IsReady() bool
	SendCommand(cmd string) error
	GetPlayer() *mcclient.Player
}

type runnerFactory func(cfg *config.Config) runner

func defaultRunnerFactory(cfg *config.Config) runner {
	return headless.New(cfg)
}

func classifyExitCategory(err error) ExitCategory {
	if err == nil {
		return ExitCategoryManualStop
	}

	switch {
	case errors.Is(err, ErrRunnerNetworkDisconnect):
		return ExitCategoryNetworkDisconnect
	case errors.Is(err, ErrRunnerAuthFailed):
		return ExitCategoryAuthFailed
	case errors.Is(err, ErrRunnerStartupTimeout):
		return ExitCategoryStartupTimeout
	default:
		return ExitCategoryUnknown
	}
}
