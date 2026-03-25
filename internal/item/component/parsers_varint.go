package component

import (
	"bytes"

	"gmcc/internal/mcclient/packet"
)

// ParseMaxStackSize 解析 max_stack_size 组件 (ID: 1, VarInt)
func ParseMaxStackSize(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	value, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return nil, err
	}
	return &ComponentResult{
		TypeID: typeID,
		Data:   value,
	}, nil
}

// ParseMaxDamage 解析 max_damage 组件 (ID: 2, VarInt)
func ParseMaxDamage(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	value, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return nil, err
	}
	return &ComponentResult{
		TypeID: typeID,
		Data:   value,
	}, nil
}

// ParseDamage 解析 damage 组件 (ID: 3, VarInt)
func ParseDamage(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	value, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return nil, err
	}
	return &ComponentResult{
		TypeID: typeID,
		Data:   value,
	}, nil
}

// ParseRepairCost 解析 repair_cost 组件 (ID: 19, VarInt)
func ParseRepairCost(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	value, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return nil, err
	}
	return &ComponentResult{
		TypeID: typeID,
		Data:   value,
	}, nil
}

// ParseEnchantable 解析 enchantable 组件 (ID: 31, VarInt)
func ParseEnchantable(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	value, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return nil, err
	}
	return &ComponentResult{
		TypeID: typeID,
		Data:   value,
	}, nil
}

// ParseMapID 解析 map_id 组件 (ID: 44, VarInt)
func ParseMapID(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	value, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return nil, err
	}
	return &ComponentResult{
		TypeID: typeID,
		Data:   value,
	}, nil
}
