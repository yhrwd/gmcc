package packet

import (
	"bytes"
	"testing"

	"gmcc/internal/constants"
)

func TestReadVarInt(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{"zero", []byte{0x00}, false},
		{"single byte", []byte{0x01}, false},
		{"two bytes", []byte{0x7f}, false},
		{"continuation", []byte{0x80, 0x01}, false},
		{"300", []byte{0xac, 0x02}, false},
		{"max mc varint", []byte{0xff, 0xff, 0xff, 0xff, 0x0f}, false},
		{"negative -1", []byte{0xff, 0xff, 0xff, 0xff, 0x1f}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.input)
			_, err := ReadVarInt(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadVarInt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncodeVarInt(t *testing.T) {
	tests := []struct {
		name  string
		value int32
	}{
		{"zero", 0},
		{"one", 1},
		{"127", 127},
		{"128", 128},
		{"300", 300},
		{"max int32", 2147483647},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeVarInt(tt.value)
			if len(got) == 0 {
				t.Errorf("EncodeVarInt() returned empty")
			}
		})
	}
}

func TestVarIntRoundTrip(t *testing.T) {
	values := []int32{0, 1, 127, 128, 300, 1000, 100000, 2147483647, -1}
	for _, v := range values {
		encoded := EncodeVarInt(v)
		r := bytes.NewReader(encoded)
		decoded, err := ReadVarInt(r)
		if err != nil {
			t.Fatalf("ReadVarInt() error = %v", err)
		}
		if decoded != v {
			t.Errorf("round trip: got %v, want %v", decoded, v)
		}
	}
}

func TestReadBytes(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		n       int
		want    []byte
		wantErr bool
	}{
		{"normal", []byte{0x01, 0x02, 0x03}, 3, []byte{0x01, 0x02, 0x03}, false},
		{"empty", []byte{}, 0, []byte{}, false},
		{"partial", []byte{0x01, 0x02}, 3, nil, true},
		{"negative", []byte{}, -1, nil, true},
		{"exceeds max", []byte{}, constants.MaxPacketSize + 1, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.input)
			got, err := ReadBytes(r, tt.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("ReadBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBool(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  bool
	}{
		{"true", []byte{0x01}, true},
		{"false", []byte{0x00}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.input)
			got, err := ReadBool(r)
			if err != nil {
				t.Errorf("ReadBool() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("ReadBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustReadBytes(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		n       int
		wantLen int
	}{
		{"normal", []byte{0x01, 0x02, 0x03}, 3, 3},
		{"empty", []byte{}, 0, 0},
		{"partial returns empty on error", []byte{0x01, 0x02}, 3, 0},
		{"negative returns empty", []byte{}, -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.input)
			got := MustReadBytes(r, tt.n, "test")
			if len(got) != tt.wantLen {
				t.Errorf("MustReadBytes() len = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}
