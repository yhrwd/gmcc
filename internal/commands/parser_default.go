package commands

import (
	"encoding/json"
	"regexp"
	"strings"
)

type DefaultParser struct {
	prefix  string
	botName string

	jsonPattern     *regexp.Regexp
	standardPattern *regexp.Regexp
}

func NewDefaultParser(prefix, botName string) *DefaultParser {
	return &DefaultParser{
		prefix:          prefix,
		botName:         botName,
		jsonPattern:     regexp.MustCompile(`^\[([^\s➥]+)\s*➥\s*([^\]]+)\]\s*(.+)$`),
		standardPattern: regexp.MustCompile(`^\[([^\]]+) -> 你?\s*\]\s*(.+)$`),
	}
}

func (p *DefaultParser) Parse(raw RawChat) *Message {
	text := p.extractPlainText(raw)
	if text == "" {
		return nil
	}

	if sender, receiver, content := p.parseCustomFormat(text); sender != "" {
		if !p.isBotReceiver(receiver) {
			return nil
		}
		if !strings.HasPrefix(content, p.prefix) {
			return nil
		}
		return &Message{
			Type:       raw.Type,
			PlainText:  content,
			RawJSON:    raw.RawJSON,
			Sender:     sender,
			SenderUUID: formatUUID(raw.SenderUUID),
			IsPrivate:  true,
			Timestamp:  raw.Timestamp,
		}
	}

	if sender, content := p.parseStandardFormat(text); sender != "" {
		if !strings.HasPrefix(content, p.prefix) {
			return nil
		}
		return &Message{
			Type:       raw.Type,
			PlainText:  content,
			RawJSON:    raw.RawJSON,
			Sender:     sender,
			SenderUUID: formatUUID(raw.SenderUUID),
			IsPrivate:  true,
			Timestamp:  raw.Timestamp,
		}
	}

	return nil
}

func (p *DefaultParser) parseCustomFormat(text string) (sender, receiver, content string) {
	matches := p.jsonPattern.FindStringSubmatch(text)
	if matches == nil || len(matches) < 4 {
		return "", "", ""
	}
	return matches[1], matches[2], matches[3]
}

func (p *DefaultParser) parseStandardFormat(text string) (sender, content string) {
	matches := p.standardPattern.FindStringSubmatch(text)
	if matches == nil || len(matches) < 3 {
		return "", ""
	}
	return matches[1], matches[2]
}

func (p *DefaultParser) isBotReceiver(receiver string) bool {
	return strings.EqualFold(receiver, p.botName)
}

func (p *DefaultParser) extractPlainText(raw RawChat) string {
	if raw.PlainText != "" {
		return raw.PlainText
	}

	if raw.RawJSON == "" {
		return ""
	}

	var jsonMsg struct {
		Extra []struct {
			Text string `json:"text"`
		} `json:"extra"`
		Text string `json:"text"`
	}

	if err := json.Unmarshal([]byte(raw.RawJSON), &jsonMsg); err != nil {
		return ""
	}

	var sb strings.Builder
	for _, e := range jsonMsg.Extra {
		sb.WriteString(e.Text)
	}
	sb.WriteString(jsonMsg.Text)
	return sb.String()
}

func formatUUID(uuid [16]byte) string {
	if uuid == [16]byte{} {
		return ""
	}
	hex := make([]byte, 32)
	for i := 0; i < 16; i++ {
		hex[i*2] = "0123456789abcdef"[uuid[i]>>4]
		hex[i*2+1] = "0123456789abcdef"[uuid[i]&0x0f]
	}
	return string(hex[0:8]) + "-" + string(hex[8:12]) + "-" + string(hex[12:16]) + "-" + string(hex[16:20]) + "-" + string(hex[20:32])
}
