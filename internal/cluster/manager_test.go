package cluster

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	authsession "gmcc/internal/auth/session"
	"gmcc/internal/logx"
	"gmcc/internal/resource"
	"gmcc/internal/state"
)

type fakeResourceManager struct {
	accountInput      resource.CreateAccountInput
	accountResult     state.AccountMeta
	accountErr        error
	accountCalled     bool
	createInput       resource.CreateInstanceInput
	createResult      state.InstanceMeta
	createErr         error
	createCalled      bool
	deleteAccountID   string
	deleteAccountErr  error
	deleteInstanceID  string
	deleteInstanceErr error
	restoreResult     resource.RestoreResourcesResult
	restoreErr        error
	restoreCalled     bool
}

func (f *fakeResourceManager) CreateAccount(in resource.CreateAccountInput) (state.AccountMeta, error) {
	f.accountCalled = true
	f.accountInput = in
	if f.accountErr != nil {
		return state.AccountMeta{}, f.accountErr
	}
	if f.accountResult.AccountID == "" {
		f.accountResult = state.AccountMeta{AccountID: in.AccountID, Enabled: in.Enabled}
	}
	return f.accountResult, nil
}

func (f *fakeResourceManager) CreateInstance(in resource.CreateInstanceInput) (state.InstanceMeta, error) {
	f.createCalled = true
	f.createInput = in
	if f.createErr != nil {
		return state.InstanceMeta{}, f.createErr
	}
	if f.createResult.InstanceID == "" {
		f.createResult = state.InstanceMeta{
			InstanceID:    in.InstanceID,
			AccountID:     in.AccountID,
			ServerAddress: in.ServerAddress,
			Enabled:       in.Enabled,
		}
	}
	return f.createResult, nil
}

func (f *fakeResourceManager) DeleteAccount(accountID string) error {
	f.deleteAccountID = accountID
	return f.deleteAccountErr
}

func (f *fakeResourceManager) DeleteInstance(instanceID string) error {
	f.deleteInstanceID = instanceID
	return f.deleteInstanceErr
}

func (f *fakeResourceManager) RestoreResources() (resource.RestoreResourcesResult, error) {
	f.restoreCalled = true
	if f.restoreErr != nil {
		return resource.RestoreResourcesResult{}, f.restoreErr
	}
	return f.restoreResult, nil
}

func TestManager_CreateInstanceUsesResourceManager(t *testing.T) {
	initTestLogger(t)

	rm := &fakeResourceManager{createResult: state.InstanceMeta{
		InstanceID:    "bot-1",
		AccountID:     "acc-main",
		ServerAddress: "mc.example.com",
		Enabled:       true,
	}}
	m := NewManager(DefaultClusterConfig(), nil)
	m.SetResourceManager(rm)

	err := m.CreateInstance("bot-1", AccountEntry{ID: "acc-main", ServerAddress: "ignored.example.com", Enabled: false})
	if err != nil {
		t.Fatalf("create instance: %v", err)
	}
	if !rm.createCalled {
		t.Fatalf("expected resource manager create to be called")
	}
	if rm.createInput.InstanceID != "bot-1" || rm.createInput.AccountID != "acc-main" {
		t.Fatalf("unexpected create input: %+v", rm.createInput)
	}
	if rm.createInput.ServerAddress != "ignored.example.com" {
		t.Fatalf("expected create input server address, got %q", rm.createInput.ServerAddress)
	}

	inst, err := m.GetInstance("bot-1")
	if err != nil {
		t.Fatalf("get instance: %v", err)
	}
	if inst.Account.ID != "acc-main" {
		t.Fatalf("expected runtime account id acc-main, got %q", inst.Account.ID)
	}
	if inst.Account.ServerAddress != "mc.example.com" {
		t.Fatalf("expected runtime server address from metadata, got %q", inst.Account.ServerAddress)
	}
	if !inst.Account.Enabled {
		t.Fatalf("expected runtime account enabled from metadata")
	}
}

