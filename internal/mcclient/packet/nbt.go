package packet

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"gmcc/internal/logx"
)

// NBT Tag types
const (
	TagEnd       = 0x00
	TagByte      = 0x01
	TagShort     = 0x02
	TagInt       = 0x03
	TagLong      = 0x04
	TagFloat     = 0x05
	TagDouble    = 0x06
	TagByteArray = 0x07
	TagString    = 0x08
	TagList      = 0x09
	TagCompound  = 0x0A
	TagIntArray  = 0x0B
	TagLongArray = 0x0C
)

// SyntaxError represents a parsing error
type SyntaxError struct {
	Offset    int64
	Message   string
	FieldPath []string
}

func (e *SyntaxError) Error() string {
	if len(e.FieldPath) > 0 {
		return e.FieldPath[len(e.FieldPath)-1] + ": " + e.Message
	}
	return e.Message
}

// Unmarshaler interface for custom decoding
type Unmarshaler interface {
	UnmarshalNBT(data []byte) error
}

// Decoder reads NBT from an io.Reader
type Decoder struct {
	r             io.Reader
	networkFormat bool
	offset        int64
	fieldPath     []string
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
		return errors.New("unexpected TAG_End")
	}

	if !d.networkFormat {
		if _, err := d.readString(); err != nil {
			return d.errorWrap(err)
		}
	}

	return d.unmarshal(val, tagType)
}

// Skip skips the current NBT tag
func (d *Decoder) Skip() error {
	tagType, err := d.readByte()
	if err != nil {
		return err
	}
	return d.skipValue(tagType)
}

// Read methods
func (d *Decoder) readByte() (byte, error) {
	var b [1]byte
	_, err := io.ReadFull(d.r, b[:])
	d.offset++
	return b[0], err
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
	return CESU8ToUTF8(buf), nil
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
		return errors.New("unexpected TAG_End")
	case TagByte:
		v, err := d.readByte()
		if err != nil {
			return err
		}
		return d.setValue(val, int8(v))
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
	default:
		return d.errorf("cannot decode list to %v", vk)
	}
}

func (d *Decoder) decodeValue(tagType byte) (any, error) {
	switch tagType {
	case TagByte:
		v, err := d.readByte()
		return int8(v), err
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
	case TagString:
		return d.readString()
	case TagCompound:
		m := make(map[string]any)
		for {
			tag, err := d.readByte()
			if err != nil {
				return nil, err
			}
			if tag == TagEnd {
				break
			}
			name, err := d.readString()
			if err != nil {
				return nil, err
			}
			v, err := d.decodeValue(tag)
			if err != nil {
				return nil, err
			}
			m[name] = v
		}
		return m, nil
	default:
		return nil, d.errorf("unknown tag type: 0x%02X", tagType)
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
		d.fieldPath = append(d.fieldPath, name)
		if err := d.handleStructField(val, name, tagType); err != nil {
			return err
		}
		d.fieldPath = d.fieldPath[:len(d.fieldPath)-1]
	}
	return nil
}

func (d *Decoder) handleStructField(val reflect.Value, name string, tagType byte) error {
	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" && !f.Anonymous {
			continue
		}
		tag := f.Tag.Get("nbt")
		if tag == "-" {
			continue
		}
		if tag == "" {
			tag = f.Name
		}
		if strings.EqualFold(tag, name) {
			return d.unmarshal(val.Field(i), tagType)
		}
	}
	return d.skipValue(tagType)
}

func (d *Decoder) readCompoundMap(val reflect.Value) error {
	m := val.Convert(reflect.TypeOf(map[string]any{})).Interface().(map[string]any)
	if m == nil {
		m = make(map[string]any)
	}
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
		v, err := d.decodeValue(tagType)
		if err != nil {
			return err
		}
		m[name] = v
	}
	val.Set(reflect.ValueOf(m))
	return nil
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
	case TagInt:
		_, err := d.readInt32()
		return err
	case TagLong:
		_, err := d.readInt64()
		return err
	case TagFloat:
		_, err := d.readInt32()
		return err
	case TagDouble:
		_, err := d.readInt64()
		return err
	case TagByteArray:
		length, err := d.readInt32()
		if err != nil {
			return err
		}
		_, err = io.CopyN(io.Discard, d.r, int64(length))
		d.offset += int64(length)
		return err
	case TagString:
		_, err := d.readString()
		return err
	case TagList:
		return d.skipList()
	case TagCompound:
		return d.skipCompound()
	case TagIntArray:
		length, err := d.readInt32()
		if err != nil {
			return err
		}
		_, err = io.CopyN(io.Discard, d.r, int64(length)*4)
		d.offset += int64(length) * 4
		return err
	case TagLongArray:
		length, err := d.readInt32()
		if err != nil {
			return err
		}
		_, err = io.CopyN(io.Discard, d.r, int64(length)*8)
		d.offset += int64(length) * 8
		return err
	default:
		return d.errorf("unknown tag type: 0x%02X", tagType)
	}
}

func (d *Decoder) skipList() error {
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
	return nil
}

func (d *Decoder) skipCompound() error {
	for {
		tagType, err := d.readByte()
		if err != nil {
			return err
		}
		if tagType == TagEnd {
			break
		}
		if _, err := d.readString(); err != nil {
			return err
		}
		if err := d.skipValue(tagType); err != nil {
			return err
		}
	}
	return nil
}

func (d *Decoder) errorf(format string, args ...any) error {
	return &SyntaxError{
		Offset:    d.offset,
		Message:   fmt.Sprintf(format, args...),
		FieldPath: append([]string{}, d.fieldPath...),
	}
}

func (d *Decoder) errorWrap(err error) error {
	if _, ok := err.(*SyntaxError); ok {
		return err
	}
	return &SyntaxError{
		Offset:    d.offset,
		Message:   err.Error(),
		FieldPath: append([]string{}, d.fieldPath...),
	}
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

// CESU8ToUTF8 converts CESU-8 (Modified UTF-8) to standard UTF-8
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

// SkipNBT skips Network NBT format (no name field)
func SkipNBT(r *bytes.Reader) error {
	if r.Len() == 0 {
		return nil
	}
	dec := NewDecoder(r).NetworkFormat(true)
	err := dec.Skip()
	if err != nil {
		errMsg := err.Error()
		if errMsg == "unexpected EOF" || strings.HasPrefix(errMsg, "unknown tag type: ") {
			logx.Warnf("SkipNBT warning: %v, %d bytes remaining", err, r.Len())
			return nil
		}
		return err
	}
	return nil
}

// ReadAnonymousNBTJSON parses Network NBT and returns JSON string
func ReadAnonymousNBTJSON(r io.Reader) (string, error) {
	dec := NewDecoder(r)
	dec.NetworkFormat(true)
	var v any
	if err := dec.Decode(&v); err != nil {
		return "", err
	}

	raw, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
