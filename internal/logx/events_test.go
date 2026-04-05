package logx

import (
	"encoding/json"
	"testing"
	"time"
)

func TestEventMarshalIncludesEventTypeAndAction(t *testing.T) {
	e := Event{
		Timestamp:  time.Date(2026, 4, 5, 8, 30, 0, 0, time.UTC),
		Level:      "info",
		EventType:  "instance.lifecycle",
		Action:     "ready",
		Message:    "instance ready",
		InstanceID: "bot-1",
		AccountID:  "acc-main",
	}

	data, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}

	if !containsJSONField(data, "event_type", "instance.lifecycle") {
		t.Fatalf("expected event_type field in %s", string(data))
	}

	if !containsJSONField(data, "action", "ready") {
		t.Fatalf("expected action field in %s", string(data))
	}

	if !containsJSONField(data, "instance_id", "bot-1") {
		t.Fatalf("expected instance_id field in %s", string(data))
	}

	if !containsJSONField(data, "account_id", "acc-main") {
		t.Fatalf("expected account_id field in %s", string(data))
	}
}

func TestEventRejectsSensitiveFields(t *testing.T) {
	e := Event{
		EventType: "auth.session",
		Action:    "auth_failed",
		Message:   "token leaked?",
		AuthError: "device_login_required",
		Fields: map[string]any{
			"access_token": "secret-token",
		},
	}

	if err := e.Validate(); err == nil {
		t.Fatalf("expected sensitive field validation failure")
	}
}

func containsJSONField(data []byte, key string, want string) bool {
	var fields map[string]any
	if err := json.Unmarshal(data, &fields); err != nil {
		return false
	}

	got, ok := fields[key]
	if !ok {
		return false
	}

	gotString, ok := got.(string)
	if !ok {
		return false
	}

	return gotString == want
}
