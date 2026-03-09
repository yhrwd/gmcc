package mcclient

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"gmcc/internal/logx"
)

// extractPlainTextFromChatJSON 从Minecraft聊天JSON中提取纯文本。
func extractPlainTextFromChatJSON(rawJSON string) string {
	if strings.TrimSpace(rawJSON) == "" {
		return ""
	}
	var node any
	if err := json.Unmarshal([]byte(rawJSON), &node); err != nil {
		logx.Debugf("解析聊天JSON失败: %v, 原始JSON: %s", err, rawJSON)
		return rawJSON // 返回原始JSON以记录数据
	}
	var parts []string
	collectChatText(node, &parts)
	text := strings.TrimSpace(strings.Join(parts, ""))
	// 过滤Minecraft颜色代码 (§[0-9a-fk-or])
	text = removeColorCodes(text)
	return text
}

// collectChatText 递归收集聊天JSON中的文本内容。
func collectChatText(node any, parts *[]string) {
	switch v := node.(type) {
	case string:
		if strings.TrimSpace(v) != "" {
			*parts = append(*parts, v)
		}
	case bool, float64, int64, int32, int16, int8, uint64, uint32, uint16, uint8:
		*parts = append(*parts, fmt.Sprint(v))
	case []any:
		for _, item := range v {
			collectChatText(item, parts)
		}
	case map[string]any:
		if text, ok := v["text"].(string); ok {
			if strings.TrimSpace(text) != "" {
				*parts = append(*parts, text)
			}
		}
		if tr, ok := v["translate"].(string); ok {
			// translate key 也保留，至少不丢语义
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
	}
}

// removeColorCodes 移除Minecraft颜色代码。
func removeColorCodes(text string) string {
	// Minecraft颜色代码: § 后面跟 [0-9a-fk-or]
	re := regexp.MustCompile(`§[0-9a-fk-or]`)
	return re.ReplaceAllString(text, "")
}
