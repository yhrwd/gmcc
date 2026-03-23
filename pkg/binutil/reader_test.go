package binutil

import (
	"bytes"
	"testing"
)

func TestReader_ReadVarInt(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected int32
		wantErr  bool
	}{
		{"zero", []byte{0x00}, 0, false},
		{"one", []byte{0x01}, 1, false},
		{"127", []byte{0x7F}, 127, false},
		{"128", []byte{0x80, 0x01}, 128, false},
		{"255", []byte{0xFF, 0x01}, 255, false},
		{"max", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x07}, 2147483647, false},
		{"empty", []byte{}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(tt.data)
			got, err := r.ReadVarInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadVarInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ReadVarInt() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestReader_ReadString(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected string
		wantErr  bool
	}{
		{"empty", []byte{0x00}, "", false},
		{"hello", []byte{0x05, 'h', 'e', 'l', 'l', 'o'}, "hello", false},
		{"chinese", []byte{0x06, 0xE4, 0xB8, 0xAD, 0xE6, 0x96, 0x87}, "中文", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(tt.data)
			got, err := r.ReadString()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ReadString() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestReader_ReadBool(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
		wantErr  bool
	}{
		{"true", []byte{0x01}, true, false},
		{"false", []byte{0x00}, false, false},
		{"nonzero", []byte{0xFF}, true, false},
		{"empty", []byte{}, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(tt.data)
			got, err := r.ReadBool()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ReadBool() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestReader_ReadInt32(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected int32
		wantErr  bool
	}{
		{"zero", []byte{0x00, 0x00, 0x00, 0x00}, 0, false},
		{"one", []byte{0x00, 0x00, 0x00, 0x01}, 1, false},
		{"max", []byte{0x7F, 0xFF, 0xFF, 0xFF}, 2147483647, false},
		{"min", []byte{0x80, 0x00, 0x00, 0x00}, -2147483648, false},
		{"empty", []byte{0x00}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(tt.data)
			got, err := r.ReadInt32()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadInt32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ReadInt32() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestReader_ReadFloat32(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected float32
		wantErr  bool
	}{
		{"zero", []byte{0x00, 0x00, 0x00, 0x00}, 0.0, false},
		{"one", []byte{0x3F, 0x80, 0x00, 0x00}, 1.0, false},
		{"minus_one", []byte{0xBF, 0x80, 0x00, 0x00}, -1.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(tt.data)
			got, err := r.ReadFloat32()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFloat32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// 使用近似比较浮点数
			if diff := got - tt.expected; diff < -1e-6 || diff > 1e-6 {
				t.Errorf("ReadFloat32() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestReader_ReadPosition(t *testing.T) {
	// 跳过位置测试，暂时不实现
	t.Skip("Position encoding not fully tested")
}

func TestReader_Len(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	r := NewReader(data)

	if r.Len() != 5 {
		t.Errorf("Len() = %v, want %v", r.Len(), 5)
	}

	r.ReadByte()
	if r.Len() != 4 {
		t.Errorf("Len() after ReadByte = %v, want %v", r.Len(), 4)
	}

	r.ReadInt32()
	if r.Len() != 0 {
		t.Errorf("Len() after ReadInt32 = %v, want %v", r.Len(), 0)
	}
}

func TestNewReaderFromBytesReader(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	br := bytes.NewReader(data)
	r := NewReaderFromBytesReader(br)

	v, err := r.ReadByte()
	if err != nil {
		t.Errorf("ReadByte() error = %v", err)
		return
	}
	if v != 0x01 {
		t.Errorf("ReadByte() = %v, want %v", v, 0x01)
	}
}
