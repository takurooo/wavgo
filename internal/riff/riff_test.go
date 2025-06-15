package riff

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRiffChunkAddSubChunk(t *testing.T) {
	riffChunk := &riffChunk{
		ID:        RIFFChunkID,
		Size:      12,
		Format:    "WAVE",
		SubChunks: make([]*Chunk, 0),
	}

	// Add a chunk
	riffChunk.AddSubChunk("test", 4, []byte{0x01, 0x02, 0x03, 0x04})

	require.Len(t, riffChunk.SubChunks, 1)
	require.Equal(t, "test", riffChunk.SubChunks[0].ID)
	require.Equal(t, uint32(4), riffChunk.SubChunks[0].Size)
	require.Equal(t, []byte{0x01, 0x02, 0x03, 0x04}, riffChunk.SubChunks[0].Data)
}

func TestRiffChunkGetFMTChunk(t *testing.T) {
	riffChunk := &riffChunk{
		ID:        RIFFChunkID,
		Size:      20,
		Format:    "WAVE",
		SubChunks: make([]*Chunk, 0),
	}

	// Add some chunks
	riffChunk.AddSubChunk("junk", 4, []byte{0x01, 0x02, 0x03, 0x04})
	riffChunk.AddSubChunk(FMTChunkID, 8, []byte{0x01, 0x00, 0x02, 0x00, 0x44, 0xAC, 0x00, 0x00})
	riffChunk.AddSubChunk("info", 4, []byte{0x05, 0x06, 0x07, 0x08})

	// Get FMT chunk
	fmtChunk, err := riffChunk.GetFMTChunk()
	require.NoError(t, err)
	require.NotNil(t, fmtChunk)
	require.Equal(t, FMTChunkID, fmtChunk.ID)
	require.Equal(t, uint32(8), fmtChunk.Size)
	require.Equal(t, []byte{0x01, 0x00, 0x02, 0x00, 0x44, 0xAC, 0x00, 0x00}, fmtChunk.Data)
}

func TestRiffChunkGetFMTChunkNotFound(t *testing.T) {
	riffChunk := &riffChunk{
		ID:        RIFFChunkID,
		Size:      12,
		Format:    "WAVE",
		SubChunks: make([]*Chunk, 0),
	}

	// Add chunks without FMT chunk
	riffChunk.AddSubChunk("junk", 4, []byte{0x01, 0x02, 0x03, 0x04})
	riffChunk.AddSubChunk("info", 4, []byte{0x05, 0x06, 0x07, 0x08})

	// Try to get FMT chunk
	fmtChunk, err := riffChunk.GetFMTChunk()
	require.Error(t, err)
	require.Nil(t, fmtChunk)
	require.EqualError(t, err, "not found FMTChunk")
}

func TestRiffChunkGetDataChunk(t *testing.T) {
	riffChunk := &riffChunk{
		ID:        RIFFChunkID,
		Size:      20,
		Format:    "WAVE",
		SubChunks: make([]*Chunk, 0),
	}

	// Add some chunks
	riffChunk.AddSubChunk(FMTChunkID, 8, []byte{0x01, 0x00, 0x02, 0x00, 0x44, 0xAC, 0x00, 0x00})
	riffChunk.AddSubChunk(DATAChunkID, 8, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08})
	riffChunk.AddSubChunk("info", 4, []byte{0x09, 0x0A, 0x0B, 0x0C})

	// Get DATA chunk
	dataChunk, err := riffChunk.GetDataChunk()
	require.NoError(t, err)
	require.NotNil(t, dataChunk)
	require.Equal(t, DATAChunkID, dataChunk.ID)
	require.Equal(t, uint32(8), dataChunk.Size)
	require.Equal(t, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}, dataChunk.Data)
}

func TestRiffChunkGetDataChunkNotFound(t *testing.T) {
	riffChunk := &riffChunk{
		ID:        RIFFChunkID,
		Size:      16,
		Format:    "WAVE",
		SubChunks: make([]*Chunk, 0),
	}

	// Add chunks without DATA chunk
	riffChunk.AddSubChunk(FMTChunkID, 8, []byte{0x01, 0x00, 0x02, 0x00, 0x44, 0xAC, 0x00, 0x00})
	riffChunk.AddSubChunk("info", 4, []byte{0x05, 0x06, 0x07, 0x08})

	// Try to get DATA chunk
	dataChunk, err := riffChunk.GetDataChunk()
	require.Error(t, err)
	require.Nil(t, dataChunk)
	require.EqualError(t, err, "not found DataChunk")
}

func TestRiffChunkMultipleSameTypeChunks(t *testing.T) {
	riffChunk := &riffChunk{
		ID:        RIFFChunkID,
		Size:      24,
		Format:    "WAVE",
		SubChunks: make([]*Chunk, 0),
	}

	// Add multiple chunks with same ID (should return first one)
	riffChunk.AddSubChunk(FMTChunkID, 4, []byte{0x01, 0x02, 0x03, 0x04})
	riffChunk.AddSubChunk(FMTChunkID, 4, []byte{0x05, 0x06, 0x07, 0x08})
	riffChunk.AddSubChunk(DATAChunkID, 4, []byte{0x09, 0x0A, 0x0B, 0x0C})
	riffChunk.AddSubChunk(DATAChunkID, 4, []byte{0x0D, 0x0E, 0x0F, 0x10})

	// Should get first FMT chunk
	fmtChunk, err := riffChunk.GetFMTChunk()
	require.NoError(t, err)
	require.Equal(t, []byte{0x01, 0x02, 0x03, 0x04}, fmtChunk.Data)

	// Should get first DATA chunk
	dataChunk, err := riffChunk.GetDataChunk()
	require.NoError(t, err)
	require.Equal(t, []byte{0x09, 0x0A, 0x0B, 0x0C}, dataChunk.Data)
}