func TestManager_StartRestoresValidatedInstancesWithoutAutoStart(t *testing.T) {
	initTestLogger(t)

	cfg := DefaultClusterConfig()
	cfg.Accounts = []AccountEntry{{ID: "acc-main", ServerAddress: "mc.example.com", Enabled: true}}
	rm := &fakeResourceManager{restoreResult: resource.RestoreResourcesResult{
		RestoredInstances: []state.InstanceMeta{{
			InstanceID:    "bot-restore",
			AccountID:     "acc-main",
			ServerAddress: "mc.example.com",
			Enabled:       true,
		}},
		RestoredCount: 1,
	}}
	m := NewManager(cfg, nil)
	m.SetResourceManager(rm)

	if err := m.Start(); err != nil {
		t.Fatalf("start manager: %v", err)
	}
	if !rm.restoreCalled {
		t.Fatalf("expected restore resources to be called")
	}

	inst, err := m.GetInstance("bot-restore")
	if err != nil {
		t.Fatalf("get restored instance: %v", err)
	}
	if got := inst.GetStatus(); got != StatusPending {
		t.Fatalf("expected restored instance to remain pending, got %s", got)
	}
	if inst.Account.ID != "acc-main" {
		t.Fatalf("expected restored runtime account id acc-main, got %q", inst.Account.ID)
	}
	if got := len(m.ListInstances()); got != 1 {
		t.Fatalf("expected one restored instance, got %d", got)
	}
}

func TestManager_GetInstanceInfoIncludesAccountID(t *testing.T) {
	initTestLogger(t)

	m := NewManager(DefaultClusterConfig(), nil)
	if err := m.CreateInstance("bot-1", AccountEntry{ID: "acc-main", ServerAddress: "mc.example.com", Enabled: true}); err != nil {
		t.Fatalf("create instance: %v", err)
	}

	info, err := m.GetInstanceInfo("bot-1")
	if err != nil {
		t.Fatalf("get instance info: %v", err)
	}
	if info.AccountID != "acc-main" {
		t.Fatalf("want account id acc-main, got %q", info.AccountID)
	}
}

func TestManager_StartDoesNotAutoLaunchEnabledAccounts(t *testing.T) {
	initTestLogger(t)

	cfg := DefaultClusterConfig()
	cfg.Accounts = []AccountEntry{{ID: "a1", Enabled: true}}
	m := NewManager(cfg, nil)
	if err := m.Start(); err != nil {
		t.Fatalf("start manager: %v", err)
	}
	if got := len(m.ListInstances()); got != 0 {
		t.Fatalf("expected 0 instances on manager start, got %d", got)
	}
}

func TestManager_DeleteTimeoutWhenRunnerBlocks(t *testing.T) {
	initTestLogger(t)

	cfg := DefaultClusterConfig()
	cfg.Global.ReconnectPolicy.Enabled = false
	m := NewManager(cfg, nil)
	if err := m.CreateInstance("a1", AccountEntry{ID: "a1"}); err != nil {
		t.Fatalf("create: %v", err)
	}
	m.deleteTimeout = 20 * time.Millisecond

	inst, err := m.GetInstance("a1")
	if err != nil {
		t.Fatalf("get instance: %v", err)
	}
	inst.startRunnerFn = func(_ uint64) error { return nil }

	if err := inst.Start(); err != nil {
		t.Fatalf("start instance: %v", err)
	}

	err = m.DeleteInstance("a1")
	if !errors.Is(err, ErrDeleteTimeout) {
		t.Fatalf("expected ErrDeleteTimeout, got %v", err)
	}

	if !m.InstanceExists("a1") {
		t.Fatalf("instance should remain after delete timeout")
	}
}

