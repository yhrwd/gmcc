package cluster

import (
	"reflect"
	"testing"
)

func TestInstance_StartTriggerRules(t *testing.T) {
	inst := newInstance("i1", AccountEntry{ID: "i1", PlayerID: "p1"}, nil)
	method, ok := reflect.TypeOf(inst).MethodByName("StartWithTrigger")
	if !ok {
		t.Fatalf("missing symbol: (*Instance).StartWithTrigger")
	}

	if method.Type.NumIn() != 2 {
		t.Fatalf("StartWithTrigger should accept trigger argument, got %d args", method.Type.NumIn()-1)
	}

	t.Fatalf("pending behavior test: start trigger rules should be enforced")
}

func TestInstance_RejectsStaleVersionEvent(t *testing.T) {
	inst := newInstance("i1", AccountEntry{ID: "i1", PlayerID: "p1"}, nil)
	instType := reflect.TypeOf(inst)

	if _, ok := instType.Elem().FieldByName("version"); !ok {
		t.Fatalf("missing symbol: Instance.version")
	}

	if _, ok := instType.MethodByName("applyExitEvent"); !ok {
		t.Fatalf("missing symbol: (*Instance).applyExitEvent")
	}

	t.Fatalf("pending behavior test: stale version exit event must be ignored")
}
