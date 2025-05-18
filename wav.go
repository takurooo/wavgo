// Package wavgo provides simple utilities for reading and writing WAV
// (RIFF) audio files.
package wavgo

// AudioFormatPCM represents the PCM audio format defined by the WAV
// specification.
const (
	AudioFormatPCM = 0x0001
)

// Format describes the basic information stored in a WAV file's `fmt ` chunk.
type Format struct {
	AudioFormat   uint16 // audio codec (1 is PCM)
	NumChannels   uint16 // number of audio channels
	SampleRate    uint32 // sampling frequency in Hz
	ByteRate      uint32 // bytes per second
	BlockAlign    uint16 // bytes per sample frame
	BitsPerSample uint16 // bits used per sample
}

// Sample holds a single sample for up to two channels.
type Sample [2]int
