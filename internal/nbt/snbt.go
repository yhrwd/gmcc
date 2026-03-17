package nbt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
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
			escaped, err := p.parseEscape()
			if err != nil {
				return "", err
			}
			result.WriteString(escaped)
		} else {
			result.WriteByte(ch)
		}
		p.pos++
	}
	return "", ErrInvalidSNBT
}

func (p *snbtParser) parseEscape() (string, error) {
	if p.pos >= len(p.input) {
		return "", ErrInvalidSNBT
	}
	ch := p.input[p.pos]
	switch ch {
	case 'b':
		return "\b", nil
	case 'f':
		return "\f", nil
	case 'n':
		return "\n", nil
	case 'r':
		return "\r", nil
	case 's':
		return " ", nil
	case 't':
		return "\t", nil
	case '\\':
		return "\\", nil
	case '\'':
		return "'", nil
	case '"':
		return "\"", nil
	case 'x':
		return p.parseHexEscape(2)
	case 'u':
		return p.parseHexEscape(4)
	case 'U':
		return p.parseHexEscape(8)
	case 'N':
		return p.parseUnicodeName()
	}
	return "", fmt.Errorf("invalid escape sequence: \\%c", ch)
}

func (p *snbtParser) parseHexEscape(digits int) (string, error) {
	if p.pos+digits > len(p.input) {
		return "", ErrInvalidSNBT
	}
	p.pos += digits
	hex := p.input[p.pos-digits : p.pos]
	code, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return "", err
	}
	return string(rune(code)), nil
}

func (p *snbtParser) parseUnicodeName() (string, error) {
	if p.pos >= len(p.input) || p.input[p.pos] != '{' {
		return "", ErrInvalidSNBT
	}
	p.pos++
	start := p.pos
	for p.pos < len(p.input) && p.input[p.pos] != '}' {
		p.pos++
	}
	if p.pos >= len(p.input) {
		return "", ErrInvalidSNBT
	}
	name := p.input[start:p.pos]
	p.pos++
	return "\\N{" + name + "}", nil
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

	if len(s) < 2 {
		return p.parseSimpleValue(s)
	}

	last := s[len(s)-1]
	secondLast := s[len(s)-2]

	if (secondLast == 's' || secondLast == 'S' || secondLast == 'u' || secondLast == 'U') &&
		(last == 'b' || last == 'B' || last == 's' || last == 'S' ||
			last == 'l' || last == 'L' || last == 'f' || last == 'F' || last == 'd' || last == 'D') {
		numPart := s[:len(s)-2]
		sign := secondLast
		suffix := last
		return p.parseNumberWithSign(numPart, sign, suffix)
	}

	if last == 'b' || last == 'B' || last == 's' || last == 'S' ||
		last == 'l' || last == 'L' || last == 'f' || last == 'F' || last == 'd' || last == 'D' {
		numPart := s[:len(s)-1]
		return p.parseNumber(numPart, last)
	}

	return p.parseSimpleValue(s)
}

func (p *snbtParser) parseSimpleValue(s string) (any, error) {
	if v, err := p.parseNumber(s, 0); err == nil {
		return v, nil
	}
	if s == "true" || s == "false" {
		return s == "true", nil
	}
	return s, nil
}

func (p *snbtParser) parseNumberWithSign(numPart string, sign byte, suffix byte) (any, error) {
	base := 10
	if strings.HasPrefix(numPart, "0x") || strings.HasPrefix(numPart, "0X") {
		base = 16
		numPart = numPart[2:]
	} else if strings.HasPrefix(numPart, "0b") || strings.HasPrefix(numPart, "0B") {
		base = 2
		numPart = numPart[2:]
	}

	isUnsigned := sign == 'u' || sign == 'U'

	switch suffix {
	case 'b', 'B':
		v, _ := strconv.ParseInt(numPart, base, 8)
		if isUnsigned {
			return uint8(v), nil
		}
		return int8(v), nil
	case 's', 'S':
		v, _ := strconv.ParseInt(numPart, base, 16)
		if isUnsigned {
			return uint16(v), nil
		}
		return int16(v), nil
	case 'l', 'L':
		v, _ := strconv.ParseInt(numPart, base, 64)
		if isUnsigned {
			return uint64(v), nil
		}
		return int64(v), nil
	}
	return p.parseNumber(string(suffix), suffix)
}

func (p *snbtParser) parseNumber(s string, suffix byte) (any, error) {
	s = strings.ReplaceAll(s, "_", "")

	base := 10
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		base = 16
		s = s[2:]
	} else if strings.HasPrefix(s, "0b") || strings.HasPrefix(s, "0B") {
		base = 2
		s = s[2:]
	}

	switch suffix {
	case 'b', 'B':
		return parseSignedInteger[int8](s, base, 8)
	case 's', 'S':
		return parseSignedInteger[int16](s, base, 16)
	case 'l', 'L':
		return parseSignedInteger[int64](s, base, 64)
	case 'f', 'F':
		return parseFloat[float32](s, 32)
	case 'd', 'D':
		return parseFloat[float64](s, 64)
	default:
		if v, err := parseIntegerAuto(s, base); err == nil {
			return v, nil
		}
		if v, err := parseScientific(s); err == nil {
			return v, nil
		}
		return nil, fmt.Errorf("invalid number: %s", s)
	}
}

func parseIntegerAuto(s string, base int) (any, error) {
	v, err := strconv.ParseInt(s, base, 64)
	if err != nil {
		return nil, err
	}
	if v >= math.MinInt8 && v <= math.MaxInt8 {
		return int8(v), nil
	}
	if v >= math.MinInt16 && v <= math.MaxInt16 {
		return int16(v), nil
	}
	if v >= math.MinInt32 && v <= math.MaxInt32 {
		return int32(v), nil
	}
	return v, nil
}

func parseScientific(s string) (any, error) {
	if !strings.ContainsAny(s, "eE") {
		return nil, fmt.Errorf("not a scientific number")
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	if math.IsInf(v, 0) || math.IsNaN(v) {
		return nil, fmt.Errorf("number out of range: %s", s)
	}
	return v, nil
}

func parseInteger[T int8 | int16 | int64](s string, bits int) (T, error) {
	v, err := strconv.ParseInt(s, 10, bits)
	return T(v), err
}

func parseSignedInteger[T int8 | int16 | int64](s string, base int, bits int) (T, error) {
	v, err := strconv.ParseInt(s, base, bits)
	if err != nil {
		return 0, err
	}
	return T(v), nil
}

func parseFloat[T float32 | float64](s string, bits int) (T, error) {
	v, err := strconv.ParseFloat(s, bits)
	return T(v), err
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
	return formatArray(b, 'B', func(v byte) string { return fmt.Sprintf("%db", int8(v)) })
}

func formatIntArray(arr []int32) string {
	return formatArray(arr, 'I', func(v int32) string { return fmt.Sprintf("%d", v) })
}

func formatLongArray(arr []int64) string {
	return formatArray(arr, 'L', func(v int64) string { return fmt.Sprintf("%dL", v) })
}

func formatArray[T any](arr []T, prefix byte, formatFunc func(T) string) string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	buf.WriteByte(prefix)
	buf.WriteByte(';')
	for i, v := range arr {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(formatFunc(v))
	}
	buf.WriteByte(']')
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
