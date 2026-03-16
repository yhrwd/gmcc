package packet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	"gmcc/internal/logx"
	"gmcc/internal/nbt"
)

type SlotData struct {
	ID    int32
	Count int32
}

func ReadVarIntFromReader(r *bytes.Reader) (int32, error) {
	return ReadVarInt(r)
}

func ReadStringFromReader(r *bytes.Reader) (string, error) {
	return ReadString(r, r)
}

func ReadBoolFromReader(r io.Reader) (bool, error) {
	return ReadBool(r)
}

func ReadInt32FromReader(r io.Reader) (int32, error) {
	return ReadInt32(r)
}

func ReadFloat64FromReader(r io.Reader) (float64, error) {
	var v float64
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

func ReadFloat32FromReader(r io.Reader) (float32, error) {
	var v float32
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

func ReadU8(r io.Reader) (byte, error) {
	var b [1]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return 0, err
	}
	return b[0], nil
}

func ReadBytes(r io.Reader, n int) []byte {
	b := make([]byte, n)
	_, _ = io.ReadFull(r, b)
	return b
}

// ReadSlotData 解析 1.21+ ItemStack 格式
// 结构: count(VarInt) -> [如果count>0] item_id(VarInt) -> components
func ReadSlotData(r *bytes.Reader) (*SlotData, error) {
	// 1. item_count (VarInt)
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil // 空物品
	}

	// 2. item_id (VarInt)
	itemID, err := ReadVarIntFromReader(r)
	if err != nil {
		return nil, err
	}

	// 3. 跳过components
	if err := SkipSlotComponents(r); err != nil {
		logx.Warnf("Slot解析失败: itemID=%d, count=%d, err=%v", itemID, count, err)
		return nil, err
	}

	return &SlotData{ID: itemID, Count: count}, nil
}

// SkipSlotComponents 跳过物品组件
// 结构: components_to_add(VarInt) -> [component_type(VarInt) + data] -> components_to_remove(VarInt) -> [component_type(VarInt)]
func SkipSlotComponents(r *bytes.Reader) error {
	// 添加的组件
	numAdd, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < numAdd; i++ {
		// 先读 component_type
		componentType, err := ReadVarIntFromReader(r)
		if err != nil {
			return err
		}
		// 再根据类型跳过数据
		if err := SkipComponentByType(r, componentType); err != nil {
			return fmt.Errorf("component type %d: %w", componentType, err)
		}
	}

	// 移除的组件 (只有 component_type)
	numRemove, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < numRemove; i++ {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
	}
	return nil
}

// SkipNBT 跳过 Network NBT 格式 (无 name 字段)
func SkipNBT(r *bytes.Reader) error {
	dec := nbt.NewDecoder(r).NetworkFormat(true)
	return dec.Skip()
}

// ReadAnonymousNBTJSON 解析 Network NBT 并返回 JSON 字符串
func ReadAnonymousNBTJSON(r io.Reader) (string, error) {
	dec := nbt.NewDecoder(r)
	dec.NetworkFormat(true)
	var v any
	if err := dec.Decode(&v); err != nil {
		return "", err
	}

	raw, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
