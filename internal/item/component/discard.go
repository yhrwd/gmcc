package component

import (
	"bytes"

	"gmcc/internal/logx"
	"gmcc/internal/mcclient/packet"
)

// makeDiscardHandler 创建丢弃处理器
func makeDiscardHandler(typeID int32) ComponentHandler {
	return func(id int32, r *bytes.Reader) (*ComponentResult, error) {
		// 根据组件类型选择跳过方式
		var err error
		switch id {
		case 0, 6, 9, 22, 35, 45, 55, 57, 64, 76, 77, 93, 95, 96, 97, 98, 99:
			// NBT 组件
			err = packet.SkipNBT(r)
		case 1, 2, 3, 7, 12, 19, 44, 61, 91, 102, 103:
			// VarInt 组件
			_, err = packet.ReadVarIntFromReader(r)
		case 4, 20:
			// 无数据组件
		case 10, 63, 69:
			// String 组件
			_, err = packet.ReadStringFromReader(r)
		case 21:
			// Bool 组件
			_, err = packet.ReadBoolFromReader(r)
		case 42, 43:
			// Int32 组件
			_, err = packet.ReadInt32FromReader(r)
		default:
			// 其他组件：尝试作为 NBT 跳过
			err = packet.SkipNBT(r)
		}

		// 忽略EOF错误，这些都是非关键数据
		if err != nil && err.Error() != "unexpected EOF" {
			logx.Debugf("跳过组件失败: typeID=%d, err=%v", id, err)
			// 不返回错误，继续处理
		}

		return &ComponentResult{TypeID: id}, nil
	}
}
