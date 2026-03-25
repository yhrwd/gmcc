package component

import (
	"bytes"
	"testing"
)

func TestParseMaxStackSize(t *testing.T) {
	tests := []struct {
		name  string
		data  []byte
		want  int32
		valid bool
	}{
		{"default stack", []byte{0x40}, 64, true},
		{"stack 16", []byte{0x10}, 16, true},
		{"stack 1", []byte{0x01}, 1, true},
		{"empty", []byte{}, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.data)
			result, err := ParseMaxStackSize(MaxStackSize, r)

			if !tt.valid {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseMaxStackSize() error = %v", err)
			}

			if result.TypeID != MaxStackSize {
				t.Errorf("TypeID = %d, want %d", result.TypeID, MaxStackSize)
			}

			if result.Data != tt.want {
				t.Errorf("Data = %v, want %v", result.Data, tt.want)
			}
		})
	}
}

func TestParseMaxDamage(t *testing.T) {
	tests := []struct {
		name  string
		data  []byte
		want  int32
		valid bool
	}{
		{"diamond tool", []byte{0x40}, 64, true},
		{"iron tool", []byte{0xF8, 0x01}, 248, true},
		{"zero", []byte{0x00}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.data)
			result, err := ParseMaxDamage(MaxDamage, r)

			if !tt.valid {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseMaxDamage() error = %v", err)
			}

			if result.TypeID != MaxDamage {
				t.Errorf("TypeID = %d, want %d", result.TypeID, MaxDamage)
			}

			if result.Data != tt.want {
				t.Errorf("Data = %v, want %v", result.Data, tt.want)
			}
		})
	}
}

func TestParseDamage(t *testing.T) {
	tests := []struct {
		name  string
		data  []byte
		want  int32
		valid bool
	}{
		{"slightly damaged", []byte{0x05}, 5, true},
		{"half damaged", []byte{0x20}, 32, true},
		{"fully damaged", []byte{0x3F}, 63, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.data)
			result, err := ParseDamage(Damage, r)

			if !tt.valid {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseDamage() error = %v", err)
			}

			if result.TypeID != Damage {
				t.Errorf("TypeID = %d, want %d", result.TypeID, Damage)
			}

			if result.Data != tt.want {
				t.Errorf("Data = %v, want %v", result.Data, tt.want)
			}
		})
	}
}

func TestParseRepairCost(t *testing.T) {
	tests := []struct {
		name  string
		data  []byte
		want  int32
		valid bool
	}{
		{"no repairs", []byte{0x00}, 0, true},
		{"one repair", []byte{0x01}, 1, true},
		{"five repairs", []byte{0x05}, 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.data)
			result, err := ParseRepairCost(RepairCost, r)

			if !tt.valid {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseRepairCost() error = %v", err)
			}

			if result.TypeID != RepairCost {
				t.Errorf("TypeID = %d, want %d", result.TypeID, RepairCost)
			}

			if result.Data != tt.want {
				t.Errorf("Data = %v, want %v", result.Data, tt.want)
			}
		})
	}
}

func TestParseEnchantable(t *testing.T) {
	tests := []struct {
		name  string
		data  []byte
		want  int32
		valid bool
	}{
		{"no enchantment", []byte{0x00}, 0, true},
		{"iron tools", []byte{0x0E}, 14, true},
		{"diamond tools", []byte{0x0A}, 10, true},
		{"gold tools", []byte{0x16}, 22, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.data)
			result, err := ParseEnchantable(Enchantable, r)

			if !tt.valid {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseEnchantable() error = %v", err)
			}

			if result.TypeID != Enchantable {
				t.Errorf("TypeID = %d, want %d", result.TypeID, Enchantable)
			}

			if result.Data != tt.want {
				t.Errorf("Data = %v, want %v", result.Data, tt.want)
			}
		})
	}
}

func TestParseMapID(t *testing.T) {
	tests := []struct {
		name  string
		data  []byte
		want  int32
		valid bool
	}{
		{"map 0", []byte{0x00}, 0, true},
		{"map 1", []byte{0x01}, 1, true},
		{"map 100", []byte{0x64}, 100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.data)
			result, err := ParseMapID(MapID, r)

			if !tt.valid {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseMapID() error = %v", err)
			}

			if result.TypeID != MapID {
				t.Errorf("TypeID = %d, want %d", result.TypeID, MapID)
			}

			if result.Data != tt.want {
				t.Errorf("Data = %v, want %v", result.Data, tt.want)
			}
		})
	}
}

func TestParseCustomModelData(t *testing.T) {
	tests := []struct {
		name  string
		data  []byte
		want  int32
		valid bool
	}{
		{"no variant", []byte{0x00, 0x00, 0x00, 0x00}, 0, true},
		{"variant 1", []byte{0x00, 0x00, 0x00, 0x01}, 1, true},
		{"variant -1", []byte{0xFF, 0xFF, 0xFF, 0xFF}, -1, true},
		{"high variant", []byte{0x12, 0x34, 0x56, 0x78}, 305419896, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.data)
			result, err := ParseCustomModelData(CustomModelData, r)

			if !tt.valid {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseCustomModelData() error = %v", err)
			}

			if result.TypeID != CustomModelData {
				t.Errorf("TypeID = %d, want %d", result.TypeID, CustomModelData)
			}

			if result.Data != tt.want {
				t.Errorf("Data = %v, want %v", result.Data, tt.want)
			}
		})
	}
}

