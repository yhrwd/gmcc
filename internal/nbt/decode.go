package nbt

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
	"sync"
	"unicode/utf16"
	"unicode/utf8"
)

// Unmarshal decodes NBT data into v
func Unmarshal(data []byte, v any) error {
	return NewDecoder(bytes.NewReader(data)).Decode(v)
}

// Decoder reads NBT from an io.Reader
type Decoder struct {
	r               io.Reader
	networkFormat   bool
	offset          int64
	fieldPath       []string
	disallowUnknown bool
}

// NewDecoder creates a new decoder
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// NetworkFormat enables network format (no root name)
func (d *Decoder) NetworkFormat(v bool) *Decoder {
	d.networkFormat = v
	return d
}

// DisallowUnknownFields prevents unknown fields
func (d *Decoder) DisallowUnknownFields() *Decoder {
	d.disallowUnknown = true
	return d
}

// Decode decodes NBT into v, returns root tag name
func (d *Decoder) Decode(v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return d.errorf("non-pointer passed to Decode")
	}

	tagType, err := d.readByte()
	if err != nil {
		return d.errorWrap(err)
	}

	if tagType == TagEnd {
		return ErrEND
	}

	if !d.networkFormat {
		if _, err := d.readString(); err != nil {
			return d.errorWrap(err)
		}
	}

	return d.unmarshal(val, tagType)
}

func (d *Decoder) unmarshal(val reflect.Value, tagType byte) error {
	u, t, val, assign := indirect(val, tagType == TagEnd)
	if assign != nil {
		defer assign()
	}
	if u != nil {
		data, err := io.ReadAll(d.r)
		if err != nil {
			return err
		}
		return u.UnmarshalNBT(data)
	}

	switch tagType {
	case TagEnd:
		return ErrEND
	case TagByte:
		v, err := d.readInt8()
		if err != nil {
			return err
		}
		return d.setValue(val, v)
	case TagShort:
		v, err := d.readInt16()
		if err != nil {
			return err
		}
		return d.setValue(val, v)
	case TagInt:
		v, err := d.readInt32()
		if err != nil {
			return err
		}
		return d.setValue(val, v)
	case TagLong:
		v, err := d.readInt64()
		if err != nil {
			return err
		}
		return d.setValue(val, v)
	case TagFloat:
		v, err := d.readInt32()
		if err != nil {
			return err
		}
		return d.setValue(val, math.Float32frombits(uint32(v)))
	case TagDouble:
		v, err := d.readInt64()
		if err != nil {
			return err
		}
		return d.setValue(val, math.Float64frombits(uint64(v)))
	case TagByteArray:
		v, err := d.readByteArray()
		if err != nil {
			return err
		}
		return d.setValue(val, v)
	case TagString:
		v, err := d.readString()
		if err != nil {
			return err
		}
		if t != nil {
			return t.UnmarshalText([]byte(v))
		}
		return d.setValue(val, v)
	case TagList:
		return d.readList(val)
	case TagCompound:
		return d.readCompound(val)
	case TagIntArray:
		v, err := d.readIntArray()
		if err != nil {
			return err
		}
		return d.setValue(val, v)
	case TagLongArray:
		v, err := d.readLongArray()
		if err != nil {
			return err
		}
		return d.setValue(val, v)
	default:
		return d.errorf("unknown tag type: 0x%02X", tagType)
	}
}

func (d *Decoder) setValue(val reflect.Value, v any) error {
	vk := val.Kind()
	if vk == reflect.Interface {
		val.Set(reflect.ValueOf(v))
		return nil
	}
	if vk == reflect.Ptr {
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		return d.setValue(val.Elem(), v)
	}
	rv := reflect.ValueOf(v)
	if rv.Type().ConvertibleTo(val.Type()) {
		val.Set(rv.Convert(val.Type()))
		return nil
	}
	return d.errorf("cannot set %v to %v", rv.Type(), val.Type())
}

// Read methods
func (d *Decoder) readByte() (byte, error) {
	var b [1]byte
	_, err := io.ReadFull(d.r, b[:])
	d.offset++
	return b[0], err
}

func (d *Decoder) readInt8() (int8, error) {
	b, err := d.readByte()
	return int8(b), err
}

func (d *Decoder) readInt16() (int16, error) {
	var b [2]byte
	_, err := io.ReadFull(d.r, b[:])
	d.offset += 2
	return int16(binary.BigEndian.Uint16(b[:])), err
}

func (d *Decoder) readInt32() (int32, error) {
	var b [4]byte
	_, err := io.ReadFull(d.r, b[:])
	d.offset += 4
	return int32(binary.BigEndian.Uint32(b[:])), err
}

