// Package wavgo provides a Go library for reading and writing WAV audio files
// with PCM format support. It offers a simple API for parsing WAV file headers,
// extracting audio samples, and creating new WAV files.
//
// The library supports common bit depths (8, 16, 24, 32 bits) and provides
// access to format information such as sample rate, number of channels, and
// bits per sample through the Format struct.
//
// Basic usage for reading:
//
//	reader := wavgo.NewReader()
//	err := reader.Open("input.wav")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer reader.Close()
//
//	err = reader.Load()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	format := reader.GetFormat()
//	samples, err := reader.GetSamples(1024)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Basic usage for writing:
//
//	format := &wavgo.Format{
//		AudioFormat:   wavgo.AudioFormatPCM,
//		NumChannels:   2,
//		SampleRate:    44100,
//		BitsPerSample: 16,
//		BlockAlign:    4,
//		ByteRate:      176400,
//	}
//
//	writer := wavgo.NewWriter(format)
//	err := writer.Open("output.wav")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer writer.Close()
//
//	samples := []wavgo.Sample{{100, 200}, {300, 400}}
//	err = writer.WriteSamples(samples)
//	if err != nil {
//		log.Fatal(err)
//	}
package wavgo

// Audio format constants as defined by the WAV specification.
const (
	// AudioFormatPCM represents the standard PCM (Pulse Code Modulation) audio format.
	// This is the most common uncompressed audio format used in WAV files.
	AudioFormatPCM = 0x0001
)

// Format describes the basic audio format information stored in a WAV file's fmt chunk.
// This structure contains all the essential parameters needed to interpret the audio data.
type Format struct {
	// AudioFormat specifies the audio compression codec. For PCM audio, this should be 1.
	AudioFormat uint16

	// NumChannels specifies the number of audio channels (1 for mono, 2 for stereo, etc.).
	NumChannels uint16

	// SampleRate specifies the sampling frequency in Hz (e.g., 44100, 48000).
	SampleRate uint32

	// ByteRate specifies the average number of bytes per second of audio data.
	// This is calculated as SampleRate * NumChannels * BitsPerSample / 8.
	ByteRate uint32

	// BlockAlign specifies the number of bytes for one complete sample frame
	// across all channels. This is calculated as NumChannels * BitsPerSample / 8.
	BlockAlign uint16

	// BitsPerSample specifies the number of bits used per audio sample
	// (typically 8, 16, 24, or 32).
	BitsPerSample uint16
}

// Sample represents a single audio sample frame that can hold data for up to two channels.
// For mono audio, only index 0 is used. For stereo audio, index 0 represents the left
// channel and index 1 represents the right channel. The values are stored as signed
// integers and their interpretation depends on the bit depth specified in the Format.
type Sample [2]int
