package wavgo

import (
	"errors"
	"os"

	"github.com/takurooo/binaryio"
	"github.com/takurooo/wavgo/internal/riff"
)

// Writer ...
type Writer struct {
	f                   *os.File
	bw                  *binaryio.Writer
	format              *Format
	headerWritten       bool
	numWrittenSamples   uint32
	headerSize          uint32
	riffChunkSizeOffset int64
	dataChunkSizeOffset int64
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
	w.bw = binaryio.NewWriter(f)
	return nil
}

// Close ...
func (w *Writer) Close() error {
	dataChunkSize := w.numWrittenSamples * uint32(w.format.BlockAlign)
	riffChunkSize := dataChunkSize + w.headerSize - 8
	w.bw.SetOffset(w.riffChunkSizeOffset)
	w.bw.WriteU32(riffChunkSize, binaryio.LittleEndian)
	w.bw.SetOffset(w.dataChunkSizeOffset)
	w.bw.WriteU32(dataChunkSize, binaryio.LittleEndian)
	if w.bw.Err() != nil {
		return w.bw.Err()
	}
	if err := w.f.Sync(); err != nil {
		return err
	}
	if err := w.f.Close(); err != nil {
		return err
	}
	return nil
}

func (w *Writer) writeHeader() error {
	// riff chunk
	w.bw.WriteS32(riff.RIFFChunkID, binaryio.BigEndian)
	w.riffChunkSizeOffset = w.bw.GetOffset()
	w.bw.WriteU32(0, binaryio.LittleEndian) // dummy write
	w.bw.WriteS32("WAVE", binaryio.BigEndian)
	// fmt chunk
	w.bw.WriteS32(riff.FMTChunkID, binaryio.BigEndian)
	w.bw.WriteU32(0x10, binaryio.LittleEndian)
	w.bw.WriteU16(w.format.AudioFormat, binaryio.LittleEndian)
	w.bw.WriteU16(w.format.NumChannels, binaryio.LittleEndian)
	w.bw.WriteU32(w.format.SampleRate, binaryio.LittleEndian)
	w.bw.WriteU32(w.format.ByteRate, binaryio.LittleEndian)
	w.bw.WriteU16(w.format.BlockAlign, binaryio.LittleEndian)
	w.bw.WriteU16(w.format.BitsPerSample, binaryio.LittleEndian)
	// data chunk
	w.bw.WriteS32(riff.DATAChunkID, binaryio.BigEndian)
	w.dataChunkSizeOffset = w.bw.GetOffset()
	w.bw.WriteU32(0, binaryio.LittleEndian) // dummy write
	if w.bw.Err() != nil {
		return w.bw.Err()
	}

	w.headerSize = uint32(w.bw.GetOffset())
	return nil
}

// WriteSamples ...
func (w *Writer) WriteSamples(samples []Sample) error {
	if !w.headerWritten {
		err := w.writeHeader()
		if err != nil {
			return err
		}
		w.headerWritten = true
	}

	var (
		numChannels   = int(w.format.NumChannels)
		bitsPerSample = int(w.format.BitsPerSample)
	)
	for _, sample := range samples {
		for ch := 0; ch < numChannels; ch++ {
			switch bitsPerSample {
			case 8:
				w.bw.WriteU8(uint8(sample[ch]))
			case 16:
				w.bw.WriteU16(uint16(sample[ch]), binaryio.LittleEndian)
			case 24:
				w.bw.WriteU24(uint32(sample[ch]), binaryio.LittleEndian)
			case 32:
				w.bw.WriteU32(uint32(sample[ch]), binaryio.LittleEndian)
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