func (d *Decoder) readInt64() (int64, error) {
	var b [8]byte
	_, err := io.ReadFull(d.r, b[:])
	d.offset += 8
	return int64(binary.BigEndian.Uint64(b[:])), err
}

func (d *Decoder) readString() (string, error) {
	length, err := d.readInt16()
	if err != nil {
		return "", err
	}
	if length < 0 {
		return "", d.errorf("string length < 0: %d", length)
	}
	if length == 0 {
		return "", nil
	}
	buf := make([]byte, length)
	_, err = io.ReadFull(d.r, buf)
	d.offset += int64(length)
	if err != nil {
		return "", err
	}
	return cesu8ToUTF8(buf), nil
}

func (d *Decoder) readByteArray() ([]byte, error) {
	length, err := d.readInt32()
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, d.errorf("byte array length < 0: %d", length)
	}
	if length == 0 {
		return []byte{}, nil
	}
	buf := make([]byte, length)
	_, err = io.ReadFull(d.r, buf)
	d.offset += int64(length)
	return buf, err
}

func (d *Decoder) readIntArray() ([]int32, error) {
	length, err := d.readInt32()
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, d.errorf("int array length < 0: %d", length)
	}
	result := make([]int32, length)
	for i := int32(0); i < length; i++ {
		result[i], err = d.readInt32()
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (d *Decoder) readLongArray() ([]int64, error) {
	length, err := d.readInt32()
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, d.errorf("long array length < 0: %d", length)
	}
	result := make([]int64, length)
	for i := int32(0); i < length; i++ {
		result[i], err = d.readInt64()
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (d *Decoder) readList(val reflect.Value) error {
	elemType, err := d.readByte()
	if err != nil {
		return err
	}
	length, err := d.readInt32()
	if err != nil {
		return err
	}
	if length < 0 {
		return d.errorf("list length < 0: %d", length)
	}

	vk := val.Kind()
	switch vk {
	case reflect.Interface:
		slice := make([]any, length)
		for i := int32(0); i < length; i++ {
			v, err := d.decodeValue(elemType)
			if err != nil {
				return err
			}
			slice[i] = v
		}
		val.Set(reflect.ValueOf(slice))
		return nil
	case reflect.Slice:
		slice := reflect.MakeSlice(val.Type(), int(length), int(length))
		for i := 0; i < int(length); i++ {
			if err := d.unmarshal(slice.Index(i), elemType); err != nil {
				return err
			}
		}
		val.Set(slice)
		return nil
	case reflect.Array:
		for i := 0; i < int(length) && i < val.Len(); i++ {
			if err := d.unmarshal(val.Index(i), elemType); err != nil {
				return err
			}
		}
		return nil
	default:
		return d.errorf("cannot decode list to %v", vk)
	}
}

func (d *Decoder) readCompound(val reflect.Value) error {
	vk := val.Kind()
	switch vk {
	case reflect.Struct:
		return d.readCompoundStruct(val)
	case reflect.Map:
		return d.readCompoundMap(val)
	case reflect.Interface:
		m := make(map[string]any)
		if err := d.readCompoundMap(reflect.ValueOf(&m).Elem()); err != nil {
			return err
		}
		val.Set(reflect.ValueOf(m))
		return nil
	default:
		return d.errorf("cannot decode compound to %v", vk)
	}
}

func (d *Decoder) readCompoundStruct(val reflect.Value) error {
	fields := cachedFields(val.Type())
	for {
		tagType, err := d.readByte()
		if err != nil {
			return err
		}
		if tagType == TagEnd {
			break
		}
		name, err := d.readString()
		if err != nil {
			return err
		}
		d.pushField(name)
		f := fields.find(name)
		if f != nil {
			fv := val
			for _, i := range f.index {
				if fv.Kind() == reflect.Ptr {
					if fv.IsNil() {
						fv.Set(reflect.New(fv.Type().Elem()))
					}
					fv = fv.Elem()
				}
				fv = fv.Field(i)
			}
			if err := d.unmarshal(fv, tagType); err != nil {
				d.popField()
				return err
			}
		} else if d.disallowUnknown {
			d.popField()
			return d.errorf("unknown field %q", name)
		} else {
			if err := d.skipValue(tagType); err != nil {
				d.popField()
				return err
			}
		}
		d.popField()
	}
	return nil
}

func (d *Decoder) readCompoundMap(val reflect.Value) error {
	if val.IsNil() {
		val.Set(reflect.MakeMap(val.Type()))
	}
	keyType := val.Type().Key()
	if keyType.Kind() != reflect.String {
		return d.errorf("map key must be string")
	}
	elemType := val.Type().Elem()
	for {
		tagType, err := d.readByte()
		if err != nil {
			return err
		}
		if tagType == TagEnd {
			break
		}
		name, err := d.readString()
		if err != nil {
			return err
		}
		d.pushField(name)
		elem := reflect.New(elemType).Elem()
		if err := d.unmarshal(elem, tagType); err != nil {
			d.popField()
			return err
		}
		val.SetMapIndex(reflect.ValueOf(name), elem)
		d.popField()
	}
	return nil
}

func (d *Decoder) decodeValue(tagType byte) (any, error) {
	switch tagType {
	case TagEnd:
		return nil, nil
	case TagByte:
		return d.readInt8()
	case TagShort:
		return d.readInt16()
	case TagInt:
		return d.readInt32()
	case TagLong:
		return d.readInt64()
	case TagFloat:
		v, err := d.readInt32()
		return math.Float32frombits(uint32(v)), err
	case TagDouble:
		v, err := d.readInt64()
		return math.Float64frombits(uint64(v)), err
	case TagByteArray:
		return d.readByteArray()
	case TagString:
		return d.readString()
	case TagList:
		return d.decodeList()
	case TagCompound:
		return d.decodeCompound()
	case TagIntArray:
		return d.readIntArray()
	case TagLongArray:
		return d.readLongArray()
	default:
		return nil, d.errorf("unknown tag type: 0x%02X", tagType)
	}
}

func (d *Decoder) decodeList() ([]any, error) {
	elemType, err := d.readByte()
	if err != nil {
		return nil, err
	}
	length, err := d.readInt32()
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, d.errorf("list length < 0: %d", length)
	}
	result := make([]any, length)
	for i := int32(0); i < length; i++ {
		result[i], err = d.decodeValue(elemType)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (d *Decoder) decodeCompound() (map[string]any, error) {
	m := make(map[string]any)
	for {
		tagType, err := d.readByte()
		if err != nil {
			return nil, err
		}
		if tagType == TagEnd {
			break
		}
		name, err := d.readString()
		if err != nil {
			return nil, err
		}
		d.pushField(name)
		v, err := d.decodeValue(tagType)
		if err != nil {
			d.popField()
			return nil, err
		}
		m[name] = v
		d.popField()
	}
	return m, nil
}

func (d *Decoder) skipValue(tagType byte) error {
	switch tagType {
	case TagEnd:
		return nil
	case TagByte:
		_, err := d.readByte()
		return err
	case TagShort:
		_, err := d.readInt16()
		return err
	case TagInt, TagFloat:
		_, err := d.readInt32()
		return err
	case TagLong, TagDouble:
		_, err := d.readInt64()
		return err
	case TagByteArray:
		length, err := d.readInt32()
		if err != nil {
			return err
		}
		for i := int32(0); i < length; i++ {
			if _, err := d.readByte(); err != nil {
				return err
			}
		}
	case TagString:
		length, err := d.readInt16()
		if err != nil {
			return err
		}
		for i := int16(0); i < length; i++ {
			if _, err := d.readByte(); err != nil {
				return err
			}
		}
	case TagList:
		elemType, err := d.readByte()
		if err != nil {
			return err
		}
		length, err := d.readInt32()
		if err != nil {
			return err
		}
		for i := int32(0); i < length; i++ {
			if err := d.skipValue(elemType); err != nil {
				return err
			}
		}
	case TagCompound:
		for {
			tt, err := d.readByte()
			if err != nil {
				return err
			}
			if tt == TagEnd {
				break
			}
			if _, err := d.readString(); err != nil {
				return err
			}
			if err := d.skipValue(tt); err != nil {
				return err
			}
		}
	case TagIntArray:
		length, err := d.readInt32()
		if err != nil {
			return err
		}
		for i := int32(0); i < length; i++ {
			if _, err := d.readInt32(); err != nil {
				return err
			}
		}
	case TagLongArray:
		length, err := d.readInt32()
		if err != nil {
			return err
		}
		for i := int32(0); i < length; i++ {
			if _, err := d.readInt64(); err != nil {
				return err
			}
		}
	default:
		return d.errorf("unknown tag type: 0x%02X", tagType)
	}
	return nil
}

func (d *Decoder) pushField(name string) {
	d.fieldPath = append(d.fieldPath, name)
}

func (d *Decoder) popField() {
	if len(d.fieldPath) > 0 {
		d.fieldPath = d.fieldPath[:len(d.fieldPath)-1]
	}
}

func (d *Decoder) errorf(format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	if len(d.fieldPath) > 0 {
		msg = strings.Join(d.fieldPath, ".") + ": " + msg
	}
	return &SyntaxError{Offset: d.offset, Message: msg, FieldPath: d.fieldPath}
}

func (d *Decoder) errorWrap(err error) error {
	if len(d.fieldPath) > 0 {
		msg := strings.Join(d.fieldPath, ".") + ": " + err.Error()
		return &SyntaxError{Offset: d.offset, Message: msg, FieldPath: d.fieldPath}
	}
	return err
}

func indirect(v reflect.Value, decodingNull bool) (Unmarshaler, encoding.TextUnmarshaler, reflect.Value, func()) {
	v0 := v
	haveAddr := false
	var assign func()

	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		haveAddr = true
		v = v.Addr()
	}

	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) {
				haveAddr = false
				v = e
				continue
			} else if v.CanSet() {
				e = reflect.New(e.Type())
				assign = func() { v0.Set(e.Elem()) }
				v = e
				continue
			}
		}
		if v.Kind() != reflect.Ptr {
			break
		}
		if decodingNull && v.CanSet() {
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Type().NumMethod() > 0 && v.CanInterface() {
			if u, ok := v.Interface().(Unmarshaler); ok {
				return u, nil, reflect.Value{}, assign
			}
			if u, ok := v.Interface().(encoding.TextUnmarshaler); ok {
				return nil, u, v, assign
			}
		}
		if haveAddr {
			v = v0
			haveAddr = false
		} else {
			v = v.Elem()
		}
	}
	return nil, nil, v, assign
}

