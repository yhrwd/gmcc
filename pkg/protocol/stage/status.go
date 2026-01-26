package stage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gmcc/internal/tui"
	"gmcc/pkg/logger"
	"gmcc/pkg/protocol/codec"
	"gmcc/pkg/protocol/connection"
	"gmcc/pkg/protocol/packet"
	"io"
	"time"
)

type ClientStatusState struct{}

func (s *ClientStatusState) Name() connection.ConnState { return connection.StateStatus }

func (s *ClientStatusState) Enter(ctx *connection.ConnContext) error {
	codec.WriteVarInt(ctx.Conn.W, 1) // Length
	codec.WriteVarInt(ctx.Conn.W, 0) // Packet ID = 0
	ctx.Conn.W.Flush()

	reader := packet.NewPacketReader(ctx.Conn.R)

	packet, err := reader.ReadPacket()
	if err != nil {
		return err
	}

	if packet.ID != 0 {
		return fmt.Errorf("unexpected packet id: %d", packet.ID)
	}

	strLen, err := codec.ReadVarInt(packet.Data)
	if err != nil {
		return err
	}

	jsonBuf := make([]byte, strLen)
	_, err = io.ReadFull(packet.Data, jsonBuf)
	if err != nil {
		return err
	}

	logger.Infof("服务器状态: %s", string(jsonBuf))

	now := time.Now().UnixMilli()

	buf := bytes.NewBuffer(nil)

	codec.WriteVarInt(buf, 1) // packet id = ping
	binary.Write(buf, binary.BigEndian, now)

	codec.WriteVarInt(ctx.Conn.W, int32(buf.Len()))
	ctx.Conn.W.Write(buf.Bytes())
	ctx.Conn.W.Flush()

	pongPkt, err := reader.ReadPacket()
	if err != nil {
		return err
	}

	if pongPkt.ID != 1 {
		return fmt.Errorf("unexpected pong packet id: %d", pongPkt.ID)
	}

	var pongTime int64
	err = binary.Read(pongPkt.Data, binary.BigEndian, &pongTime)
	if err != nil {
		return err
	}

	latency := time.Now().UnixMilli() - pongTime
	logger.Infof("延迟: %d ms", latency)

	tui.Push(tui.MsgUpdateStatus{
		Latency: fmt.Sprintf("%d ms", latency),
		Addr:    "",
	})

	return nil
}

func (s *ClientStatusState) HandlePacket(ctx *connection.ConnContext, packet any) error {
	return nil
}

func (s *ClientStatusState) Exit(ctx *connection.ConnContext) error {
	return nil
}
