package binio

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	// Test data: various binary values
	data := []byte{
		0x01,       // uint8: 1
		0x02, 0x03, // uint16 LE: 770
		0x04, 0x05, 0x06, // uint24 LE: 395268
		0x07, 0x08, 0x09, 0x0A, // uint32 LE: 168496135
		'T', 'E', 'S', 'T', // string: "TEST"
	}

	reader := NewReader(bytes.NewReader(data))

	// Test ReadU8
	u8 := reader.ReadU8()
	require.NoError(t, reader.Err())
	require.Equal(t, uint8(0x01), u8)

	// Test ReadU16 Little Endian
	u16 := reader.ReadU16(binary.LittleEndian)
	require.NoError(t, reader.Err())
	require.Equal(t, uint16(0x0302), u16)

	// Test ReadU24 Little Endian
	u24 := reader.ReadU24(binary.LittleEndian)
	require.NoError(t, reader.Err())
	require.Equal(t, uint32(0x060504), u24)

	// Test ReadU32 Little Endian
	u32 := reader.ReadU32(binary.LittleEndian)
	require.NoError(t, reader.Err())
	require.Equal(t, uint32(0x0A090807), u32)

	// Test ReadS32
	s32 := reader.ReadS32(binary.LittleEndian)
	require.NoError(t, reader.Err())
	require.Equal(t, "TEST", s32)
}

func TestReaderBigEndian(t *testing.T) {
	data := []byte{
		0x01, 0x02, // uint16 BE: 258
		0x03, 0x04, 0x05, // uint24 BE: 197637
		0x06, 0x07, 0x08, 0x09, // uint32 BE: 101124105
	}

	reader := NewReader(bytes.NewReader(data))

	// Test ReadU16 Big Endian
	u16 := reader.ReadU16(binary.BigEndian)
	require.NoError(t, reader.Err())
	require.Equal(t, uint16(0x0102), u16)

	// Test ReadU24 Big Endian
	u24 := reader.ReadU24(binary.BigEndian)
	require.NoError(t, reader.Err())
	require.Equal(t, uint32(0x030405), u24)

	// Test ReadU32 Big Endian
	u32 := reader.ReadU32(binary.BigEndian)
	require.NoError(t, reader.Err())
	require.Equal(t, uint32(0x06070809), u32)
}

func TestReaderRaw(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	reader := NewReader(bytes.NewReader(data))

	raw := reader.ReadRaw(3)
	require.NoError(t, reader.Err())
	require.Equal(t, []byte{0x01, 0x02, 0x03}, raw)

	raw2 := reader.ReadRaw(2)
	require.NoError(t, reader.Err())
	require.Equal(t, []byte{0x04, 0x05}, raw2)
}

func TestReaderErrorHandling(t *testing.T) {
	// Test EOF error
	data := []byte{0x01}
	reader := NewReader(bytes.NewReader(data))

	// Read successful
	u8 := reader.ReadU8()
	require.NoError(t, reader.Err())
	require.Equal(t, uint8(0x01), u8)

	// Read beyond EOF
	u16 := reader.ReadU16(binary.LittleEndian)
	require.Error(t, reader.Err())
	require.Equal(t, uint16(0), u16)

	// Subsequent reads should return zero values
	u32 := reader.ReadU32(binary.LittleEndian)
	require.Error(t, reader.Err())
	require.Equal(t, uint32(0), u32)

	s32 := reader.ReadS32(binary.LittleEndian)
	require.Error(t, reader.Err())
	require.Equal(t, "", s32)
}

func TestReaderErrorPropagation(t *testing.T) {
	// Create a reader that always fails
	failingReader := &failingReaderAt{}
	reader := NewReader(failingReader)

	u8 := reader.ReadU8()
	require.Error(t, reader.Err())
	require.Equal(t, uint8(0), u8)

	// Error should persist
	u16 := reader.ReadU16(binary.LittleEndian)
	require.Error(t, reader.Err())
	require.Equal(t, uint16(0), u16)
}

// failingReaderAt is a mock that always returns an error
type failingReaderAt struct{}

func (f *failingReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}
