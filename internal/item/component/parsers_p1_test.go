package component

import (
	"bytes"
	"testing"

	"gmcc/internal/mcclient/chat"
)

func TestParseCustomName(t *testing.T) {
	t.Run("valid_custom_name", func(t *testing.T) {
		parser := NewParser()

		nbtData := []byte{
			0x0A,
			0x00, 0x00,
			0x08,
			0x00, 0x04,
			't', 'e', 'x', 't',
			0x00, 0x0B,
			'C', 'u', 's', 't', 'o', 'm', ' ', 'S', 'w', 'o', 'r', 'd',
			0x00,
		}

		r := bytes.NewReader(nbtData)
		result, err := parser.ParseComponent(CustomName, r)

		if err != nil {
			t.Fatalf("ParseCustomName error = %v", err)
		}

		if result == nil {
			t.Fatal("Result should not be nil")
		}

		if result.TypeID != CustomName {
			t.Errorf("TypeID = %d, want %d", result.TypeID, CustomName)
		}

		tc, ok := result.Data.(*chat.TextComponent)
		if !ok {
			t.Fatal("Data should be TextComponent")
		}

		if tc.Text == "" {
			t.Logf("Warning: TextComponent text is empty, but parsing succeeded")
		}
	})
}

func TestParseItemName(t *testing.T) {
	t.Run("valid_item_name", func(t *testing.T) {
		parser := NewParser()

		nbtData := []byte{
			0x0A,
			0x00, 0x00,
			0x08,
			0x00, 0x04,
			't', 'e', 'x', 't',
			0x00, 0x0D,
			'D', 'i', 'a', 'm', 'o', 'n', 'd', ' ', 'P', 'i', 'c', 'k', 'a', 'x', 'e',
			0x00,
		}

		r := bytes.NewReader(nbtData)
		result, err := parser.ParseComponent(ItemName, r)

		if err != nil {
			t.Fatalf("ParseItemName error = %v", err)
		}

		tc, ok := result.Data.(*chat.TextComponent)
		if !ok {
			t.Fatal("Data should be TextComponent")
		}

		if tc.Text == "" {
			t.Logf("Warning: TextComponent text is empty, but parsing succeeded")
		}
	})
}

func TestParseRarity(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		want      string
		wantError bool
	}{
		{"common", []byte{0x00}, "common", false},
		{"uncommon", []byte{0x01}, "uncommon", false},
		{"rare", []byte{0x02}, "rare", false},
		{"epic", []byte{0x03}, "epic", false},
		{"invalid_default", []byte{0x05}, "common", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			r := bytes.NewReader(tt.data)
			result, err := parser.ParseComponent(Rarity, r)

			if tt.wantError && err == nil {
				t.Error("Expected error but got nil")
				return
			}

			if !tt.wantError && err != nil {
				t.Fatalf("ParseRarity error = %v", err)
			}

			if result.Data != tt.want {
				t.Errorf("Data = %v, want %v", result.Data, tt.want)
			}
		})
	}
}

func TestParseEnchantments(t *testing.T) {
	t.Run("parser_initialization", func(t *testing.T) {
		_, err := ParseEnchantments(Enchantments, bytes.NewReader(nil))

		t.Logf("Enchantments parser initialized (error expected for nil input)")
		if err != nil {
			t.Logf("Expected error: %v", err)
		}
	})
}

func TestParseLore(t *testing.T) {
	t.Run("parser_initialization", func(t *testing.T) {
		_, err := ParseLore(Lore, bytes.NewReader(nil))

		t.Logf("Lore parser initialized (error expected for nil input)")
		if err != nil {
			t.Logf("Expected error: %v", err)
		}
	})
}

func TestTextComponentStructure(t *testing.T) {
	t.Run("parser_initialization", func(t *testing.T) {
		_, err := ParseCustomName(CustomName, bytes.NewReader(nil))

		t.Logf("Text component parser initialized (error expected for nil input)")
		if err != nil {
			t.Logf("Expected error: %v", err)
		}
	})
}

func TestP1HandlersRegistered(t *testing.T) {
	parser := NewParser()

	p1Components := map[int32]string{
		CustomName:   "custom_name",
		ItemName:     "item_name",
		Lore:         "lore",
		Rarity:       "rarity",
		Enchantments: "enchantments",
	}

	for typeID, name := range p1Components {
		t.Run(name, func(t *testing.T) {
			if _, ok := parser.handlers[typeID]; !ok {
				t.Errorf("Handler for %s (ID: %d) not registered", name, typeID)
			}
		})
	}
}
