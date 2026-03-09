package mcclient

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type CharacterMapping struct {
	Description string `json:"description"`
	ReplaceWith string `json:"replace_with"`
}

type CharacterMapConfig struct {
	EnableReplace   bool                        `json:"enable_replace"`
	ShowUnicodeInfo bool                        `json:"show_unicode_info"`
	Mappings        map[string]CharacterMapping `json:"mappings"`
}

type CharacterAnalyzer struct {
	mu         sync.RWMutex
	config     *CharacterMapConfig
	configPath string
}

func NewCharacterAnalyzer(configPath string) *CharacterAnalyzer {
	analyzer := &CharacterAnalyzer{
		configPath: configPath,
	}
	analyzer.loadConfig()
	return analyzer
}

func (ca *CharacterAnalyzer) loadConfig() error {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	if ca.configPath == "" {
		ca.configPath = "charmap.json"
	}

	data, err := os.ReadFile(ca.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			ca.config = ca.getDefaultConfig()
			return ca.saveConfig()
		}
		return fmt.Errorf("读取字符映射配置失败: %w", err)
	}

	ca.config = &CharacterMapConfig{}
	if err := json.Unmarshal(data, ca.config); err != nil {
		return fmt.Errorf("解析字符映射配置失败: %w", err)
	}

	if ca.config.Mappings == nil {
		ca.config.Mappings = make(map[string]CharacterMapping)
	}

	return nil
}

func (ca *CharacterAnalyzer) saveConfig() error {
	if ca.config == nil {
		ca.config = ca.getDefaultConfig()
	}

	data, err := json.MarshalIndent(ca.config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化字符映射配置失败: %w", err)
	}

	tmpPath := ca.configPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("写入临时文件失败: %w", err)
	}

	if err := os.Rename(tmpPath, ca.configPath); err != nil {
		_ = os.Remove(ca.configPath)
		if err2 := os.Rename(tmpPath, ca.configPath); err2 != nil {
			_ = os.Remove(tmpPath)
			return fmt.Errorf("重命名配置文件失败: %w", err2)
		}
	}

	return nil
}

func (ca *CharacterAnalyzer) getDefaultConfig() *CharacterMapConfig {
	return &CharacterMapConfig{
		EnableReplace:   false,
		ShowUnicodeInfo: true,
		Mappings: map[string]CharacterMapping{
			"\\uE000": {
				Description: "Minecraft 私用区字符示例 1",
				ReplaceWith: "█",
			},
			"\\uE001": {
				Description: "Minecraft 私用区字符示例 2",
				ReplaceWith: "▓",
			},
			"\\uE002": {
				Description: "Minecraft 私用区字符示例 3",
				ReplaceWith: "▒",
			},
		},
	}
}

func (ca *CharacterAnalyzer) AnalyzeText(text string) string {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	if !ca.config.ShowUnicodeInfo || strings.TrimSpace(text) == "" {
		return ""
	}

	var specialChars []string
	seenChars := make(map[rune]bool)

	for i, r := range text {
		if r > 127 || r < 32 || r == 0xFFFD {
			if !seenChars[r] {
				seenChars[r] = true
				charInfo := fmt.Sprintf("'%c'(U+%04X,pos=%d)", r, r, i)
				specialChars = append(specialChars, charInfo)
			}
		}
	}

	if len(specialChars) > 0 {
		sort.Strings(specialChars)
		return "特殊字符: " + strings.Join(specialChars, ", ")
	}

	return ""
}

func (ca *CharacterAnalyzer) ReplaceText(text string) string {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	if !ca.config.EnableReplace || len(ca.config.Mappings) == 0 {
		return text
	}

	result := text
	for unicodeStr, mapping := range ca.config.Mappings {
		char := parseUnicodeEscape(unicodeStr)
		if char != 0 {
			result = strings.ReplaceAll(result, string(char), mapping.ReplaceWith)
		}
	}

	return result
}

func parseUnicodeEscape(s string) rune {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "\\u") {
		return 0
	}

	hexStr := strings.TrimPrefix(s, "\\u")
	var codePoint uint32
	_, err := fmt.Sscanf(hexStr, "%X", &codePoint)
	if err != nil {
		return 0
	}

	return rune(codePoint)
}

