package wavgo

import (
	"bytes"
	"errors"
	"os"

	bio "github.com/takurooo/binaryio"
	"github.com/takurooo/wavgo/riff"
)

// Reader ...
type Reader struct {
	f          *os.File
	format     *Format
	numSamples uint32
	br         *bio.Reader
}

// NewReader ...
func NewReader() *Reader {
	return &Reader{f: nil, format: nil, numSamples: 0, br: nil}
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

	// fmt.Println("AudioFormat    : ", r.format.AudioFormat)
	// fmt.Println("NumChannels    : ", r.format.NumChannels)
	// fmt.Println("SampleRate     : ", r.format.SampleRate)
	// fmt.Println("ByteRate       : ", r.format.ByteRate)
	// fmt.Println("BlockAlign     : ", r.format.BlockAlign)
	// fmt.Println("BitsPerSample  : ", r.format.BitsPerSample)

	// ----------------------------
	// Data Chunk
	// ----------------------------
	dataChunk, err := riffChunk.GetDataChunk()
	if err != nil {
		return err
	}
	r.numSamples = dataChunk.Size / uint32(r.format.BlockAlign)

	r.br = bio.NewReader(bytes.NewReader(dataChunk.Data))

	return nil
}

// GetFormat ...
func (r *Reader) GetFormat(format *Format) {
	*format = *(r.format)
}

// GetNumSamples ...
func (r *Reader) GetNumSamples() uint32 {
	return r.numSamples
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
				v = int(r.br.ReadU16(bio.LittleEndian))
			case 24:
				v = int(r.br.ReadU24(bio.LittleEndian))
			case 32:
				v = int(r.br.ReadU32(bio.LittleEndian))
			default:
				return nil, errors.New("not supported BitsPerSample")
			}

			if r.br.Err() != nil {
				return samples, r.br.Err()
			}

			samples[i][ch] = v
		}
	}
	return samples, nil
}

func parseFormatChunkData(fmtChunk *riff.Chunk) (*Format, error) {
	br := bio.NewReader(bytes.NewReader(fmtChunk.Data))
	format := &Format{
		AudioFormat:   br.ReadU16(bio.LittleEndian),
		NumChannels:   br.ReadU16(bio.LittleEndian),
		SampleRate:    br.ReadU32(bio.LittleEndian),
		ByteRate:      br.ReadU32(bio.LittleEndian),
		BlockAlign:    br.ReadU16(bio.LittleEndian),
		BitsPerSample: br.ReadU16(bio.LittleEndian),
	}
	if br.Err() != nil {
		return nil, br.Err()
	}
	return format, nil
}