func TestParseDyedColor(t *testing.T) {
	tests := []struct {
		name  string
		data  []byte
		want  int32
		valid bool
	}{
		{"red", []byte{0x55, 0xAA, 0xFF, 0x00}, 1437269760, true},
		{"blue", []byte{0x00, 0x00, 0xFF, 0x00}, 65280, true},
		{"green", []byte{0x00, 0x80, 0x00, 0x00}, 8388608, true},
		{"white", []byte{0x00, 0xFF, 0xFF, 0xFF}, 16777215, true},
		{"black", []byte{0x00, 0x00, 0x00, 0x00}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.data)
			result, err := ParseDyedColor(DyedColor, r)

			if !tt.valid {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseDyedColor() error = %v", err)
			}

			if result.TypeID != DyedColor {
				t.Errorf("TypeID = %d, want %d", result.TypeID, DyedColor)
			}

			if result.Data != tt.want {
				t.Errorf("Data = %v, want %v", result.Data, tt.want)
			}
		})
	}
}

func TestParseMapColor(t *testing.T) {
	tests := []struct {
		name  string
		data  []byte
		want  int32
		valid bool
	}{
		{"default", []byte{0x00, 0x00, 0x00, 0x00}, 0, true},
		{"white", []byte{0xFF, 0xFF, 0xFF, 0x00}, -256, true},
		{"black", []byte{0x00, 0x00, 0x00, 0x00}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.data)
			result, err := ParseMapColor(MapColor, r)

			if !tt.valid {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseMapColor() error = %v", err)
			}

			if result.TypeID != MapColor {
				t.Errorf("TypeID = %d, want %d", result.TypeID, MapColor)
			}

			if result.Data != tt.want {
				t.Errorf("Data = %v, want %v", result.Data, tt.want)
			}
		})
	}
}

func TestParseUnbreakable(t *testing.T) {
	r := bytes.NewReader([]byte{})
	result, err := ParseUnbreakable(Unbreakable, r)

	if err != nil {
		t.Fatalf("ParseUnbreakable() error = %v", err)
	}

	if result.TypeID != Unbreakable {
		t.Errorf("TypeID = %d, want %d", result.TypeID, Unbreakable)
	}

	if result.Data != true {
		t.Errorf("Data = %v, want true", result.Data)
	}
}

func TestParseEnchantmentGlintOverride(t *testing.T) {
	tests := []struct {
		name  string
		data  byte
		want  bool
		valid bool
	}{
		{"enable glint", 0x01, true, true},
		{"disable glint", 0x00, false, true},
		{"explicit false", 0x00, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader([]byte{tt.data})
			result, err := ParseEnchantmentGlintOverride(EnchantmentGlintOverride, r)

			if !tt.valid {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseEnchantmentGlintOverride() error = %v", err)
			}

			if result.TypeID != EnchantmentGlintOverride {
				t.Errorf("TypeID = %d, want %d", result.TypeID, EnchantmentGlintOverride)
			}

			if result.Data != tt.want {
				t.Errorf("Data = %v, want %v", result.Data, tt.want)
			}
		})
	}
}

func TestGenericParsers(t *testing.T) {
	t.Run("ParseInt32", func(t *testing.T) {
		data := []byte{0x00, 0x00, 0x00, 0x01}
		r := bytes.NewReader(data)
		result, err := ParseInt32(99, r)

		if err != nil {
			t.Fatalf("ParseInt32() error = %v", err)
		}

		if result.Data != int32(1) {
			t.Errorf("Data = %v, want 1", result.Data)
		}
	})

	t.Run("ParseBool", func(t *testing.T) {
		data := []byte{0x01}
		r := bytes.NewReader(data)
		result, err := ParseBool(100, r)

		if err != nil {
			t.Fatalf("ParseBool() error = %v", err)
		}

		if result.TypeID != 100 {
			t.Errorf("TypeID = %d, want 100", result.TypeID)
		}

		if result.Data != true {
			t.Errorf("Data = %v, want true", result.Data)
		}
	})
}

func TestParserIntegration(t *testing.T) {
	parser := NewParser()

	components := []struct {
		typeID int32
		data   []byte
		want   any
	}{
		{MaxStackSize, []byte{0x40}, int32(64)},
		{MaxDamage, []byte{0xF8, 0x01}, int32(248)},
		{Damage, []byte{0x05}, int32(5)},
		{Unbreakable, []byte{}, true},
		{CustomModelData, []byte{0x00, 0x00, 0x00, 0x01}, int32(1)},
		{RepairCost, []byte{0x03}, int32(3)},
		{EnchantmentGlintOverride, []byte{0x01}, true},
		{Enchantable, []byte{0x0A}, int32(10)},
		{DyedColor, []byte{0xFF, 0x00, 0x00, 0x00}, int32(-16777216)},
		{MapColor, []byte{0x00, 0x00, 0xFF, 0x00}, int32(65280)},
		{MapID, []byte{0x01}, int32(1)},
	}

	for _, tc := range components {
		t.Run("TypeID_"+string(rune(tc.typeID)), func(t *testing.T) {
			r := bytes.NewReader(tc.data)
			result, err := parser.ParseComponent(tc.typeID, r)

			if err != nil {
				t.Fatalf("ParseComponent(%d) error = %v", tc.typeID, err)
			}

			if result.TypeID != tc.typeID {
				t.Errorf("TypeID = %d, want %d", result.TypeID, tc.typeID)
			}

			if result.Data != tc.want {
				t.Errorf("Data = %v, want %v, type: %T vs %T", result.Data, tc.want, result.Data, tc.want)
			} else if result.Data != tc.want {
				t.Logf("Data matches but comparison failed: %v == %v", result.Data, tc.want)
			}
		})
	}
}
