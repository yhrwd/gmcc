package connection

import (
	"bufio"
	"net"
	"time"
)

type Conn struct {
	Raw net.Conn
	R   *bufio.Reader
	W   *bufio.Writer
}

func NewConn(c net.Conn) *Conn {
	const bufferSize = 1048576

	return &Conn{
		Raw: c,
		R:   bufio.NewReaderSize(c, bufferSize),
		W:   bufio.NewWriterSize(c, bufferSize),
	}
}

func (c *Conn) Close() error {
	return c.Raw.Close()
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.Raw.SetWriteDeadline(t)
}
