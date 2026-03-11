package mcclient

import (
	"encoding/json"
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
	Action      string         `json:"action"`
	Value       *TextComponent `json:"value,omitempty"`
	ValueString string         `json:"value"`
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

var colorToMotd = map[string]string{
	"black":         "§0",
	"dark_blue":     "§1",
	"dark_green":    "§2",
	"dark_aqua":     "§3",
	"dark_red":      "§4",
	"dark_purple":   "§5",
	"gold":          "§6",
	"gray":          "§7",
	"dark_gray":     "§8",
	"blue":          "§9",
	"green":         "§a",
	"aqua":          "§b",
	"red":           "§c",
	"light_purple":  "§d",
	"yellow":        "§e",
	"white":         "§f",
	"obfuscated":    "§k",
	"bold":          "§l",
	"strikethrough": "§m",
	"underlined":    "§n",
	"italic":        "§o",
	"reset":         "§r",
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
	Color         *int
	Bold          bool
	Italic        bool
	Underlined    bool
	Strikethrough bool
	Obfuscated    bool
}

func (c *TextComponent) render(sb *strings.Builder, parent *Style) {
	style := &Style{}
	if parent != nil {
		*style = *parent
	}

	if c.Color != "" {
		if code, ok := colorMap[c.Color]; ok {
			style.Color = &code
		} else if strings.HasPrefix(c.Color, "#") {
			if r, g, b, ok := parseHex(c.Color); ok {
				c256 := nearestColor(r, g, b)
				style.Color = &c256
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

func nearestColor(r, g, b int) int {
	colors := []struct{ r, g, b, code int }{
		{0, 0, 0, 30},
		{0, 0, 170, 34},
		{0, 170, 0, 32},
		{0, 170, 170, 36},
		{170, 0, 0, 31},
		{170, 0, 170, 35},
		{255, 170, 0, 33},
		{170, 170, 170, 37},
		{85, 85, 85, 90},
		{85, 85, 255, 94},
		{85, 255, 85, 92},
		{85, 255, 255, 96},
		{255, 85, 85, 91},
		{255, 85, 255, 95},
		{255, 255, 85, 93},
		{255, 255, 255, 97},
	}

	minDist := 1<<31 - 1
	bestCode := 37
	for _, c := range colors {
		dr, dg, db := r-c.r, g-c.g, b-c.b
		dist := dr*dr + dg*dg + db*db
		if dist < minDist {
			minDist = dist
			bestCode = c.code
		}
	}
	return bestCode
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
		colorName := motdColorName(*s.Color)
		if colorName != "" {
			codes = append(codes, colorName)
		}
	}

	return strings.Join(codes, "")
}

func motdColorName(code int) string {
	switch code {
	case 30:
		return "§0"
	case 31:
		return "§4"
	case 32:
		return "§2"
	case 33:
		return "§6"
	case 34:
		return "§1"
	case 35:
		return "§5"
	case 36:
		return "§3"
	case 37:
		return "§7"
	case 90:
		return "§8"
	case 91:
		return "§c"
	case 92:
		return "§a"
	case 93:
		return "§e"
	case 94:
		return "§9"
	case 95:
		return "§d"
	case 96:
		return "§b"
	case 97:
		return "§f"
	default:
		return ""
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
		if code, ok := colorMap[c.Color]; ok {
			style.Color = &code
		} else if strings.HasPrefix(c.Color, "#") {
			if r, g, b, ok := parseHex(c.Color); ok {
				c256 := nearestColor(r, g, b)
				style.Color = &c256
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
