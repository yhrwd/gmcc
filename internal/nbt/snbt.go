package nbt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidSNBT = errors.New("invalid SNBT syntax")
)

// ParseSNBT parses a stringified NBT string into a value
func ParseSNBT(snbt string) (any, error) {
	p := &snbtParser{input: snbt}
	return p.parse()
}

// UnmarshalSNBT parses SNBT into v
func UnmarshalSNBT(snbt string, v any) error {
	val, err := ParseSNBT(snbt)
	if err != nil {
		return err
	}

	// Use json.Marshal/Unmarshal for reflection
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

type snbtParser struct {
	input string
	pos   int
}

func (p *snbtParser) parse() (any, error) {
	p.skipWhitespace()
	if p.pos >= len(p.input) {
		return nil, ErrInvalidSNBT
	}

	ch := p.input[p.pos]
	switch ch {
	case '{':
		return p.parseCompound()
	case '[':
		return p.parseList()
	case '"', '\'':
		s, err := p.parseString()
		if err != nil {
			return nil, err
		}
		// Check for number suffixes
		if p.pos < len(p.input) {
			suffix := p.input[p.pos]
			if suffix == 'b' || suffix == 'B' || suffix == 's' || suffix == 'S' ||
				suffix == 'l' || suffix == 'L' || suffix == 'f' || suffix == 'F' ||
				suffix == 'd' || suffix == 'D' {
				p.pos++
				return p.parseNumber(s, suffix)
			}
		}
		return s, nil
	default:
		// Try to parse as unquoted string or number
		return p.parseUnquoted()
	}
}

func (p *snbtParser) parseCompound() (map[string]any, error) {
	if p.input[p.pos] != '{' {
		return nil, ErrInvalidSNBT
	}
	p.pos++

	result := make(map[string]any)
	p.skipWhitespace()

	if p.pos < len(p.input) && p.input[p.pos] == '}' {
		p.pos++
		return result, nil
	}

	for {
		p.skipWhitespace()
		if p.pos >= len(p.input) {
			return nil, ErrInvalidSNBT
		}

		// Parse key
		var key string
		var err error
		if p.input[p.pos] == '"' || p.input[p.pos] == '\'' {
			key, err = p.parseString()
		} else {
			key, err = p.parseUnquotedKey()
		}
		if err != nil {
			return nil, err
		}

		p.skipWhitespace()
		if p.pos >= len(p.input) || p.input[p.pos] != ':' {
			return nil, ErrInvalidSNBT
		}
		p.pos++

		p.skipWhitespace()
		value, err := p.parse()
		if err != nil {
			return nil, err
		}
		result[key] = value

		p.skipWhitespace()
		if p.pos >= len(p.input) {
			return nil, ErrInvalidSNBT
		}

		if p.input[p.pos] == '}' {
			p.pos++
			return result, nil
		}
		if p.input[p.pos] != ',' {
			return nil, ErrInvalidSNBT
		}
		p.pos++
	}
}

func (p *snbtParser) parseList() ([]any, error) {
	if p.input[p.pos] != '[' {
		return nil, ErrInvalidSNBT
	}
	p.pos++

	p.skipWhitespace()
	if p.pos >= len(p.input) {
		return nil, ErrInvalidSNBT
	}

	// Check for type prefix (B;, I;, L;)
	var listType byte
	if p.pos+1 < len(p.input) && p.input[p.pos+1] == ';' {
		listType = p.input[p.pos]
		p.pos += 2
		p.skipWhitespace()
	}

	if p.pos < len(p.input) && p.input[p.pos] == ']' {
		p.pos++
		return []any{}, nil
	}

	var result []any
	for {
		p.skipWhitespace()
		value, err := p.parse()
		if err != nil {
			return nil, err
		}

		// Apply type conversion for typed arrays
		if listType != 0 {
			value, err = p.convertArrayType(value, listType)
			if err != nil {
				return nil, err
			}
		}
		result = append(result, value)

		p.skipWhitespace()
		if p.pos >= len(p.input) {
			return nil, ErrInvalidSNBT
		}

		if p.input[p.pos] == ']' {
			p.pos++
			return result, nil
		}
		if p.input[p.pos] != ',' {
			return nil, ErrInvalidSNBT
		}
		p.pos++
	}
}

func (p *snbtParser) parseString() (string, error) {
	if p.pos >= len(p.input) {
		return "", ErrInvalidSNBT
	}
	quote := p.input[p.pos]
	if quote != '"' && quote != '\'' {
		return "", ErrInvalidSNBT
	}
	p.pos++

	var result strings.Builder
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if ch == quote {
			p.pos++
			return result.String(), nil
		}
		if ch == '\\' && p.pos+1 < len(p.input) {
			p.pos++
			escaped := p.input[p.pos]
			switch escaped {
			case 'n':
				result.WriteByte('\n')
			case 'r':
				result.WriteByte('\r')
			case 't':
				result.WriteByte('\t')
			case '\\', '"', '\'':
				result.WriteByte(escaped)
			default:
				result.WriteByte('\\')
				result.WriteByte(escaped)
			}
		} else {
			result.WriteByte(ch)
		}
		p.pos++
	}
	return "", ErrInvalidSNBT
}

