package mcclient

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
	"gmcc/internal/nbt"
)

type packet struct {
	ID   int32
	Data []byte
}

const maxFrameLength = 2 * 1024 * 1024

type packetConn struct {
	conn net.Conn
	br   *bufio.Reader
	w    io.Writer

	mu                   sync.Mutex
	compressionThreshold int
}

func newPacketConn(conn net.Conn) *packetConn {
	return &packetConn{
		conn:                 conn,
		br:                   bufio.NewReader(conn),
		w:                    conn,
		compressionThreshold: -1,
	}
}

func (c *packetConn) Close() error {
	return c.conn.Close()
}

func (c *packetConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *packetConn) SetCompressionThreshold(threshold int) {
	c.compressionThreshold = threshold
}

func (c *packetConn) EnableEncryption(secret []byte) error {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return fmt.Errorf("create AES cipher failed: %w", err)
	}
	decrypter := newCFB8(block, secret, false)
	encrypter := newCFB8(block, secret, true)
	c.br = bufio.NewReader(&cipher.StreamReader{S: decrypter, R: c.conn})
	c.w = &cipher.StreamWriter{S: encrypter, W: c.conn}
	logx.Debugf("已启用连接加密: mode=AES/CFB8 secretLen=%d", len(secret))
	return nil
}

func (c *packetConn) ReadPacket() (packet, error) {
	frameLen, err := readVarInt(c.br)
	if err != nil {
		logx.Debugf("读取 frame 长度失败: %v", err)
		return packet{}, err
	}
	if frameLen < 0 {
		return packet{}, fmt.Errorf("invalid frame length %d", frameLen)
	}
	if frameLen > maxFrameLength {
		return packet{}, fmt.Errorf("invalid frame length %d > %d (可能是加密/压缩流失步)", frameLen, maxFrameLength)
	}

	frame := make([]byte, frameLen)
	if _, err := io.ReadFull(c.br, frame); err != nil {
		logx.Debugf("读取 frame 内容失败: frameLen=%d err=%v", frameLen, err)
		return packet{}, err
	}

	data := frame
	if c.compressionThreshold >= 0 {
		r := bytes.NewReader(frame)
		uncompressedLen, err := readVarInt(r)
		if err != nil {
			return packet{}, fmt.Errorf("read compressed length failed: %w", err)
		}
		rest, err := io.ReadAll(r)
		if err != nil {
			return packet{}, err
		}

		if uncompressedLen == 0 {
			data = rest
		} else {
			zr, err := zlib.NewReader(bytes.NewReader(rest))
			if err != nil {
				return packet{}, fmt.Errorf("create zlib reader failed: %w", err)
			}
			uncompressed, err := io.ReadAll(zr)
			_ = zr.Close()
			if err != nil {
				return packet{}, fmt.Errorf("decompress packet failed: %w", err)
			}
			if len(uncompressed) != int(uncompressedLen) {
				return packet{}, fmt.Errorf("invalid decompressed size: want %d got %d", uncompressedLen, len(uncompressed))
			}
			data = uncompressed
		}

		logx.PacketLogf(
			"收包 frame: frameLen=%d compressed=%t threshold=%d uncompressedLen=%d framePreview=%s",
			frameLen,
			uncompressedLen != 0,
			c.compressionThreshold,
			len(data),
			rawPreview(frame),
		)
	}

	r := bytes.NewReader(data)
	packetID, err := readVarInt(r)
	if err != nil {
		return packet{}, fmt.Errorf("read packet id failed: %w", err)
	}
	payload, err := io.ReadAll(r)
	if err != nil {
		return packet{}, err
	}

	if c.compressionThreshold < 0 {
		logx.PacketLogf(
			"收包 frame: frameLen=%d compressed=false threshold=%d uncompressedLen=%d framePreview=%s",
			frameLen,
			c.compressionThreshold,
			len(data),
			rawPreview(frame),
		)
	}

	return packet{ID: packetID, Data: payload}, nil
}

