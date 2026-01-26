package stage

import (
    "gmcc/internal/protocol/connection"
    "gmcc/internal/protocol/codec"
    "gmcc/internal/protocol/packet"
    "gmcc/internal/logger"
    "gmcc/internal/tui"
    "bytes"
    "time"
    "encoding/binary"
    "fmt"
    "io"
)

type ClientStatusState struct{}

func (s *ClientStatusState) Name() connection.ConnState { return connection.StateStatus }

func (s *ClientStatusState) Enter(ctx *connection.ConnContext) error {
    // 使用 ctx.Conn.W (bufio.Writer)
    codec.WriteVarInt(ctx.Conn.W, 1) // Length
    codec.WriteVarInt(ctx.Conn.W, 0) // Packet ID = 0
    ctx.Conn.W.Flush()         // Flush！

    // 2. 读取响应
    // ctx.Conn.R (bufio.Reader) 传给 PacketReader
    reader := packet.NewPacketReader(ctx.Conn.R) 
    
    packet, err := reader.ReadPacket()
    if err != nil {
        return err
    }

    if packet.ID != 0 {
        return fmt.Errorf("unexpected packet id: %d", packet.ID)
    }

    // 解析 JSON
    strLen, err := codec.ReadVarInt(packet.Data)
    if err != nil { return err }
    
    jsonBuf := make([]byte, strLen)
    _, err = io.ReadFull(packet.Data, jsonBuf)
    if err != nil { return err }
    
    logger.Infof("服务器状态: %s", string(jsonBuf))
    
    // ===== Ping / Pong =====
    now := time.Now().UnixMilli()

    buf := bytes.NewBuffer(nil)
    
    // packet data
    codec.WriteVarInt(buf, 1) // packet id = ping
    binary.Write(buf, binary.BigEndian, now)
    
    // packet length
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
        Addr: "忘记在状态机留接口了...",
    })

    return nil
}

func (s *ClientStatusState) HandlePacket(ctx *connection.ConnContext, packet any) error {
	return nil
}

func (s *ClientStatusState) Exit(ctx *connection.ConnContext) error {
	return nil
}
