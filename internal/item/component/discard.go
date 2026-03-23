package component

// makeDiscardHandler 创建丢弃处理器
func makeDiscardHandler(typeID int32) ComponentHandler {
	// 当前阶段只返回包含TypeID的结果，跳过逻辑在slot.go中实现
	return func(id int32, r any) (*ComponentResult, error) {
		return &ComponentResult{TypeID: id}, nil
	}
}
