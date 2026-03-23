package component

// ComponentResult 组件解析结果
type ComponentResult struct {
	TypeID int32 // 组件类型ID
	Data   any   // 解析后的数据（当前阶段可能为nil）
}

// ComponentHandler 组件处理器函数类型
type ComponentHandler func(typeID int32, r any) (*ComponentResult, error)

// Parser 组件解析器
type Parser struct {
	handlers map[int32]ComponentHandler
}

// NewParser 创建默认解析器
func NewParser() *Parser {
	return &Parser{
		handlers: defaultHandlers(),
	}
}

// RegisterHandler 注册组件处理器
func (p *Parser) RegisterHandler(typeID int32, handler ComponentHandler) {
	p.handlers[typeID] = handler
}

// ParseComponent 解析单个组件
func (p *Parser) ParseComponent(typeID int32, r any) (*ComponentResult, error) {
	if handler, ok := p.handlers[typeID]; ok {
		return handler(typeID, r)
	}
	// 未知组件：返回仅包含TypeID的结果
	return &ComponentResult{TypeID: typeID}, nil
}
