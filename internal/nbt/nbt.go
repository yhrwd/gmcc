// Package nbt implements Minecraft's Named Binary Tag format.
package nbt

import (
	"errors"
)

// Tag types
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

// TagNames maps tag types to strings
var TagNames = map[byte]string{
	TagEnd:       "TAG_End",
	TagByte:      "TAG_Byte",
	TagShort:     "TAG_Short",
	TagInt:       "TAG_Int",
	TagLong:      "TAG_Long",
	TagFloat:     "TAG_Float",
	TagDouble:    "TAG_Double",
	TagByteArray: "TAG_ByteArray",
	TagString:    "TAG_String",
	TagList:      "TAG_List",
	TagCompound:  "TAG_Compound",
	TagIntArray:  "TAG_IntArray",
	TagLongArray: "TAG_LongArray",
}

// ErrEND is returned for unexpected TAG_End
var ErrEND = errors.New("unexpected TAG_End")

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

// Marshaler interface for custom encoding
type Marshaler interface {
	MarshalNBT() ([]byte, error)
}

// Unmarshaler interface for custom decoding
type Unmarshaler interface {
	UnmarshalNBT(data []byte) error
}

// RawMessage stores raw NBT for delayed decoding
type RawMessage struct {
	Type byte
	Data []byte
}

func (m RawMessage) TagType() byte  { return m.Type }
func (m RawMessage) String() string { return string(m.Data) }
