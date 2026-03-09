package mcclient

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestNBTDecoder_String(t *testing.T) {
	// NBT String: tag type (0x08), length (0x00 0x04), data "test"
	data := []byte{0x08, 0x00, 0x04, 't', 'e', 's', 't'}
	dec := newNBTDecoder(bytes.NewReader(data), true)

	v, err := dec.decodeRoot()
	if err != nil {
		t.Fatalf("decodeRoot failed: %v", err)
	}

	s, ok := v.(string)
	if !ok {
		t.Fatalf("expected string, got %T", v)
	}
	if s != "test" {
		t.Errorf("expected 'test', got %q", s)
	}
}

func TestNBTDecoder_CESU8String(t *testing.T) {
	// NBT String with CESU-8 encoded emoji: 😀 (U+1F600)
	// CESU-8: ED A0 BD ED B8 80 (6 bytes)
	data := []byte{
		0x08,       // TagString
		0x00, 0x06, // length = 6
		0xED, 0xA0, 0xBD, 0xED, 0xB8, 0x80, // CESU-8 emoji
	}

	dec := newNBTDecoder(bytes.NewReader(data), true)
	v, err := dec.decodeRoot()
	if err != nil {
		t.Fatalf("decodeRoot failed: %v", err)
	}

	s, ok := v.(string)
	if !ok {
		t.Fatalf("expected string, got %T", v)
	}
	if s != "😀" {
		t.Errorf("expected '😀', got %q (len=%d)", s, len(s))
	}
}

func TestNBTDecoder_Compound(t *testing.T) {
	// NBT Compound with one string field:
	// TagCompound, TagString, name "name", value "test", TagEnd
	data := []byte{
		0x0A,                           // TagCompound
		0x08,                           // TagString
		0x00, 0x04, 'n', 'a', 'm', 'e', // name
		0x00, 0x04, 't', 'e', 's', 't', // value
		0x00, // TagEnd
	}

	dec := newNBTDecoder(bytes.NewReader(data), true)
	v, err := dec.decodeRoot()
	if err != nil {
		t.Fatalf("decodeRoot failed: %v", err)
	}

	m, ok := v.(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T", v)
	}

	if s, ok := m["name"].(string); !ok || s != "test" {
		t.Errorf("expected name='test', got %v", m["name"])
	}
}

func TestNBTDecoder_List(t *testing.T) {
	// NBT List of strings: ["a", "b", "c"]
	data := []byte{
		0x09,                   // TagList
		0x08,                   // element type: TagString
		0x00, 0x00, 0x00, 0x03, // length = 3
		0x00, 0x01, 'a',
		0x00, 0x01, 'b',
		0x00, 0x01, 'c',
	}

	dec := newNBTDecoder(bytes.NewReader(data), true)
	v, err := dec.decodeRoot()
	if err != nil {
		t.Fatalf("decodeRoot failed: %v", err)
	}

	list, ok := v.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", v)
	}
	if len(list) != 3 {
		t.Errorf("expected 3 elements, got %d", len(list))
	}
}

func TestNBTDecoder_Int(t *testing.T) {
	// NBT Int: 12345 (0x00003039)
	data := []byte{
		0x03,                   // TagInt
		0x00, 0x00, 0x30, 0x39, // value = 12345
	}

	dec := newNBTDecoder(bytes.NewReader(data), true)
	v, err := dec.decodeRoot()
	if err != nil {
		t.Fatalf("decodeRoot failed: %v", err)
	}

	n, ok := v.(int32)
	if !ok {
		t.Fatalf("expected int32, got %T", v)
	}
	if n != 12345 {
		t.Errorf("expected 12345, got %d", n)
	}
}

func TestNBTDecoder_ByteArray(t *testing.T) {
	// NBT ByteArray: [1, 2, 3]
	data := []byte{
		0x07,                   // TagByteArray
		0x00, 0x00, 0x00, 0x03, // length = 3
		0x01, 0x02, 0x03,
	}

	dec := newNBTDecoder(bytes.NewReader(data), true)
	v, err := dec.decodeRoot()
	if err != nil {
		t.Fatalf("decodeRoot failed: %v", err)
	}

	ba, ok := v.([]byte)
	if !ok {
		t.Fatalf("expected []byte, got %T", v)
	}
	if len(ba) != 3 {
		t.Errorf("expected 3 bytes, got %d", len(ba))
	}
}

func TestReadAnonymousNBTJSON(t *testing.T) {
	// Minecraft chat message format (network NBT)
	// Root compound with text="Hello"
	data := []byte{
		0x0A,                           // TagCompound
		0x08,                           // TagString
		0x00, 0x04, 't', 'e', 'x', 't', // name: "text"
		0x00, 0x05, 'H', 'e', 'l', 'l', 'o', // value: "Hello"
		0x00, // TagEnd (root)
	}

	jsonStr, err := readAnonymousNBTJSON(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("readAnonymousNBTJSON failed: %v", err)
	}

	// Verify JSON is valid
	var parsed map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Errorf("invalid JSON: %v\nJSON: %s", err, jsonStr)
	}

	if parsed["text"] != "Hello" {
		t.Errorf("expected text='Hello', got %v", parsed["text"])
	}
}

func TestNBTDecoder_CESU8Chinese(t *testing.T) {
	// 中文 "测试" UTF-8 编码
	chinese := []byte{0xE6, 0xB5, 0x8B, 0xE8, 0xAF, 0x95}
	data := append([]byte{
		0x08,       // TagString
		0x00, 0x06, // length = 6
	}, chinese...)

	dec := newNBTDecoder(bytes.NewReader(data), true)
	v, err := dec.decodeRoot()
	if err != nil {
		t.Fatalf("decodeRoot failed: %v", err)
	}

	s, ok := v.(string)
	if !ok {
		t.Fatalf("expected string, got %T", v)
	}
	if s != "测试" {
		t.Errorf("expected '测试', got %q", s)
	}
}
