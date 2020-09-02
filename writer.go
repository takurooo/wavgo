package wav

import (
	"errors"
	"os"

	bio "github.com/takurooo/binaryio"
	"github.com/takurooo/wavgo/riff"
)

// Writer ...
type Writer struct {
	f                 *os.File
	bw                *bio.Writer
	format            *Format
	headerWritten     bool
	numWrittenSamples uint32
}

// NewWriter ...
func NewWriter(format *Format) *Writer {
	return &Writer{f: nil, bw: nil, format: format, headerWritten: false, numWrittenSamples: 0}
}

// Open ...
func (w *Writer) Open(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	w.f = f
	w.bw = bio.NewWriter(f)
	return nil
}

// Close ...
func (w *Writer) Close() {
	wavHeaderSize := uint32(44)
	dataChunkSize := w.numWrittenSamples * uint32(w.format.BlockAlign)
	riffChunkSize := dataChunkSize + wavHeaderSize - 8
	w.bw.SetOffset(4)
	w.bw.WriteU32(riffChunkSize, bio.LittleEndian)
	w.bw.SetOffset(40)
	w.bw.WriteU32(dataChunkSize, bio.LittleEndian)
	w.f.Close()
}

func (w *Writer) writeHeader() {
	// riff chunk
	w.bw.WriteS32("RIFF", bio.BigEndian)
	w.bw.WriteU32(0, bio.LittleEndian) // dummy write offset 4
	w.bw.WriteS32("WAVE", bio.BigEndian)
	// fmt chunk
	w.bw.WriteS32(riff.FMT, bio.BigEndian)
	w.bw.WriteU32(0x10, bio.LittleEndian)
	w.bw.WriteU16(w.format.AudioFormat, bio.LittleEndian)
	w.bw.WriteU16(w.format.NumChannels, bio.LittleEndian)
	w.bw.WriteU32(w.format.SampleRate, bio.LittleEndian)
	w.bw.WriteU32(w.format.ByteRate, bio.LittleEndian)
	w.bw.WriteU16(w.format.BlockAlign, bio.LittleEndian)
	w.bw.WriteU16(w.format.BitsPerSample, bio.LittleEndian)
	// data chunk
	w.bw.WriteS32(riff.DATA, bio.BigEndian)
	w.bw.WriteU32(0, bio.LittleEndian) // dummy write offset 40
}

// WriteSamples ...
func (w *Writer) WriteSamples(samples []Sample) error {
	if !w.headerWritten {
		w.writeHeader()
	}

	numChannels := int(w.format.NumChannels)
	bitsPerSample := int(w.format.BitsPerSample)

	for _, sample := range samples {
		for ch := 0; ch < numChannels; ch++ {
			switch bitsPerSample {
			case 8:
				w.bw.WriteU8(uint8(sample[ch]))
			case 16:
				w.bw.WriteU16(uint16(sample[ch]), bio.LittleEndian)
			case 24:
				w.bw.WriteU24(uint32(sample[ch]), bio.LittleEndian)
			case 32:
				w.bw.WriteU32(uint32(sample[ch]), bio.LittleEndian)
			default:
				return errors.New("not supported BitsPerSample")
			}

			if w.bw.Err() != nil {
				return w.bw.Err()
			}
		}
		w.numWrittenSamples++
	}
	return nil
}