func TestManager_LocalSimulationReconnectPolicy(t *testing.T) {
	initTestLogger(t)

	tests := []struct {
		name            string
		outcome         runnerOutcome
		waitFor         time.Duration
		wantReconnect   bool
		wantFinalStatus InstanceStatus
	}{
		{"network reconnect", runnerOutcome{err: ErrRunnerNetworkDisconnect, ready: true, runDelay: 150 * time.Millisecond}, 220 * time.Millisecond, true, StatusReconnecting},
		{"auth no reconnect", runnerOutcome{err: ErrRunnerAuthFailed}, 40 * time.Millisecond, false, StatusError},
		{"timeout no reconnect", runnerOutcome{err: ErrRunnerStartupTimeout}, 40 * time.Millisecond, false, StatusError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultClusterConfig()
			cfg.Global.ReconnectPolicy.Enabled = true
			cfg.Global.ReconnectPolicy.BaseDelay = time.Second
			cfg.Global.ReconnectPolicy.MaxDelay = time.Second
			cfg.Global.ReconnectPolicy.Multiplier = 1
			cfg.Global.ReconnectPolicy.MaxRetries = 1

			m := NewManager(cfg, nil)
			if err := m.CreateInstance("a1", AccountEntry{ID: "a1"}); err != nil {
				t.Fatalf("create: %v", err)
			}

			inst, err := m.GetInstance("a1")
			if err != nil {
				t.Fatalf("get instance: %v", err)
			}

			factory := newScriptedRunnerFactory(tt.outcome)
			inst.runnerFactory = factory.Build

			if err := inst.Start(); err != nil {
				t.Fatalf("start instance: %v", err)
			}

			time.Sleep(tt.waitFor)

			if got := inst.GetStatus(); got != tt.wantFinalStatus {
				t.Fatalf("status mismatch, got=%s want=%s", got, tt.wantFinalStatus)
			}

			m.supervisionMu.Lock()
			hasReconnect := m.supervising[inst.ID]
			m.supervisionMu.Unlock()
			if hasReconnect != tt.wantReconnect {
				t.Fatalf("reconnect mismatch, got=%v want=%v", hasReconnect, tt.wantReconnect)
			}

			if tt.wantReconnect && factory.BuildCount() != 1 {
				t.Fatalf("expected reconnect to be pending before retry, builds=%d", factory.BuildCount())
			}

			_ = m.Stop()
		})
	}
}

func TestClassifyExitCategoryRecognizesWrappedNetworkErrors(t *testing.T) {
	tests := []error{
		errors.New("读取数据包失败 (state=play): EOF"),
		errors.New("连接服务器失败: connection refused"),
		errors.New("服务器已关闭连接 (state=play): use of closed network connection"),
	}

	for _, err := range tests {
		if got := classifyExitCategory(err); got != ExitCategoryNetworkDisconnect {
			t.Fatalf("expected network disconnect for %v, got %s", err, got)
		}
	}
}

func TestManager_AuthSessionErrorsDoNotReconnect(t *testing.T) {
	initTestLogger(t)

	cfg := DefaultClusterConfig()
	cfg.Global.ReconnectPolicy.Enabled = true
	m := NewManager(cfg, nil)
	if err := m.CreateInstance("a1", AccountEntry{ID: "a1"}); err != nil {
		t.Fatalf("create: %v", err)
	}
	inst, err := m.GetInstance("a1")
	if err != nil {
		t.Fatalf("get instance: %v", err)
	}
	factory := newScriptedRunnerFactory(runnerOutcome{err: authsession.ErrDeviceLoginRequired})
	inst.runnerFactory = factory.Build

	if err := inst.Start(); err != nil {
		t.Fatalf("start instance: %v", err)
	}
	time.Sleep(50 * time.Millisecond)

	if got := inst.GetStatus(); got != StatusError {
		t.Fatalf("expected auth failure to end in error, got %s", got)
	}
	m.supervisionMu.Lock()
	hasReconnect := m.supervising[inst.ID]
	m.supervisionMu.Unlock()
	if hasReconnect {
		t.Fatalf("auth failure should not trigger reconnect supervision")
	}
}

