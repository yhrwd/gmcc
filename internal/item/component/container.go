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

	for i := int32(0); i < length; i++ {
		if err := packet.SkipSlotData(r); err != nil {
			logx.Debugf("跳过容器物品槽失败: index=%d, err=%v", i, err)
			// 继续处理，不中断
		}
	}

	return &ComponentResult{TypeID: typeID}, nil
}
