package cluster

import (
	"errors"
	"testing"
	"time"
)

func TestManager_StartDoesNotAutoLaunchEnabledAccounts(t *testing.T) {
	initTestLogger(t)

	cfg := DefaultClusterConfig()
	cfg.Accounts = []AccountEntry{{ID: "a1", PlayerID: "p1", Enabled: true}}
	m := NewManager(cfg)
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
	m := NewManager(cfg)
	if err := m.CreateInstance("a1", AccountEntry{ID: "a1", PlayerID: "p1"}); err != nil {
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

			m := NewManager(cfg)
			if err := m.CreateInstance("a1", AccountEntry{ID: "a1", PlayerID: "p1"}); err != nil {
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
