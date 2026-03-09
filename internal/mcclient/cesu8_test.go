package mcclient

import (
	"testing"
)

func TestCESU8ToUTF8(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "ASCII only",
			input:    []byte("Hello World"),
			expected: "Hello World",
		},
		{
			name:     "Chinese characters (U+4E2D U+6587)",
			input:    []byte{0xE4, 0xB8, 0xAD, 0xE6, 0x96, 0x87}, // UTF-8 编码的"中文"
			expected: "中文",
		},
		{
			name: "Emoji (U+1F600) - CESU-8 encoded",
			// Emoji 😀 (U+1F600) 在 CESU-8 中编码为代理对:
			// U+D83D (高代理) -> ED A0 BD
			// U+DE00 (低代理) -> ED B8 80
			input:    []byte{0xED, 0xA0, 0xBD, 0xED, 0xB8, 0x80},
			expected: "😀",
		},
		{
			name:     "Mixed ASCII and Chinese",
			input:    []byte("Player \xE6\xB5\x8B\xE8\xAF\x95 joined"), // "Player 测试 joined"
			expected: "Player 测试 joined",
		},
		{
			name: "Multiple emoji - CESU-8",
			// 🎮 U+1F3AE -> UTF-16: D83C DFAE -> CESU-8: ED A0 BC ED BE AE
			// 🎯 U+1F3AF -> UTF-16: D83C DFAF -> CESU-8: ED A0 BC ED BE AF
			input:    []byte{0xED, 0xA0, 0xBC, 0xED, 0xBE, 0xAE, 0xED, 0xA0, 0xBC, 0xED, 0xBE, 0xAF},
			expected: "🎮🎯",
		},
		{
			name:     "Empty input",
			input:    []byte{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cesu8ToUTF8(tt.input)
			if result != tt.expected {
				t.Errorf("cesu8ToUTF8() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestValidUTF8OrReplace(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "Valid UTF-8",
			input:    []byte("Hello 中文"),
			expected: "Hello 中文",
		},
		{
			name:     "Invalid UTF-8 byte",
			input:    []byte{0xFF, 0xFE},
			expected: "??",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validUTF8OrReplace(tt.input)
			if result != tt.expected {
				t.Errorf("validUTF8OrReplace() = %q, want %q", result, tt.expected)
			}
		})
	}
}
