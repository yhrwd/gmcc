package packet

import (
	"bytes"
	"gmcc/pkg/protocol/codec"
	"io"
)

// Packet 代表一个 Minecraft 数据包
type Packet struct {
	ID   int32         // Packet ID
	Data *bytes.Reader // Payload 数据的读取器
}

// PacketReader 负责从连接中读取 Minecraft 数据包
type PacketReader struct {
	r io.Reader
}

// NewPacketReader 创建一个新的 PacketReader
func NewPacketReader(r io.Reader) *PacketReader {
	return &PacketReader{
		r: r,
	}
}

// ReadPacket 读取一个完整的数据包
// 返回: *Packet, error
func (pr *PacketReader) ReadPacket() (*Packet, error) {
	// 1. 读取数据包总长度 (VarInt)
	length, err := codec.ReadVarInt(pr.r)
	if err != nil {
		return nil, err
	}

	// 2. 读取 Packet ID + Payload 的全部字节
	// 这里必须用 ReadFull，确保读取到完整的 length 字节
	packetData := make([]byte, length)
	_, err = io.ReadFull(pr.r, packetData)
	if err != nil {
		return nil, err
	}

	// 3. 使用 bytes.Reader 来解析 Packet ID
	r := bytes.NewReader(packetData)
	packetID, err := codec.ReadVarInt(r)
	if err != nil {
		return nil, err
	}

	// 4. 返回 Packet 对象
	// 此时 r 已经指向了 Payload 的起始位置，且 r.Len() 就是 Payload 的长度
	return &Packet{
		ID:   packetID,
		Data: r,
	}, nil
}
