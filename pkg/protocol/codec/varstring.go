package codec

import (
	"io"
)

func WriteString(w io.Writer, s string) error {
	if err := WriteVarInt(w, int32(len(s))); err != nil {
		return err
	}
	_, err := w.Write([]byte(s))
	return err
}
