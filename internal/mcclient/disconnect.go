package mcclient

import (
	"bytes"
	"strings"

	chatjson "gmcc/internal/mcclient/chat"
	"gmcc/internal/mcclient/packet"
)

func disconnectReasonFromNBT(data []byte) string {
	rawJSON, err := packet.ReadAnonymousNBTJSON(bytes.NewReader(data))
	if err != nil {
		return packet.RawPreview(data)
	}

	plain := strings.TrimSpace(chatjson.ExtractPlainTextFromChatJSON(rawJSON))
	if plain != "" {
		return plain
	}
	return rawJSON
}

func disconnectReasonFromJSON(data []byte) string {
	r := bytes.NewReader(data)
	rawJSON, err := packet.ReadString(r, r)
	if err != nil {
		return packet.RawPreview(data)
	}

	plain := strings.TrimSpace(chatjson.ExtractPlainTextFromChatJSON(rawJSON))
	if plain != "" {
		return plain
	}
	return rawJSON
}
