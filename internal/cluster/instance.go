package cluster

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gmcc/internal/config"
	"gmcc/internal/headless"
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

	// 运行时
	runner  *headless.Runner
	cancel  context.CancelFunc
	errChan chan error

	// 父级管理器
	manager *Manager
}

// newInstance 创建新实例（内部使用）
func newInstance(id string, account AccountEntry, manager *Manager) *Instance {
	return &Instance{
		ID:      id,
		Account: account,
		status:  StatusPending,
		errChan: make(chan error, 1),
		manager: manager,
	}
}

// Start 启动实例
func (i *Instance) Start() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	// 检查当前状态
	if i.status == StatusRunning || i.status == StatusStarting {
		return fmt.Errorf("instance %s is already running", i.ID)
	}

	if i.status == StatusReconnecting {
		return fmt.Errorf("instance %s is reconnecting", i.ID)
	}

	// 重置状态
	i.errorMsg = ""
	if i.status != StatusReconnecting {
		// 非重连情况下才重置重连计数
		i.reconnectCount = 0
	}

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
	}

	// 创建runner
	i.runner = headless.New(cfg)

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	i.cancel = cancel
	i.errChan = make(chan error, 1)

	// 启动goroutine运行实例
	i.status = StatusStarting
	i.startTime = time.Now()

	go func() {
		err := i.runner.Run(ctx)
		if err != nil {
			logx.Errorf("实例 %s 运行错误: %v", i.ID, err)
		}
		i.errChan <- err

		// 更新状态
		i.mu.Lock()
		if i.status != StatusStopped {
			if err != nil {
				i.status = StatusError
				i.errorMsg = err.Error()
			} else {
				i.status = StatusStopped
			}
		}
		i.mu.Unlock()

		// 通知管理器
		if i.manager != nil {
			i.manager.handleInstanceStopped(i.ID, err)
		}
	}()

	// 等待实例就绪（异步）
	go i.waitForReady()

	logx.Infof("实例 %s 已启动: player_id=%s, server=%s", i.ID, i.Account.PlayerID, i.Account.ServerAddress)
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

	i.status = StatusStopped
	logx.Infof("实例 %s 已停止", i.ID)
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
			i.status = StatusRunning
			i.reconnectCount = 0 // 成功后重置重连计数
			i.mu.Unlock()
			logx.Infof("实例 %s 已就绪", i.ID)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	// 超时
	i.mu.Lock()
	if i.status == StatusStarting || i.status == StatusReconnecting {
		i.status = StatusError
		i.errorMsg = "timeout waiting for ready"
	}
	i.mu.Unlock()
	logx.Warnf("实例 %s 启动超时", i.ID)
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
