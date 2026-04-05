package logx

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var forbiddenEventFields = map[string]struct{}{
	"access_token":  {},
	"refresh_token": {},
	"device_code":   {},
	"xsts_token":    {},
}

// Event 表示结构化可观测事件。
type Event struct {
	Timestamp  time.Time      `json:"ts"`
	Level      string         `json:"level,omitempty"`
	EventType  string         `json:"event_type"`
	Action     string         `json:"action"`
	Message    string         `json:"message,omitempty"`
	InstanceID string         `json:"instance_id,omitempty"`
	AccountID  string         `json:"account_id,omitempty"`
	PlayerID   string         `json:"player_id,omitempty"`
	Reason     string         `json:"reason,omitempty"`
	AuthError  string         `json:"auth_error,omitempty"`
	Result     string         `json:"result,omitempty"`
	Fields     map[string]any `json:"-"`
}

// Validate 校验事件必填字段与敏感字段约束。
func (e Event) Validate() error {
	if strings.TrimSpace(e.EventType) == "" || strings.TrimSpace(e.Action) == "" {
		return fmt.Errorf("event_type and action are required")
	}

	for key := range e.Fields {
		normalizedKey := strings.ToLower(strings.TrimSpace(key))
		if _, blocked := forbiddenEventFields[normalizedKey]; blocked {
			return fmt.Errorf("sensitive event field forbidden: %s", key)
		}
	}

	return nil
}

// MarshalJSON 输出包含扩展字段的结构化事件。
func (e Event) MarshalJSON() ([]byte, error) {
	payload := make(map[string]any, len(e.Fields)+10)

	ts := e.Timestamp
	if ts.IsZero() {
		ts = time.Now().UTC()
	} else {
		ts = ts.UTC()
	}
	payload["ts"] = ts

	level := strings.ToLower(strings.TrimSpace(e.Level))
	if level == "" {
		level = "info"
	}
	payload["level"] = level
	payload["event_type"] = e.EventType
	payload["action"] = e.Action

	if e.Message != "" {
		payload["message"] = e.Message
	}
	if e.InstanceID != "" {
		payload["instance_id"] = e.InstanceID
	}
	if e.AccountID != "" {
		payload["account_id"] = e.AccountID
	}
	if e.PlayerID != "" {
		payload["player_id"] = e.PlayerID
	}
	if e.Reason != "" {
		payload["reason"] = e.Reason
	}
	if e.AuthError != "" {
		payload["auth_error"] = e.AuthError
	}
	if e.Result != "" {
		payload["result"] = e.Result
	}

	for key, value := range e.Fields {
		normalizedKey := strings.TrimSpace(key)
		if normalizedKey == "" {
			continue
		}
		if _, exists := payload[normalizedKey]; exists {
			continue
		}
		payload[normalizedKey] = value
	}

	return json.Marshal(payload)
}

// NewLifecycleEvent 创建实例生命周期事件。
func NewLifecycleEvent(level, action, message, instanceID, accountID string) Event {
	return Event{
		Timestamp:  time.Now().UTC(),
		Level:      level,
		EventType:  "instance.lifecycle",
		Action:     action,
		Message:    message,
		InstanceID: instanceID,
		AccountID:  accountID,
	}
}

// NewReconnectEvent 创建实例重连事件。
func NewReconnectEvent(level, action, message, instanceID, accountID, reason string) Event {
	return Event{
		Timestamp:  time.Now().UTC(),
		Level:      level,
		EventType:  "instance.reconnect",
		Action:     action,
		Message:    message,
		InstanceID: instanceID,
		AccountID:  accountID,
		Reason:     reason,
	}
}

// NewAuthEvent 创建认证会话事件。
func NewAuthEvent(level, action, message, instanceID, accountID, authError, result string) Event {
	return Event{
		Timestamp:  time.Now().UTC(),
		Level:      level,
		EventType:  "auth.session",
		Action:     action,
		Message:    message,
		InstanceID: instanceID,
		AccountID:  accountID,
		AuthError:  authError,
		Result:     result,
	}
}
