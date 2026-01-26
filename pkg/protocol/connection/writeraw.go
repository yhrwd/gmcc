package connection

import (
	"bytes"
	"errors"
	"fmt"
	"gmcc/pkg/protocol/codec"
	"net"
	"time"
)

func (c *ConnContext) WriteRaw(data []byte) error {
	c.WriteMu.Lock()
	defer c.WriteMu.Unlock()

	if err := c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return fmt.Errorf("set write deadline: %w", err)
	}

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

	if c.SharedSecret != nil {
		if c.EncryptStream == nil {
			return errors.New("encryption enabled but encryptStream is nil")
		}
		c.EncryptStream.XORKeyStream(payload, payload)
	}

	lengthBuf := &bytes.Buffer{}
	if err := codec.WriteVarInt(lengthBuf, int32(len(payload))); err != nil {
		return fmt.Errorf("write length: %w", err)
	}

	buffers := net.Buffers{lengthBuf.Bytes(), payload}
	_, err := buffers.WriteTo(c.Conn.Raw)
	if err != nil {
		return fmt.Errorf("write to conn: %w", err)
	}

	return nil
}
