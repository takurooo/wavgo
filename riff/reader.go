package riff

import (
	"errors"
	"io"

	bio "github.com/takurooo/binaryio"
)

func ReadRIFFChunk(r io.ReaderAt) (*riffChunk, error) {
	breader := bio.NewReader(r)
	// ----------------------------
	// Read RIFF Chunk
	// ----------------------------
	chunkID := breader.ReadS32(bio.BigEndian)
	if chunkID != RIFF {
		return nil, errors.New("not found riff chunk")
	}

	chunkSize := breader.ReadU32(bio.LittleEndian)
	format := breader.ReadS32(bio.BigEndian)
	riffChunk := &riffChunk{chunkID, chunkSize, format, make([]*Chunk, 0)}

	// ----------------------------
	// Read SubChunks
	// ----------------------------
	numBytesLeft := riffChunk.Size - 4
	for 0 < numBytesLeft {
		subChunkID := breader.ReadS32(bio.BigEndian)
		subChunkSize := breader.ReadU32(bio.LittleEndian)
		chunkData := breader.ReadRaw(uint64(subChunkSize))
		if breader.Err() != nil {
			return nil, breader.Err()
		}
		riffChunk.AddSubChunk(subChunkID, subChunkSize, chunkData)
		numBytesLeft -= subChunkSize + 8
	}
	return riffChunk, nil
}
