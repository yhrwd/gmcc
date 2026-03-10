package mcclient

import (
	"testing"

	"gmcc/internal/nbt"
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
			input:    []byte{0xE4, 0xB8, 0xAD, 0xE6, 0x96, 0x87},
			expected: "中文",
		},
		{
			name:     "Emoji (U+1F600) - CESU-8 encoded",
			input:    []byte{0xED, 0xA0, 0xBD, 0xED, 0xB8, 0x80},
			expected: "😀",
		},
		{
			name:     "Mixed ASCII and Chinese",
			input:    []byte("Player \xE6\xB5\x8B\xE8\xAF\x95 joined"),
			expected: "Player 测试 joined",
		},
		{
			name:     "Multiple emoji - CESU-8",
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
			result := nbt.CESU8ToUTF8(tt.input)
			if result != tt.expected {
				t.Errorf("CESU8ToUTF8() = %q, want %q", result, tt.expected)
			}
		})
	}
}
