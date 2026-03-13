package chat

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
	Obfuscated    any             `json:"obfuscated,omitempty"`
	Translate     string          `json:"translate,omitempty"`
	With          []TextComponent `json:"with,omitempty"`
	Selector      string          `json:"selector,omitempty"`
	Keybind       string          `json:"keybind,omitempty"`
	ClickEvent    *ClickEvent     `json:"clickEvent,omitempty"`
	HoverEvent    *HoverEvent     `json:"hoverEvent,omitempty"`
	Insertion     string          `json:"insertion,omitempty"`
}

type ClickEvent struct {
	Action string `json:"action"`
	Value  string `json:"value"`
}

type HoverEvent struct {
	Action string         `json:"action"`
	Value  *TextComponent `json:"value,omitempty"`
}

var basicColorMap = map[string]ColorRGB{
	"black":        {0, 0, 0},
	"dark_blue":    {0, 0, 170},
	"dark_green":   {0, 170, 0},
	"dark_aqua":    {0, 170, 170},
	"dark_red":     {170, 0, 0},
	"dark_purple":  {170, 0, 170},
	"gold":         {255, 170, 0},
	"gray":         {170, 170, 170},
	"dark_gray":    {85, 85, 85},
	"blue":         {85, 85, 255},
	"green":        {85, 255, 85},
	"aqua":         {85, 255, 255},
	"red":          {255, 85, 85},
	"light_purple": {255, 85, 255},
	"yellow":       {255, 255, 85},
	"white":        {255, 255, 255},
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
	result := sb.String()
	return result
}

type Style struct {
	Color         *ColorRGB
	Bold          bool
	Italic        bool
	Underlined    bool
	Strikethrough bool
	Obfuscated    bool
}

type ColorRGB struct {
	R, G, B int
}

func (c *TextComponent) render(sb *strings.Builder, parent *Style) {
	style := &Style{}
	if parent != nil {
		*style = *parent
	}

	if c.Color != "" {
		if code, ok := basicColorMap[c.Color]; ok {
			style.Color = &ColorRGB{R: code.R, G: code.G, B: code.B}
		} else if strings.HasPrefix(c.Color, "#") {
			if r, g, b, ok := parseHex(c.Color); ok {
				style.Color = &ColorRGB{R: r, G: g, B: b}
			}
		}
	}

	if toBool(c.Bold) {
		style.Bold = true
	}
	if toBool(c.Italic) {
		style.Italic = true
	}
	if toBool(c.Underlined) {
		style.Underlined = true
	}

	sb.WriteString(style.toANSI())

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

	for i := range c.Extra {
		c.Extra[i].render(sb, style)
	}

	if parent != nil {
		sb.WriteString(parent.toANSI())
	}
}

func (s *Style) toANSI() string {
	var sb strings.Builder

	if s.Bold {
		sb.WriteString("\033[1m")
	}
	if s.Italic {
		sb.WriteString("\033[3m")
	}
	if s.Underlined {
		sb.WriteString("\033[4m")
	}
	if s.Strikethrough {
		sb.WriteString("\033[9m")
	}
	if s.Obfuscated {
		sb.WriteString("\033[8m")
	}
	if s.Color != nil {
		sb.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm", s.Color.R, s.Color.G, s.Color.B))
	}

	return sb.String()
}

func (s *Style) toMotd() string {
	var codes []string

	if s.Obfuscated {
		codes = append(codes, "§k")
	}
	if s.Bold {
		codes = append(codes, "§l")
	}
	if s.Strikethrough {
		codes = append(codes, "§m")
	}
	if s.Underlined {
		codes = append(codes, "§n")
	}
	if s.Italic {
		codes = append(codes, "§o")
	}
	if s.Color != nil {
		colorName := rgbToMotd(s.Color)
		if colorName != "" {
			codes = append(codes, colorName)
		}
	}

	return strings.Join(codes, "")
}

