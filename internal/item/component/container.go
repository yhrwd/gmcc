package component

import (
	"bytes"

	"gmcc/internal/logx"
	"gmcc/internal/mcclient/packet"
)

// containerCallback 全局回调变量
var containerCallback func(size int32) error

// SetContainerCallback 注册容器回调
func SetContainerCallback(callback func(size int32) error) {
	containerCallback = callback
}

// ContainerComponentHandler 容器组件特殊处理器
func ContainerComponentHandler(typeID int32, r *bytes.Reader) (*ComponentResult, error) {
	// 1. 读取容器大小
	size, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return nil, err
	}

	// 2. 触发回调（如已注册）
	if containerCallback != nil {
		if err := containerCallback(size); err != nil {
			return nil, err
		}
	}

	// 3. 继续读取容器内容
	length, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return nil, err
	}

	// 读取容器中的物品槽（使用内部简化版）
	for i := int32(0); i < length; i++ {
		if err := skipSlotDataInternal(r); err != nil {
			logx.Debugf("跳过容器物品槽失败: index=%d, err=%v", i, err)
			// 继续处理，不中断
		}
	}

	return &ComponentResult{TypeID: typeID}, nil
}

// skipSlotDataInternal 内部使用的简化版物品槽跳过函数
// 避免与 optimized.go 产生初始化循环
func skipSlotDataInternal(r *bytes.Reader) error {
	// 读取数量
	count, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	if count == 0 {
		return nil // 空物品
	}

	// 读取物品ID
	_, err = packet.ReadVarIntFromReader(r)
	if err != nil {
		return err
	}

	// 跳过组件
	return skipComponentsInternal(r)
}

// skipComponentsInternal 内部使用的组件跳过函数
func skipComponentsInternal(r *bytes.Reader) error {
	// 读取添加的组件数量
	numAdd, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return err
	}

	// 读取移除的组件数量
	numRemove, err := packet.ReadVarIntFromReader(r)
	if err != nil {
		return err
	}

	// 跳过添加的组件
	for i := int32(0); i < numAdd; i++ {
		// 读取组件类型
		typeID, err := packet.ReadVarIntFromReader(r)
		if err != nil {
			return err
		}
		// 直接创建新解析器（避免初始化循环）
		parser := NewParser()
		_, err = parser.ParseComponent(typeID, r)
		if err != nil {
			return err
		}
	}

	// 跳过移除的组件
	for i := int32(0); i < numRemove; i++ {
		if _, err := packet.ReadVarIntFromReader(r); err != nil {
			return err
		}
	}

	return nil
}
