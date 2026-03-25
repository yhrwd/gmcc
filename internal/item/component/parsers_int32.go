package component

import (
	"bytes"

	"gmcc/internal/mcclient/packet"
)

// ParseInt32 解析 Int32 类型的组件 (通用)
func ParseInt32(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	value, err := packet.ReadInt32FromReader(r)
	if err != nil {
		return nil, err
	}
	return &ComponentResult{
		TypeID: typeID,
		Data:   value,
	}, nil
}

// ParseCustomModelData 解析 custom_model_data 组件 (ID: 17, Int32)
func ParseCustomModelData(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	value, err := packet.ReadInt32FromReader(r)
	if err != nil {
		return nil, err
	}
	return &ComponentResult{
		TypeID: typeID,
		Data:   value,
	}, nil
}

// ParseDyedColor 解析 dyed_color 组件 (ID: 42, Int32)
func ParseDyedColor(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	value, err := packet.ReadInt32FromReader(r)
	if err != nil {
		return nil, err
	}
	return &ComponentResult{
		TypeID: typeID,
		Data:   value,
	}, nil
}

// ParseMapColor 解析 map_color 组件 (ID: 43, Int32)
func ParseMapColor(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	value, err := packet.ReadInt32FromReader(r)
	if err != nil {
		return nil, err
	}
	return &ComponentResult{
		TypeID: typeID,
		Data:   value,
	}, nil
}
