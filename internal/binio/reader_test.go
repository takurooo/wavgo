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

	t.Run("ReadU8", func(t *testing.T) {
		reader := NewReader(bytes.NewReader(data))
		u8 := reader.ReadU8()
		require.NoError(t, reader.Err())
		require.Equal(t, uint8(0x01), u8)
	})

	t.Run("ReadU16_LittleEndian", func(t *testing.T) {
		reader := NewReader(bytes.NewReader(data))
		reader.ReadU8() // Skip first byte
		u16 := reader.ReadU16(binary.LittleEndian)
		require.NoError(t, reader.Err())
		require.Equal(t, uint16(0x0302), u16)
	})

	t.Run("ReadU24_LittleEndian", func(t *testing.T) {
		reader := NewReader(bytes.NewReader(data))
		reader.ReadU8()                     // Skip first byte
		reader.ReadU16(binary.LittleEndian) // Skip next 2 bytes
		u24 := reader.ReadU24(binary.LittleEndian)
		require.NoError(t, reader.Err())
		require.Equal(t, uint32(0x060504), u24)
	})

	t.Run("ReadU32_LittleEndian", func(t *testing.T) {
		reader := NewReader(bytes.NewReader(data))
		reader.ReadU8()                     // Skip first byte
		reader.ReadU16(binary.LittleEndian) // Skip next 2 bytes
		reader.ReadU24(binary.LittleEndian) // Skip next 3 bytes
		u32 := reader.ReadU32(binary.LittleEndian)
		require.NoError(t, reader.Err())
		require.Equal(t, uint32(0x0A090807), u32)
	})

	t.Run("ReadS32", func(t *testing.T) {
		reader := NewReader(bytes.NewReader(data))
		reader.ReadU8()                     // Skip first byte
		reader.ReadU16(binary.LittleEndian) // Skip next 2 bytes
		reader.ReadU24(binary.LittleEndian) // Skip next 3 bytes
		reader.ReadU32(binary.LittleEndian) // Skip next 4 bytes
		s32 := reader.ReadS32(binary.LittleEndian)
		require.NoError(t, reader.Err())
		require.Equal(t, "TEST", s32)
	})
}

func TestReaderBigEndian(t *testing.T) {
	data := []byte{
		0x01, 0x02, // uint16 BE: 258
		0x03, 0x04, 0x05, // uint24 BE: 197637
		0x06, 0x07, 0x08, 0x09, // uint32 BE: 101124105
	}

	t.Run("ReadU16_BigEndian", func(t *testing.T) {
		reader := NewReader(bytes.NewReader(data))
		u16 := reader.ReadU16(binary.BigEndian)
		require.NoError(t, reader.Err())
		require.Equal(t, uint16(0x0102), u16)
	})

	t.Run("ReadU24_BigEndian", func(t *testing.T) {
		reader := NewReader(bytes.NewReader(data))
		reader.ReadU16(binary.BigEndian) // Skip first 2 bytes
		u24 := reader.ReadU24(binary.BigEndian)
		require.NoError(t, reader.Err())
		require.Equal(t, uint32(0x030405), u24)
	})

	t.Run("ReadU32_BigEndian", func(t *testing.T) {
		reader := NewReader(bytes.NewReader(data))
		reader.ReadU16(binary.BigEndian) // Skip first 2 bytes
		reader.ReadU24(binary.BigEndian) // Skip next 3 bytes
		u32 := reader.ReadU32(binary.BigEndian)
		require.NoError(t, reader.Err())
		require.Equal(t, uint32(0x06070809), u32)
	})
}

func TestReaderRaw(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05}

	t.Run("ReadRaw_FirstChunk", func(t *testing.T) {
		reader := NewReader(bytes.NewReader(data))
		raw := reader.ReadRaw(3)
		require.NoError(t, reader.Err())
		require.Equal(t, []byte{0x01, 0x02, 0x03}, raw)
	})

	t.Run("ReadRaw_Sequential", func(t *testing.T) {
		reader := NewReader(bytes.NewReader(data))
		reader.ReadRaw(3) // Read first 3 bytes
		raw2 := reader.ReadRaw(2)
		require.NoError(t, reader.Err())
		require.Equal(t, []byte{0x04, 0x05}, raw2)
	})
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
