package commands

import (
	"testing"
	"time"
)

func TestDefaultParser_Parse(t *testing.T) {
	tests := []struct {
		name       string
		prefix     string
		botName    string
		raw        RawChat
		wantNil    bool
		wantSender string
		wantText   string
	}{
		{
			name:    "custom format ride command",
			prefix:  "!",
			botName: "TestBot",
			raw: RawChat{
				PlainText: "[Player1 ➥ TestBot] !ride",
			},
			wantNil:    false,
			wantSender: "Player1",
			wantText:   "!ride",
		},
		{
			name:    "custom format not for bot",
			prefix:  "!",
			botName: "TestBot",
			raw: RawChat{
				PlainText: "[Player1 ➥ OtherBot] !ride",
			},
			wantNil: true,
		},
		{
			name:    "custom format no prefix",
			prefix:  "!",
			botName: "TestBot",
			raw: RawChat{
				PlainText: "[Player1 ➥ TestBot] hello",
			},
			wantNil: true,
		},
		{
			name:    "standard format tell",
			prefix:  "!",
			botName: "TestBot",
			raw: RawChat{
				PlainText: "[Player2 -> 你] !ride",
			},
			wantNil:    false,
			wantSender: "Player2",
			wantText:   "!ride",
		},
		{
			name:    "standard format no prefix",
			prefix:  "!",
			botName: "TestBot",
			raw: RawChat{
				PlainText: "[Player2 -> 你] hello",
			},
			wantNil: true,
		},
		{
			name:    "empty message",
			prefix:  "!",
			botName: "TestBot",
			raw: RawChat{
				PlainText: "",
			},
			wantNil: true,
		},
		{
			name:    "unknown format",
			prefix:  "!",
			botName: "TestBot",
			raw: RawChat{
				PlainText: "Player1: !ride",
			},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewDefaultParser(tt.prefix, tt.botName)
			msg := parser.Parse(tt.raw)

			if tt.wantNil {
				if msg != nil {
					t.Errorf("Parse() should return nil, got %+v", msg)
				}
				return
			}

			if msg == nil {
				t.Error("Parse() returned nil, expected non-nil")
				return
			}

			if msg.Sender != tt.wantSender {
				t.Errorf("Parse().Sender = %q, want %q", msg.Sender, tt.wantSender)
			}

			if msg.PlainText != tt.wantText {
				t.Errorf("Parse().PlainText = %q, want %q", msg.PlainText, tt.wantText)
			}

			if !msg.IsPrivate {
				t.Error("Parse().IsPrivate should be true")
			}
		})
	}
}

func TestDefaultParser_ExtractPlainText(t *testing.T) {
	parser := NewDefaultParser("!", "Bot")

	tests := []struct {
		name     string
		raw      RawChat
		expected string
	}{
		{
			name: "from plain text",
			raw: RawChat{
				PlainText: "[Player ➥ Bot] hello",
			},
			expected: "[Player ➥ Bot] hello",
		},
		{
			name: "from JSON",
			raw: RawChat{
				RawJSON: `{"extra":[{"text":"["},{"text":"Player "},{"text":"➥ "},{"text":"Bot"},{"text":"] "},{"text":"hello"}],"text":""}`,
			},
			expected: "[Player ➥ Bot] hello",
		},
		{
			name:     "empty",
			raw:      RawChat{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := parser.extractPlainText(tt.raw)
			if text != tt.expected {
				t.Errorf("extractPlainText() = %q, want %q", text, tt.expected)
			}
		})
	}
}

func TestNewRawChat(t *testing.T) {
	typ := "player_chat"
	plainText := "[Player ➥ Bot] test"
	rawJSON := `{"text":""}`
	senderName := "TestPlayer"
	senderUUID := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	raw := NewRawChat(typ, plainText, rawJSON, senderName, senderUUID)

	if raw.Type != typ {
		t.Errorf("Type = %q, want %q", raw.Type, typ)
	}
	if raw.PlainText != plainText {
		t.Errorf("PlainText = %q, want %q", raw.PlainText, plainText)
	}
	if raw.SenderName != senderName {
		t.Errorf("SenderName = %q, want %q", raw.SenderName, senderName)
	}
	if raw.SenderUUID != senderUUID {
		t.Error("SenderUUID mismatch")
	}
	if raw.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}

	t2 := time.Now().Add(-time.Hour)
	raw = NewRawChatWithTime(typ, plainText, rawJSON, senderName, senderUUID, t2)
	if !raw.Timestamp.Equal(t2) {
		t.Error("Timestamp should match provided time")
	}
}
