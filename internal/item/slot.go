package item

import (
	"bytes"
	"fmt"

	"gmcc/internal/item/component"
	"gmcc/internal/logx"
	"gmcc/internal/mcclient/packet"
)

// SlotData 表示物品槽中的物品数据
type SlotData struct {
	ID         int32                        // 物品ID
	Count      int32                        // 数量
	Components []*component.ComponentResult // 组件列表
}

// IsEmpty 检查槽位是否为空
func (s *SlotData) IsEmpty() bool {
	return s == nil || s.Count <= 0
}

// ReadSlotData 从 Reader 读取物品槽数据
func ReadSlotData(r *bytes.Reader) (*SlotData, error) {
	// 读取数量 (VarInt)
	count, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, nil // 空物品
	}

	// 读取物品ID (VarInt)
	itemID, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return nil, err
	}

	// 读取组件
	components, err := readComponents(r)
	if err != nil {
		logx.Warnf("Slot组件解析失败: itemID=%d, count=%d, err=%v", itemID, count, err)
		return nil, err
	}

	return &SlotData{
		ID:         itemID,
		Count:      count,
		Components: components,
	}, nil
}

// readComponents 读取物品组件列表
func readComponents(r *bytes.Reader) ([]*component.ComponentResult, error) {
	// 读取添加的组件数量
	numAdd, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("read add component count: %w", err)
	}

	// 读取移除的组件数量
	numRemove, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("read remove component count: %w", err)
	}

	// 优化的批量解析 - 预分配容量，避免动态扩容
	components := make([]*component.ComponentResult, 0, numAdd)

	// 从池获取解析器
	parser := component.Acquire()
	defer component.Release(parser)

	// 解析添加的组件
	for i := int32(0); i < numAdd; i++ {
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

		components = append(components, result)
	}

	// 跳过移除的组件
	for i := int32(0); i < numRemove; i++ {
		if _, err := packet.ReadVarIntFromReader(r); err != nil {
			return nil, fmt.Errorf("skip removed component %d: %w", i, err)
		}
	}

	return components, nil
}
