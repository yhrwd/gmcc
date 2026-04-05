package cluster

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	authsession "gmcc/internal/auth/session"
	appconfig "gmcc/internal/config"
	"gmcc/internal/logx"
	"gmcc/internal/resource"
	"gmcc/internal/state"
)

var (
	ErrInstanceNotFound = errors.New("instance not found")
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

	// 监督与删除控制
	supervisionMu sync.Mutex
	supervising   map[string]bool
	deleteTimeout time.Duration

	// 上下文
	ctx    context.Context
	cancel context.CancelFunc

	authManager *authsession.AuthManager
	resourceMgr resourceService
}

type resourceService interface {
	CreateAccount(in resource.CreateAccountInput) (state.AccountMeta, error)
	CreateInstance(in resource.CreateInstanceInput) (state.InstanceMeta, error)
	DeleteAccount(accountID string) error
	DeleteInstance(instanceID string) error
	RestoreResources() (resource.RestoreResourcesResult, error)
}

// NewManager 创建集群管理器
func NewManager(config ClusterConfig, authManager *authsession.AuthManager, configPath ...string) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	m := &Manager{
		config:        config,
		instances:     make(map[string]*Instance),
		startTime:     time.Now(),
		supervising:   make(map[string]bool),
		deleteTimeout: 10 * time.Second,
		ctx:           ctx,
		cancel:        cancel,
		authManager:   authManager,
	}

	// 如果传入了配置路径，保存它
	if len(configPath) > 0 && configPath[0] != "" {
		m.configPath = configPath[0]
	}

	return m
}

// SetResourceManager 设置资源管理器
func (m *Manager) SetResourceManager(resourceManager resourceService) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.resourceMgr = resourceManager
}

// Start 启动集群管理器
func (m *Manager) Start() error {
	logx.Infof("集群管理器启动: max_instances=%d", m.config.Global.MaxInstances)
	if m.resourceMgr != nil {
		if err := m.RestoreInstances(); err != nil {
			return err
		}
	}
	// 设计约束：启动管理器不自动拉起实例，仅通过 API/命令显式启动
	return nil
}

func (m *Manager) beginSupervision(id string) bool {
	m.supervisionMu.Lock()
	defer m.supervisionMu.Unlock()

	if m.supervising[id] {
		return false
	}

	m.supervising[id] = true
	return true
}

func (m *Manager) endSupervision(id string) {
	m.supervisionMu.Lock()
	delete(m.supervising, id)
	m.supervisionMu.Unlock()
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
	if m.resourceMgr != nil {
		meta, err := m.resourceMgr.CreateInstance(resource.CreateInstanceInput{
			InstanceID:    id,
			AccountID:     account.ID,
			ServerAddress: account.ServerAddress,
			Enabled:       account.Enabled,
		})
		if err != nil {
			return err
		}

		account.ID = meta.AccountID
		account.ServerAddress = meta.ServerAddress
		account.Enabled = meta.Enabled
		return m.createRuntimeInstance(meta.InstanceID, account)
	}

	return m.createRuntimeInstance(id, account)
}

func (m *Manager) createRuntimeInstance(id string, account AccountEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查是否已存在
	if _, exists := m.instances[id]; exists {
		return fmt.Errorf("%w: %s", resource.ErrInstanceAlreadyExists, id)
	}

	// 检查最大实例数
	if m.config.Global.MaxInstances > 0 && len(m.instances) >= m.config.Global.MaxInstances {
		return ErrMaxInstances
	}

	// 创建实例
	inst := newInstance(id, account, m)
	inst.authManager = m.authManager
	m.instances[id] = inst

	logx.Infof("实例已创建: %s (account_id=%s)", id, account.ID)
	return nil
}