func cesu8ToUTF8(data []byte) string {
	return CESU8ToUTF8(data)
}

// CESU8ToUTF8 将 CESU-8（Modified UTF-8）转换为标准 UTF-8。
// Minecraft Java 使用 CESU-8 编码字符串。
func CESU8ToUTF8(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	hasCESU8 := false
	for i := 0; i < len(data)-5; i++ {
		if data[i] == 0xED && (data[i+1]&0xF0) == 0xA0 {
			hasCESU8 = true
			break
		}
	}
	if !hasCESU8 {
		return string(data)
	}
	var result []byte
	i := 0
	for i < len(data) {
		if i+5 < len(data) && data[i] == 0xED {
			b1, b2 := data[i+1], data[i+2]
			b3, b4, b5 := data[i+3], data[i+4], data[i+5]
			if (b1&0xF0) == 0xA0 && (b2&0xC0) == 0x80 && b3 == 0xED && (b4&0xF0) == 0xB0 && (b5&0xC0) == 0x80 {
				high := uint16(0xD800) + uint16(b1&0x0F)<<6 + uint16(b2&0x3F)
				low := uint16(0xDC00) + uint16(b4&0x0F)<<6 + uint16(b5&0x3F)
				decoded := utf16.Decode([]uint16{high, low})
				if len(decoded) > 0 {
					var buf [4]byte
					n := utf8.EncodeRune(buf[:], decoded[0])
					result = append(result, buf[:n]...)
					i += 6
					continue
				}
			}
		}
		result = append(result, data[i])
		i++
	}
	return string(result)
}

