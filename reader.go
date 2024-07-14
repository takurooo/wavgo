package wavgo

import (
	"bytes"
	"errors"
	"os"

	"github.com/takurooo/binaryio"
	"github.com/takurooo/wavgo/internal/riff"
)

// Reader ...
type Reader struct {
	f              *os.File
	format         Format
	numSamples     uint32
	numSamplesLeft uint32
	br             *binaryio.Reader
}

// NewReader ...
func NewReader() *Reader {
	return &Reader{}
}

// Open ...
func (r *Reader) Open(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	r.f = f
	return nil
}

// Close ...
func (r *Reader) Close() {
	r.f.Close()
}

// ReadOnMemory ...
func (r *Reader) ReadOnMemory() error {
	// ----------------------------
	// RIFF Chunk
	// ----------------------------
	riffChunk, err := riff.ReadRIFFChunk(r.f)
	if err != nil {
		return err
	}
	// ----------------------------
	// Format Chunk
	// ----------------------------
	fmtChunk, err := riffChunk.GetFMTChunk()
	if err != nil {
		return err
	}

	r.format, err = parseFormatChunkData(fmtChunk)
	if err != nil {
		return err
	}
	// ----------------------------
	// Data Chunk
	// ----------------------------
	dataChunk, err := riffChunk.GetDataChunk()
	if err != nil {
		return err
	}
	r.numSamples = dataChunk.Size / uint32(r.format.BlockAlign)
	r.numSamplesLeft = r.numSamples
	r.br = binaryio.NewReader(bytes.NewReader(dataChunk.Data))
	return nil
}

// GetFormat ...
func (r *Reader) GetFormat() Format {
	return r.format
}

// GetNumSamples ...
func (r *Reader) GetNumSamples() uint32 {
	return r.numSamples
}

// GetNumSamples ...
func (r *Reader) GetNumSamplesLeft() uint32 {
	return r.numSamplesLeft
}

// GetSamples ...
func (r *Reader) GetSamples(numSamples int) ([]Sample, error) {
	samples := make([]Sample, numSamples)
	bitsPerSample := int(r.format.BitsPerSample)
	numChannels := int(r.format.NumChannels)

	for i := 0; i < numSamples; i++ {
		for ch := 0; ch < numChannels; ch++ {
			var v int
			switch bitsPerSample {
			case 8:
				v = int(r.br.ReadU8())
			case 16:
				v = int(r.br.ReadU16(binaryio.LittleEndian))
			case 24:
				v = int(r.br.ReadU24(binaryio.LittleEndian))
			case 32:
				v = int(r.br.ReadU32(binaryio.LittleEndian))
			default:
				return nil, errors.New("not supported BitsPerSample")
			}

			if r.br.Err() != nil {
				return nil, r.br.Err()
			}

			samples[i][ch] = v
		}
		r.numSamplesLeft -= 1
	}
	return samples, nil
}

func parseFormatChunkData(fmtChunk *riff.Chunk) (Format, error) {
	br := binaryio.NewReader(bytes.NewReader(fmtChunk.Data))
	format := Format{
		AudioFormat:   br.ReadU16(binaryio.LittleEndian),
		NumChannels:   br.ReadU16(binaryio.LittleEndian),
		SampleRate:    br.ReadU32(binaryio.LittleEndian),
		ByteRate:      br.ReadU32(binaryio.LittleEndian),
		BlockAlign:    br.ReadU16(binaryio.LittleEndian),
		BitsPerSample: br.ReadU16(binaryio.LittleEndian),
	}
	if br.Err() != nil {
		return Format{}, br.Err()
	}
	return format, nil
}