// RestoreInstances 从资源元数据恢复内存中的实例对象（不自动启动）
func (m *Manager) RestoreInstances() error {
	m.mu.RLock()
	resourceMgr := m.resourceMgr
	m.mu.RUnlock()

	if resourceMgr == nil {
		return nil
	}

	result, err := resourceMgr.RestoreResources()
	if err != nil {
		return fmt.Errorf("restore resources: %w", err)
	}

	for _, meta := range result.RestoredInstances {
		account := m.accountForRestoredInstance(meta)
		if err := m.createRuntimeInstance(meta.InstanceID, account); err != nil {
			if errors.Is(err, resource.ErrInstanceAlreadyExists) {
				continue
			}
			return fmt.Errorf("materialize restored instance %q: %w", meta.InstanceID, err)
		}
	}

	if result.RestoredCount > 0 || result.SkippedCount > 0 {
		logx.Infof("集群实例恢复完成: restored=%d skipped=%d", result.RestoredCount, result.SkippedCount)
	}
	return nil
}

func (m *Manager) accountForRestoredInstance(meta state.InstanceMeta) AccountEntry {
	account := AccountEntry{
		ID:            meta.AccountID,
		ServerAddress: meta.ServerAddress,
		Enabled:       meta.Enabled,
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, configured := range m.config.Accounts {
		if configured.ID == meta.AccountID {
			if configured.ServerAddress != "" {
				account.ServerAddress = configured.ServerAddress
			}
			return account
		}
	}

	return account
}

// DeleteInstance 删除实例
func (m *Manager) DeleteInstance(id string) error {
	inst, err := m.GetInstance(id)
	if err != nil {
		return err
	}

	status := inst.GetStatus()
	if status == StatusRunning || status == StatusStarting || status == StatusReconnecting {
		_ = inst.Stop()

		timeout := m.deleteTimeout
		if timeout <= 0 {
			timeout = 10 * time.Second
		}

		ctx, cancel := context.WithTimeout(m.ctx, timeout)
		defer cancel()
		if err := inst.waitExit(ctx); err != nil {
			inst.markError("delete timeout")
			return ErrDeleteTimeout
		}
	}

	inst.markDeleted()

	if m.resourceMgr != nil {
		if err := m.resourceMgr.DeleteInstance(id); err != nil && !errors.Is(err, resource.ErrInstanceNotFound) {
			return err
		}
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
	// 这里可以触发重连逻辑
	inst, getErr := m.GetInstance(id)
	if getErr != nil {
		return
	}

	// 检查是否需要自动重连
	if m.config.Global.ReconnectPolicy.Enabled {
		if classifyExitCategory(err) != ExitCategoryNetworkDisconnect {
			return
		}

		if !m.beginSupervision(inst.ID) {
			return
		}
		go m.handleReconnect(inst)
	}
}

// handleReconnect 处理自动重连
func (m *Manager) handleReconnect(inst *Instance) {
	defer m.endSupervision(inst.ID)

	policy := m.config.Global.ReconnectPolicy

	// 首次进入重连态并记录尝试次数
	inst.mu.Lock()
	_ = inst.transitionToLocked(StatusReconnecting, "network disconnect")
	attempt := inst.bumpReconnectAttemptsLocked() - 1 // 本次尝试的索引（从0开始）
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

		m.emitReconnectEvent("warn", "scheduled", "instance reconnect scheduled", inst, "network_disconnect", map[string]any{
			"attempt":     attempt + 1,
			"backoff_ms":  delay.Milliseconds(),
			"max_retries": policy.MaxRetries,
		})
		logx.Summaryf("warn", "实例 %s 断线，计划在 %s 后重连（第 %d 次）", inst.ID, delay, attempt+1)

		select {
		case <-m.ctx.Done():
			return
		case <-inst.reconnectStopCh:
			return
		case <-time.After(delay):
		}

		if inst.shouldAbortReconnect() {
			return
		}

		// 尝试启动（仅自动重连触发）
		if err := inst.StartWithTrigger(StartTriggerAutoReconnect); err == nil {
			m.emitReconnectEvent("info", "succeeded", "instance reconnect succeeded", inst, "network_disconnect", map[string]any{
				"attempt":     attempt + 1,
				"max_retries": policy.MaxRetries,
			})
			logx.Summaryf("info", "实例 %s 重连成功", inst.ID)
			inst.resetReconnectAttempts()
			return
		} else {
			m.emitReconnectEvent("warn", "failed", "instance reconnect attempt failed", inst, reconnectFailureReason(err), map[string]any{
				"attempt":     attempt + 1,
				"max_retries": policy.MaxRetries,
			})
			logx.Warnf("实例 %s 重连失败（第 %d 次）: %v", inst.ID, attempt+1, err)
			attempt = inst.bumpReconnectAttempts() - 1
		}
	}

	m.emitReconnectEvent("error", "exhausted", "instance reconnect retries exhausted", inst, "network_disconnect", map[string]any{
		"attempt":     attempt,
		"max_retries": policy.MaxRetries,
	})
	logx.Errorf("实例 %s 重连次数已耗尽", inst.ID)
}

func reconnectFailureReason(err error) string {
	if err == nil {
		return "unknown"
	}
	category := classifyExitCategory(err)
	if category != ExitCategoryUnknown {
		return string(category)
	}
	return "start_failed"
}

func (m *Manager) emitReconnectEvent(level, action, message string, inst *Instance, reason string, fields map[string]any) {
	event := logx.NewReconnectEvent(level, action, message, inst.ID, inst.Account.ID, reason)
	event.Fields = fields
	logx.Emit(event)
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
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查ID是否已存在
	for _, acc := range m.config.Accounts {
		if acc.ID == account.ID {
			return fmt.Errorf("account with ID '%s' already exists", account.ID)
		}
	}

	// 添加账号
	m.config.Accounts = append(m.config.Accounts, account)
	logx.Infof("账号已添加到集群配置: %s", account.ID)

	// 如果配置了自动保存，保存配置
	if m.configPath != "" {
		if err := m.saveConfigLocked(); err != nil {
			logx.Warnf("保存集群配置失败: %v", err)
			// 不返回错误，继续操作
		}
	}

	return nil
}

// RemoveAccount 从集群配置中移除账号
func (m *Manager) RemoveAccount(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.resourceMgr != nil {
		if err := m.resourceMgr.DeleteAccount(id); err != nil && !errors.Is(err, resource.ErrAccountNotFound) {
			return err
		}
	}

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
		if err := m.saveConfigLocked(); err != nil {
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

	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.saveConfigLocked()
}

// SetConfigPath 设置配置文件路径
func (m *Manager) SetConfigPath(path string) {
	m.configPath = path
}

func (m *Manager) saveConfigLocked() error {
	cfg, err := appconfig.Load(m.configPath)
	if err != nil {
		return fmt.Errorf("load app config: %w", err)
	}

	cfg.Cluster = appconfig.ClusterConfig{
		Global: appconfig.GlobalConfig{
			MaxInstances: m.config.Global.MaxInstances,
			ReconnectPolicy: appconfig.ReconnectPolicy{
				Enabled:    m.config.Global.ReconnectPolicy.Enabled,
				MaxRetries: m.config.Global.ReconnectPolicy.MaxRetries,
				BaseDelay:  m.config.Global.ReconnectPolicy.BaseDelay,
				MaxDelay:   m.config.Global.ReconnectPolicy.MaxDelay,
				Multiplier: m.config.Global.ReconnectPolicy.Multiplier,
			},
		},
		Accounts: convertAccountsToConfig(m.config.Accounts),
	}

	if err := appconfig.Save(m.configPath, *cfg); err != nil {
		return fmt.Errorf("save app config: %w", err)
	}
	return nil
}

func convertAccountsToConfig(accounts []AccountEntry) []appconfig.AccountEntry {
	result := make([]appconfig.AccountEntry, len(accounts))
	for i, account := range accounts {
		result[i] = appconfig.AccountEntry{
			ID:            account.ID,
			ServerAddress: account.ServerAddress,
			Enabled:       account.Enabled,
		}
	}
	return result
}
