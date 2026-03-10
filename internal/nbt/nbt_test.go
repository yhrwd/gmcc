package nbt

import (
	"bytes"
	"reflect"
	"testing"
)

func TestMarshalUnmarshal_Compound(t *testing.T) {
	type Player struct {
		Name  string `nbt:"name"`
		HP    int32  `nbt:"Health"`
		Level int16  `nbt:"Level"`
	}

	player := Player{
		Name:  "Steve",
		HP:    20,
		Level: 100,
	}

	data, err := Marshal(&player)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var result Player
	if err := Unmarshal(data, &result); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if result.Name != player.Name {
		t.Errorf("Name: got %q, want %q", result.Name, player.Name)
	}
	if result.HP != player.HP {
		t.Errorf("HP: got %d, want %d", result.HP, player.HP)
	}
	if result.Level != player.Level {
		t.Errorf("Level: got %d, want %d", result.Level, player.Level)
	}
}

func TestMarshalUnmarshal_Types(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{"byte", int8(127)},
		{"short", int16(32767)},
		{"int", int32(2147483647)},
		{"long", int64(9223372036854775807)},
		{"float", float32(3.14159)},
		{"double", float64(3.141592653589793)},
		{"string", "Hello, World!"},
		{"byte_array", []byte{1, 2, 3, 4, 5}},
		{"int_array", []int32{100, 200, 300}},
		{"long_array", []int64{1000, 2000, 3000}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.value)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			if len(data) == 0 {
				t.Fatal("Marshal returned empty data")
			}
		})
	}
}

func TestSNBTParse_Compound(t *testing.T) {
	snbt := `{name:"Steve",Health:20,Level:100}`

	result, err := ParseSNBT(snbt)
	if err != nil {
		t.Fatalf("ParseSNBT failed: %v", err)
	}

	m, ok := result.(map[string]any)
	if !ok {
		t.Fatal("Expected map[string]any")
	}

	if m["name"] != "Steve" {
		t.Errorf("name: got %v, want Steve", m["name"])
	}
}

func TestSNBTParse_List(t *testing.T) {
	snbt := `[1,2,3,4,5]`

	result, err := ParseSNBT(snbt)
	if err != nil {
		t.Fatalf("ParseSNBT failed: %v", err)
	}

	list, ok := result.([]any)
	if !ok {
		t.Fatal("Expected []any")
	}

	if len(list) != 5 {
		t.Errorf("list length: got %d, want 5", len(list))
	}
}

func TestSNBTParse_ByteArray(t *testing.T) {
	snbt := `[B;1b,2b,3b]`

	result, err := ParseSNBT(snbt)
	if err != nil {
		t.Fatalf("ParseSNBT failed: %v", err)
	}

	list, ok := result.([]any)
	if !ok {
		t.Fatal("Expected []any")
	}

	if len(list) != 3 {
		t.Errorf("list length: got %d, want 3", len(list))
	}
}

func TestFormatSNBT(t *testing.T) {
	m := map[string]any{
		"name":   "Steve",
		"Health": int32(20),
		"Level":  int16(100),
	}

	snbt := FormatSNBT(m)
	if snbt == "" {
		t.Fatal("FormatSNBT returned empty string")
	}

	// Round-trip test
	result, err := ParseSNBT(snbt)
	if err != nil {
		t.Fatalf("ParseSNBT failed: %v", err)
	}

	m2, ok := result.(map[string]any)
	if !ok {
		t.Fatal("Expected map[string]any")
	}

	if m2["name"] != "Steve" {
		t.Errorf("name: got %v, want Steve", m2["name"])
	}
}

func TestCESU8(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"ASCII", []byte("Hello"), "Hello"},
		{"Chinese", []byte{0xE4, 0xB8, 0xAD, 0xE6, 0x96, 0x87}, "中文"},
		{"Emoji_CESU8", []byte{0xED, 0xA0, 0xBD, 0xED, 0xB8, 0x80}, "😀"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal with CESU-8 encoded string
			data := append([]byte{0x08, 0x00, byte(len(tt.input))}, tt.input...)

			// Unmarshal
			decoder := NewDecoder(bytes.NewReader(data))
			decoder.NetworkFormat(true)
			var result string
			if err := decoder.Decode(&result); err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

// decodeAny decodes any NBT value
func decodeAny(r *bytes.Reader) any {
	d := NewDecoder(r)
	d.NetworkFormat(true)
	_, _ = d.readByte() // tag type
	var result any
	_ = d.unmarshal(reflect.ValueOf(&result).Elem(), TagEnd)
	return result
}
