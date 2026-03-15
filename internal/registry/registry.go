package registry

import (
	"sync"

	"gmcc/internal/i18n"
)

type ItemInfo struct {
	ID          int32
	Name        string
	DisplayName string
	StackSize   int32
}

type ItemRegistry struct {
	mu     sync.RWMutex
	byID   map[int32]*ItemInfo
	byName map[string]*ItemInfo
}

var defaultRegistry *ItemRegistry
var once sync.Once

func GetItemRegistry() *ItemRegistry {
	once.Do(func() {
		defaultRegistry = newItemRegistry()
	})
	return defaultRegistry
}

func newItemRegistry() *ItemRegistry {
	r := &ItemRegistry{
		byID:   make(map[int32]*ItemInfo),
		byName: make(map[string]*ItemInfo),
	}
	r.loadItems()
	return r
}

func (r *ItemRegistry) GetByID(id int32) *ItemInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.byID[id]
}

func (r *ItemRegistry) GetByName(name string) *ItemInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.byName[name]
}

func (r *ItemRegistry) DisplayName(id int32) string {
	if item := r.GetByID(id); item != nil {
		return i18n.ItemName(item.Name)
	}
	return "Unknown"
}

func (r *ItemRegistry) NameToID(name string) int32 {
	if item := r.GetByName(name); item != nil {
		return item.ID
	}
	return -1
}

func (r *ItemRegistry) IDToName(id int32) string {
	if item := r.GetByID(id); item != nil {
		return item.Name
	}
	return "unknown"
}

func (r *ItemRegistry) LocalizedName(id int32) string {
	if item := r.GetByID(id); item != nil {
		translated := i18n.ItemName(item.Name)
		if translated != "item.minecraft."+item.Name && translated != "block.minecraft."+item.Name {
			return translated
		}
		return item.DisplayName
	}
	return "Unknown"
}
