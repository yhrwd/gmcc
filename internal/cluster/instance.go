package cluster

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	authsession "gmcc/internal/auth/session"
	"gmcc/internal/config"
	"gmcc/internal/logx"
)

// InstanceStatus 实例状态
type InstanceStatus string

const (
	StatusPending      InstanceStatus = "pending"
	StatusStarting     InstanceStatus = "starting"
	StatusRunning      InstanceStatus = "running"
	StatusReconnecting InstanceStatus = "reconnecting"
	StatusStopped      InstanceStatus = "stopped"
	StatusError        InstanceStatus = "error"
)

// InstanceInfo 实例信息（可导出）
type InstanceInfo struct {
	ID             string         `json:"id"`
	PlayerID       string         `json:"player_id"`
	ServerAddress  string         `json:"server_address"`
	Status         InstanceStatus `json:"status"`
	OnlineDuration time.Duration  `json:"online_duration"`
	LastActive     time.Time      `json:"last_active"`
	ReconnectCount int            `json:"reconnect_count"`
	Error          string         `json:"error,omitempty"`
	Health         float32        `json:"health,omitempty"`
	Food           int32          `json:"food,omitempty"`
	Position       *Position      `json:"position,omitempty"`
}

// Position 位置信息
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Instance 实例（单个Minecraft客户端）
type Instance struct {
	// 元数据
	ID      string
	Account AccountEntry

	// 状态
	mu             sync.RWMutex
	status         InstanceStatus
	startTime      time.Time
	lastActive     time.Time
	reconnectCount int
	errorMsg       string
	version        uint64
	runVersion     uint64
	deleted        bool
	exitCh         chan struct{}

	// 运行时
	runner  runner
	cancel  context.CancelFunc
	errChan chan error

	// 父级管理器
	manager *Manager

	startRunnerFn func(runVersion uint64) error
	runnerFactory runnerFactory
	authManager   *authsession.AuthManager
}

// newInstance 创建新实例（内部使用）
func newInstance(id string, account AccountEntry, manager *Manager) *Instance {
	return &Instance{
		ID:            id,
		Account:       account,
		status:        StatusPending,
		errChan:       make(chan error, 1),
		exitCh:        make(chan struct{}, 1),
		manager:       manager,
		runnerFactory: defaultRunnerFactory,
	}
}

// Start 启动实例
func (i *Instance) Start() error {
	return i.StartWithTrigger(StartTriggerManualStart)
}

// StartWithTrigger 按触发类型启动实例
func (i *Instance) StartWithTrigger(trigger StartTrigger) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	switch i.status {
	case StatusStarting, StatusRunning:
		return ErrInstanceRunningLike
	case StatusReconnecting:
		if trigger != StartTriggerAutoReconnect {
			return ErrInstanceRunningLike
		}
	}

	i.errorMsg = ""
	if trigger != StartTriggerAutoReconnect {
		i.resetReconnectAttemptsLocked()
	}
	i.version++
	i.runVersion = i.version
	i.exitCh = make(chan struct{}, 1)

	if err := i.transitionToLocked(StatusStarting, ""); err != nil {
		return err
	}

	if i.startRunnerFn != nil {
		return i.startRunnerFn(i.runVersion)
	}

	if err := i.startRunnerLocked(trigger, i.runVersion); err != nil {
		_ = i.transitionToLocked(StatusError, err.Error())
		i.emitLifecycleEvent("error", "error", "instance start failed", err.Error(), map[string]any{
			"status":  string(StatusError),
			"trigger": string(trigger),
		})
		return err
	}

	return nil
}

