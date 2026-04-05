package logx

import (
	"fmt"
	"os"
	"path/filepath"
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

type structuredWriter struct {
	file *os.File
}

func newStructuredWriter(activePath string, _ int64, _ int) (*structuredWriter, error) {
	if err := os.MkdirAll(filepath.Dir(activePath), 0o755); err != nil {
		return nil, fmt.Errorf("create structured log directory: %w", err)
	}

	file, err := os.OpenFile(activePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open structured log file: %w", err)
	}

	return &structuredWriter{file: file}, nil
}

func (w *structuredWriter) Write(p []byte) (int, error) {
	if w == nil || w.file == nil {
		return 0, fmt.Errorf("structured writer is closed")
	}

	return w.file.Write(p)
}

func (w *structuredWriter) Close() error {
	if w == nil || w.file == nil {
		return nil
	}

	err := w.file.Close()
	w.file = nil
	return err
}