func (c *packetConn) WritePacket(packetID int32, payload []byte) error {
	body := append(encodeVarInt(packetID), payload...)

	var frame []byte
	if c.compressionThreshold >= 0 {
		if len(body) >= c.compressionThreshold {
			var compressed bytes.Buffer
			zw := zlib.NewWriter(&compressed)
			if _, err := zw.Write(body); err != nil {
				_ = zw.Close()
				return err
			}
			if err := zw.Close(); err != nil {
				return err
			}
			frame = append(encodeVarInt(int32(len(body))), compressed.Bytes()...)
		} else {
			frame = append(encodeVarInt(0), body...)
		}
	} else {
		frame = body
	}

	raw := append(encodeVarInt(int32(len(frame))), frame...)

	logx.PacketLogf(
		"发包: id=0x%02X payloadLen=%d bodyLen=%d frameLen=%d compressed=%t threshold=%d payloadPreview=%s",
		packetID,
		len(payload),
		len(body),
		len(frame),
		c.compressionThreshold >= 0 && len(body) >= c.compressionThreshold,
		c.compressionThreshold,
		rawPreview(payload),
	)

	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.w.Write(raw)
	if err != nil {
		logx.PacketLogf("发包失败: id=0x%02X err=%v", packetID, err)
	}
	return err
}

func readVarInt(r io.ByteReader) (int32, error) {
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

func encodeVarInt(v int32) []byte {
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

func readBool(r io.Reader) (bool, error) {
	var b [1]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return false, err
	}
	return b[0] != 0, nil
}

func encodeBool(v bool) []byte {
	if v {
		return []byte{1}
	}
	return []byte{0}
}

func readString(r io.ByteReader, rr io.Reader) (string, error) {
	n, err := readVarInt(r)
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

func encodeString(s string) []byte {
	b := []byte(s)
	out := make([]byte, 0, len(b)+5)
	out = append(out, encodeVarInt(int32(len(b)))...)
	out = append(out, b...)
	return out
}

func readByteArray(r io.ByteReader, rr io.Reader) ([]byte, error) {
	n, err := readVarInt(r)
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

func encodeByteArray(b []byte) []byte {
	out := make([]byte, 0, len(b)+5)
	out = append(out, encodeVarInt(int32(len(b)))...)
	out = append(out, b...)
	return out
}

func readInt64(r io.Reader) (int64, error) {
	var v int64
	if err := binary.Read(r, binary.BigEndian, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func encodeInt64(v int64) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(v))
	return b[:]
}

func readInt32(r io.Reader) (int32, error) {
	var v int32
	if err := binary.Read(r, binary.BigEndian, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func encodeInt32(v int32) []byte {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(v))
	return b[:]
}

func readUUID(r io.Reader) ([16]byte, error) {
	var id [16]byte
	_, err := io.ReadFull(r, id[:])
	return id, err
}

type cfb8 struct {
	b        cipher.Block
	next     []byte
	tmp      []byte
	encrypt  bool
	blockLen int
}

func newCFB8(block cipher.Block, iv []byte, encrypt bool) cipher.Stream {
	bs := block.BlockSize()
	next := make([]byte, bs)
	copy(next, iv)
	return &cfb8{
		b:        block,
		next:     next,
		tmp:      make([]byte, bs),
		encrypt:  encrypt,
		blockLen: bs,
	}
}

func (x *cfb8) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic("cfb8 output smaller than input")
	}

	for i := 0; i < len(src); i++ {
		in := src[i]
		x.b.Encrypt(x.tmp, x.next)
		out := in ^ x.tmp[0]
		dst[i] = out

		copy(x.next, x.next[1:])
		if x.encrypt {
			x.next[x.blockLen-1] = out
		} else {
			x.next[x.blockLen-1] = in
		}
	}
}