func (i *Instance) startRunnerLocked(trigger StartTrigger, runVersion uint64) error {

	// 创建配置
	cfg := &config.Config{
		Account: config.AccountConfig{
			PlayerID:        i.Account.PlayerID,
			UseOfficialAuth: i.Account.UseOfficialAuth,
		},
		Server: config.ServerConfig{
			Address: i.Account.ServerAddress,
		},
		Log: config.LogConfig{
			LogDir:     "logs",
			MaxSize:    512,
			Debug:      false,
			EnableFile: true,
		},
		ClusterRuntime: config.ClusterRuntimeConfig{
			AccountID:   i.Account.ID,
			AuthManager: i.authManager,
		},
	}

	// 创建runner
	factory := i.runnerFactory
	if factory == nil {
		factory = defaultRunnerFactory
	}
	i.runner = factory(cfg, i.authManager)

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	i.cancel = cancel
	i.errChan = make(chan error, 1)

	// 启动goroutine运行实例
	i.startTime = time.Now()

	go func() {
		err := i.runner.Run(ctx)
		i.errChan <- err

		category := classifyExitCategory(err)
		if errors.Is(err, authsession.ErrDeviceLoginRequired) ||
			errors.Is(err, authsession.ErrRefreshTokenInvalid) ||
			errors.Is(err, authsession.ErrProviderUnavailable) ||
			errors.Is(err, authsession.ErrRefreshUpstream) ||
			errors.Is(err, authsession.ErrOwnershipFailed) ||
			errors.Is(err, authsession.ErrProfileInvalid) ||
			errors.Is(err, authsession.ErrXSTSDenied) {
			category = ExitCategoryAuthFailed
		}
		handled := i.applyExitEvent(runVersion, category, err)
		if handled {
			i.signalExit()
		}

		// 通知管理器
		if handled && i.manager != nil {
			i.manager.handleInstanceStopped(i.ID, err)
		}
	}()

	// 等待实例就绪（异步）
	go i.waitForReady()

	i.emitLifecycleEvent("info", "start", "instance started", "", map[string]any{
		"status":         string(StatusStarting),
		"trigger":        string(trigger),
		"server_address": i.Account.ServerAddress,
	})
	return nil
}

// Stop 停止实例
func (i *Instance) Stop() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.status == StatusStopped {
		return nil // 已经停止
	}

	if i.cancel != nil {
		i.cancel()
		i.cancel = nil
	}

	_ = i.transitionToLocked(StatusStopped, "")
	i.emitLifecycleEvent("info", "stop", "instance stopped", "", map[string]any{
		"status": string(StatusStopped),
	})
	return nil
}

// Restart 重启实例
func (i *Instance) Restart() error {
	if err := i.Stop(); err != nil {
		return err
	}

	// 等待实例完全停止
	time.Sleep(500 * time.Millisecond)

	return i.Start()
}

// GetStatus 获取实例状态
func (i *Instance) GetStatus() InstanceStatus {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.status
}

// GetInfo 获取实例信息
func (i *Instance) GetInfo() InstanceInfo {
	i.mu.RLock()
	defer i.mu.RUnlock()

	info := InstanceInfo{
		ID:             i.ID,
		PlayerID:       i.Account.PlayerID,
		ServerAddress:  i.Account.ServerAddress,
		Status:         i.status,
		ReconnectCount: i.reconnectCount,
		Error:          i.errorMsg,
	}

	// 计算在线时长
	if i.status == StatusRunning && !i.startTime.IsZero() {
		info.OnlineDuration = time.Since(i.startTime)
	}

	info.LastActive = i.lastActive

	// 获取玩家状态
	if i.runner != nil && i.runner.IsReady() {
		if player := i.runner.GetPlayer(); player != nil {
			health, _, food, _ := player.GetHealth()
			info.Health = health
			info.Food = food
			x, y, z := player.GetPosition()
			info.Position = &Position{X: x, Y: y, Z: z}
		}
	}

	return info
}

// IsReady 检查实例是否就绪
func (i *Instance) IsReady() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.status == StatusRunning && i.runner != nil && i.runner.IsReady()
}

// SendCommand 发送命令
func (i *Instance) SendCommand(cmd string) error {
	if !i.IsReady() {
		return fmt.Errorf("instance %s is not ready", i.ID)
	}

	return i.runner.SendCommand(cmd)
}

