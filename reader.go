package wavgo

import (
	"bytes"
	"encoding/binary"
	"os"

	"github.com/takurooo/wavgo/internal/binio"
	"github.com/takurooo/wavgo/internal/riff"
)

// Reader provides access to the samples and metadata contained in a WAV file.
type Reader struct {
	f              *os.File
	format         Format
	numSamples     uint32
	numSamplesLeft uint32
	br             *binio.Reader
}

// NewReader creates an empty Reader.
func NewReader() *Reader {
	return &Reader{}
}

// Open opens the specified WAV file for reading.
func (r *Reader) Open(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	r.f = f
	return nil
}

// Close closes the underlying file descriptor.
func (r *Reader) Close() error {
	if r.f == nil {
		return nil
	}
	return r.f.Close()
}

// Load reads and parses the WAV file into memory.
func (r *Reader) Load() error {
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
	r.br = binio.NewReader(bytes.NewReader(dataChunk.Data))
	return nil
}

// GetFormat returns the file's format information.
func (r *Reader) GetFormat() Format {
	return r.format
}

// GetNumSamples returns the total number of samples in the file.
func (r *Reader) GetNumSamples() uint32 {
	return r.numSamples
}

// GetNumSamplesLeft returns the number of samples remaining.
func (r *Reader) GetNumSamplesLeft() uint32 {
	return r.numSamplesLeft
}

// GetSamples reads the next numSamples samples from the data chunk.
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
				v = int(r.br.ReadU16(binary.LittleEndian))
			case 24:
				v = int(r.br.ReadU24(binary.LittleEndian))
			case 32:
				v = int(r.br.ReadU32(binary.LittleEndian))
			default:
				return nil, ErrUnsupportedBitsPerSample
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
	br := binio.NewReader(bytes.NewReader(fmtChunk.Data))
	format := Format{
		AudioFormat:   br.ReadU16(binary.LittleEndian),
		NumChannels:   br.ReadU16(binary.LittleEndian),
		SampleRate:    br.ReadU32(binary.LittleEndian),
		ByteRate:      br.ReadU32(binary.LittleEndian),
		BlockAlign:    br.ReadU16(binary.LittleEndian),
		BitsPerSample: br.ReadU16(binary.LittleEndian),
	}
	if br.Err() != nil {
		return Format{}, br.Err()
	}
	return format, nil
}
