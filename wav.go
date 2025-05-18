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
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
}

// Sample holds a single sample for up to two channels.
type Sample [2]int