// waitForReady 等待实例就绪
func (i *Instance) waitForReady() {
	// 检查是否已经在等待中
	i.mu.Lock()
	if i.status != StatusStarting && i.status != StatusReconnecting {
		i.mu.Unlock()
		return
	}
	i.mu.Unlock()

	// 最多等待30秒
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		// 检查状态是否已改变（可能被停止）
		i.mu.RLock()
		currentStatus := i.status
		i.mu.RUnlock()

		if currentStatus == StatusStopped {
			return
		}

		if i.runner != nil && i.runner.IsReady() {
			i.mu.Lock()
			// 再次检查状态，防止重复设置
			if i.status == StatusRunning {
				i.mu.Unlock()
				return
			}
			if err := i.transitionToLocked(StatusRunning, ""); err != nil {
				i.mu.Unlock()
				return
			}
			i.resetReconnectAttemptsLocked() // 成功后重置重连计数
			i.emitLifecycleEvent("info", "ready", "instance ready", "", map[string]any{
				"status": string(StatusRunning),
			})
			i.mu.Unlock()
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	// 超时
	i.mu.Lock()
	if i.status == StatusStarting || i.status == StatusReconnecting {
		_ = i.transitionToLocked(StatusError, "timeout waiting for ready")
		i.emitLifecycleEvent("warn", "error", "instance startup timed out", "timeout waiting for ready", map[string]any{
			"status":        string(StatusError),
			"exit_category": string(ExitCategoryStartupTimeout),
		})
	}
	i.mu.Unlock()
}

func (i *Instance) transitionToLocked(next InstanceStatus, reason string) error {
	if !canTransition(i.status, next) {
		return ErrInvalidTransition
	}

	i.status = next
	i.lastActive = time.Now()
	if reason != "" {
		i.errorMsg = reason
	}

	return nil
}

func (i *Instance) applyExitEvent(runVersion uint64, category ExitCategory, err error) bool {
	i.mu.Lock()
	defer i.mu.Unlock()

	if runVersion != i.version || i.deleted {
		return false
	}

	if category == ExitCategoryNetworkDisconnect {
		i.emitLifecycleEvent("warn", "disconnect", "instance disconnected", string(category), map[string]any{
			"status":        string(StatusReconnecting),
			"exit_category": string(category),
		})
		_ = i.transitionToLocked(StatusReconnecting, "network disconnect")
		return true
	}

	next := StatusStopped
	reason := ""
	if category != ExitCategoryManualStop && err != nil {
		next = StatusError
		reason = err.Error()
	}
	_ = i.transitionToLocked(next, reason)
	if category == ExitCategoryManualStop && err != nil {
		return true
	}
	if next == StatusError {
		i.emitLifecycleEvent("error", "error", "instance exited with error", reason, map[string]any{
			"status":        string(next),
			"exit_category": string(category),
		})
	} else {
		i.emitLifecycleEvent("info", "stop", "instance stopped", "", map[string]any{
			"status":        string(next),
			"exit_category": string(category),
		})
	}

	return true
}

func (i *Instance) emitLifecycleEvent(level, action, message, reason string, fields map[string]any) {
	event := logx.NewLifecycleEvent(level, action, message, i.ID, i.Account.ID)
	event.PlayerID = i.Account.PlayerID
	if strings.TrimSpace(reason) != "" {
		event.Reason = reason
	}
	event.Fields = fields
	logx.Emit(event)
}
func (i *Instance) signalExit() {
	select {
	case i.exitCh <- struct{}{}:
	default:
	}
}

func (i *Instance) waitExit(ctx context.Context) error {
	select {
	case <-i.exitCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (i *Instance) markDeleted() {
	i.mu.Lock()
	i.deleted = true
	i.mu.Unlock()
}

func (i *Instance) markError(msg string) {
	i.mu.Lock()
	i.errorMsg = msg
	i.lastActive = time.Now()
	i.mu.Unlock()
}

func (i *Instance) bumpReconnectAttempts() int {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.bumpReconnectAttemptsLocked()
}

func (i *Instance) bumpReconnectAttemptsLocked() int {
	i.reconnectCount++
	return i.reconnectCount
}

func (i *Instance) resetReconnectAttempts() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.resetReconnectAttemptsLocked()
}

func (i *Instance) resetReconnectAttemptsLocked() {
	i.reconnectCount = 0
}

// updateLastActive 更新最后活动时间
func (i *Instance) updateLastActive() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.lastActive = time.Now()
}

// String 返回实例字符串表示
func (s InstanceStatus) String() string {
	return string(s)
}
