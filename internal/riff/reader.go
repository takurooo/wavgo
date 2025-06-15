package riff

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/takurooo/wavgo/internal/binio"
)

func ReadRIFFChunk(r io.ReaderAt) (*riffChunk, error) {
	breader := binio.NewReader(r)
	// ----------------------------
	// Read RIFF Chunk
	// ----------------------------
	var (
		chunkID   = breader.ReadS32(binary.BigEndian)
		chunkSize = breader.ReadU32(binary.LittleEndian)
		format    = breader.ReadS32(binary.BigEndian)
	)
	if breader.Err() != nil {
		return nil, breader.Err()
	}
	if chunkID != RIFFChunkID {
		return nil, errors.New("not found riff chunk")
	}
	riffChunk := &riffChunk{chunkID, chunkSize, format, make([]*Chunk, 0)}
	// ----------------------------
	// Read SubChunks
	// ----------------------------
	numBytesLeft := riffChunk.Size - 4
	for 0 < numBytesLeft {
		var (
			subChunkID   = breader.ReadS32(binary.BigEndian)
			subChunkSize = breader.ReadU32(binary.LittleEndian)
			chunkData    = breader.ReadRaw(uint64(subChunkSize))
		)
		if breader.Err() != nil {
			return nil, breader.Err()
		}

		chunkOverhead := uint32(8) // 4 bytes ID + 4 bytes size
		if subChunkSize+chunkOverhead > numBytesLeft {
			return nil, errors.New("invalid chunk size: exceeds remaining bytes")
		}

		riffChunk.AddSubChunk(subChunkID, subChunkSize, chunkData)
		numBytesLeft -= subChunkSize + chunkOverhead
	}
	return riffChunk, nil
}
