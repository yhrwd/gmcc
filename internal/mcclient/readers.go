package mcclient

import (
	"bytes"
	"encoding/binary"
	"io"

	"gmcc/internal/nbt"
)

type SlotData struct {
	ID    int32
	Count int32
}

func readVarIntFromReader(r *bytes.Reader) (int32, error) {
	return readVarInt(r)
}

func readStringFromReader(r *bytes.Reader) (string, error) {
	return readString(r, r)
}

func readBoolFromReader(r io.Reader) (bool, error) {
	return readBool(r)
}

func readInt32FromReader(r io.Reader) (int32, error) {
	return readInt32(r)
}

func readFloat64FromReader(r io.Reader) (float64, error) {
	var v float64
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

func readFloat32FromReader(r io.Reader) (float32, error) {
	var v float32
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

func readU8(r io.Reader) (byte, error) {
	var b [1]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return 0, err
	}
	return b[0], nil
}

func readBytes(r io.Reader, n int) []byte {
	b := make([]byte, n)
	_, _ = io.ReadFull(r, b)
	return b
}

func readSlotData(r *bytes.Reader) (*SlotData, error) {
	count, err := readVarIntFromReader(r)
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		return nil, nil
	}

	itemID, err := readVarIntFromReader(r)
	if err != nil {
		return nil, nil
	}

	numComponentsToAdd, err := readVarIntFromReader(r)
	if err != nil {
		return nil, nil
	}

	for i := int32(0); i < numComponentsToAdd; i++ {
		if err := skipComponentData(r); err != nil {
			return nil, nil
		}
	}

	numComponentsToRemove, err := readVarIntFromReader(r)
	if err != nil {
		return nil, nil
	}

	for i := int32(0); i < numComponentsToRemove; i++ {
		if _, err := readVarIntFromReader(r); err != nil {
			return nil, nil
		}
	}

	return &SlotData{ID: itemID, Count: count}, nil
}

func skipComponentData(r *bytes.Reader) error {
	componentType, err := readVarIntFromReader(r)
	if err != nil {
		return err
	}
	return skipComponentByType(r, componentType)
}

func skipNBT(r *bytes.Reader) error {
	dec := nbt.NewDecoder(r).NetworkFormat(true)
	return dec.Skip()
}
