package logx

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestLogTokenCache(t *testing.T) {
	// Setup
	var buf bytes.Buffer
	consoleLogger = log.New(&buf, "", 0)

	// Test with all parameters
	LogTokenCache("minecraft", "player1", "uuid-123")

	output := buf.String()
	if !strings.Contains(output, "[INFO] 使用缓存的 minecraft token: player1 (uuid-123)") {
		t.Errorf("Expected output to contain token info, got: %s", output)
	}

	// Test with minimal parameters
	buf.Reset()
	LogTokenCache("microsoft", "", "")

	output = buf.String()
	if !strings.Contains(output, "[INFO] 使用缓存的 microsoft token") {
		t.Errorf("Expected output to contain token type only, got: %s", output)
	}
}

func TestLogTokenExpired(t *testing.T) {
	// Setup
	var buf bytes.Buffer
	consoleLogger = log.New(&buf, "", 0)

	// Test with error
	LogTokenExpired("microsoft", os.ErrNotExist)

	output := buf.String()
	if !strings.Contains(output, "缓存 microsoft token 已失效") {
		t.Errorf("Expected output to contain expired token message, got: %s", output)
	}

	// Test without error
	buf.Reset()
	LogTokenExpired("minecraft", nil)

	output = buf.String()
	if !strings.Contains(output, "缓存 minecraft token 已失效") {
		t.Errorf("Expected output to contain expired token message without refresh, got: %s", output)
	}
}

func TestPacketErrorWithContext(t *testing.T) {
	// Setup
	var buf bytes.Buffer
	consoleLogger = log.New(&buf, "", 0)

	testData := []byte{0x01, 0x02, 0x03}

	// Test with context and error
	PacketErrorWithContext("test_packet", testData, os.ErrNotExist, "读取 test 数据失败")

	output := buf.String()
	if !strings.Contains(output, "Packet解析失败: test_packet") {
		t.Errorf("Expected output to contain packet error, got: %s", output)
	}
	if !strings.Contains(output, "读取 test 数据失败") {
		t.Errorf("Expected output to contain context, got: %s", output)
	}

	// Test with context but no error
	buf.Reset()
	PacketErrorWithContext("test_packet", testData, nil, "处理数据")

	output = buf.String()
	if !strings.Contains(output, "处理数据") {
		t.Errorf("Expected output to contain context only, got: %s", output)
	}

	// Test with no context
	buf.Reset()
	PacketErrorWithContext("test_packet", testData, os.ErrNotExist, "")

	output = buf.String()
	if !strings.Contains(output, "file does not exist") {
		t.Errorf("Expected output to contain error message, got: %s", output)
	}
}
