package riff

import (
	"errors"
	"io"

	bio "github.com/takurooo/binaryio"
)

// Reader ...
type Reader struct {
	br *bio.Reader
}

// NewReader ...
func NewReader(r io.ReaderAt) *Reader {
	return &Reader{bio.NewReader(r)}
}

// Read ...
func (r *Reader) Read() (riffChunk *RIFFChunk, err error) {

	// ----------------------------
	// Read RIFF Chunk
	// ----------------------------
	chunkID := r.br.ReadS32(bio.BigEndian)
	if chunkID != "RIFF" {
		return nil, errors.New("not found riff chunk")
	}
	chunkSize := r.br.ReadU32(bio.LittleEndian)
	format := r.br.ReadS32(bio.BigEndian)

	riffChunk = &RIFFChunk{chunkID, chunkSize, format, make([]*Chunk, 0)}

	// ----------------------------
	// Read SubChunks
	// ----------------------------
	numBytesLeft := riffChunk.Size - 4
	for 0 < numBytesLeft {
		subChunkID := r.br.ReadS32(bio.BigEndian)
		subChunkSize := r.br.ReadU32(bio.LittleEndian)
		chunkData := r.br.ReadRaw(uint64(subChunkSize))

		riffChunk.SubChunks = append(riffChunk.SubChunks,
			&Chunk{subChunkID, subChunkSize, chunkData})

		if r.br.Err() != nil {
			return nil, r.br.Err()
		}

		numBytesLeft -= subChunkSize + 8
	}

	return riffChunk, nil
}

// GetChunk ...
func (r *Reader) GetChunk(riffChunk *RIFFChunk, chunkID string) (chunk *Chunk) {
	chunk = nil
	for _, c := range riffChunk.SubChunks {
		if c.ID == chunkID {
			chunk = c
			break
		}
	}
	return chunk
}
