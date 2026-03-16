package i18n

import (
	"strings"
	"sync"
)

type I18n struct {
	mu       sync.RWMutex
	lang     string
	messages map[string]string
}

var defaultI18n *I18n
var once sync.Once

func GetI18n() *I18n {
	once.Do(func() {
		defaultI18n = newI18n()
	})
	return defaultI18n
}

func newI18n() *I18n {
	i := &I18n{
		lang:     "en_us",
		messages: make(map[string]string),
	}
	i.loadLanguage("en_us")
	return i
}

func (i *I18n) SetLanguage(lang string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.lang = strings.ToLower(lang)
	i.messages = make(map[string]string)
	i.loadLanguage(i.lang)
}

func (i *I18n) GetLanguage() string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.lang
}

func (i *I18n) Translate(key string, args ...any) string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if msg, ok := i.messages[key]; ok {
		return i.format(msg, args...)
	}
	return key
}

func (i *I18n) TranslateWithFallback(key string, fallback string, args ...any) string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if msg, ok := i.messages[key]; ok {
		return i.format(msg, args...)
	}
	return i.format(fallback, args...)
}

func (i *I18n) Has(key string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	_, ok := i.messages[key]
	return ok
}

func (i *I18n) format(template string, args ...any) string {
	if len(args) == 0 {
		return template
	}

	result := template
	for idx, arg := range args {
		old := result
		key := "%" + string(rune('1'+idx))
		if strings.Contains(result, key) {
			result = strings.ReplaceAll(result, key, formatArg(arg))
			continue
		}
		key = "%" + string(rune('0'+idx+1))
		if strings.Contains(result, key) {
			result = strings.ReplaceAll(result, key, formatArg(arg))
			continue
		}
		if strings.Contains(old, "%s") {
			result = strings.Replace(old, "%s", formatArg(arg), 1)
		} else if strings.Contains(old, "%d") {
			result = strings.Replace(old, "%d", formatArg(arg), 1)
		}
	}
	return result
}

func formatArg(arg any) string {
	switch v := arg.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64:
		return formatInt(v)
	case uint, uint8, uint16, uint32, uint64:
		return formatUint(v)
	case float32, float64:
		return formatFloat(v)
	default:
		return ""
	}
}

func formatInt(v any) string {
	switch n := v.(type) {
	case int:
		return itoa(int64(n))
	case int8:
		return itoa(int64(n))
	case int16:
		return itoa(int64(n))
	case int32:
		return itoa(int64(n))
	case int64:
		return itoa(n)
	}
	return ""
}

func formatUint(v any) string {
	switch n := v.(type) {
	case uint:
		return utoa(uint64(n))
	case uint8:
		return utoa(uint64(n))
	case uint16:
		return utoa(uint64(n))
	case uint32:
		return utoa(uint64(n))
	case uint64:
		return utoa(n)
	}
	return ""
}

func formatFloat(v any) string {
	switch n := v.(type) {
	case float32:
		return ftoa(float64(n), 2)
	case float64:
		return ftoa(n, 2)
	}
	return ""
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}

func utoa(n uint64) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}

func ftoa(n float64, precision int) string {
	intPart := int64(n)
	fracPart := n - float64(intPart)
	if fracPart < 0 {
		fracPart = -fracPart
	}

	result := itoa(intPart)
	if precision > 0 {
		mul := float64(10)
		for i := 1; i < precision; i++ {
			mul *= 10
		}
		fracInt := int64(fracPart*mul + 0.5)
		fracStr := utoa(uint64(fracInt))
		for len(fracStr) < precision {
			fracStr = "0" + fracStr
		}
		result += "." + fracStr
	}
	return result
}

func Translate(key string, args ...any) string {
	return GetI18n().Translate(key, args...)
}

func TranslateWithFallback(key string, fallback string, args ...any) string {
	return GetI18n().TranslateWithFallback(key, fallback, args...)
}

func (i *I18n) ItemName(name string) string {
	key := "item.minecraft." + name
	if msg := i.Translate(key); msg != key {
		return msg
	}
	key = "block.minecraft." + name
	return i.Translate(key)
}

func (i *I18n) BlockName(name string) string {
	return i.Translate("block.minecraft." + name)
}

func (i *I18n) EntityName(name string) string {
	return i.Translate("entity.minecraft." + name)
}

func (i *I18n) EffectName(name string) string {
	return i.Translate("effect.minecraft." + name)
}

func (i *I18n) EnchantmentName(name string) string {
	return i.Translate("enchantment.minecraft." + name)
}

func ItemName(name string) string {
	return GetI18n().ItemName(name)
}

func BlockName(name string) string {
	return GetI18n().BlockName(name)
}

func EntityName(name string) string {
	return GetI18n().EntityName(name)
}

func EffectName(name string) string {
	return GetI18n().EffectName(name)
}

func EnchantmentName(name string) string {
	return GetI18n().EnchantmentName(name)
}
