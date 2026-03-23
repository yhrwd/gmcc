package component

import (
	"testing"
)

func TestNewParser(t *testing.T) {
	parser := NewParser()
	if parser == nil {
		t.Error("NewParser() = nil, want non-nil")
	}
	if parser.handlers == nil {
		t.Error("NewParser().handlers = nil, want non-nil")
	}
}

func TestParseComponent(t *testing.T) {
	parser := NewParser()

	// 测试解析一个组件
	result, err := parser.ParseComponent(0, nil)
	if err != nil {
		t.Errorf("ParseComponent() error = %v", err)
	}
	if result == nil {
		t.Error("ParseComponent() = nil, want non-nil")
	}
	if result.TypeID != 0 {
		t.Errorf("ParseComponent() TypeID = %d, want 0", result.TypeID)
	}
}
