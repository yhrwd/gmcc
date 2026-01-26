下面给你一个 完整、可直接使用、生产级别的 WriteRaw 实现。
它包含：
 
- 并发安全（writeMu）
- 自动处理压缩
- 自动处理加密
- 自动写入 Packet Length（VarInt）
- 超时控制
- 错误处理
- 与你的 ConnContext 结构完全兼容
 

 
一、WriteRaw 完整代码（推荐版）
 
go
  
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

 

 
二、你需要补充的字段（ConnContext）
 
因为上面的 WriteRaw 使用了加密流，你需要在 ConnContext 里加一个字段：
 
go
  
type ConnContext struct {
	// ...

	encryptStream cipher.Stream // 加密流（登录成功后初始化）
}

 

 
三、你需要补充的辅助功能
 
1. 压缩功能（codec.Compress）
 
go
  
package codec

import (
	"bytes"
	"compress/zlib"
	"io"
)

func Compress(data []byte, threshold int) ([]byte, error) {
	if len(data) < threshold {
		// 不压缩，直接返回原数据
		return data, nil
	}

	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

 
2. 加密初始化（登录成功后）
 
go
  
// 在登录成功后调用
func (c *ConnContext) EnableEncryption(sharedSecret []byte) error {
	block, err := aes.NewCipher(sharedSecret)
	if err != nil {
		return err
	}

	iv := make([]byte, block.BlockSize())
	// Minecraft 使用全零 IV
	c.encryptStream = cipher.NewCTR(block, iv)

	return nil
}

 

 
四、WriteRaw 的使用示例
 
go
  
// 发送握手包
func (s *ClientHandshakeState) DoHandshake(ctx *connection.ConnContext) error {
	packet := &bytes.Buffer{}
	codec.WriteVarInt(packet, 0x00) // Handshake Packet ID
	// ... 写入其他字段 ...

	return ctx.WriteRaw(packet.Bytes())
}

 

 
五、总结：WriteRaw 需要完善的内容
 
你需要确保以下内容都已实现：
 
1.writeMu 加锁（你已有）
2.压缩逻辑（根据 CompressionThreshold）
3.加密逻辑（使用 AES-CTR）
4.Packet Length VarInt 写入
5.写超时
6.使用 net.Buffers 进行零拷贝写入（高性能）
7.错误包装（fmt.Errorf + %w）
 

 
如果你愿意，我还能帮你：
 
- 实现 ReadPacket（完整带解密、解压）
- 帮你把状态机全部串联起来
- 帮你写一个完整的 Minecraft 客户端连接流程（握手 → Status → Login → Play）
 
你想继续哪一部分？