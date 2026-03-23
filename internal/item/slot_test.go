package item

import (
	"bytes"
	"testing"
)

func TestReadSlotData(t *testing.T) {
	// 测试空物品
	data := []byte{0x00} // count = 0
	r := bytes.NewReader(data)

	slot, err := ReadSlotData(r)
	if err != nil {
		t.Errorf("ReadSlotData(empty) error = %v", err)
	}
	if slot != nil {
		t.Errorf("ReadSlotData(empty) = %v, want nil", slot)
	}

	// 测试非空物品（简单物品）
	data = []byte{0x01, 0x01, 0x00, 0x00} // count=1, itemID=1, addCount=0, removeCount=0
	r = bytes.NewReader(data)

	slot, err = ReadSlotData(r)
	if err != nil {
		t.Errorf("ReadSlotData(simple) error = %v", err)
	}
	if slot == nil {
		t.Errorf("ReadSlotData(simple) = nil, want slot")
		return
	}
	if slot.ID != 1 || slot.Count != 1 {
		t.Errorf("ReadSlotData(simple) = {ID:%d, Count:%d}, want {ID:1, Count:1}", slot.ID, slot.Count)
	}
	if len(slot.Components) != 0 {
		t.Errorf("ReadSlotData(simple) Components len = %d, want 0", len(slot.Components))
	}
}

func TestSlotData_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		slot *SlotData
		want bool
	}{
		{"nil", nil, true},
		{"zero_count", &SlotData{Count: 0}, true},
		{"negative_count", &SlotData{Count: -1}, true},
		{"valid", &SlotData{Count: 1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.slot.IsEmpty(); got != tt.want {
				t.Errorf("SlotData.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadComponents(t *testing.T) {
	// 测试无组件
	data := []byte{0x00, 0x00} // addCount=0, removeCount=0
	r := bytes.NewReader(data)

	components, err := readComponents(r)
	if err != nil {
		t.Errorf("readComponents() error = %v", err)
	}
	if components == nil {
		t.Errorf("readComponents() = nil, want non-nil slice")
	}
	if len(components) != 0 {
		t.Errorf("readComponents() len = %d, want 0", len(components))
	}
}
