package binio

import (
	"encoding/binary"
	"io"
)

type Writer struct {
	w      io.WriteSeeker
	offset int64
	err    error
}

func NewWriter(w io.WriteSeeker) *Writer {
	return &Writer{w: w}
}

func (bw *Writer) write(b []byte) {
	if bw.err != nil {
		return
	}
	n, err := bw.w.Write(b)
	if err != nil {
		bw.err = err
		return
	}
	if n != len(b) {
		bw.err = io.ErrShortWrite
		return
	}
	bw.offset += int64(n)
}

func (bw *Writer) WriteU8(v uint8) {
	bw.write([]byte{v})
}

func (bw *Writer) WriteU16(v uint16, order binary.ByteOrder) {
	buf := make([]byte, 2)
	order.PutUint16(buf, v)
	bw.write(buf)
}

func (bw *Writer) WriteU24(v uint32, order binary.ByteOrder) {
	buf := make([]byte, 3)
	if order == binary.LittleEndian {
		buf[0] = byte(v)
		buf[1] = byte(v >> 8)
		buf[2] = byte(v >> 16)
	} else {
		buf[2] = byte(v)
		buf[1] = byte(v >> 8)
		buf[0] = byte(v >> 16)
	}
	bw.write(buf)
}

func (bw *Writer) WriteU32(v uint32, order binary.ByteOrder) {
	buf := make([]byte, 4)
	order.PutUint32(buf, v)
	bw.write(buf)
}

func (bw *Writer) WriteS32(s string, order binary.ByteOrder) {
	_ = order // order is ignored for strings
	if len(s) != 4 {
		bw.err = io.ErrShortWrite
		return
	}
	bw.write([]byte(s))
}

func (bw *Writer) SetOffset(off int64) {
	if bw.err != nil {
		return
	}
	_, err := bw.w.Seek(off, io.SeekStart)
	if err != nil {
		bw.err = err
		return
	}
	bw.offset = off
}

func (bw *Writer) GetOffset() int64 {
	return bw.offset
}

func (bw *Writer) Err() error {
	return bw.err
}
