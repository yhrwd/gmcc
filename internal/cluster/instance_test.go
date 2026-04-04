package cluster

import (
	"testing"
)

func TestInstance_StartTriggerRules(t *testing.T) {
	tests := []struct {
		name      string
		status    InstanceStatus
		trigger   StartTrigger
		wantError bool
	}{
		{"manual from stopped", StatusStopped, StartTriggerManualStart, false},
		{"manual from reconnecting denied", StatusReconnecting, StartTriggerManualStart, true},
		{"auto reconnect from reconnecting allowed", StatusReconnecting, StartTriggerAutoReconnect, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := newInstance("i1", AccountEntry{ID: "i1", PlayerID: "p1"}, nil)
			inst.status = tt.status
			inst.startRunnerFn = func(_ int64) error {
				return nil
			}

			err := inst.StartWithTrigger(tt.trigger)
			if (err != nil) != tt.wantError {
				t.Fatalf("wantError=%v, err=%v", tt.wantError, err)
			}
		})
	}
}

func TestInstance_RejectsStaleVersionEvent(t *testing.T) {
	inst := newInstance("i1", AccountEntry{ID: "i1", PlayerID: "p1"}, nil)
	inst.version = 7
	inst.runVersion = 7
	inst.status = StatusRunning

	changed := inst.applyExitEvent(6, ExitCategoryNetworkDisconnect, errFakeNetworkEOF)
	if changed {
		t.Fatalf("expected stale event not to mutate state")
	}

	if got := inst.GetStatus(); got != StatusRunning {
		t.Fatalf("expected status to remain running, got %s", got)
	}
}
