package component

// defaultHandlers 返回默认处理器映射
func defaultHandlers() map[int32]ComponentHandler {
	handlers := make(map[int32]ComponentHandler)

	// ID 范围 MinComponentID-MaxComponentID
	for typeID := MinComponentID; typeID <= MaxComponentID; typeID++ {
		if typeID == Container {
			// 容器组件特殊处理
			handlers[typeID] = ContainerComponentHandler
		} else {
			handlers[typeID] = makeDiscardHandler(typeID)
		}
	}

	return handlers
}
