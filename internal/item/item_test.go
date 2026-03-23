package item

import (
	"bytes"
	"testing"
)

func TestReadSlotData_WithComponents(t *testing.T) {
	// 测试带组件的物品
	// 一个带VarInt组件的物品
	data := []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x05, 0x00, 0x00}
	// count=1, itemID=1, addCount=1, componentID=1(componentID VarInt), value=5, removeCount=0
	r := bytes.NewReader(data)

	slot, err := ReadSlotData(r)
	if err != nil {
		t.Errorf("ReadSlotData(with component) error = %v", err)
	}
	if slot == nil {
		t.Errorf("ReadSlotData(with component) = nil, want slot")
		return
	}

	if len(slot.Components) != 1 {
		t.Errorf("ReadSlotData(with component) Components len = %d, want 1", len(slot.Components))
	}

	if slot.Components[0] != nil {
		if slot.Components[0].TypeID != 1 {
			t.Errorf("ReadSlotData(with component) Component[0].TypeID = %d, want 1", slot.Components[0].TypeID)
		}
	}
}

func TestReadComponents_Container(t *testing.T) {
	// 跳过容器测试，需要更多数据
	t.Skip("Container component test requires more data")
}
