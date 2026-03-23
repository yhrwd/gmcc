package packet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	"gmcc/internal/constants"
	"gmcc/internal/logx"
	"gmcc/internal/nbt"
)

// SlotData 兼容类型定义（避免循环依赖）
type SlotData struct {
	ID    int32
	Count int32
}

// MustReadBytes 读取字节，错误不终止但记录日志
// 用于解析非关键数据时忽略错误
func MustReadBytes(r io.Reader, n int, name string) []byte {
	b, err := ReadBytes(r, n)
	if err != nil {
		logx.PacketWarn(name, err)
	}
	return b
}

type byteReaderWrapper struct {
	r io.Reader
}

func newByteReaderWrapper(r io.Reader) *byteReaderWrapper {
	return &byteReaderWrapper{r: r}
}

func (b *byteReaderWrapper) Read(p []byte) (n int, err error) {
	return b.r.Read(p)
}

func (b *byteReaderWrapper) ReadByte() (byte, error) {
	var buf [1]byte
	if _, err := io.ReadFull(b.r, buf[:]); err != nil {
		return 0, err
	}
	return buf[0], nil
}

// MustReadVarInt 读取 VarInt，错误不终止但记录日志
func MustReadVarInt(r io.Reader, name string) int32 {
	v, err := ReadVarIntFromReader(r)
	if err != nil {
		logx.PacketWarn(name, err)
	}
	return v
}

// MustReadString 读取字符串，错误不终止但记录日志
func MustReadString(r io.Reader, name string) string {
	s, err := ReadStringFromReader(r)
	if err != nil {
		logx.PacketWarn(name, err)
	}
	return s
}

// MustReadBool 读取布尔值，错误不终止但记录日志
func MustReadBool(r io.Reader, name string) bool {
	v, err := ReadBool(r)
	if err != nil {
		logx.PacketWarn(name, err)
	}
	return v
}

// MustReadU8 读取 uint8，错误不终止但记录日志
func MustReadU8(r io.Reader, name string) byte {
	v, err := ReadU8(r)
	if err != nil {
		logx.PacketWarn(name, err)
	}
	return v
}

// MustReadInt32 读取 int32，错误不终止但记录日志
func MustReadInt32(r io.Reader, name string) int32 {
	v, err := ReadInt32(r)
	if err != nil {
		logx.PacketWarn(name, err)
	}
	return v
}

// MustReadFloat64 读取 float64，错误不终止但记录日志
func MustReadFloat64(r io.Reader, name string) float64 {
	v, err := ReadFloat64FromReader(r)
	if err != nil {
		logx.PacketWarn(name, err)
	}
	return v
}

func ReadVarIntFromReader(r io.Reader) (int32, error) {
	if br, ok := r.(io.ByteReader); ok {
		return ReadVarInt(br)
	}
	return ReadVarInt(newByteReaderWrapper(r))
}

func ReadStringFromReader(r io.Reader) (string, error) {
	if br, ok := r.(io.ByteReader); ok {
		return ReadString(br, r)
	}
	return ReadString(newByteReaderWrapper(r), r)
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

func ReadBytes(r io.Reader, n int) ([]byte, error) {
	if n < 0 {
		return nil, fmt.Errorf("negative read length: %d", n)
	}
	if n > constants.MaxPacketSize {
		return nil, fmt.Errorf("read length exceeds max allowed: %d > %d", n, constants.MaxPacketSize)
	}
	b := make([]byte, n)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, fmt.Errorf("read %d bytes: %w", n, err)
	}
	return b, nil
}

// ReadSlotData 使用新的内部实现
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

	// 3. 跳过components (暂时使用旧逻辑)
	if err := SkipSlotComponents(r); err != nil {
		logx.Warnf("Slot解析失败: itemID=%d, count=%d, err=%v", itemID, count, err)
		return nil, err
	}

	return &SlotData{ID: itemID, Count: count}, nil
}

// 内部实现，将在阶段3后完全移除
// func internalReadSlotData(r *bytes.Reader) (*SlotData, error) {
// 	// 这里将在阶段3中实现新的组件解析系统集成
// 	return nil, nil
// }

// SkipSlotComponents 跳过物品组件 (1.21.11)
// 结构: addComponentPatchesCount(VarInt) -> removeComponentPatchesCount(VarInt) ->
//
//	addComponentPatches(Array) -> removeComponentPatches(Array)
func SkipSlotComponents(r *bytes.Reader) error {
	// 添加的组件数量
	numAdd, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}

	// 移除的组件数量
	numRemove, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}

	// 添加的组件数组 (component_type(VarInt) + data)
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

	// 移除的组件数组 (只有 component_type)
	for i := int32(0); i < numRemove; i++ {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
	}
	return nil
}

// SkipNBT 跳过 Network NBT 格式 (无 name 字段)
func SkipNBT(r *bytes.Reader) error {
	if r.Len() == 0 {
		return nil
	}
	dec := nbt.NewDecoder(r).NetworkFormat(true)
	err := dec.Skip()
	if err != nil && err.Error() == "unexpected EOF" {
		return nil
	}
	return err
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
