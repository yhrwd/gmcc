package cluster

import (
	"reflect"
	"testing"
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

	mgrType := reflect.TypeOf(m)
	if _, ok := mgrType.Elem().FieldByName("deleteTimeout"); !ok {
		t.Fatalf("missing symbol: Manager.deleteTimeout")
	}

	t.Fatalf("pending behavior test: delete should timeout when runner blocks")
}
