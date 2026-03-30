package mcclient

import (
	"bytes"
	"testing"

	"gmcc/internal/mcclient/packet"
)

func TestBuildChatSignableBodyUsesMillisecondTimestamp(t *testing.T) {
	var playerUUID [16]byte
	var sessionID [16]byte

	body := buildChatSignableBody(playerUUID, sessionID, 7, "hello", 1234567890123, 99, nil)
	want := packet.EncodeInt64(1234567890123)

	if !bytes.Contains(body, want) {
		t.Fatalf("signable body does not contain millisecond timestamp %x", want)
	}
	if bytes.Contains(body, packet.EncodeInt64(1234567890)) {
		t.Fatalf("signable body should not contain second-level timestamp")
	}
}
