package component

import (
	"bytes"
	"fmt"
	"sync"

	"gmcc/internal/mcclient/packet"
)

// parserPool 解析器对象池
var parserPool = sync.Pool{
	New: func() any {
		return NewParser()
	},
}

// Acquire 从池中获取 Parser
func Acquire() *Parser {
	return parserPool.Get().(*Parser)
}

// Release 将 Parser 归还到池中
func Release(p *Parser) {
	parserPool.Put(p)
}

// BatchParse 批量解析组件列表（性能优化版本）
func BatchParse(r *bytes.Reader) ([]*ComponentResult, error) {
	// 读取组件数量
	count, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return nil, err
	}

	// 预分配容量，避免动态扩容
	results := make([]*ComponentResult, 0, int(count))

	// 从池获取解析器
	parser := Acquire()
	defer Release(parser)

	// 解析组件
	for i := int32(0); i < count; i++ {
		// 读取组件类型ID
		typeID, err := packet.ReadVarIntFromReader(r)
		if err != nil {
			return nil, fmt.Errorf("read component type %d: %w", i, err)
		}

		// 解析组件
		result, err := parser.ParseComponent(typeID, r)
		if err != nil {
			return nil, fmt.Errorf("parse component %d: %w", typeID, err)
		}

		results = append(results, result)
	}

	return results, nil
}
