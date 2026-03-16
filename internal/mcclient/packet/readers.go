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

func ReadSlotData(r *bytes.Reader) (*SlotData, error) {
	startPos := r.Len()

	// 物品ID: VarInt, 0 表示空槽
	itemID, err := ReadVarIntFromReader(r)
	if err != nil {
		return nil, err
	}
	if itemID == 0 {
		return nil, nil
	}

	// 数量: VarInt
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return nil, err
	}

	// 组件数据
	if err := SkipSlotComponents(r); err != nil {
		remaining := r.Len()
		logx.Warnf("Slot解析失败: itemID=%d, count=%d, startPos=%d, remaining=%d, err=%v",
			itemID, count, startPos, remaining, err)
		return nil, err
	}

	return &SlotData{ID: itemID, Count: count}, nil
}

func SkipSlotComponents(r *bytes.Reader) error {
	numAdd, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < numAdd; i++ {
		if err := SkipComponentData(r); err != nil {
			return err
		}
	}

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

func SkipComponentData(r *bytes.Reader) error {
	componentType, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	if err := SkipComponentByType(r, componentType); err != nil {
		return fmt.Errorf("component type %d: %w", componentType, err)
	}
	return nil
}

func SkipNBT(r *bytes.Reader) error {
	dec := nbt.NewDecoder(r).NetworkFormat(true)
	return dec.Skip()
}

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
