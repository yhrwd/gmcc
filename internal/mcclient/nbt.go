package mcclient

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

const (
	tagEnd       = 0x00
	tagByte      = 0x01
	tagShort     = 0x02
	tagInt       = 0x03
	tagLong      = 0x04
	tagFloat     = 0x05
	tagDouble    = 0x06
	tagByteArray = 0x07
	tagString    = 0x08
	tagList      = 0x09
	tagCompound  = 0x0A
	tagIntArray  = 0x0B
	tagLongArray = 0x0C
)

type nbtDecoder struct {
	r             io.Reader
	networkFormat bool
}

func newNBTDecoder(r io.Reader, networkFormat bool) *nbtDecoder {
	return &nbtDecoder{r: r, networkFormat: networkFormat}
}

func (d *nbtDecoder) decodeRoot() (any, error) {
	tagType, err := d.readByte()
	if err != nil {
		return nil, err
	}

	if tagType == tagEnd {
		return nil, nil
	}

	if d.networkFormat {
		return d.decodeValue(tagType)
	}

	_, err = d.readString()
	if err != nil {
		return nil, err
	}

	return d.decodeValue(tagType)
}

func (d *nbtDecoder) decodeValue(tagType byte) (any, error) {
	switch tagType {
	case tagEnd:
		return nil, nil
	case tagByte:
		v, err := d.readByte()
		return int8(v), err
	case tagShort:
		v, err := d.readInt16()
		return v, err
	case tagInt:
		return d.readInt32()
	case tagLong:
		return d.readInt64()
	case tagFloat:
		v, err := d.readInt32()
		if err != nil {
			return nil, err
		}
		return float32frombits(uint32(v)), nil
	case tagDouble:
		v, err := d.readInt64()
		if err != nil {
			return nil, err
		}
		return float64frombits(uint64(v)), nil
	case tagByteArray:
		return d.readByteArray()
	case tagString:
		return d.readString()
	case tagList:
		return d.readList()
	case tagCompound:
		return d.readCompound()
	case tagIntArray:
		return d.readIntArray()
	case tagLongArray:
		return d.readLongArray()
	default:
		return nil, fmt.Errorf("unknown NBT tag type: 0x%02X", tagType)
	}
}

func (d *nbtDecoder) readByte() (byte, error) {
	var b [1]byte
	_, err := io.ReadFull(d.r, b[:])
	return b[0], err
}

func (d *nbtDecoder) readInt16() (int16, error) {
	var b [2]byte
	_, err := io.ReadFull(d.r, b[:])
	if err != nil {
		return 0, err
	}
	return int16(binary.BigEndian.Uint16(b[:])), nil
}

func (d *nbtDecoder) readInt32() (int32, error) {
	var b [4]byte
	_, err := io.ReadFull(d.r, b[:])
	if err != nil {
		return 0, err
	}
	return int32(binary.BigEndian.Uint32(b[:])), nil
}

func (d *nbtDecoder) readInt64() (int64, error) {
	var b [8]byte
	_, err := io.ReadFull(d.r, b[:])
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(b[:])), nil
}

func (d *nbtDecoder) readString() (string, error) {
	length, err := d.readInt16()
	if err != nil {
		return "", err
	}
	if length < 0 {
		return "", fmt.Errorf("NBT string length < 0: %d", length)
	}
	if length == 0 {
		return "", nil
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(d.r, buf)
	if err != nil {
		return "", err
	}

	return cesu8ToUTF8(buf), nil
}

func (d *nbtDecoder) readByteArray() ([]byte, error) {
	length, err := d.readInt32()
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, fmt.Errorf("NBT byte array length < 0: %d", length)
	}
	if length == 0 {
		return []byte{}, nil
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(d.r, buf)
	return buf, err
}

func (d *nbtDecoder) readList() ([]any, error) {
	elemType, err := d.readByte()
	if err != nil {
		return nil, err
	}

	length, err := d.readInt32()
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, fmt.Errorf("NBT list length < 0: %d", length)
	}
	if length == 0 {
		return []any{}, nil
	}

	result := make([]any, 0, length)
	for i := int32(0); i < length; i++ {
		v, err := d.decodeValue(elemType)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}

	return result, nil
}

func (d *nbtDecoder) readCompound() (map[string]any, error) {
	result := make(map[string]any)

	for {
		tagType, err := d.readByte()
		if err != nil {
			return nil, err
		}
		if tagType == tagEnd {
			break
		}

		name, err := d.readString()
		if err != nil {
			return nil, err
		}

		value, err := d.decodeValue(tagType)
		if err != nil {
			return nil, err
		}

		result[name] = value
	}

	return result, nil
}

func (d *nbtDecoder) readIntArray() ([]int32, error) {
	length, err := d.readInt32()
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, fmt.Errorf("NBT int array length < 0: %d", length)
	}
	if length == 0 {
		return []int32{}, nil
	}

	result := make([]int32, 0, length)
	for i := int32(0); i < length; i++ {
		v, err := d.readInt32()
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}

	return result, nil
}

func (d *nbtDecoder) readLongArray() ([]int64, error) {
	length, err := d.readInt32()
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, fmt.Errorf("NBT long array length < 0: %d", length)
	}
	if length == 0 {
		return []int64{}, nil
	}

	result := make([]int64, 0, length)
	for i := int32(0); i < length; i++ {
		v, err := d.readInt64()
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}

	return result, nil
}

func float32frombits(b uint32) float32 {
	return math.Float32frombits(b)
}

func float64frombits(b uint64) float64 {
	return math.Float64frombits(b)
}