type fieldCache struct {
	mu    sync.RWMutex
	cache map[reflect.Type]*fields
}

var globalCache = &fieldCache{cache: make(map[reflect.Type]*fields)}

type fields struct {
	list      []field
	nameIndex map[string]int
}

type field struct {
	name  string
	index []int
}

func cachedFields(t reflect.Type) *fields {
	globalCache.mu.RLock()
	f, ok := globalCache.cache[t]
	globalCache.mu.RUnlock()
	if ok {
		return f
	}
	globalCache.mu.Lock()
	defer globalCache.mu.Unlock()
	if f, ok := globalCache.cache[t]; ok {
		return f
	}
	f = newFields(t)
	globalCache.cache[t] = f
	return f
}

func newFields(t reflect.Type) *fields {
	fs := &fields{nameIndex: make(map[string]int)}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" && !f.Anonymous {
			continue
		}
		name := f.Tag.Get("nbt")
		if name == "-" {
			continue
		}
		if name == "" {
			name = f.Name
		}
		fs.list = append(fs.list, field{name: name, index: []int{i}})
		fs.nameIndex[name] = len(fs.list) - 1
	}
	return fs
}

func (fs *fields) find(name string) *field {
	if i, ok := fs.nameIndex[name]; ok {
		return &fs.list[i]
	}
	for i := range fs.list {
		if strings.EqualFold(fs.list[i].name, name) {
			return &fs.list[i]
		}
	}
	return nil
}
