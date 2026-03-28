package commands

import (
	"testing"
)

func TestAuthManager_Check(t *testing.T) {
	tests := []struct {
		name      string
		allowAll  bool
		whitelist []string
		msg       *Message
		expected  bool
	}{
		{
			name:      "allow all",
			allowAll:  true,
			whitelist: nil,
			msg: &Message{
				Sender:     "Player1",
				SenderUUID: "uuid-123",
			},
			expected: true,
		},
		{
			name:      "whitelist by name",
			allowAll:  false,
			whitelist: []string{"Player1", "Player2"},
			msg: &Message{
				Sender:     "Player1",
				SenderUUID: "uuid-123",
			},
			expected: true,
		},
		{
			name:      "whitelist by UUID",
			allowAll:  false,
			whitelist: []string{"uuid-456", "uuid-789"},
			msg: &Message{
				Sender:     "Player1",
				SenderUUID: "uuid-456",
			},
			expected: true,
		},
		{
			name:      "not in whitelist",
			allowAll:  false,
			whitelist: []string{"Player1", "Player2"},
			msg: &Message{
				Sender:     "Stranger",
				SenderUUID: "uuid-999",
			},
			expected: false,
		},
		{
			name:      "case insensitive name",
			allowAll:  false,
			whitelist: []string{"Player1"},
			msg: &Message{
				Sender:     "player1",
				SenderUUID: "uuid-123",
			},
			expected: true,
		},
		{
			name:      "empty whitelist",
			allowAll:  false,
			whitelist: []string{},
			msg: &Message{
				Sender:     "Player1",
				SenderUUID: "uuid-123",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := NewAuthManager()
			auth.SetAllowAll(tt.allowAll)
			auth.SetWhitelist(tt.whitelist)

			result := auth.Check(tt.msg)
			if result != tt.expected {
				t.Errorf("Check() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestAuthManager_AddRemove(t *testing.T) {
	auth := NewAuthManager()

	auth.AddToWhitelist("Player1", "Player2")
	if !auth.IsAllowed("Player1", "") {
		t.Error("Player1 should be allowed after AddToWhitelist")
	}

	auth.RemoveFromWhitelist("Player1")
	if auth.IsAllowed("Player1", "") {
		t.Error("Player1 should not be allowed after RemoveFromWhitelist")
	}
	if !auth.IsAllowed("Player2", "") {
		t.Error("Player2 should still be allowed")
	}
}

func TestAuthManager_GetWhitelist(t *testing.T) {
	auth := NewAuthManager()
	auth.SetWhitelist([]string{"Player1", "Player2", "Player3"})

	list := auth.GetWhitelist()
	if len(list) != 3 {
		t.Errorf("GetWhitelist() returned %d items, expected 3", len(list))
	}
}