func (p *snbtParser) parseUnquotedKey() (string, error) {
	var result strings.Builder
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if ch == ':' || ch == ',' || ch == '}' || ch == ']' || unicode.IsSpace(rune(ch)) {
			break
		}
		result.WriteByte(ch)
		p.pos++
	}
	return result.String(), nil
}

func (p *snbtParser) parseUnquoted() (any, error) {
	var result strings.Builder
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if ch == ',' || ch == '}' || ch == ']' || unicode.IsSpace(rune(ch)) {
			break
		}
		result.WriteByte(ch)
		p.pos++
	}

	s := result.String()

	// Check for number suffixes
	if len(s) > 1 {
		last := s[len(s)-1]
		numPart := s[:len(s)-1]
		switch last {
		case 'b', 'B':
			return p.parseNumber(numPart, last)
		case 's', 'S':
			return p.parseNumber(numPart, last)
		case 'l', 'L':
			return p.parseNumber(numPart, last)
		case 'f', 'F':
			return p.parseNumber(numPart, last)
		case 'd', 'D':
			return p.parseNumber(numPart, last)
		}
	}

	// Try to parse as number
	if v, err := p.parseNumber(s, 0); err == nil {
		return v, nil
	}

	// Check for boolean
	if s == "true" || s == "false" {
		return s == "true", nil
	}

	return s, nil
}

func (p *snbtParser) parseNumber(s string, suffix byte) (any, error) {
	// Remove underscore separators
	s = strings.ReplaceAll(s, "_", "")

	switch suffix {
	case 'b', 'B':
		v, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return nil, err
		}
		return int8(v), nil
	case 's', 'S':
		v, err := strconv.ParseInt(s, 10, 16)
		if err != nil {
			return nil, err
		}
		return int16(v), nil
	case 'l', 'L':
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		return v, nil
	case 'f', 'F':
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, err
		}
		return float32(v), nil
	case 'd', 'D':
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		return v, nil
	default:
		// Try int first
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			if v >= int64(-2147483648) && v <= int64(2147483647) {
				return int32(v), nil
			}
			return v, nil
		}
		// Try float
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			return v, nil
		}
		return nil, fmt.Errorf("invalid number: %s", s)
	}
}

func (p *snbtParser) convertArrayType(value any, listType byte) (any, error) {
	switch listType {
	case 'B':
		if f, ok := value.(float64); ok {
			return int8(f), nil
		}
	case 'I':
		if f, ok := value.(float64); ok {
			return int32(f), nil
		}
	case 'L':
		if f, ok := value.(float64); ok {
			return int64(f), nil
		}
	}
	return value, nil
}

func (p *snbtParser) skipWhitespace() {
	for p.pos < len(p.input) {
		if !unicode.IsSpace(rune(p.input[p.pos])) {
			break
		}
		p.pos++
	}
}

// FormatSNBT converts a value to SNBT string
func FormatSNBT(v any) string {
	return formatValue(v, true)
}

func formatValue(v any, topLevel bool) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case bool:
		if val {
			return "true"
		}
		return "false"
	case int8:
		return fmt.Sprintf("%db", val)
	case int16:
		return fmt.Sprintf("%ds", val)
	case int32:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%dL", val)
	case float32:
		return fmt.Sprintf("%gf", val)
	case float64:
		return fmt.Sprintf("%gd", val)
	case string:
		if topLevel {
			return fmt.Sprintf(`"%s"`, escapeString(val))
		}
		return fmt.Sprintf(`"%s"`, escapeString(val))
	case []byte:
		return formatByteArray(val)
	case []int32:
		return formatIntArray(val)
	case []int64:
		return formatLongArray(val)
	case []any:
		return formatList(val)
	case map[string]any:
		return formatCompound(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func escapeString(s string) string {
	var result strings.Builder
	for _, ch := range s {
		switch ch {
		case '"':
			result.WriteString(`\"`)
		case '\\':
			result.WriteString(`\\`)
		case '\n':
			result.WriteString(`\n`)
		case '\r':
			result.WriteString(`\r`)
		case '\t':
			result.WriteString(`\t`)
		default:
			result.WriteRune(ch)
		}
	}
	return result.String()
}

func formatByteArray(b []byte) string {
	var buf bytes.Buffer
	buf.WriteString("[B;")
	for i, v := range b {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(fmt.Sprintf("%db", int8(v)))
	}
	buf.WriteString("]")
	return buf.String()
}

func formatIntArray(arr []int32) string {
	var buf bytes.Buffer
	buf.WriteString("[I;")
	for i, v := range arr {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(fmt.Sprintf("%d", v))
	}
	buf.WriteString("]")
	return buf.String()
}

func formatLongArray(arr []int64) string {
	var buf bytes.Buffer
	buf.WriteString("[L;")
	for i, v := range arr {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(fmt.Sprintf("%dL", v))
	}
	buf.WriteString("]")
	return buf.String()
}

func formatList(list []any) string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, v := range list {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(formatValue(v, false))
	}
	buf.WriteByte(']')
	return buf.String()
}

func formatCompound(m map[string]any) string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	first := true
	for k, v := range m {
		if !first {
			buf.WriteByte(',')
		}
		first = false
		buf.WriteString(fmt.Sprintf(`"%s":`, escapeString(k)))
		buf.WriteString(formatValue(v, false))
	}
	buf.WriteByte('}')
	return buf.String()
}