func TestManager_StopCancelsPendingReconnect(t *testing.T) {
	initTestLogger(t)

	cfg := DefaultClusterConfig()
	cfg.Global.ReconnectPolicy.Enabled = true
	cfg.Global.ReconnectPolicy.BaseDelay = 200 * time.Millisecond
	cfg.Global.ReconnectPolicy.MaxDelay = 200 * time.Millisecond
	cfg.Global.ReconnectPolicy.Multiplier = 1
	cfg.Global.ReconnectPolicy.MaxRetries = 1

	m := NewManager(cfg, nil)
	if err := m.CreateInstance("a1", AccountEntry{ID: "a1"}); err != nil {
		t.Fatalf("create: %v", err)
	}
	inst, err := m.GetInstance("a1")
	if err != nil {
		t.Fatalf("get instance: %v", err)
	}
	factory := newScriptedRunnerFactory(runnerOutcome{err: ErrRunnerNetworkDisconnect, ready: true, runDelay: 10 * time.Millisecond}, runnerOutcome{})
	inst.runnerFactory = factory.Build

	if err := inst.Start(); err != nil {
		t.Fatalf("start instance: %v", err)
	}
	time.Sleep(40 * time.Millisecond)
	if err := m.StopInstance("a1"); err != nil {
		t.Fatalf("stop instance: %v", err)
	}
	time.Sleep(260 * time.Millisecond)

	if got := factory.BuildCount(); got != 1 {
		t.Fatalf("expected no restart after stop, builds=%d", got)
	}
}

func TestManager_DeleteCancelsPendingReconnect(t *testing.T) {
	initTestLogger(t)

	cfg := DefaultClusterConfig()
	cfg.Global.ReconnectPolicy.Enabled = true
	cfg.Global.ReconnectPolicy.BaseDelay = 200 * time.Millisecond
	cfg.Global.ReconnectPolicy.MaxDelay = 200 * time.Millisecond
	cfg.Global.ReconnectPolicy.Multiplier = 1
	cfg.Global.ReconnectPolicy.MaxRetries = 1

	m := NewManager(cfg, nil)
	if err := m.CreateInstance("a1", AccountEntry{ID: "a1"}); err != nil {
		t.Fatalf("create: %v", err)
	}
	inst, err := m.GetInstance("a1")
	if err != nil {
		t.Fatalf("get instance: %v", err)
	}
	factory := newScriptedRunnerFactory(runnerOutcome{err: ErrRunnerNetworkDisconnect, ready: true, runDelay: 10 * time.Millisecond}, runnerOutcome{})
	inst.runnerFactory = factory.Build

	if err := inst.Start(); err != nil {
		t.Fatalf("start instance: %v", err)
	}
	time.Sleep(40 * time.Millisecond)
	if err := m.DeleteInstance("a1"); err != nil {
		t.Fatalf("delete instance: %v", err)
	}
	time.Sleep(260 * time.Millisecond)

	if got := factory.BuildCount(); got != 1 {
		t.Fatalf("expected no restart after delete, builds=%d", got)
	}
	if m.InstanceExists("a1") {
		t.Fatalf("instance should remain deleted")
	}
}

