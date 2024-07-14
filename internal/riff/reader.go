package riff

import (
	"errors"
	"io"

	"github.com/takurooo/binaryio"
)

func ReadRIFFChunk(r io.ReaderAt) (*riffChunk, error) {
	breader := binaryio.NewReader(r)
	// ----------------------------
	// Read RIFF Chunk
	// ----------------------------
	var (
		chunkID   = breader.ReadS32(binaryio.BigEndian)
		chunkSize = breader.ReadU32(binaryio.LittleEndian)
		format    = breader.ReadS32(binaryio.BigEndian)
	)
	if breader.Err() != nil {
		return nil, breader.Err()
	}
	if chunkID != RIFF {
		return nil, errors.New("not found riff chunk")
	}
	riffChunk := &riffChunk{chunkID, chunkSize, format, make([]*Chunk, 0)}
	// ----------------------------
	// Read SubChunks
	// ----------------------------
	numBytesLeft := riffChunk.Size - 4
	for 0 < numBytesLeft {
		var (
			subChunkID   = breader.ReadS32(binaryio.BigEndian)
			subChunkSize = breader.ReadU32(binaryio.LittleEndian)
			chunkData    = breader.ReadRaw(uint64(subChunkSize))
		)
		if breader.Err() != nil {
			return nil, breader.Err()
		}
		riffChunk.AddSubChunk(subChunkID, subChunkSize, chunkData)
		numBytesLeft -= subChunkSize + 8
	}
	return riffChunk, nil
}
