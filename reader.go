package wav

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
	riffReader := riff.NewReader(r.f)
	// ----------------------------
	// RIFF Chunk
	// ----------------------------
	riffChunk, err := riffReader.Read()
	if err != nil {
		return nil
	}

	// ----------------------------
	// Format Chunk
	// ----------------------------
	fmtChunk := riffReader.GetChunk(riffChunk, riff.FMT)
	if fmtChunk == nil {
		return errors.New("not found FmtChunk")
	}

	r.format, err = readFormat(fmtChunk)
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
	dataChunk := riffReader.GetChunk(riffChunk, riff.DATA)
	if dataChunk == nil {
		return errors.New("not found DataChunk")
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
func (r *Reader) GetSamples(numSamples int) (samples []Sample, err error) {

	samples = make([]Sample, numSamples)

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

func readFormat(fmtChunk *riff.Chunk) (format *Format, err error) {
	format = &Format{}
	br := bio.NewReader(bytes.NewReader(fmtChunk.Data))
	format.AudioFormat = br.ReadU16(bio.LittleEndian)
	format.NumChannels = br.ReadU16(bio.LittleEndian)
	format.SampleRate = br.ReadU32(bio.LittleEndian)
	format.ByteRate = br.ReadU32(bio.LittleEndian)
	format.BlockAlign = br.ReadU16(bio.LittleEndian)
	format.BitsPerSample = br.ReadU16(bio.LittleEndian)
	if br.Err() != nil {
		return nil, br.Err()
	}
	return format, nil
}
