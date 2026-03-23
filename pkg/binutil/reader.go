package binutil

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"gmcc/internal/constants"
)

// Reader 提供二进制数据读取功能
type Reader struct {
	data []byte
	pos  int
}

// NewReader 从字节数据创建 Reader
func NewReader(data []byte) *Reader {
	return &Reader{data: data}
}

// NewReaderFromBytesReader 从 bytes.Reader 创建 Reader
func NewReaderFromBytesReader(r *bytes.Reader) *Reader {
	data, _ := io.ReadAll(r)
	return &Reader{data: data}
}

// ReadVarInt 读取 VarInt 编码的整数
func (r *Reader) ReadVarInt() (int32, error) {
	var result int32
	var shift uint
	for {
		if r.pos >= len(r.data) {
			return 0, fmt.Errorf("unexpected EOF reading VarInt")
		}
		b := r.data[r.pos]
		r.pos++
		result |= int32(b&0x7F) << shift
		shift += 7
		if (b & 0x80) == 0 {
			break
		}
		if shift >= 32 {
			return 0, fmt.Errorf("VarInt too large")
		}
	}
	return result, nil
}

// ReadVarLong 读取 VarLong 编码的整数
func (r *Reader) ReadVarLong() (int64, error) {
	var result int64
	var shift uint
	for {
		if r.pos >= len(r.data) {
			return 0, fmt.Errorf("unexpected EOF reading VarLong")
		}
		b := r.data[r.pos]
		r.pos++
		result |= int64(b&0x7F) << shift
		shift += 7
		if (b & 0x80) == 0 {
			break
		}
		if shift >= 64 {
			return 0, fmt.Errorf("VarLong too large")
		}
	}
	return result, nil
}

// ReadString 读取字符串（带长度前缀）
func (r *Reader) ReadString() (string, error) {
	length, err := r.ReadVarInt()
	if err != nil {
		return "", err
	}
	if length < 0 {
		return "", fmt.Errorf("negative string length: %d", length)
	}
	if int(length) > constants.MaxPacketSize {
		return "", fmt.Errorf("string length exceeds max: %d", length)
	}
	if r.pos+int(length) > len(r.data) {
		return "", fmt.Errorf("insufficient data for string of length %d", length)
	}
	str := string(r.data[r.pos : r.pos+int(length)])
	r.pos += int(length)
	return str, nil
}

// ReadBool 读取布尔值
func (r *Reader) ReadBool() (bool, error) {
	if r.pos >= len(r.data) {
		return false, fmt.Errorf("unexpected EOF reading bool")
	}
	b := r.data[r.pos]
	r.pos++
	return b != 0, nil
}

// ReadByte 读取单个字节
func (r *Reader) ReadByte() (byte, error) {
	if r.pos >= len(r.data) {
		return 0, fmt.Errorf("unexpected EOF reading byte")
	}
	b := r.data[r.pos]
	r.pos++
	return b, nil
}

// ReadInt16 读取 int16
func (r *Reader) ReadInt16() (int16, error) {
	if r.pos+2 > len(r.data) {
		return 0, fmt.Errorf("unexpected EOF reading int16")
	}
	v := binary.BigEndian.Uint16(r.data[r.pos : r.pos+2])
	r.pos += 2
	return int16(v), nil
}

// ReadInt32 读取 int32
func (r *Reader) ReadInt32() (int32, error) {
	if r.pos+4 > len(r.data) {
		return 0, fmt.Errorf("unexpected EOF reading int32")
	}
	v := binary.BigEndian.Uint32(r.data[r.pos : r.pos+4])
	r.pos += 4
	return int32(v), nil
}

// ReadInt64 读取 int64
func (r *Reader) ReadInt64() (int64, error) {
	if r.pos+8 > len(r.data) {
		return 0, fmt.Errorf("unexpected EOF reading int64")
	}
	v := binary.BigEndian.Uint64(r.data[r.pos : r.pos+8])
	r.pos += 8
	return int64(v), nil
}

// ReadFloat32 读取 float32
func (r *Reader) ReadFloat32() (float32, error) {
	if r.pos+4 > len(r.data) {
		return 0, fmt.Errorf("unexpected EOF reading float32")
	}
	v := binary.BigEndian.Uint32(r.data[r.pos : r.pos+4])
	r.pos += 4
	return math.Float32frombits(v), nil
}

// ReadFloat64 读取 float64
func (r *Reader) ReadFloat64() (float64, error) {
	if r.pos+8 > len(r.data) {
		return 0, fmt.Errorf("unexpected EOF reading float64")
	}
	v := binary.BigEndian.Uint64(r.data[r.pos : r.pos+8])
	r.pos += 8
	return math.Float64frombits(v), nil
}

// ReadBytes 读取指定长度的字节
func (r *Reader) ReadBytes(n int) ([]byte, error) {
	if n < 0 {
		return nil, fmt.Errorf("negative read length: %d", n)
	}
	if n > constants.MaxPacketSize {
		return nil, fmt.Errorf("read length exceeds max: %d", n)
	}
	if r.pos+n > len(r.data) {
		return nil, fmt.Errorf("insufficient data: need %d, have %d", n, len(r.data)-r.pos)
	}
	b := make([]byte, n)
	copy(b, r.data[r.pos:r.pos+n])
	r.pos += n
	return b, nil
}

// ReadUUID 读取 UUID（16字节）
func (r *Reader) ReadUUID() ([16]byte, error) {
	var uuid [16]byte
	if r.pos+16 > len(r.data) {
		return uuid, fmt.Errorf("unexpected EOF reading UUID")
	}
	copy(uuid[:], r.data[r.pos:r.pos+16])
	r.pos += 16
	return uuid, nil
}

// ReadPosition 读取位置（64位编码的坐标）
func (r *Reader) ReadPosition() (x, y, z int32, err error) {
	val, err := r.ReadInt64()
	if err != nil {
		return 0, 0, 0, err
	}
	x = int32(val >> 38)
	y = int32((val << 52) >> 52)
	z = int32((val << 26) >> 38)
	return x, y, z, nil
}

// ReadBitSet 读取位集合
func (r *Reader) ReadBitSet() ([]int64, error) {
	length, err := r.ReadVarInt()
	if err != nil {
		return nil, err
	}
	bits := make([]int64, length)
	for i := int32(0); i < length; i++ {
		bits[i], err = r.ReadInt64()
		if err != nil {
			return nil, err
		}
	}
	return bits, nil
}

// ReadRemaining 读取剩余所有数据
func (r *Reader) ReadRemaining() []byte {
	data := r.data[r.pos:]
	r.pos = len(r.data)
	return data
}

// Len 返回剩余字节数
func (r *Reader) Len() int {
	return len(r.data) - r.pos
}

// Position 返回当前位置
func (r *Reader) Position() int {
	return r.pos
}

// Seek 设置读取位置
func (r *Reader) Seek(pos int) error {
	if pos < 0 || pos > len(r.data) {
		return fmt.Errorf("invalid seek position: %d", pos)
	}
	r.pos = pos
	return nil
}
