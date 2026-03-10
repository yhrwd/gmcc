package mcclient

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type TextComponent struct {
	Text          string          `json:"text,omitempty"`
	Extra         []TextComponent `json:"extra,omitempty"`
	Color         string          `json:"color,omitempty"`
	Bold          bool            `json:"bold,omitempty"`
	Italic        bool            `json:"italic,omitempty"`
	Underlined    bool            `json:"underlined,omitempty"`
	Strikethrough bool            `json:"strikethrough,omitempty"`
	Obfuscated    bool            `json:"obfuscated,omitempty"`
	Translate     string          `json:"translate,omitempty"`
	With          []any           `json:"with,omitempty"`
	Selector      string          `json:"selector,omitempty"`
	Keybind       string          `json:"keybind,omitempty"`
	Insertion     string          `json:"insertion,omitempty"`
}

var ansiColors = map[string]string{
	"black":        "\033[30m",
	"dark_blue":    "\033[34m",
	"dark_green":   "\033[32m",
	"dark_aqua":    "\033[36m",
	"dark_red":     "\033[31m",
	"dark_purple":  "\033[35m",
	"gold":         "\033[33m",
	"gray":         "\033[37m",
	"dark_gray":    "\033[90m",
	"blue":         "\033[94m",
	"green":        "\033[92m",
	"aqua":         "\033[96m",
	"red":          "\033[91m",
	"light_purple": "\033[95m",
	"yellow":       "\033[93m",
	"white":        "\033[97m",
}

const ansiReset = "\033[0m"

func ParseTextComponent(rawJSON string) (*TextComponent, error) {
	var comp TextComponent
	if err := json.Unmarshal([]byte(rawJSON), &comp); err != nil {
		return nil, err
	}
	return &comp, nil
}

func (c *TextComponent) ToANSI() string {
	var sb strings.Builder
	c.renderANSI(&sb, StyleState{})
	return sb.String()
}

type StyleState struct {
	Color         string
	Bold          bool
	Italic        bool
	Underlined    bool
	Strikethrough bool
}

func (c *TextComponent) renderANSI(sb *strings.Builder, parent StyleState) {
	style := parent
	changed := false

	if c.Color != "" {
		style.Color = c.Color
		changed = true
	}
	if c.Bold {
		style.Bold = true
		changed = true
	}
	if c.Italic {
		style.Italic = true
		changed = true
	}
	if c.Underlined {
		style.Underlined = true
		changed = true
	}
	if c.Strikethrough {
		style.Strikethrough = true
		changed = true
	}

	if changed {
		sb.WriteString(formatStyle(style))
	}

	if c.Text != "" {
		sb.WriteString(c.Text)
	}

	if c.Translate != "" {
		switch c.Translate {
		case "command.unknown.argument":
			sb.WriteString("未知命令参数")
		case "command.context.here":
			sb.WriteString(" <--[此处]")
		default:
			sb.WriteString("[")
			sb.WriteString(c.Translate)
			sb.WriteString("]")
		}
	}

	if c.Selector != "" {
		sb.WriteString(c.Selector)
	}

	if c.Keybind != "" {
		sb.WriteString(c.Keybind)
	}

	for _, extra := range c.Extra {
		extra.renderANSI(sb, style)
	}

	if changed {
		sb.WriteString(formatStyle(parent))
	}
}

func formatStyle(s StyleState) string {
	var codes []string

	if s.Bold {
		codes = append(codes, "1")
	}
	if s.Italic {
		codes = append(codes, "3")
	}
	if s.Underlined {
		codes = append(codes, "4")
	}
	if s.Strikethrough {
		codes = append(codes, "9")
	}

	var result strings.Builder

	if len(codes) > 0 {
		result.WriteString("\033[")
		result.WriteString(strings.Join(codes, ";"))
		result.WriteString("m")
	}

	if s.Color != "" {
		if ansi, ok := ansiColors[s.Color]; ok {
			result.WriteString(ansi)
		} else if strings.HasPrefix(s.Color, "#") {
			if fg, err := hexToANSI(s.Color); err == nil {
				result.WriteString(fg)
			}
		}
	}

	return result.String()
}

func hexToANSI(hex string) (string, error) {
	if len(hex) != 7 || hex[0] != '#' {
		return "", fmt.Errorf("invalid hex color")
	}
	r, err := strconv.ParseInt(hex[1:3], 16, 32)
	if err != nil {
		return "", err
	}
	g, err := strconv.ParseInt(hex[3:5], 16, 32)
	if err != nil {
		return "", err
	}
	b, err := strconv.ParseInt(hex[5:7], 16, 32)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b), nil
}

func (c *TextComponent) ToPlain() string {
	var sb strings.Builder
	c.renderPlain(&sb)
	return sb.String()
}

func (c *TextComponent) renderPlain(sb *strings.Builder) {
	sb.WriteString(c.Text)

	if c.Translate != "" {
		sb.WriteString("[")
		sb.WriteString(c.Translate)
		sb.WriteString("]")
	}

	if c.Selector != "" {
		sb.WriteString(c.Selector)
	}

	if c.Keybind != "" {
		sb.WriteString(c.Keybind)
	}

	for _, extra := range c.Extra {
		extra.renderPlain(sb)
	}
}
