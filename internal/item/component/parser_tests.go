package component

import (
	"bytes"
	"testing"
)

func TestParseComponent_NBT(t *testing.T) {
	parser := NewParser()

	// 测试解析NBT组件
	// 最小的有效NBT数据 (TAG_END)
	r := bytes.NewReader([]byte{0x00})
	result, err := parser.ParseComponent(0, r)
	if err != nil {
		t.Errorf("ParseComponent(NBT) error = %v", err)
	}
	if result == nil {
		t.Error("ParseComponent(NBT) = nil, want non-nil")
	}
	if result.TypeID != 0 {
		t.Errorf("ParseComponent(NBT) TypeID = %d, want 0", result.TypeID)
	}
}

func TestParseComponent_VarInt(t *testing.T) {
	parser := NewParser()

	// 测试解析VarInt组件
	r := bytes.NewReader([]byte{0xFF, 0x01}) // 255
	result, err := parser.ParseComponent(1, r)
	if err != nil {
		t.Errorf("ParseComponent(VarInt) error = %v", err)
	}
	if result == nil {
		t.Error("ParseComponent(VarInt) = nil, want non-nil")
	}
	if result.TypeID != 1 {
		t.Errorf("ParseComponent(VarInt) TypeID = %d, want 1", result.TypeID)
	}
}

func TestParseComponent_String(t *testing.T) {
	parser := NewParser()

	// 测试解析String组件
	r := bytes.NewReader([]byte{0x05, 'h', 'e', 'l', 'l', 'o'}) // "hello"
	result, err := parser.ParseComponent(10, r)
	if err != nil {
		t.Errorf("ParseComponent(String) error = %v", err)
	}
	if result == nil {
		t.Error("ParseComponent(String) = nil, want non-nil")
	}
	if result.TypeID != 10 {
		t.Errorf("ParseComponent(String) TypeID = %d, want 10", result.TypeID)
	}
}

func TestParseComponent_Bool(t *testing.T) {
	parser := NewParser()

	// 测试解析Bool组件
	r := bytes.NewReader([]byte{0x01}) // true
	result, err := parser.ParseComponent(21, r)
	if err != nil {
		t.Errorf("ParseComponent(Bool) error = %v", err)
	}
	if result == nil {
		t.Error("ParseComponent(Bool) = nil, want non-nil")
	}
	if result.TypeID != 21 {
		t.Errorf("ParseComponent(Bool) TypeID = %d, want 21", result.TypeID)
	}
}

func TestRegisterHandler(t *testing.T) {
	parser := NewParser()

	// 注册自定义处理器
	customHandler := func(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
		return &ComponentResult{
			TypeID: typeID,
			Data:   "custom",
		}, nil
	}

	parser.RegisterHandler(999, customHandler)

	// 测试自定义处理器
	r := bytes.NewReader([]byte{})
	result, err := parser.ParseComponent(999, r)
	if err != nil {
		t.Errorf("ParseComponent(custom) error = %v", err)
	}
	if result == nil {
		t.Error("ParseComponent(custom) = nil, want non-nil")
	}
	if result.TypeID != 999 {
		t.Errorf("ParseComponent(custom) TypeID = %d, want 999", result.TypeID)
	}
	if result.Data != "custom" {
		t.Errorf("ParseComponent(custom) Data = %v, want custom", result.Data)
	}
}
