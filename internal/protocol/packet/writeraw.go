package connection

import (
	"bytes"
	"context"
	"crypto/cipher"
	"errors"
	"fmt"
	"gmcc/internal/protocol/codec"
	"io"
	"net"
	"sync"
	"time"
)

// WriteRaw 发送原始数据包内容（Packet ID + Payload）
// 自动处理：
// 1. 加锁（writeMu）
// 2. 压缩（如果启用）
// 3. 加密（如果启用）
// 4. 写入 Packet Length（VarInt）
// 5. 超时控制
func (c *ConnContext) WriteRaw(data []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	// 设置写超时
	if err := c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return fmt.Errorf("set write deadline: %w", err)
	}

	// 1. 压缩（如果启用）
	var payload []byte
	if c.CompressionThreshold >= 0 {
		compressed, err := codec.Compress(data, c.CompressionThreshold)
		if err != nil {
			return fmt.Errorf("compress: %w", err)
		}
		payload = compressed
	} else {
		payload = data
	}

	// 2. 加密（如果启用）
	if c.SharedSecret != nil {
		// 假设你有一个 cipher.Stream
		// 这里需要你提前初始化 encryptStream
		// 例如：c.encryptStream = cipher.NewCTR(block, iv)
		if c.encryptStream == nil {
			return errors.New("encryption enabled but encryptStream is nil")
		}
		c.encryptStream.XORKeyStream(payload, payload)
	}

	// 3. 写入 Packet Length (VarInt)
	lengthBuf := &bytes.Buffer{}
	if err := codec.WriteVarInt(lengthBuf, int32(len(payload))); err != nil {
		return fmt.Errorf("write length: %w", err)
	}

	// 4. 写入数据
	buffers := net.Buffers{lengthBuf.Bytes(), payload}
	_, err := buffers.WriteTo(c.Conn)
	if err != nil {
		return fmt.Errorf("write to conn: %w", err)
	}

	return nil
}
