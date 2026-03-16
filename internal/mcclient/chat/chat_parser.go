package chat

import (
	"encoding/json"
	"regexp"
	"strings"

	"gmcc/internal/logx"
)

// ExtractPlainTextFromChatJSON 从Minecraft聊天JSON中提取纯文本。
func ExtractPlainTextFromChatJSON(rawJSON string) string {
	if strings.TrimSpace(rawJSON) == "" {
		return ""
	}

	var node any
	if err := json.Unmarshal([]byte(rawJSON), &node); err != nil {
		logx.Debugf("解析聊天JSON失败: %v, 原始JSON: %s", err, rawJSON)
		return rawJSON
	}

	var parts []string
	CollectChatText(node, &parts)
	text := strings.TrimSpace(strings.Join(parts, ""))
	text = RemoveColorCodes(text)
	return text
}

// CollectChatText 递归收集聊天JSON中的文本内容。
func CollectChatText(node any, parts *[]string) {
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
			if fallback, ok := v["fallback"].(string); ok && strings.TrimSpace(fallback) != "" {
				*parts = append(*parts, fallback)
			} else {
				*parts = append(*parts, "["+tr+"]")
			}
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
				CollectChatText(val, parts)
			} else if name, ok := score["name"]; ok {
				CollectChatText(name, parts)
			}
		}
		if with, ok := v["with"].([]any); ok {
			for _, item := range with {
				CollectChatText(item, parts)
			}
		}
		if extra, ok := v["extra"].([]any); ok {
			for _, item := range extra {
				CollectChatText(item, parts)
			}
		}
		if content, ok := v["content"]; ok {
			CollectChatText(content, parts)
		}
		if separator, ok := v["separator"]; ok {
			CollectChatText(separator, parts)
		}
	case []any:
		for _, item := range v {
			CollectChatText(item, parts)
		}
	}
}

// RemoveColorCodes 移除Minecraft颜色代码。
func RemoveColorCodes(text string) string {
	// Minecraft颜色代码: § 后面跟 [0-9a-fk-or]
	re := regexp.MustCompile(`§[0-9a-fk-or]`)
	return RemoveHexColorCodes(re.ReplaceAllString(text, ""))
}

// RemoveHexColorCodes 移除 1.16+ 的十六进制颜色代码 §#xxxxxx
func RemoveHexColorCodes(text string) string {
	re := regexp.MustCompile(`§#[0-9a-fA-F]{6}`)
	return re.ReplaceAllString(text, "")
}
