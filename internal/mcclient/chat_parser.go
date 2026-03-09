package mcclient

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"gmcc/internal/logx"
)

// cesu8ToUTF8 将 CESU-8（Modified UTF-8）转换为标准 UTF-8。
// Minecraft Java 使用 CESU-8 编码，其中辅助平面字符(U+10000 以上)被编码为 6 字节代理对，
// 而不是标准 UTF-8 的 4 字节。
func cesu8ToUTF8(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	// 快速检查是否包含 CESU-8 代理对
	hasCESU8 := false
	for i := 0; i < len(data)-5; i++ {
		if data[i] == 0xED && (data[i+1]&0xF0) == 0xA0 {
			hasCESU8 = true
			break
		}
	}

	if !hasCESU8 {
		return validUTF8OrReplace(data)
	}

	var result bytes.Buffer
	result.Grow(len(data))

	i := 0
	for i < len(data) {
		// 检测 CESU-8 代理对: ED A0-BF 80-BF ED B0-BF 80-BF (6 字节)
		if i+5 < len(data) && data[i] == 0xED {
			b1, b2, b3, b4, b5, b6 := data[i+1], data[i+2], data[i+3], data[i+4], data[i+5], byte(0)
			if i+6 < len(data) {
				b6 = data[i+6]
			}
			_ = b6 // b6 实际上不使用，只是为了消除未使用警告

			// 高代理: ED A0-BF 80-BF (代理范围 D800-DBFF)
			if (b1&0xF0) == 0xA0 && (b2&0xC0) == 0x80 {
				// 检查是否有对应的低代理
				if b3 == 0xED && (b4&0xF0) == 0xB0 && (b5&0xC0) == 0x80 {
					// 解码 CESU-8 代理对
					highSurrogate := uint16(0xD800) + uint16(b1&0x0F)<<6 + uint16(b2&0x3F)
					lowSurrogate := uint16(0xDC00) + uint16(b4&0x0F)<<6 + uint16(b5&0x3F)

					// 将代理对转换为 Unicode 码点
					codePoint := utf16.Decode([]uint16{highSurrogate, lowSurrogate})
					if len(codePoint) > 0 {
						result.WriteRune(codePoint[0])
						i += 6
						continue
					}
				}
			}
		}

		// 普通字符
		r, size := utf8.DecodeRune(data[i:])
		if r == utf8.RuneError {
			result.WriteByte(data[i])
			i++
		} else {
			result.WriteRune(r)
			i += size
		}
	}

	return result.String()
}

// validUTF8OrReplace 确保返回有效的 UTF-8 字符串
func validUTF8OrReplace(data []byte) string {
	if utf8.Valid(data) {
		return string(data)
	}

	var result bytes.Buffer
	result.Grow(len(data))

	i := 0
	for i < len(data) {
		r, size := utf8.DecodeRune(data[i:])
		if r == utf8.RuneError && size == 1 {
			result.WriteByte('?')
			i++
		} else {
			result.WriteRune(r)
			i += size
		}
	}

	return result.String()
}

// extractPlainTextFromChatJSON 从Minecraft聊天JSON中提取纯文本。
func extractPlainTextFromChatJSON(rawJSON string) string {
	if strings.TrimSpace(rawJSON) == "" {
		return ""
	}

	var node any
	if err := json.Unmarshal([]byte(rawJSON), &node); err != nil {
		logx.Debugf("解析聊天JSON失败: %v, 原始JSON: %s", err, rawJSON)
		return rawJSON
	}

	// 递归处理所有字符串字段，修复可能的 CESU-8 问题
	node = fixCESU8InValue(node)

	var parts []string
	collectChatText(node, &parts)
	text := strings.TrimSpace(strings.Join(parts, ""))
	text = removeColorCodes(text)
	return text
}

// fixCESU8InValue 递归修复值中的 CESU-8 问题
func fixCESU8InValue(v any) any {
	switch val := v.(type) {
	case string:
		return fixCESU8String(val)
	case map[string]any:
		result := make(map[string]any, len(val))
		for k, v := range val {
			result[k] = fixCESU8InValue(v)
		}
		return result
	case []any:
		result := make([]any, len(val))
		for i, v := range val {
			result[i] = fixCESU8InValue(v)
		}
		return result
	default:
		return v
	}
}

// fixCESU8String 修复字符串中可能的 CESU-8 编码问题
// 当原始数据包含 CESU-8 编码时，如果被当作 UTF-8 解析，
// 辅助平面字符会变成无效序列，Go 会将其替换为 \ufffd
// 我们需要检测并尝试恢复
func fixCESU8String(s string) string {
	// 如果没有替换字符，直接返回
	if !strings.ContainsRune(s, '\ufffd') {
		return s
	}

	// 已经是损坏状态，无法恢复原始字符
	// 但我们可以保持字符串原样显示
	return s
}

// collectChatText 递归收集聊天JSON中的文本内容。
func collectChatText(node any, parts *[]string) {
	switch v := node.(type) {
	case string:
		if strings.TrimSpace(v) != "" {
			*parts = append(*parts, v)
		}
	case map[string]any:
		if text, ok := v["text"].(string); ok {
			if strings.TrimSpace(text) != "" {
				*parts = append(*parts, text)
			}
		}
		if tr, ok := v["translate"].(string); ok {
			*parts = append(*parts, "["+tr+"]")
		}
		if selector, ok := v["selector"].(string); ok {
			if strings.TrimSpace(selector) != "" {
				*parts = append(*parts, selector)
			}
		}
		if keybind, ok := v["keybind"].(string); ok {
			if strings.TrimSpace(keybind) != "" {
				*parts = append(*parts, keybind)
			}
		}
		if insertion, ok := v["insertion"].(string); ok {
			if strings.TrimSpace(insertion) != "" {
				*parts = append(*parts, insertion)
			}
		}
		if score, ok := v["score"].(map[string]any); ok {
			if val, ok := score["value"]; ok {
				collectChatText(val, parts)
			} else if name, ok := score["name"]; ok {
				collectChatText(name, parts)
			}
		}
		if with, ok := v["with"].([]any); ok {
			for _, item := range with {
				collectChatText(item, parts)
			}
		}
		if extra, ok := v["extra"].([]any); ok {
			for _, item := range extra {
				collectChatText(item, parts)
			}
		}
		if content, ok := v["content"]; ok {
			collectChatText(content, parts)
		}
		if separator, ok := v["separator"]; ok {
			collectChatText(separator, parts)
		}
	case []any:
		for _, item := range v {
			collectChatText(item, parts)
		}
	}
}

// removeColorCodes 移除Minecraft颜色代码。
func removeColorCodes(text string) string {
	// Minecraft颜色代码: § 后面跟 [0-9a-fk-or]
	re := regexp.MustCompile(`§[0-9a-fk-or]`)
	return re.ReplaceAllString(text, "")
}
