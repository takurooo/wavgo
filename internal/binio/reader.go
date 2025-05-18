package binio

import (
	"encoding/binary"
	"io"
)

type Reader struct {
	r   io.ReaderAt
	off int64
	err error
}

func NewReader(r io.ReaderAt) *Reader {
	return &Reader{r: r}
}

func (br *Reader) read(n int) []byte {
	if br.err != nil {
		return nil
	}
	buf := make([]byte, n)
	_, err := br.r.ReadAt(buf, br.off)
	if err != nil {
		br.err = err
		return nil
	}
	br.off += int64(n)
	return buf
}

func (br *Reader) ReadRaw(n uint64) []byte {
	return br.read(int(n))
}

func (br *Reader) ReadU8() uint8 {
	b := br.read(1)
	if br.err != nil {
		return 0
	}
	return b[0]
}

func (br *Reader) ReadU16(order binary.ByteOrder) uint16 {
	b := br.read(2)
	if br.err != nil {
		return 0
	}
	return order.Uint16(b)
}

func (br *Reader) ReadU24(order binary.ByteOrder) uint32 {
	b := br.read(3)
	if br.err != nil {
		return 0
	}
	if order == binary.LittleEndian {
		return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
	}
	return uint32(b[2]) | uint32(b[1])<<8 | uint32(b[0])<<16
}

func (br *Reader) ReadU32(order binary.ByteOrder) uint32 {
	b := br.read(4)
	if br.err != nil {
		return 0
	}
	return order.Uint32(b)
}

func (br *Reader) ReadS32(order binary.ByteOrder) string {
	_ = order // order is ignored for strings
	b := br.read(4)
	if br.err != nil {
		return ""
	}
	return string(b)
}

func (br *Reader) Err() error {
	return br.err
}
