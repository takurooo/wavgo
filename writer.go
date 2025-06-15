package wavgo

import (
	"encoding/binary"
	"os"

	"github.com/takurooo/wavgo/internal/binio"
	"github.com/takurooo/wavgo/internal/riff"
)

// Writer provides functionality to create and write WAV audio files.
// It uses the provided Format configuration to structure the output file
// and supports writing audio samples with automatic header generation.
type Writer struct {
	f                   *os.File
	bw                  *binio.Writer
	format              *Format
	headerWritten       bool
	numWrittenSamples   uint32
	headerSize          uint32
	riffChunkSizeOffset int64
	dataChunkSizeOffset int64
}

// NewWriter creates a new WAV file writer configured with the specified Format.
// The format parameter defines the audio characteristics such as sample rate,
// bit depth, and channel configuration. The writer must be opened with Open()
// before writing samples.
func NewWriter(format *Format) *Writer {
	return &Writer{f: nil, bw: nil, format: format, headerWritten: false, numWrittenSamples: 0}
}

// Open creates the destination WAV file at the specified path. This method
// prepares the file for writing but does not write the WAV header yet.
// The header is written automatically on the first call to WriteSamples().
func (w *Writer) Open(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	w.f = f
	w.bw = binio.NewWriter(f)
	return nil
}

// Close finalizes the WAV file by updating the RIFF and data chunk sizes
// in the header, syncing the file to disk, and closing the file handle.
// This method must be called to ensure the WAV file is properly formatted
// and all data is written to disk.
func (w *Writer) Close() error {
	dataChunkSize := w.numWrittenSamples * uint32(w.format.BlockAlign)
	riffChunkSize := dataChunkSize + w.headerSize - 8
	w.bw.SetOffset(w.riffChunkSizeOffset)
	w.bw.WriteU32(riffChunkSize, binary.LittleEndian)
	if w.bw.Err() != nil {
		return w.bw.Err()
	}
	w.bw.SetOffset(w.dataChunkSizeOffset)
	w.bw.WriteU32(dataChunkSize, binary.LittleEndian)
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
	w.bw.WriteS32(riff.RIFFChunkID, binary.BigEndian)
	w.riffChunkSizeOffset = w.bw.GetOffset()
	w.bw.WriteU32(0, binary.LittleEndian) // dummy write
	w.bw.WriteS32("WAVE", binary.BigEndian)
	// fmt chunk
	w.bw.WriteS32(riff.FMTChunkID, binary.BigEndian)
	w.bw.WriteU32(0x10, binary.LittleEndian)
	w.bw.WriteU16(w.format.AudioFormat, binary.LittleEndian)
	w.bw.WriteU16(w.format.NumChannels, binary.LittleEndian)
	w.bw.WriteU32(w.format.SampleRate, binary.LittleEndian)
	w.bw.WriteU32(w.format.ByteRate, binary.LittleEndian)
	w.bw.WriteU16(w.format.BlockAlign, binary.LittleEndian)
	w.bw.WriteU16(w.format.BitsPerSample, binary.LittleEndian)
	// data chunk
	w.bw.WriteS32(riff.DATAChunkID, binary.BigEndian)
	w.dataChunkSizeOffset = w.bw.GetOffset()
	w.bw.WriteU32(0, binary.LittleEndian) // dummy write
	if w.bw.Err() != nil {
		return w.bw.Err()
	}

	w.headerSize = uint32(w.bw.GetOffset())
	return nil
}

// WriteSamples writes the provided audio samples to the WAV file. On the first
// call, this method automatically writes the WAV header before writing sample data.
// Each Sample in the slice should contain data for all channels defined in the Format.
// The method handles the conversion of sample data to the appropriate bit depth
// and byte order as specified in the format configuration.
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
				w.bw.WriteU16(uint16(sample[ch]), binary.LittleEndian)
			case 24:
				w.bw.WriteU24(uint32(sample[ch]), binary.LittleEndian)
			case 32:
				w.bw.WriteU32(uint32(sample[ch]), binary.LittleEndian)
			default:
				return ErrUnsupportedBitsPerSample
			}

			if w.bw.Err() != nil {
				return w.bw.Err()
			}
		}
		w.numWrittenSamples++
	}
	return nil
}
