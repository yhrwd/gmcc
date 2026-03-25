package component

// defaultHandlers 返回默认处理器映射
func defaultHandlers() map[int32]ComponentHandler {
	handlers := make(map[int32]ComponentHandler)

	// P0: 基础类型 - VarInt
	handlers[MaxStackSize] = ParseMaxStackSize
	handlers[MaxDamage] = ParseMaxDamage
	handlers[Damage] = ParseDamage
	handlers[RepairCost] = ParseRepairCost
	handlers[Enchantable] = ParseEnchantable
	handlers[MapID] = ParseMapID

	// P0: 基础类型 - Int32
	handlers[CustomModelData] = ParseCustomModelData
	handlers[DyedColor] = ParseDyedColor
	handlers[MapColor] = ParseMapColor

	// P0: 基础类型 - Bool
	handlers[Unbreakable] = ParseUnbreakable
	handlers[EnchantmentGlintOverride] = ParseEnchantmentGlintOverride

	// P1: 常用显示组件
	handlers[CustomName] = ParseCustomName
	handlers[ItemName] = ParseItemName
	handlers[Lore] = ParseLore
	handlers[Rarity] = ParseRarity
	handlers[Enchantments] = ParseEnchantments

	// ID 范围 MinComponentID-MaxComponentID - 其他使用丢弃处理器
	for typeID := MinComponentID; typeID <= MaxComponentID; typeID++ {
		if _, exists := handlers[typeID]; !exists {
			if typeID == Container {
				handlers[typeID] = ContainerComponentHandler
			} else {
				handlers[typeID] = makeDiscardHandler(typeID)
			}
		}
	}

	return handlers
}
