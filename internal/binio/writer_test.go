package binio

import (
	"encoding/binary"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriter(t *testing.T) {
	buf := &seekableBuffer{}
	writer := NewWriter(buf)

	// Test WriteU8
	writer.WriteU8(0x01)
	require.NoError(t, writer.Err())

	// Test WriteU16 Little Endian
	writer.WriteU16(0x0302, binary.LittleEndian)
	require.NoError(t, writer.Err())

	// Test WriteU24 Little Endian
	writer.WriteU24(0x060504, binary.LittleEndian)
	require.NoError(t, writer.Err())

	// Test WriteU32 Little Endian
	writer.WriteU32(0x0A090807, binary.LittleEndian)
	require.NoError(t, writer.Err())

	// Test WriteS32
	writer.WriteS32("TEST", binary.LittleEndian)
	require.NoError(t, writer.Err())

	// Verify written data
	expected := []byte{
		0x01,                   // uint8: 1
		0x02, 0x03,             // uint16 LE: 770
		0x04, 0x05, 0x06,       // uint24 LE: 395268
		0x07, 0x08, 0x09, 0x0A, // uint32 LE: 168496135
		'T', 'E', 'S', 'T',     // string: "TEST"
	}
	require.Equal(t, expected, buf.Bytes())
	require.Equal(t, int64(len(expected)), writer.GetOffset())
}

func TestWriterBigEndian(t *testing.T) {
	buf := &seekableBuffer{}
	writer := NewWriter(buf)

	// Test WriteU16 Big Endian
	writer.WriteU16(0x0102, binary.BigEndian)
	require.NoError(t, writer.Err())

	// Test WriteU24 Big Endian
	writer.WriteU24(0x030405, binary.BigEndian)
	require.NoError(t, writer.Err())

	// Test WriteU32 Big Endian
	writer.WriteU32(0x06070809, binary.BigEndian)
	require.NoError(t, writer.Err())

	expected := []byte{
		0x01, 0x02,             // uint16 BE: 258
		0x03, 0x04, 0x05,       // uint24 BE: 197637
		0x06, 0x07, 0x08, 0x09, // uint32 BE: 101124105
	}
	require.Equal(t, expected, buf.Bytes())
}

func TestWriterSeekAndOffset(t *testing.T) {
	buf := &seekableBuffer{}
	writer := NewWriter(buf)

	// Write some data
	writer.WriteU32(0x12345678, binary.LittleEndian)
	require.NoError(t, writer.Err())
	require.Equal(t, int64(4), writer.GetOffset())

	// Seek to beginning
	writer.SetOffset(0)
	require.NoError(t, writer.Err())
	require.Equal(t, int64(0), writer.GetOffset())

	// Overwrite first two bytes
	writer.WriteU16(0xABCD, binary.LittleEndian)
	require.NoError(t, writer.Err())
	require.Equal(t, int64(2), writer.GetOffset())

	// Verify we can write and seek
	require.Greater(t, len(buf.Bytes()), 0)
}

func TestWriterStringValidation(t *testing.T) {
	buf := &seekableBuffer{}
	writer := NewWriter(buf)

	// Test valid 4-character string
	writer.WriteS32("RIFF", binary.LittleEndian)
	require.NoError(t, writer.Err())
	require.Equal(t, []byte("RIFF"), buf.Bytes())

	// Test invalid string length (too short)
	buf.Reset()
	writer = NewWriter(buf)
	writer.WriteS32("ABC", binary.LittleEndian)
	require.Error(t, writer.Err())
	require.EqualError(t, writer.Err(), "string must be exactly 4 characters")
	require.Empty(t, buf.Bytes())

	// Test invalid string length (too long)
	buf.Reset()
	writer = NewWriter(buf)
	writer.WriteS32("TOOLONG", binary.LittleEndian)
	require.Error(t, writer.Err())
	require.EqualError(t, writer.Err(), "string must be exactly 4 characters")
	require.Empty(t, buf.Bytes())

	// Test empty string
	buf.Reset()
	writer = NewWriter(buf)
	writer.WriteS32("", binary.LittleEndian)
	require.Error(t, writer.Err())
	require.EqualError(t, writer.Err(), "string must be exactly 4 characters")
	require.Empty(t, buf.Bytes())
}

func TestWriterErrorHandling(t *testing.T) {
	// Test with failing writer
	failingWriter := &failingWriteSeeker{}
	writer := NewWriter(failingWriter)

	writer.WriteU8(0x01)
	require.Error(t, writer.Err())

	// Error should persist
	writer.WriteU16(0x0102, binary.LittleEndian)
	require.Error(t, writer.Err())
}

func TestWriterSeekError(t *testing.T) {
	// Test seek failure
	failingSeeker := &failingSeeker{}
	writer := NewWriter(failingSeeker)

	// First write should succeed
	writer.WriteU8(0x01)
	require.NoError(t, writer.Err())

	// Seek should fail
	writer.SetOffset(0)
	require.Error(t, writer.Err())

	// Subsequent operations should be no-ops
	originalOffset := writer.GetOffset()
	writer.WriteU8(0x02)
	require.Error(t, writer.Err())
	require.Equal(t, originalOffset, writer.GetOffset())
}

func TestWriterU24EdgeCases(t *testing.T) {
	buf := &seekableBuffer{}
	writer := NewWriter(buf)

	// Test maximum 24-bit value
	writer.WriteU24(0xFFFFFF, binary.LittleEndian)
	require.NoError(t, writer.Err())
	require.Equal(t, []byte{0xFF, 0xFF, 0xFF}, buf.Bytes())

	// Test with big endian
	buf.Reset()
	writer = NewWriter(buf)
	writer.WriteU24(0xFFFFFF, binary.BigEndian)
	require.NoError(t, writer.Err())
	require.Equal(t, []byte{0xFF, 0xFF, 0xFF}, buf.Bytes())

	// Test with mixed endianness
	buf.Reset()
	writer = NewWriter(buf)
	writer.WriteU24(0x123456, binary.LittleEndian)
	require.NoError(t, writer.Err())
	require.Equal(t, []byte{0x56, 0x34, 0x12}, buf.Bytes())

	buf.Reset()
	writer = NewWriter(buf)
	writer.WriteU24(0x123456, binary.BigEndian)
	require.NoError(t, writer.Err())
	require.Equal(t, []byte{0x12, 0x34, 0x56}, buf.Bytes())
}

// Mock types for testing error conditions

// seekableBuffer implements io.WriteSeeker using bytes.Buffer
type seekableBuffer struct {
	data   []byte
	offset int64
}

func (s *seekableBuffer) Write(p []byte) (n int, err error) {
	// Expand buffer if needed
	needed := int(s.offset) + len(p)
	if needed > len(s.data) {
		newData := make([]byte, needed)
		copy(newData, s.data)
		s.data = newData
	}
	
	// Write data at current offset
	copy(s.data[s.offset:], p)
	s.offset += int64(len(p))
	return len(p), nil
}

func (s *seekableBuffer) Seek(offset int64, whence int) (int64, error) {
	var newOffset int64
	switch whence {
	case io.SeekStart:
		newOffset = offset
	case io.SeekCurrent:
		newOffset = s.offset + offset
	case io.SeekEnd:
		newOffset = int64(len(s.data)) + offset
	default:
		return 0, errors.New("invalid whence")
	}
	
	if newOffset < 0 {
		return 0, errors.New("negative seek position")
	}
	
	s.offset = newOffset
	return newOffset, nil
}

func (s *seekableBuffer) Bytes() []byte {
	return s.data[:s.offset]
}

func (s *seekableBuffer) Reset() {
	s.data = nil
	s.offset = 0
}

func (s *seekableBuffer) Len() int {
	return int(s.offset)
}

type failingWriteSeeker struct{}

func (f *failingWriteSeeker) Write(p []byte) (n int, err error) {
	return 0, io.ErrClosedPipe
}

func (f *failingWriteSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, io.ErrClosedPipe
}

type failingSeeker struct {
	buf *seekableBuffer
}

func (f *failingSeeker) Write(p []byte) (n int, err error) {
	if f.buf == nil {
		f.buf = &seekableBuffer{}
	}
	return f.buf.Write(p)
}

func (f *failingSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, errors.New("seek failed")
}