func rgbToMotd(c *ColorRGB) string {
	switch {
	case c.R == 0 && c.G == 0 && c.B == 0:
		return "§0"
	case c.R == 170 && c.G == 0 && c.B == 0:
		return "§4"
	case c.R == 0 && c.G == 170 && c.B == 0:
		return "§2"
	case c.R == 0 && c.G == 0 && c.B == 170:
		return "§1"
	case c.R == 170 && c.G == 0 && c.B == 170:
		return "§5"
	case c.R == 0 && c.G == 170 && c.B == 170:
		return "§3"
	case c.R == 255 && c.G == 170 && c.B == 0:
		return "§6"
	case c.R == 170 && c.G == 170 && c.B == 170:
		return "§7"
	case c.R == 85 && c.G == 85 && c.B == 85:
		return "§8"
	case c.R == 85 && c.G == 85 && c.B == 255:
		return "§9"
	case c.R == 85 && c.G == 255 && c.B == 85:
		return "§a"
	case c.R == 85 && c.G == 255 && c.B == 255:
		return "§b"
	case c.R == 255 && c.G == 85 && c.B == 85:
		return "§c"
	case c.R == 255 && c.G == 85 && c.B == 255:
		return "§d"
	case c.R == 255 && c.G == 255 && c.B == 85:
		return "§e"
	case c.R == 255 && c.G == 255 && c.B == 255:
		return "§f"
	default:
		return fmt.Sprintf("§#%02x%02x%02x", c.R, c.G, c.B)
	}
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

func (c *TextComponent) ToMotd() string {
	var sb strings.Builder
	c.renderMotd(&sb, nil)
	return sb.String()
}

func (c *TextComponent) renderMotd(sb *strings.Builder, parent *Style) {
	style := &Style{}
	if parent != nil {
		*style = *parent
	}

	if c.Color != "" {
		if code, ok := basicColorMap[c.Color]; ok {
			style.Color = &ColorRGB{R: code.R, G: code.G, B: code.B}
		} else if strings.HasPrefix(c.Color, "#") {
			if r, g, b, ok := parseHex(c.Color); ok {
				style.Color = &ColorRGB{R: r, G: g, B: b}
			}
		}
	}

	if toBool(c.Bold) {
		style.Bold = true
	}
	if toBool(c.Italic) {
		style.Italic = true
	}
	if toBool(c.Underlined) {
		style.Underlined = true
	}
	if toBool(c.Strikethrough) {
		style.Strikethrough = true
	}
	if toBool(c.Obfuscated) {
		style.Obfuscated = true
	}

	sb.WriteString(style.toMotd())

	if c.Text != "" {
		sb.WriteString(c.Text)
	}

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

	for i := range c.Extra {
		c.Extra[i].renderMotd(sb, style)
	}

	if parent != nil {
		sb.WriteString(parent.toMotd())
	}
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

// MessageBuilder helps build complex Minecraft text components
type MessageBuilder struct {
	root TextComponent
	curr *TextComponent
}

func NewMessageBuilder(text string) *MessageBuilder {
	b := &MessageBuilder{}
	b.root.Text = text
	b.curr = &b.root
	return b
}

func (b *MessageBuilder) Color(color string) *MessageBuilder {
	b.curr.Color = color
	return b
}

func (b *MessageBuilder) Bold(v bool) *MessageBuilder {
	b.curr.Bold = v
	return b
}

func (b *MessageBuilder) Italic(v bool) *MessageBuilder {
	b.curr.Italic = v
	return b
}

func (b *MessageBuilder) Add(text string) *MessageBuilder {
	b.root.Extra = append(b.root.Extra, TextComponent{Text: text})
	b.curr = &b.root.Extra[len(b.root.Extra)-1]
	return b
}

func (b *MessageBuilder) Build() TextComponent {
	return b.root
}

func (b *MessageBuilder) BuildJSON() string {
	data, _ := json.Marshal(b.root)
	return string(data)
}
