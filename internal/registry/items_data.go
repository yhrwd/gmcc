package registry

type itemEntry struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	StackSize   int32  `json:"stackSize"`
}

func (r *ItemRegistry) loadItems() {
	items := getItemsData()
	for i := range items {
		item := &items[i]
		r.byID[item.ID] = &ItemInfo{
			ID:          item.ID,
			Name:        item.Name,
			DisplayName: item.DisplayName,
			StackSize:   item.StackSize,
		}
		r.byName[item.Name] = r.byID[item.ID]
	}
}
