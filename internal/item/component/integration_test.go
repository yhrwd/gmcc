package component

import (
	"bytes"
	"testing"
)

func TestRealWorldComponentParsing(t *testing.T) {
	parser := NewParser()

	t.Run("diamond_pickaxe_components", func(t *testing.T) {
		components := []struct {
			typeID int32
			data   []byte
			want   any
		}{
			{MaxStackSize, []byte{0x01}, int32(1)},
			{MaxDamage, []byte{0x98, 0x01}, int32(152)},
			{Unbreakable, []byte{}, true},
			{RepairCost, []byte{0x02}, int32(2)},
		}

		for _, tc := range components {
			t.Run("component_validation", func(t *testing.T) {
				r := bytes.NewReader(tc.data)
				result, err := parser.ParseComponent(tc.typeID, r)

				if err != nil {
					t.Fatalf("ParseComponent(%d) error = %v", tc.typeID, err)
				}

				if result.TypeID != tc.typeID {
					t.Errorf("TypeID = %d, want %d", result.TypeID, tc.typeID)
				}

				if result.Data != tc.want {
					t.Errorf("Data = %v, want %v", result.Data, tc.want)
				}
			})
		}
	})

	t.Run("map_items", func(t *testing.T) {
		mapTests := []struct {
			name   string
			typeID int32
			data   []byte
			want   any
		}{
			{"map_id_1", MapID, []byte{0x01}, int32(1)},
			{"map_id_0", MapID, []byte{0x00}, int32(0)},
			{"map_id_255", MapID, []byte{0xFF, 0x01}, int32(255)},
		}

		for _, tt := range mapTests {
			t.Run(tt.name, func(t *testing.T) {
				r := bytes.NewReader(tt.data)
				result, err := parser.ParseComponent(tt.typeID, r)

				if err != nil {
					t.Fatalf("ParseComponent(%d) error = %v", tt.typeID, err)
				}

				if result.Data != tt.want {
					t.Errorf("Data = %v, want %v", result.Data, tt.want)
				}
			})
		}
	})
}

func TestHandlerRegistration(t *testing.T) {
	p := NewParser()

	t.Run("verify_p0_handlers_registered", func(t *testing.T) {
		p0Components := map[int32]string{
			MaxStackSize:             "max_stack_size",
			MaxDamage:                "max_damage",
			Damage:                   "damage",
			Unbreakable:              "unbreakable",
			CustomModelData:          "custom_model_data",
			RepairCost:               "repair_cost",
			EnchantmentGlintOverride: "enchantment_glint_override",
			Enchantable:              "enchantable",
			DyedColor:                "dyed_color",
			MapColor:                 "map_color",
			MapID:                    "map_id",
		}

		testData := map[int32][]byte{
			MaxStackSize:             {0x01},
			MaxDamage:                {0x01},
			Damage:                   {0x00},
			Unbreakable:              {},
			CustomModelData:          {0x00, 0x00, 0x00, 0x00},
			RepairCost:               {0x00},
			EnchantmentGlintOverride: {0x00},
			Enchantable:              {0x00},
			DyedColor:                {0x00, 0x00, 0x00, 0x00},
			MapColor:                 {0x00, 0x00, 0x00, 0x00},
			MapID:                    {0x00},
		}

		for typeID, name := range p0Components {
			t.Run(name, func(t *testing.T) {
				if _, ok := p.handlers[typeID]; !ok {
					t.Errorf("Handler for %s (ID: %d) not registered", name, typeID)
				}
			})
		}

		t.Run("call_all_handlers", func(t *testing.T) {
			for typeID, data := range testData {
				name := p0Components[typeID]
				t.Run(name, func(t *testing.T) {
					handler, ok := p.handlers[typeID]
					if !ok {
						t.Fatalf("Handler not found for %s (ID: %d)", name, typeID)
					}
					if handler == nil {
						t.Fatalf("Handler is nil for %s (ID: %d)", name, typeID)
					}
					r := bytes.NewReader(data)
					_, err := handler(typeID, r)
					if err != nil {
						t.Errorf("Handler for %s (ID: %d) failed with error: %v", name, typeID, err)
					}
				})
			}
		})
	})
}

func TestComponentResultStructure(t *testing.T) {
	t.Run("verify_result_structure", func(t *testing.T) {
		parser := NewParser()

		data := []byte{0x40}
		r := bytes.NewReader(data)
		result, err := parser.ParseComponent(MaxStackSize, r)

		if err != nil {
			t.Fatalf("ParseComponent failed: %v", err)
		}

		if result.TypeID <= 0 {
			t.Errorf("TypeID should be positive, got %d", result.TypeID)
		}

		if result.Data == nil {
			t.Error("Data should not be nil for parsed components")
		}
	})
}

func TestAllP0ComponentsFunctional(t *testing.T) {
	parser := NewParser()

	testCases := map[int32][]byte{
		MaxStackSize:             {0x40},
		MaxDamage:                {0x98, 0x01},
		Damage:                   {0x05},
		Unbreakable:              {},
		CustomModelData:          {0x00, 0x00, 0x00, 0x01},
		RepairCost:               {0x01},
		EnchantmentGlintOverride: {0x01},
		Enchantable:              {0x0A},
		DyedColor:                {0xFF, 0x55, 0x00, 0x00},
		MapColor:                 {0x00, 0x00, 0xFF, 0x00},
		MapID:                    {0x06},
	}

	for typeID, data := range testCases {
		t.Run(string(rune(typeID)), func(t *testing.T) {
			r := bytes.NewReader(data)
			result, err := parser.ParseComponent(typeID, r)

			if err != nil {
				t.Fatalf("ParseComponent(%d) failed: %v", typeID, err)
			}

			if result == nil {
				t.Fatal("Result should not be nil")
			}

			if result.TypeID != typeID {
				t.Errorf("TypeID = %d, want %d", result.TypeID, typeID)
			}

			if result.Data == nil {
				t.Error("Data should not be nil")
			}
		})
	}
}