func TestManager_ReconnectPolicyEmitsStructuredEvent(t *testing.T) {
	logDir := t.TempDir()
	if err := logx.Init(logDir, true, 1024, false); err != nil {
		t.Fatalf("init logger: %v", err)
	}
	t.Cleanup(func() {
		_ = logx.Close()
	})

	cfg := DefaultClusterConfig()
	cfg.Global.ReconnectPolicy.Enabled = true
	cfg.Global.ReconnectPolicy.BaseDelay = 10 * time.Millisecond
	cfg.Global.ReconnectPolicy.MaxDelay = 10 * time.Millisecond
	cfg.Global.ReconnectPolicy.Multiplier = 1
	cfg.Global.ReconnectPolicy.MaxRetries = 1

	m := NewManager(cfg, nil)
	if err := m.CreateInstance("a1", AccountEntry{ID: "a1"}); err != nil {
		t.Fatalf("create: %v", err)
	}

	inst, err := m.GetInstance("a1")
	if err != nil {
		t.Fatalf("get instance: %v", err)
	}
	factory := newScriptedRunnerFactory(
		runnerOutcome{err: ErrRunnerNetworkDisconnect, ready: true, runDelay: 20 * time.Millisecond},
		runnerOutcome{},
	)
	inst.runnerFactory = factory.Build

	if err := inst.Start(); err != nil {
		t.Fatalf("start instance: %v", err)
	}

	eventsPath := filepath.Join(logDir, "gmcc-events.jsonl")
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		data, readErr := os.ReadFile(eventsPath)
		if readErr == nil {
			content := string(data)
			if strings.Contains(content, `"event_type":"instance.reconnect"`) &&
				strings.Contains(content, `"action":"scheduled"`) &&
				strings.Contains(content, `"action":"succeeded"`) {
				break
			}
		}
		time.Sleep(20 * time.Millisecond)
	}

	data, err := os.ReadFile(eventsPath)
	if err != nil {
		t.Fatalf("read events: %v", err)
	}

	var scheduled, succeeded bool
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var event map[string]any
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatalf("unmarshal event: %v", err)
		}
		if event["event_type"] != "instance.reconnect" {
			continue
		}
		switch event["action"] {
		case "scheduled":
			scheduled = true
			if got := event["reason"]; got != "network_disconnect" {
				t.Fatalf("scheduled reason mismatch: %v", got)
			}
			if got := event["attempt"]; got != float64(1) {
				t.Fatalf("scheduled attempt mismatch: %v", got)
			}
		case "succeeded":
			succeeded = true
			if got := event["attempt"]; got != float64(1) {
				t.Fatalf("succeeded attempt mismatch: %v", got)
			}
		}
	}

	if !scheduled {
		t.Fatalf("expected scheduled reconnect event")
	}
	if !succeeded {
		t.Fatalf("expected succeeded reconnect event")
	}

	_ = m.Stop()
}

func TestManager_InstanceLifecycleEmitsStructuredEvents(t *testing.T) {
	logDir := t.TempDir()
	if err := logx.Init(logDir, true, 1024, false); err != nil {
		t.Fatalf("init logger: %v", err)
	}
	t.Cleanup(func() {
		_ = logx.Close()
	})

	cfg := DefaultClusterConfig()
	cfg.Global.ReconnectPolicy.Enabled = false

	m := NewManager(cfg, nil)
	if err := m.CreateInstance("a1", AccountEntry{ID: "a1", ServerAddress: "example.org:25565"}); err != nil {
		t.Fatalf("create: %v", err)
	}

	inst, err := m.GetInstance("a1")
	if err != nil {
		t.Fatalf("get instance: %v", err)
	}
	factory := newScriptedRunnerFactory(runnerOutcome{ready: true})
	inst.runnerFactory = factory.Build

	if err := inst.Start(); err != nil {
		t.Fatalf("start instance: %v", err)
	}
	time.Sleep(150 * time.Millisecond)
	if err := inst.Stop(); err != nil {
		t.Fatalf("stop instance: %v", err)
	}
	time.Sleep(50 * time.Millisecond)

	data, err := os.ReadFile(filepath.Join(logDir, "gmcc-events.jsonl"))
	if err != nil {
		t.Fatalf("read events: %v", err)
	}

	wantActions := map[string]bool{"start": false, "ready": false, "stop": false}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var event map[string]any
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatalf("unmarshal event: %v", err)
		}
		if event["event_type"] != "instance.lifecycle" {
			continue
		}
		action, _ := event["action"].(string)
		if _, ok := wantActions[action]; ok {
			wantActions[action] = true
			if got := event["instance_id"]; got != "a1" {
				t.Fatalf("instance id mismatch for %s: %v", action, got)
			}
			if got := event["account_id"]; got != "a1" {
				t.Fatalf("account id mismatch for %s: %v", action, got)
			}
		}
	}

	for action, seen := range wantActions {
		if !seen {
			t.Fatalf("expected lifecycle action %s in structured events", action)
		}
	}
}