func (ca *CharacterAnalyzer) AddMapping(unicodeStr, description, replaceWith string) error {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	if ca.config.Mappings == nil {
		ca.config.Mappings = make(map[string]CharacterMapping)
	}

	ca.config.Mappings[unicodeStr] = CharacterMapping{
		Description: description,
		ReplaceWith: replaceWith,
	}

	return ca.saveConfig()
}

func (ca *CharacterAnalyzer) RemoveMapping(unicodeStr string) error {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	if ca.config.Mappings == nil {
		return nil
	}

	delete(ca.config.Mappings, unicodeStr)
	return ca.saveConfig()
}

func (ca *CharacterAnalyzer) SetEnableReplace(enable bool) error {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	ca.config.EnableReplace = enable
	return ca.saveConfig()
}

func (ca *CharacterAnalyzer) SetShowUnicodeInfo(show bool) error {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	ca.config.ShowUnicodeInfo = show
	return ca.saveConfig()
}

func (ca *CharacterAnalyzer) GetConfig() *CharacterMapConfig {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	configCopy := *ca.config
	configCopy.Mappings = make(map[string]CharacterMapping)
	for k, v := range ca.config.Mappings {
		configCopy.Mappings[k] = v
	}
	return &configCopy
}

func (ca *CharacterAnalyzer) GenerateMappingTemplate(outputPath string) error {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	template := CharacterMapConfig{
		EnableReplace:   true,
		ShowUnicodeInfo: true,
		Mappings: map[string]CharacterMapping{
			"\\uE000": {
				Description: "方块图标 - 左上角",
				ReplaceWith: "┌",
			},
			"\\uE001": {
				Description: "方块图标 - 横线",
				ReplaceWith: "─",
			},
			"\\uE002": {
				Description: "方块图标 - 右上角",
				ReplaceWith: "┐",
			},
			"\\uE003": {
				Description: "方块图标 - 竖线",
				ReplaceWith: "│",
			},
			"\\uE004": {
				Description: "方块图标 - 左下角",
				ReplaceWith: "└",
			},
			"\\uE005": {
				Description: "方块图标 - 右下角",
				ReplaceWith: "┘",
			},
			"\\uE010": {
				Description: "实心方块",
				ReplaceWith: "█",
			},
			"\\uE011": {
				Description: "深色阴影",
				ReplaceWith: "▓",
			},
			"\\uE012": {
				Description: "中等阴影",
				ReplaceWith: "▒",
			},
			"\\uE013": {
				Description: "浅色阴影",
				ReplaceWith: "░",
			},
			"\\uE020": {
				Description: "货币符号 - 金币",
				ReplaceWith: "●",
			},
			"\\uE021": {
				Description: "货币符号 - 银币",
				ReplaceWith: "○",
			},
			"\\uE030": {
				Description: "图标 - 心形",
				ReplaceWith: "♥",
			},
			"\\uE031": {
				Description: "图标 - 星形",
				ReplaceWith: "★",
			},
			"\\uE032": {
				Description: "图标 - 对勾",
				ReplaceWith: "✓",
			},
			"\\uE033": {
				Description: "图标 - 叉号",
				ReplaceWith: "✗",
			},
			"\\uE040": {
				Description: "箭头 - 右",
				ReplaceWith: "→",
			},
			"\\uE041": {
				Description: "箭头 - 左",
				ReplaceWith: "←",
			},
			"\\uE042": {
				Description: "箭头 - 上",
				ReplaceWith: "↑",
			},
			"\\uE043": {
				Description: "箭头 - 下",
				ReplaceWith: "↓",
			},
		},
	}

	if outputPath == "" {
		outputPath = "charmap_template.json"
	}

	data, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化模板失败: %w", err)
	}

	return os.WriteFile(outputPath, data, 0644)
}

func InitializeCharacterMap(configDir string) (*CharacterAnalyzer, error) {
	configPath := filepath.Join(configDir, "charmap.json")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("创建配置目录失败: %w", err)
	}

	analyzer := NewCharacterAnalyzer(configPath)

	return analyzer, nil
}
