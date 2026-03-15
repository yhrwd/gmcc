package packet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"

	"gmcc/internal/nbt"
)

type SlotData struct {
	ID    int32
	Count int32
}

func ReadVarIntFromReader(r *bytes.Reader) (int32, error) {
	return ReadVarInt(r)
}

func ReadStringFromReader(r *bytes.Reader) (string, error) {
	return ReadString(r, r)
}

func ReadBoolFromReader(r io.Reader) (bool, error) {
	return ReadBool(r)
}

func ReadInt32FromReader(r io.Reader) (int32, error) {
	return ReadInt32(r)
}

func ReadFloat64FromReader(r io.Reader) (float64, error) {
	var v float64
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

func ReadFloat32FromReader(r io.Reader) (float32, error) {
	var v float32
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

func ReadU8(r io.Reader) (byte, error) {
	var b [1]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return 0, err
	}
	return b[0], nil
}

func ReadBytes(r io.Reader, n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ReadSlotData(r *bytes.Reader) (*SlotData, error) {
	count, err := ReadVarIntFromReader(r)
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		return nil, nil
	}

	itemID, err := ReadVarIntFromReader(r)
	if err != nil {
		return nil, nil
	}

	if err := SkipSlotComponents(r); err != nil {
		return nil, err
	}

	return &SlotData{ID: itemID, Count: count}, nil
}

func SkipSlotComponents(r *bytes.Reader) error {
	numAdd, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < numAdd; i++ {
		if err := SkipComponentData(r); err != nil {
			return err
		}
	}

	numRemove, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	for i := int32(0); i < numRemove; i++ {
		if _, err := ReadVarIntFromReader(r); err != nil {
			return err
		}
	}
	return nil
}

func SkipComponentData(r *bytes.Reader) error {
	componentType, err := ReadVarIntFromReader(r)
	if err != nil {
		return err
	}
	return SkipComponentByType(r, componentType)
}

func SkipNBT(r *bytes.Reader) error {
	dec := nbt.NewDecoder(r).NetworkFormat(true)
	return dec.Skip()
}

func ReadAnonymousNBTJSON(r io.Reader) (string, error) {
	dec := nbt.NewDecoder(r)
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
