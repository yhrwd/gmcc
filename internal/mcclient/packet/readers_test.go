package packet

import (
	"bytes"
	"testing"
)

func TestReadSlotDataSkipsCustomModelData(t *testing.T) {
	data := make([]byte, 0, 32)
	data = append(data, EncodeVarInt(1)...)  // count
	data = append(data, EncodeVarInt(5)...)  // itemID
	data = append(data, EncodeVarInt(1)...)  // addedComponentCount
	data = append(data, EncodeVarInt(0)...)  // removedComponentCount
	data = append(data, EncodeVarInt(17)...) // custom_model_data
	data = append(data, EncodeVarInt(0)...)  // floats len
	data = append(data, EncodeVarInt(0)...)  // flags len
	data = append(data, EncodeVarInt(0)...)  // strings len
	data = append(data, EncodeVarInt(0)...)  // colors len

	r := bytes.NewReader(data)
	slot, err := ReadSlotData(r)
	if err != nil {
		t.Fatalf("ReadSlotData() error = %v", err)
	}
	if slot == nil || slot.ID != 5 || slot.Count != 1 {
		t.Fatalf("ReadSlotData() = %#v", slot)
	}
	if r.Len() != 0 {
		t.Fatalf("ReadSlotData() 未消费完整数据，剩余 %d 字节", r.Len())
	}
}

func TestReadSlotDataSkipsCanBreakExactMatchers(t *testing.T) {
	data := make([]byte, 0, 64)
	data = append(data, EncodeVarInt(1)...)  // count
	data = append(data, EncodeVarInt(5)...)  // itemID
	data = append(data, EncodeVarInt(1)...)  // addedComponentCount
	data = append(data, EncodeVarInt(0)...)  // removedComponentCount
	data = append(data, EncodeVarInt(15)...) // can_break
	data = append(data, EncodeVarInt(1)...)  // predicate count
	data = append(data, 0x00)                // no blockSet
	data = append(data, 0x00)                // no properties
	data = append(data, 0x00)                // no NBT
	data = append(data, EncodeVarInt(1)...)  // exactMatchers len
	data = append(data, EncodeVarInt(17)...) // custom_model_data
	data = append(data, EncodeVarInt(0)...)  // floats len
	data = append(data, EncodeVarInt(0)...)  // flags len
	data = append(data, EncodeVarInt(0)...)  // strings len
	data = append(data, EncodeVarInt(0)...)  // colors len
	data = append(data, EncodeVarInt(1)...)  // partialMatchers len
	data = append(data, EncodeVarInt(31)...) // enchantable

	r := bytes.NewReader(data)
	slot, err := ReadSlotData(r)
	if err != nil {
		t.Fatalf("ReadSlotData() error = %v", err)
	}
	if slot == nil || slot.ID != 5 || slot.Count != 1 {
		t.Fatalf("ReadSlotData() = %#v", slot)
	}
	if r.Len() != 0 {
		t.Fatalf("ReadSlotData() 未消费完整数据，剩余 %d 字节", r.Len())
	}
}

func TestReadSlotDataSkipsRepairableDirectIDSet(t *testing.T) {
	data := make([]byte, 0, 32)
	data = append(data, EncodeVarInt(1)...)   // count
	data = append(data, EncodeVarInt(962)...) // itemID: chainmail_helmet
	data = append(data, EncodeVarInt(2)...)   // addedComponentCount
	data = append(data, EncodeVarInt(0)...)   // removedComponentCount
	data = append(data, EncodeVarInt(33)...)  // repairable
	data = append(data, EncodeVarInt(1)...)   // items count
	data = append(data, EncodeVarInt(886)...) // iron_ingot-like item id
	data = append(data, EncodeVarInt(31)...)  // enchantable
	data = append(data, EncodeVarInt(10)...)  // enchantable value

	r := bytes.NewReader(data)
	slot, err := ReadSlotData(r)
	if err != nil {
		t.Fatalf("ReadSlotData() error = %v", err)
	}
	if slot == nil || slot.ID != 962 || slot.Count != 1 {
		t.Fatalf("ReadSlotData() = %#v", slot)
	}
	if r.Len() != 0 {
		t.Fatalf("ReadSlotData() 未消费完整数据，剩余 %d 字节", r.Len())
	}
}

func TestReadSlotDataSkipsRepairableNamedIDSet(t *testing.T) {
	tag := "#minecraft:repairs_chain_armor"

	data := make([]byte, 0, 64)
	data = append(data, EncodeVarInt(1)...)   // count
	data = append(data, EncodeVarInt(962)...) // itemID: chainmail_helmet
	data = append(data, EncodeVarInt(2)...)   // addedComponentCount
	data = append(data, EncodeVarInt(0)...)   // removedComponentCount
	data = append(data, EncodeVarInt(33)...)  // repairable
	data = append(data, EncodeVarInt(int32(len(tag)))...)
	data = append(data, []byte(tag)...)
	data = append(data, EncodeVarInt(31)...) // enchantable
	data = append(data, EncodeVarInt(10)...) // enchantable value

	r := bytes.NewReader(data)
	slot, err := ReadSlotData(r)
	if err != nil {
		t.Fatalf("ReadSlotData() error = %v", err)
	}
	if slot == nil || slot.ID != 962 || slot.Count != 1 {
		t.Fatalf("ReadSlotData() = %#v", slot)
	}
	if r.Len() != 0 {
		t.Fatalf("ReadSlotData() 未消费完整数据，剩余 %d 字节", r.Len())
	}
}
