package stage

import (
	"bytes"
	"encoding/binary"
	"gmcc/pkg/logger"
	"gmcc/pkg/protocol/codec"
	"gmcc/pkg/protocol/connection"
)

type ClientHandshakeState struct {
	Addr string
	Port uint16
}

func (s *ClientHandshakeState) Name() connection.ConnState {
	return connection.StateHandshake
}

// Enter: 这里只做初始化，不发数据
func (s *ClientHandshakeState) Enter(ctx *connection.ConnContext) error {
	logger.Debugf("[HANDSHAKE] Entering handshake state for %s:%d", s.Addr, s.Port)
	// 可以在这里设置压缩阈值为 -1 等初始化工作
	return nil
}

// HandlePacket: 握手状态下，客户端通常不会收到服务器的包，除非是错误
func (s *ClientHandshakeState) HandlePacket(ctx *connection.ConnContext, packet any) error {
	logger.Warnf("[HANDSHAKE] Unexpected packet received: %v", packet)

	payload := &bytes.Buffer{}

	// 1. 构造 Handshake Payload
	codec.WriteVarInt(payload, int32(ctx.ProtocolVersion))
	codec.WriteString(payload, s.Addr) // 注意：通常有专门的 WriteString 方法（VarInt长度 + 字符串）
	binary.Write(payload, binary.BigEndian, s.Port)
	codec.WriteVarInt(payload, 1) // Next State: 1 (Status)

	// 2. 包装 Packet ID (0x00)
	pkt := &bytes.Buffer{}
	codec.WriteVarInt(pkt, 0x00) // Packet ID
	pkt.Write(payload.Bytes())

	// 3. 安全发送 (使用 WriteMu 保护)
	if err := ctx.WriteRaw(pkt.Bytes()); err != nil {
		return err
	}

	logger.Debugf("[HANDSHAKE] Handshake sent successfully")

	// 4. 发送完握手包后，立即切换状态
	return ctx.SM.Switch(connection.StateStatus)
}

func (s *ClientHandshakeState) Exit(ctx *connection.ConnContext) error {
	return nil
}
