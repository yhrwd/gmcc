package packet

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"gmcc/internal/logx"
	"gmcc/internal/mcclient/crypto"
	"gmcc/internal/nbt"
)

type Packet struct {
	ID   int32
	Data []byte
}

const MaxFrameLength = 2 * 1024 * 1024

type PacketConn struct {
	Conn net.Conn
	Br   *bufio.Reader
	W    io.Writer

	Mu                   sync.Mutex
	CompressionThreshold int
}

func NewPacketConn(conn net.Conn) *PacketConn {
	return &PacketConn{
		Conn:                 conn,
		Br:                   bufio.NewReader(conn),
		W:                    conn,
		CompressionThreshold: -1,
	}
}

func (c *PacketConn) Close() error {
	return c.Conn.Close()
}

func (c *PacketConn) SetReadDeadline(t time.Time) error {
	return c.Conn.SetReadDeadline(t)
}

func (c *PacketConn) SetCompressionThreshold(threshold int) {
	c.CompressionThreshold = threshold
}

func (c *PacketConn) EnableEncryption(secret []byte) error {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return fmt.Errorf("create AES cipher failed: %w", err)
	}
	decrypter := crypto.NewCFB8(block, secret, false)
	encrypter := crypto.NewCFB8(block, secret, true)
	c.Br = bufio.NewReader(&cipher.StreamReader{S: decrypter, R: c.Conn})
	c.W = &cipher.StreamWriter{S: encrypter, W: c.Conn}
	logx.Debugf("已启用连接加密: mode=AES/CFB8 secretLen=%d", len(secret))
	return nil
}

func (c *PacketConn) ReadPacket() (Packet, error) {
	frame, err := c.readFrame()
	if err != nil {
		return Packet{}, err
	}

	data, err := c.decompressFrame(frame)
	if err != nil {
		return Packet{}, err
	}

	return ParsePacketData(data)
}

func (c *PacketConn) readFrame() ([]byte, error) {
	frameLen, err := ReadVarInt(c.Br)
	if err != nil {
		logx.Debugf("读取 frame 长度失败: %v", err)
		return nil, err
	}
	if frameLen < 0 || frameLen > MaxFrameLength {
		return nil, fmt.Errorf("invalid frame length %d", frameLen)
	}

	frame := make([]byte, frameLen)
	if _, err := io.ReadFull(c.Br, frame); err != nil {
		logx.Debugf("读取 frame 内容失败: frameLen=%d err=%v", frameLen, err)
		return nil, err
	}
	return frame, nil
}

func (c *PacketConn) decompressFrame(frame []byte) ([]byte, error) {
	if c.CompressionThreshold < 0 {
		logx.PacketLogf("收包 frame: frameLen=%d compressed=false threshold=%d uncompressedLen=%d framePreview=%s",
			len(frame), c.CompressionThreshold, len(frame), RawPreview(frame))
		return frame, nil
	}

	r := bytes.NewReader(frame)
	uncompressedLen, err := ReadVarInt(r)
	if err != nil {
		return nil, fmt.Errorf("read compressed length failed: %w", err)
	}

	rest, _ := io.ReadAll(r)
	var data []byte

	if uncompressedLen == 0 {
		data = rest
	} else {
		zr, err := zlib.NewReader(bytes.NewReader(rest))
		if err != nil {
			return nil, fmt.Errorf("create zlib reader failed: %w", err)
		}
		data, err = io.ReadAll(zr)
		_ = zr.Close()
		if err != nil {
			return nil, fmt.Errorf("decompress packet failed: %w", err)
		}
		if len(data) != int(uncompressedLen) {
			return nil, fmt.Errorf("invalid decompressed size: want %d got %d", uncompressedLen, len(data))
		}
	}

	logx.PacketLogf("收包 frame: frameLen=%d compressed=%t threshold=%d uncompressedLen=%d framePreview=%s",
		len(frame), uncompressedLen != 0, c.CompressionThreshold, len(data), RawPreview(frame))
	return data, nil
}

func ParsePacketData(data []byte) (Packet, error) {
	r := bytes.NewReader(data)
	packetID, err := ReadVarInt(r)
	if err != nil {
		return Packet{}, fmt.Errorf("read packet id failed: %w", err)
	}
	payload, _ := io.ReadAll(r)
	return Packet{ID: packetID, Data: payload}, nil
}

