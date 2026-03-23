package component

// defaultHandlers 返回默认处理器映射
func defaultHandlers() map[int32]ComponentHandler {
	handlers := make(map[int32]ComponentHandler)

	// ID 范围 0-103
	for typeID := int32(0); typeID <= 103; typeID++ {
		if typeID == 73 {
			// 容器组件特殊处理
			handlers[typeID] = ContainerComponentHandler
		} else {
			handlers[typeID] = makeDiscardHandler(typeID)
		}
	}

	return handlers
}
