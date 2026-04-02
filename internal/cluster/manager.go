package cluster

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"gmcc/internal/logx"
)

var (
	ErrInstanceNotFound = errors.New("instance not found")
	ErrInstanceRunning  = errors.New("instance already running")
	ErrMaxInstances     = errors.New("max instances reached")
)

// Manager 集群管理器
type Manager struct {
	config     ClusterConfig
	configPath string // 配置文件路径

	// 实例管理
	mu        sync.RWMutex
	instances map[string]*Instance
	startTime time.Time

	// 上下文
	ctx    context.Context
	cancel context.CancelFunc
}

// NewManager 创建集群管理器
func NewManager(config ClusterConfig, configPath ...string) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	m := &Manager{
		config:    config,
		instances: make(map[string]*Instance),
		startTime: time.Now(),
		ctx:       ctx,
		cancel:    cancel,
	}

	// 如果传入了配置路径，保存它
	if len(configPath) > 0 && configPath[0] != "" {
		m.configPath = configPath[0]
	}

	return m
}

// Start 启动集群管理器
func (m *Manager) Start() error {
	logx.Infof("集群管理器启动: max_instances=%d", m.config.Global.MaxInstances)

	// 自动启动所有启用的账号
	for _, account := range m.config.Accounts {
		if !account.Enabled {
			continue
		}

		if err := m.CreateInstance(account.ID, account); err != nil {
			logx.Warnf("创建实例失败 %s: %v", account.ID, err)
			continue
		}

		if err := m.StartInstance(account.ID); err != nil {
			logx.Warnf("启动实例失败 %s: %v", account.ID, err)
		}
	}

	return nil
}

// Stop 停止集群管理器
func (m *Manager) Stop() error {
	logx.Infof("集群管理器停止中...")

	// 取消上下文
	if m.cancel != nil {
		m.cancel()
	}

	// 停止所有实例
	m.mu.Lock()
	instances := make([]*Instance, 0, len(m.instances))
	for _, inst := range m.instances {
		instances = append(instances, inst)
	}
	m.mu.Unlock()

	// 并行停止所有实例
	var wg sync.WaitGroup
	for _, inst := range instances {
		wg.Add(1)
		go func(i *Instance) {
			defer wg.Done()
			if err := i.Stop(); err != nil {
				logx.Warnf("停止实例失败 %s: %v", i.ID, err)
			}
		}(inst)
	}
	wg.Wait()

	logx.Infof("集群管理器已停止")
	return nil
}

// CreateInstance 创建实例
func (m *Manager) CreateInstance(id string, account AccountEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查是否已存在
	if _, exists := m.instances[id]; exists {
		return fmt.Errorf("instance %s already exists", id)
	}

	// 检查最大实例数
	if m.config.Global.MaxInstances > 0 && len(m.instances) >= m.config.Global.MaxInstances {
		return ErrMaxInstances
	}

	// 创建实例
	inst := newInstance(id, account, m)
	m.instances[id] = inst

	logx.Infof("实例已创建: %s (player_id=%s)", id, account.PlayerID)
	return nil
}

// DeleteInstance 删除实例
func (m *Manager) DeleteInstance(id string) error {
	m.mu.Lock()
	inst, exists := m.instances[id]
	m.mu.Unlock()

	if !exists {
		return ErrInstanceNotFound
	}

	// 先停止实例
	if inst.GetStatus() == StatusRunning || inst.GetStatus() == StatusStarting {
		if err := inst.Stop(); err != nil {
			return fmt.Errorf("failed to stop instance: %w", err)
		}
		// 等待停止
		time.Sleep(500 * time.Millisecond)
	}

	m.mu.Lock()
	delete(m.instances, id)
	m.mu.Unlock()

	logx.Infof("实例已删除: %s", id)
	return nil
}

// StartInstance 启动实例
func (m *Manager) StartInstance(id string) error {
	inst, err := m.GetInstance(id)
	if err != nil {
		return err
	}

	return inst.Start()
}

// StopInstance 停止实例
func (m *Manager) StopInstance(id string) error {
	inst, err := m.GetInstance(id)
	if err != nil {
		return err
	}

	return inst.Stop()
}

// RestartInstance 重启实例
func (m *Manager) RestartInstance(id string) error {
	inst, err := m.GetInstance(id)
	if err != nil {
		return err
	}

	return inst.Restart()
}

// GetInstance 获取实例
func (m *Manager) GetInstance(id string) (*Instance, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	inst, exists := m.instances[id]
	if !exists {
		return nil, ErrInstanceNotFound
	}

	return inst, nil
}

// ListInstances 列出所有实例
func (m *Manager) ListInstances() []InstanceInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	infos := make([]InstanceInfo, 0, len(m.instances))
	for _, inst := range m.instances {
		infos = append(infos, inst.GetInfo())
	}

	return infos
}