func (c *PacketConn) WritePacket(packetID int32, payload []byte) error {
	body := append(EncodeVarInt(packetID), payload...)
	frame := c.compressBody(body)
	return c.writeFrame(packetID, payload, body, frame)
}

func (c *PacketConn) compressBody(body []byte) []byte {
	if c.CompressionThreshold < 0 {
		return body
	}

	if len(body) >= c.CompressionThreshold {
		var compressed bytes.Buffer
		zw := zlib.NewWriter(&compressed)
		if _, err := zw.Write(body); err != nil {
			logx.Warnf("compress: zlib 写入失败: %v", err)
			_ = zw.Close()
			return append(EncodeVarInt(0), body...)
		}
		if err := zw.Close(); err != nil {
			logx.Warnf("compress: zlib 关闭失败: %v", err)
			return append(EncodeVarInt(0), body...)
		}
		return append(EncodeVarInt(int32(len(body))), compressed.Bytes()...)
	}
	return append(EncodeVarInt(0), body...)
}

func (c *PacketConn) writeFrame(packetID int32, payload, body, frame []byte) error {
	raw := append(EncodeVarInt(int32(len(frame))), frame...)

	logx.PacketLogf("发包: id=0x%02X payloadLen=%d bodyLen=%d frameLen=%d compressed=%t threshold=%d payloadPreview=%s",
		packetID, len(payload), len(body), len(frame),
		c.CompressionThreshold >= 0 && len(body) >= c.CompressionThreshold,
		c.CompressionThreshold, RawPreview(payload))

	c.Mu.Lock()
	defer c.Mu.Unlock()
	_, err := c.W.Write(raw)
	return err
}

func ReadVarInt(r io.ByteReader) (int32, error) {
	var result int32
	var shift uint
	for i := 0; i < 5; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		result |= int32(b&0x7F) << shift
		if b&0x80 == 0 {
			return result, nil
		}
		shift += 7
	}
	return 0, fmt.Errorf("varint too large")
}

func EncodeVarInt(v int32) []byte {
	x := uint32(v)
	out := make([]byte, 0, 5)
	for {
		if x&^uint32(0x7F) == 0 {
			out = append(out, byte(x))
			return out
		}
		out = append(out, byte((x&0x7F)|0x80))
		x >>= 7
	}
}

func ReadBool(r io.Reader) (bool, error) {
	var b [1]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return false, err
	}
	return b[0] != 0, nil
}

func EncodeBool(v bool) []byte {
	if v {
		return []byte{1}
	}
	return []byte{0}
}

func ReadString(r io.ByteReader, rr io.Reader) (string, error) {
	n, err := ReadVarInt(r)
	if err != nil {
		return "", err
	}
	if n < 0 {
		return "", fmt.Errorf("invalid string length %d", n)
	}
	buf := make([]byte, n)
	if _, err := io.ReadFull(rr, buf); err != nil {
		return "", err
	}
	return nbt.CESU8ToUTF8(buf), nil
}

func EncodeString(s string) []byte {
	b := []byte(s)
	out := make([]byte, 0, len(b)+5)
	out = append(out, EncodeVarInt(int32(len(b)))...)
	out = append(out, b...)
	return out
}

func ReadByteArray(r io.ByteReader, rr io.Reader) ([]byte, error) {
	n, err := ReadVarInt(r)
	if err != nil {
		return nil, err
	}
	if n < 0 {
		return nil, fmt.Errorf("invalid byte array length %d", n)
	}
	buf := make([]byte, n)
	if _, err := io.ReadFull(rr, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func EncodeByteArray(b []byte) []byte {
	out := make([]byte, 0, len(b)+5)
	out = append(out, EncodeVarInt(int32(len(b)))...)
	out = append(out, b...)
	return out
}

func ReadInt64(r io.Reader) (int64, error) {
	var v int64
	if err := binary.Read(r, binary.BigEndian, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func EncodeInt64(v int64) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(v))
	return b[:]
}

func ReadInt32(r io.Reader) (int32, error) {
	var v int32
	if err := binary.Read(r, binary.BigEndian, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func EncodeInt32(v int32) []byte {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(v))
	return b[:]
}

func ReadUUID(r io.Reader) ([16]byte, error) {
	var id [16]byte
	_, err := io.ReadFull(r, id[:])
	return id, err
}

func DiscardN(r io.Reader, n int) error {
	_, err := io.CopyN(io.Discard, r, int64(n))
	return err
}
