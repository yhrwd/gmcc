package cluster

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"

	authsession "gmcc/internal/auth/session"
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

type runnerFactory func(cfg *config.Config, authManager *authsession.AuthManager) runner

func defaultRunnerFactory(cfg *config.Config, authManager *authsession.AuthManager) runner {
	if cfg != nil {
		cfg.ClusterRuntime.AuthManager = authManager
	}
	return headless.New(cfg)
}

func classifyExitCategory(err error) ExitCategory {
	if err == nil {
		return ExitCategoryManualStop
	}
	if errors.Is(err, context.Canceled) {
		return ExitCategoryManualStop
	}

	switch {
	case errors.Is(err, ErrRunnerNetworkDisconnect):
		return ExitCategoryNetworkDisconnect
	case errors.Is(err, ErrRunnerAuthFailed):
		return ExitCategoryAuthFailed
	case errors.Is(err, ErrRunnerStartupTimeout):
		return ExitCategoryStartupTimeout
	case isNetworkDisconnectError(err):
		return ExitCategoryNetworkDisconnect
	default:
		return ExitCategoryUnknown
	}
}

func isNetworkDisconnectError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, net.ErrClosed) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "eof") ||
		strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "closed network connection") ||
		strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "i/o timeout") ||
		strings.Contains(msg, "no connection could be made")
}
