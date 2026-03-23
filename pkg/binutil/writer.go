package binutil

import (
	"bytes"
	"encoding/binary"
)

// Writer 提供二进制数据写入功能
type Writer struct {
	buf bytes.Buffer
}

// NewWriter 创建新的 Writer
func NewWriter() *Writer {
	return &Writer{}
}

// WriteVarInt 写入 VarInt 编码的整数
func (w *Writer) WriteVarInt(v int32) error {
	for {
		b := byte(v & 0x7F)
		v >>= 7
		if v != 0 {
			b |= 0x80
		}
		w.buf.WriteByte(b)
		if v == 0 {
			break
		}
	}
	return nil
}

// WriteVarLong 写入 VarLong 编码的整数
func (w *Writer) WriteVarLong(v int64) error {
	for {
		b := byte(v & 0x7F)
		v >>= 7
		if v != 0 {
			b |= 0x80
		}
		w.buf.WriteByte(b)
		if v == 0 {
			break
		}
	}
	return nil
}

// WriteString 写入字符串（带长度前缀）
func (w *Writer) WriteString(s string) error {
	if err := w.WriteVarInt(int32(len(s))); err != nil {
		return err
	}
	_, err := w.buf.WriteString(s)
	return err
}

// WriteBool 写入布尔值
func (w *Writer) WriteBool(v bool) error {
	if v {
		return w.buf.WriteByte(1)
	}
	return w.buf.WriteByte(0)
}

// WriteByte 写入单个字节
func (w *Writer) WriteByte(b byte) error {
	return w.buf.WriteByte(b)
}

// WriteInt16 写入 int16
func (w *Writer) WriteInt16(v int16) error {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(v))
	_, err := w.buf.Write(b)
	return err
}

// WriteInt32 写入 int32
func (w *Writer) WriteInt32(v int32) error {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(v))
	_, err := w.buf.Write(b)
	return err
}

// WriteInt64 写入 int64
func (w *Writer) WriteInt64(v int64) error {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	_, err := w.buf.Write(b)
	return err
}

// WriteFloat32 写入 float32
func (w *Writer) WriteFloat32(v float32) error {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(v))
	_, err := w.buf.Write(b)
	return err
}

// WriteFloat64 写入 float64
func (w *Writer) WriteFloat64(v float64) error {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	_, err := w.buf.Write(b)
	return err
}

// WriteBytes 写入字节切片
func (w *Writer) WriteBytes(b []byte) error {
	_, err := w.buf.Write(b)
	return err
}

// WriteUUID 写入 UUID（16字节）
func (w *Writer) WriteUUID(uuid [16]byte) error {
	_, err := w.buf.Write(uuid[:])
	return err
}

// WritePosition 写入位置（64位编码的坐标）
func (w *Writer) WritePosition(x, y, z int32) error {
	val := (int64(x) << 38) | ((int64(z) & 0x3FFFFFF) << 12) | (int64(y) & 0xFFF)
	return w.WriteInt64(val)
}

// Bytes 返回写入的字节数据
func (w *Writer) Bytes() []byte {
	return w.buf.Bytes()
}

// Len 返回当前缓冲区长度
func (w *Writer) Len() int {
	return w.buf.Len()
}

// Reset 清空缓冲区
func (w *Writer) Reset() {
	w.buf.Reset()
}
