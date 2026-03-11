package player

import (
	"sync"
)

type Item struct {
	ID           string
	Count        int32
	Damage       int32
	DisplayName  string
	Lore         []string
	Enchantments map[string]int32
	NBT          map[string]any
}

type Inventory struct {
	mu    sync.RWMutex
	slots map[int8]*Item
}

func NewInventory() *Inventory {
	return &Inventory{
		slots: make(map[int8]*Item),
	}
}

func (i *Inventory) SetSlot(slot int8, item *Item) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if item == nil || item.Count == 0 {
		delete(i.slots, slot)
	} else {
		i.slots[slot] = item
	}
}

func (i *Inventory) GetSlot(slot int8) *Item {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.slots[slot]
}

func (i *Inventory) GetAll() map[int8]*Item {
	i.mu.RLock()
	defer i.mu.RUnlock()
	result := make(map[int8]*Item, len(i.slots))
	for k, v := range i.slots {
		result[k] = v
	}
	return result
}

func (i *Inventory) GetHotbar() []*Item {
	i.mu.RLock()
	defer i.mu.RUnlock()
	hotbar := make([]*Item, 9)
	for j := int8(0); j < 9; j++ {
		hotbar[j] = i.slots[j]
	}
	return hotbar
}

func (i *Inventory) GetMainInventory() []*Item {
	i.mu.RLock()
	defer i.mu.RUnlock()
	main := make([]*Item, 27)
	for j := int8(0); j < 27; j++ {
		main[j] = i.slots[j+9]
	}
	return main
}

func (i *Inventory) GetArmor() []*Item {
	i.mu.RLock()
	defer i.mu.RUnlock()
	armor := make([]*Item, 4)
	for j := int8(0); j < 4; j++ {
		armor[j] = i.slots[j+36]
	}
	return armor
}

func (i *Inventory) GetOffhand() *Item {
	return i.GetSlot(40)
}

func (i *Inventory) Clear() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.slots = make(map[int8]*Item)
}
