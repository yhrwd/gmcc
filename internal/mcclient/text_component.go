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
	Bold          any             `json:"bold,omitempty"`
	Italic        any             `json:"italic,omitempty"`
	Underlined    any             `json:"underlined,omitempty"`
	Strikethrough any             `json:"strikethrough,omitempty"`
	Translate     string          `json:"translate,omitempty"`
}

var colorMap = map[string]int{
	"black":        30,
	"dark_blue":    34,
	"dark_green":   32,
	"dark_aqua":    36,
	"dark_red":     31,
	"dark_purple":  35,
	"gold":         33,
	"gray":         37,
	"dark_gray":    90,
	"blue":         94,
	"green":        92,
	"aqua":         96,
	"red":          91,
	"light_purple": 95,
	"yellow":       93,
	"white":        97,
}

func ParseTextComponent(rawJSON string) (*TextComponent, error) {
	var comp TextComponent
	if err := json.Unmarshal([]byte(rawJSON), &comp); err != nil {
		return nil, err
	}
	return &comp, nil
}

func (c *TextComponent) ToANSI() string {
	var sb strings.Builder
	c.render(&sb, nil)
	sb.WriteString("\033[0m")
	return sb.String()
}

type Style struct {
	Color      *int
	Bold       bool
	Italic     bool
	Underlined bool
}

func (c *TextComponent) render(sb *strings.Builder, parent *Style) {
	style := &Style{}
	if parent != nil {
		*style = *parent
	}

	hasChanges := false

	if c.Color != "" {
		if code, ok := colorMap[c.Color]; ok {
			style.Color = &code
			hasChanges = true
		} else if strings.HasPrefix(c.Color, "#") {
			if r, g, b, ok := parseHex(c.Color); ok {
				style.Color = nil
				sb.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b))
				hasChanges = true
			}
		}
	}

	if toBool(c.Bold) && !style.Bold {
		style.Bold = true
		hasChanges = true
	}
	if toBool(c.Italic) && !style.Italic {
		style.Italic = true
		hasChanges = true
	}
	if toBool(c.Underlined) && !style.Underlined {
		style.Underlined = true
		hasChanges = true
	}

	if hasChanges {
		sb.WriteString(style.toANSI())
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

	for _, extra := range c.Extra {
		extra.render(sb, style)
	}

	if hasChanges && parent != nil {
		sb.WriteString(parent.toANSI())
	}
}

func (s *Style) toANSI() string {
	var codes []int

	if s.Bold {
		codes = append(codes, 1)
	}
	if s.Italic {
		codes = append(codes, 3)
	}
	if s.Underlined {
		codes = append(codes, 4)
	}
	if s.Color != nil {
		codes = append(codes, *s.Color)
	}

	if len(codes) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\033[")
	for i, code := range codes {
		if i > 0 {
			sb.WriteString(";")
		}
		sb.WriteString(strconv.Itoa(code))
	}
	sb.WriteString("m")
	return sb.String()
}

func toBool(v any) bool {
	switch val := v.(type) {
	case bool:
		return val
	case float64:
		return val != 0
	case string:
		return val == "true" || val == "1"
	default:
		return false
	}
}

func parseHex(hex string) (r, g, b int, ok bool) {
	if len(hex) != 7 || hex[0] != '#' {
		return 0, 0, 0, false
	}
	rr, err := strconv.ParseInt(hex[1:3], 16, 32)
	if err != nil {
		return 0, 0, 0, false
	}
	gg, err := strconv.ParseInt(hex[3:5], 16, 32)
	if err != nil {
		return 0, 0, 0, false
	}
	bb, err := strconv.ParseInt(hex[5:7], 16, 32)
	if err != nil {
		return 0, 0, 0, false
	}
	return int(rr), int(gg), int(bb), true
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
	for _, extra := range c.Extra {
		extra.renderPlain(sb)
	}
}
