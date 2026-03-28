package commands

import (
	"time"
)

type RawChat struct {
	Type       string
	PlainText  string
	RawJSON    string
	SenderName string
	SenderUUID [16]byte
	Timestamp  time.Time
}

type MessageParser interface {
	Parse(raw RawChat) *Message
}

func NewRawChat(typ, plainText, rawJSON, senderName string, senderUUID [16]byte) RawChat {
	return RawChat{
		Type:       typ,
		PlainText:  plainText,
		RawJSON:    rawJSON,
		SenderName: senderName,
		SenderUUID: senderUUID,
		Timestamp:  time.Now(),
	}
}

func NewRawChatWithTime(typ, plainText, rawJSON, senderName string, senderUUID [16]byte, timestamp time.Time) RawChat {
	return RawChat{
		Type:       typ,
		PlainText:  plainText,
		RawJSON:    rawJSON,
		SenderName: senderName,
		SenderUUID: senderUUID,
		Timestamp:  timestamp,
	}
}