// GetClusterStatus 获取集群状态
func (m *Manager) GetClusterStatus() ClusterStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	total := len(m.instances)
	running := 0

	for _, inst := range m.instances {
		if inst.GetStatus() == StatusRunning {
			running++
		}
	}

	status := "running"
	if running == 0 && total > 0 {
		status = "stopped"
	} else if running < total {
		status = "partial"
	}

	return ClusterStatus{
		Status:           status,
		TotalInstances:   total,
		RunningInstances: running,
		Uptime:           time.Since(m.startTime),
	}
}

// GetInstanceInfo 获取实例信息
func (m *Manager) GetInstanceInfo(id string) (InstanceInfo, error) {
	inst, err := m.GetInstance(id)
	if err != nil {
		return InstanceInfo{}, err
	}

	return inst.GetInfo(), nil
}

// InstanceExists 检查实例是否存在
func (m *Manager) InstanceExists(id string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.instances[id]
	return exists
}

// handleInstanceStopped 处理实例停止事件（由实例调用）
func (m *Manager) handleInstanceStopped(id string, err error) {
	logx.Debugf("实例已停止: %s, err=%v", id, err)

	// 这里可以触发重连逻辑
	inst, err := m.GetInstance(id)
	if err != nil {
		return
	}

	// 检查是否需要自动重连
	if m.config.Global.ReconnectPolicy.Enabled {
		go m.handleReconnect(inst)
	}
}

// handleReconnect 处理自动重连
func (m *Manager) handleReconnect(inst *Instance) {
	policy := m.config.Global.ReconnectPolicy

	// 增加重连计数
	inst.mu.Lock()
	inst.reconnectCount++
	attempt := inst.reconnectCount - 1 // 本次尝试的索引（从0开始）
	inst.status = StatusReconnecting
	inst.mu.Unlock()

	for attempt < policy.MaxRetries || policy.MaxRetries == 0 {
		// 计算退避延迟
		delay := policy.BaseDelay
		for i := 0; i < attempt; i++ {
			delay = time.Duration(float64(delay) * policy.Multiplier)
			if delay > policy.MaxDelay {
				delay = policy.MaxDelay
				break
			}
		}

		logx.Infof("实例 %s 将在 %v 后重连 (尝试 %d/%d)",
			inst.ID, delay, attempt+1, policy.MaxRetries)

		select {
		case <-m.ctx.Done():
			return
		case <-time.After(delay):
		}

		// 尝试启动
		if err := inst.Start(); err == nil {
			logx.Infof("实例 %s 重连成功", inst.ID)
			return
		} else {
			logx.Warnf("实例 %s 重连失败: %v", inst.ID, err)
			// 增加重连计数
			inst.mu.Lock()
			inst.reconnectCount++
			attempt = inst.reconnectCount - 1
			inst.mu.Unlock()
		}
	}

	logx.Errorf("实例 %s 重连次数用尽", inst.ID)
}

// GetConfig 获取配置
func (m *Manager) GetConfig() ClusterConfig {
	return m.config
}

// UpdateConfig 更新配置（部分更新）
func (m *Manager) UpdateConfig(config ClusterConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config = config
}

// AddAccount 添加账号到集群配置
func (m *Manager) AddAccount(account AccountEntry) error {
	// 检查ID是否已存在
	for _, acc := range m.config.Accounts {
		if acc.ID == account.ID {
			return fmt.Errorf("account with ID '%s' already exists", account.ID)
		}
	}

	// 添加账号
	m.config.Accounts = append(m.config.Accounts, account)
	logx.Infof("账号已添加到集群配置: %s (%s)", account.ID, account.PlayerID)

	// 如果配置了自动保存，保存配置
	if m.configPath != "" {
		if err := SaveClusterConfig(m.configPath, m.config); err != nil {
			logx.Warnf("保存集群配置失败: %v", err)
			// 不返回错误，继续操作
		}
	}

	return nil
}

// RemoveAccount 从集群配置中移除账号
func (m *Manager) RemoveAccount(id string) error {
	found := false
	for i, acc := range m.config.Accounts {
		if acc.ID == id {
			// 删除账号
			m.config.Accounts = append(m.config.Accounts[:i], m.config.Accounts[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return ErrInstanceNotFound
	}

	logx.Infof("账号已从集群配置移除: %s", id)

	// 如果配置了自动保存，保存配置
	if m.configPath != "" {
		if err := SaveClusterConfig(m.configPath, m.config); err != nil {
			logx.Warnf("保存集群配置失败: %v", err)
		}
	}

	return nil
}

// SaveConfig 保存当前配置到文件
func (m *Manager) SaveConfig() error {
	if m.configPath == "" {
		return fmt.Errorf("no config path set")
	}

	return SaveClusterConfig(m.configPath, m.config)
}

// SetConfigPath 设置配置文件路径
func (m *Manager) SetConfigPath(path string) {
	m.configPath = path
}
