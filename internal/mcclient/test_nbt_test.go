package mcclient

import (
	"bytes"
	"encoding/json"
	"testing"

	"gmcc/internal/mcclient/packet"
	"gmcc/internal/nbt"
)

func TestNBTDecoder_String(t *testing.T) {
	data := []byte{0x08, 0x00, 0x04, 't', 'e', 's', 't'}
	dec := nbt.NewDecoder(bytes.NewReader(data))
	dec.NetworkFormat(true)
	var result string
	if err := dec.Decode(&result); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if result != "test" {
		t.Errorf("expected 'test', got %q", result)
	}
}

func TestNBTDecoder_CESU8String(t *testing.T) {
	data := []byte{
		0x08,
		0x00, 0x06,
		0xED, 0xA0, 0xBD, 0xED, 0xB8, 0x80,
	}

	dec := nbt.NewDecoder(bytes.NewReader(data))
	dec.NetworkFormat(true)
	var result string
	if err := dec.Decode(&result); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if result != "😀" {
		t.Errorf("expected '😀', got %q (len=%d)", result, len(result))
	}
}

func TestNBTDecoder_Compound(t *testing.T) {
	data := []byte{
		0x0A,
		0x08,
		0x00, 0x04, 'n', 'a', 'm', 'e',
		0x00, 0x04, 't', 'e', 's', 't',
		0x00,
	}

	dec := nbt.NewDecoder(bytes.NewReader(data))
	dec.NetworkFormat(true)
	var result map[string]any
	if err := dec.Decode(&result); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if s, ok := result["name"].(string); !ok || s != "test" {
		t.Errorf("expected name='test', got %v", result["name"])
	}
}

func TestNBTDecoder_List(t *testing.T) {
	data := []byte{
		0x09,
		0x08,
		0x00, 0x00, 0x00, 0x03,
		0x00, 0x01, 'a',
		0x00, 0x01, 'b',
		0x00, 0x01, 'c',
	}

	dec := nbt.NewDecoder(bytes.NewReader(data))
	dec.NetworkFormat(true)
	var result []any
	if err := dec.Decode(&result); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 elements, got %d", len(result))
	}
}

func TestNBTDecoder_Int(t *testing.T) {
	data := []byte{
		0x03,
		0x00, 0x00, 0x30, 0x39,
	}

	dec := nbt.NewDecoder(bytes.NewReader(data))
	dec.NetworkFormat(true)
	var result int32
	if err := dec.Decode(&result); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if result != 12345 {
		t.Errorf("expected 12345, got %d", result)
	}
}

func TestNBTDecoder_ByteArray(t *testing.T) {
	data := []byte{
		0x07,
		0x00, 0x00, 0x00, 0x03,
		0x01, 0x02, 0x03,
	}

	dec := nbt.NewDecoder(bytes.NewReader(data))
	dec.NetworkFormat(true)
	var result []byte
	if err := dec.Decode(&result); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 bytes, got %d", len(result))
	}
}

func TestReadAnonymousNBTJSON(t *testing.T) {
	data := []byte{
		0x0A,
		0x08,
		0x00, 0x04, 't', 'e', 'x', 't',
		0x00, 0x05, 'H', 'e', 'l', 'l', 'o',
		0x00,
	}

	jsonStr, err := packet.ReadAnonymousNBTJSON(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("ReadAnonymousNBTJSON failed: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Errorf("invalid JSON: %v\nJSON: %s", err, jsonStr)
	}

	if parsed["text"] != "Hello" {
		t.Errorf("expected text='Hello', got %v", parsed["text"])
	}
}

func TestNBTDecoder_CESU8Chinese(t *testing.T) {
	chinese := []byte{0xE6, 0xB5, 0x8B, 0xE8, 0xAF, 0x95}
	data := append([]byte{
		0x08,
		0x00, 0x06,
	}, chinese...)

	dec := nbt.NewDecoder(bytes.NewReader(data))
	dec.NetworkFormat(true)
	var result string
	if err := dec.Decode(&result); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if result != "测试" {
		t.Errorf("expected '测试', got %q", result)
	}
}
