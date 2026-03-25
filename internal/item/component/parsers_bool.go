package component

import (
	"bytes"

	"gmcc/internal/mcclient/packet"
)

// ParseBool 解析 Bool 类型的组件 (通用)
func ParseBool(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	value, err := packet.ReadBoolFromReader(r)
	if err != nil {
		return nil, err
	}
	return &ComponentResult{
		TypeID: typeID,
		Data:   value,
	}, nil
}

// ParseUnbreakable 解析 unbreakable 组件 (ID: 4, 无数据)
func ParseUnbreakable(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	return &ComponentResult{
		TypeID: typeID,
		Data:   true,
	}, nil
}

// ParseEnchantmentGlintOverride 解析 enchantment_glint_override 组件 (ID: 21, Bool)
func ParseEnchantmentGlintOverride(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	value, err := packet.ReadBoolFromReader(r)
	if err != nil {
		return nil, err
	}
	return &ComponentResult{
		TypeID: typeID,
		Data:   value,
	}, nil
}
