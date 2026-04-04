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
