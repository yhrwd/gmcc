package nbt

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"io"
	"math"
	"reflect"
)

// Marshal encodes v into NBT format
func Marshal(v any) ([]byte, error) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).Encode(v, ""); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Encoder writes NBT to an io.Writer
type Encoder struct {
	w             io.Writer
	networkFormat bool
}

// NewEncoder creates a new encoder
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// NetworkFormat enables network format
func (e *Encoder) NetworkFormat(v bool) *Encoder {
	e.networkFormat = v
	return e
}

// Encode encodes v into NBT with the given name
func (e *Encoder) Encode(v any, name string) error {
	val := reflect.ValueOf(v)
	tagType := e.tagTypeOf(val)

	if tagType == TagEnd {
		return e.writeTag(TagEnd, name)
	}

	if !e.networkFormat {
		if err := e.writeTag(tagType, name); err != nil {
			return err
		}
	}

	return e.encodeValue(val, tagType)
}

func (e *Encoder) encodeValue(val reflect.Value, tagType byte) error {
	// Handle interface and pointer indirection
	for val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	if m, ok := val.Interface().(Marshaler); ok {
		data, err := m.MarshalNBT()
		if err != nil {
			return err
		}
		return e.writeBytes(data)
	}

	if tm, ok := val.Interface().(encoding.TextMarshaler); ok && tagType == TagString {
		data, err := tm.MarshalText()
		if err != nil {
			return err
		}
		return e.writeString(string(data))
	}

	switch tagType {
	case TagByte:
		var v int8
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v = int8(val.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v = int8(val.Uint())
		case reflect.Bool:
			if val.Bool() {
				v = 1
			}
		}
		return e.writeInt8(v)
	case TagShort:
		return e.writeInt16(int16(val.Int()))
	case TagInt:
		return e.writeInt32(int32(val.Int()))
	case TagLong:
		return e.writeInt64(val.Int())
	case TagFloat:
		return e.writeInt32(int32(math.Float32bits(float32(val.Float()))))
	case TagDouble:
		return e.writeInt64(int64(math.Float64bits(val.Float())))
	case TagByteArray:
		return e.writeByteArray(val.Bytes())
	case TagString:
		return e.writeString(val.String())
	case TagList:
		return e.encodeList(val)
	case TagCompound:
		return e.encodeCompound(val)
	case TagIntArray:
		return e.encodeIntArray(val)
	case TagLongArray:
		return e.encodeLongArray(val)
	default:
		return nil
	}
}

func (e *Encoder) tagTypeOf(val reflect.Value) byte {
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return TagEnd
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Bool:
		return TagByte
	case reflect.Int, reflect.Int8:
		return TagByte
	case reflect.Int16:
		return TagShort
	case reflect.Int32:
		return TagInt
	case reflect.Int64:
		return TagLong
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return TagLong
	case reflect.Float32:
		return TagFloat
	case reflect.Float64:
		return TagDouble
	case reflect.String:
		return TagString
	case reflect.Slice:
		if val.Type().Elem().Kind() == reflect.Uint8 {
			return TagByteArray
		}
		if val.Type().Elem().Kind() == reflect.Int32 {
			return TagIntArray
		}
		if val.Type().Elem().Kind() == reflect.Int64 {
			return TagLongArray
		}
		return TagList
	case reflect.Array:
		return e.tagTypeOf(val.Slice(0, val.Len()))
	case reflect.Map:
		return TagCompound
	case reflect.Struct:
		return TagCompound
	default:
		return TagEnd
	}
}

func (e *Encoder) encodeList(val reflect.Value) error {
	if val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	length := val.Len()
	if length == 0 {
		if err := e.writeByte(TagEnd); err != nil {
			return err
		}
		return e.writeInt32(0)
	}

	elemType := e.tagTypeOf(val.Index(0))
	if err := e.writeByte(elemType); err != nil {
		return err
	}
	if err := e.writeInt32(int32(length)); err != nil {
		return err
	}

	for i := 0; i < length; i++ {
		if err := e.encodeValue(val.Index(i), elemType); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeCompound(val reflect.Value) error {
	if val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Map:
		for _, key := range val.MapKeys() {
			name := key.String()
			elem := val.MapIndex(key)
			tagType := e.tagTypeOf(elem)
			if err := e.writeTag(tagType, name); err != nil {
				return err
			}
			if err := e.encodeValue(elem, tagType); err != nil {
				return err
			}
		}
	case reflect.Struct:
		fs := cachedFields(val.Type())
		for _, f := range fs.list {
			fv := val
			for _, i := range f.index {
				if fv.Kind() == reflect.Ptr {
					if fv.IsNil() {
						break
					}
					fv = fv.Elem()
				}
				fv = fv.Field(i)
			}
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				continue
			}
			tagType := e.tagTypeOf(fv)
			if tagType == TagEnd {
				continue
			}
			if err := e.writeTag(tagType, f.name); err != nil {
				return err
			}
			if err := e.encodeValue(fv, tagType); err != nil {
				return err
			}
		}
	}
	return e.writeByte(TagEnd)
}

func (e *Encoder) encodeIntArray(val reflect.Value) error {
	length := val.Len()
	if err := e.writeInt32(int32(length)); err != nil {
		return err
	}
	for i := 0; i < length; i++ {
		if err := e.writeInt32(int32(val.Index(i).Int())); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeLongArray(val reflect.Value) error {
	length := val.Len()
	if err := e.writeInt32(int32(length)); err != nil {
		return err
	}
	for i := 0; i < length; i++ {
		if err := e.writeInt64(val.Index(i).Int()); err != nil {
			return err
		}
	}
	return nil
}

// Write methods
func (e *Encoder) writeByte(b byte) error {
	_, err := e.w.Write([]byte{b})
	return err
}

func (e *Encoder) writeBytes(b []byte) error {
	_, err := e.w.Write(b)
	return err
}

func (e *Encoder) writeInt8(v int8) error {
	return e.writeByte(byte(v))
}

func (e *Encoder) writeInt16(v int16) error {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], uint16(v))
	return e.writeBytes(b[:])
}

func (e *Encoder) writeInt32(v int32) error {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(v))
	return e.writeBytes(b[:])
}

func (e *Encoder) writeInt64(v int64) error {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(v))
	return e.writeBytes(b[:])
}

func (e *Encoder) writeString(s string) error {
	b := []byte(s)
	if err := e.writeInt16(int16(len(b))); err != nil {
		return err
	}
	return e.writeBytes(b)
}

func (e *Encoder) writeByteArray(b []byte) error {
	if err := e.writeInt32(int32(len(b))); err != nil {
		return err
	}
	return e.writeBytes(b)
}

func (e *Encoder) writeTag(tagType byte, name string) error {
	if err := e.writeByte(tagType); err != nil {
		return err
	}
	if tagType != TagEnd {
		return e.writeString(name)
	}
	return nil
}
