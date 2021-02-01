package binary

import (
	"encoding/binary"
	"time"
)

type Encoder struct {
	order binary.ByteOrder
}

func NewEncoder() *Encoder {
	return &Encoder{
		order: binary.LittleEndian,
	}
}

func (c *Encoder) Decode(b []byte) time.Time {
	seconds := c.order.Uint64(b)
	t := time.Unix(int64(seconds), 0)

	return t
}

func (c *Encoder) Encode(t time.Time) []byte {
	message := make([]byte, 8)
	seconds := uint64(t.Unix())
	c.order.PutUint64(message[:8], seconds)

	return message
}
