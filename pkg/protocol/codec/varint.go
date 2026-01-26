package codec

import (
	"fmt"
	"io"
)

func ReadVarInt(r io.Reader) (int32, error) {
	var num int32
	var shift uint
	buf := make([]byte, 1)
	for i := 0; i < 5; i++ {
		if _, err := r.Read(buf); err != nil {
			return 0, err
		}
		b := buf[0]
		num |= int32(b&0x7F) << shift
		if b&0x80 == 0 {
			return num, nil
		}
		shift += 7
	}
	return 0, fmt.Errorf("varint too big")
}

func WriteVarInt(w io.Writer, num int32) error {
	buf := make([]byte, 0, 5)
	for {
		if (num & ^0x7F) == 0 {
			buf = append(buf, byte(num))
			break
		}
		buf = append(buf, byte(num&0x7F|0x80))
		num >>= 7
	}
	_, err := w.Write(buf)
	return err
}
