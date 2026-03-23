package component

import (
	"bytes"
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

	// 测试解析VarInt组件
	r := bytes.NewReader([]byte{0x01})
	result, err := parser.ParseComponent(1, r)
	if err != nil {
		t.Errorf("ParseComponent() error = %v", err)
	}
	if result == nil {
		t.Error("ParseComponent() = nil, want non-nil")
	}
	if result.TypeID != 1 {
		t.Errorf("ParseComponent() TypeID = %d, want 1", result.TypeID)
	}
}
