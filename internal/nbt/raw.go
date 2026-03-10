package nbt

import (
	"bytes"
	"encoding"
	"reflect"
)

func init() {
	// Register RawMessage handling
}

// UnmarshalNBT implements Unmarshaler for RawMessage
func (m *RawMessage) UnmarshalNBT(data []byte) error {
	m.Data = data
	return nil
}

// MarshalNBT implements Marshaler for RawMessage
func (m RawMessage) MarshalNBT() ([]byte, error) {
	return m.Data, nil
}

// Decode decodes the raw NBT data into v
func (m RawMessage) Decode(v any) error {
	return Unmarshal(m.Data, v)
}

// decodeRawMessage handles RawMessage during unmarshaling
func (d *Decoder) decodeRawMessage(tagType byte) (RawMessage, error) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	// Write tag type and read the value
	if err := enc.writeByte(tagType); err != nil {
		return RawMessage{}, err
	}

	// Read the tag data from decoder and write to buffer
	if err := d.copyTag(&buf, tagType); err != nil {
		return RawMessage{}, err
	}

	return RawMessage{Type: tagType, Data: buf.Bytes()}, nil
}

func (d *Decoder) copyTag(w *bytes.Buffer, tagType byte) error {
	switch tagType {
	case TagEnd:
		return nil
	case TagByte:
		b, err := d.readByte()
		if err != nil {
			return err
		}
		return w.WriteByte(b)
	case TagShort:
		v, err := d.readInt16()
		if err != nil {
			return err
		}
		var buf [2]byte
		buf[0] = byte(v >> 8)
		buf[1] = byte(v)
		_, err = w.Write(buf[:])
		return err
	case TagInt, TagFloat:
		v, err := d.readInt32()
		if err != nil {
			return err
		}
		var buf [4]byte
		buf[0] = byte(v >> 24)
		buf[1] = byte(v >> 16)
		buf[2] = byte(v >> 8)
		buf[3] = byte(v)
		_, err = w.Write(buf[:])
		return err
	case TagLong, TagDouble:
		v, err := d.readInt64()
		if err != nil {
			return err
		}
		var buf [8]byte
		for i := 7; i >= 0; i-- {
			buf[7-i] = byte(v >> (i * 8))
		}
		_, err = w.Write(buf[:])
		return err
	case TagByteArray:
		length, err := d.readInt32()
		if err != nil {
			return err
		}
		var lbuf [4]byte
		lbuf[0] = byte(length >> 24)
		lbuf[1] = byte(length >> 16)
		lbuf[2] = byte(length >> 8)
		lbuf[3] = byte(length)
		w.Write(lbuf[:])
		data := make([]byte, length)
		if _, err := d.r.Read(data); err != nil {
			return err
		}
		_, err = w.Write(data)
		return err
	case TagString:
		s, err := d.readString()
		if err != nil {
			return err
		}
		data := []byte(s)
		var lbuf [2]byte
		lbuf[0] = byte(len(data) >> 8)
		lbuf[1] = byte(len(data))
		w.Write(lbuf[:])
		_, err = w.Write(data)
		return err
	case TagList:
		elemType, err := d.readByte()
		if err != nil {
			return err
		}
		if err := w.WriteByte(elemType); err != nil {
			return err
		}
		length, err := d.readInt32()
		if err != nil {
			return err
		}
		var lbuf [4]byte
		lbuf[0] = byte(length >> 24)
		lbuf[1] = byte(length >> 16)
		lbuf[2] = byte(length >> 8)
		lbuf[3] = byte(length)
		w.Write(lbuf[:])
		for i := int32(0); i < length; i++ {
			if err := d.copyTag(w, elemType); err != nil {
				return err
			}
		}
		return nil
	case TagCompound:
		for {
			tt, err := d.readByte()
			if err != nil {
				return err
			}
			if err := w.WriteByte(tt); err != nil {
				return err
			}
			if tt == TagEnd {
				break
			}
			name, err := d.readString()
			if err != nil {
				return err
			}
			data := []byte(name)
			var lbuf [2]byte
			lbuf[0] = byte(len(data) >> 8)
			lbuf[1] = byte(len(data))
			w.Write(lbuf[:])
			w.Write(data)
			if err := d.copyTag(w, tt); err != nil {
				return err
			}
		}
		return nil
	case TagIntArray:
		length, err := d.readInt32()
		if err != nil {
			return err
		}
		var lbuf [4]byte
		lbuf[0] = byte(length >> 24)
		lbuf[1] = byte(length >> 16)
		lbuf[2] = byte(length >> 8)
		lbuf[3] = byte(length)
		w.Write(lbuf[:])
		for i := int32(0); i < length; i++ {
			v, err := d.readInt32()
			if err != nil {
				return err
			}
			var buf [4]byte
			buf[0] = byte(v >> 24)
			buf[1] = byte(v >> 16)
			buf[2] = byte(v >> 8)
			buf[3] = byte(v)
			w.Write(buf[:])
		}
		return nil
	case TagLongArray:
		length, err := d.readInt32()
		if err != nil {
			return err
		}
		var lbuf [4]byte
		lbuf[0] = byte(length >> 24)
		lbuf[1] = byte(length >> 16)
		lbuf[2] = byte(length >> 8)
		lbuf[3] = byte(length)
		w.Write(lbuf[:])
		for i := int32(0); i < length; i++ {
			v, err := d.readInt64()
			if err != nil {
				return err
			}
			var buf [8]byte
			for j := 7; j >= 0; j-- {
				buf[7-j] = byte(v >> (j * 8))
			}
			w.Write(buf[:])
		}
		return nil
	default:
		return d.errorf("unknown tag type: 0x%02X", tagType)
	}
}

// indirect handles pointer and interface indirection
func indirectRaw(val reflect.Value) (Marshaler, encoding.TextMarshaler, reflect.Value, func()) {
	// Simplified version for RawMessage - just return the value
	return nil, nil, val, nil
}
