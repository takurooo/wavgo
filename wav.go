package wavgo

// AudioFormatPCM ...
const (
	AudioFormatPCM = 0x0001
)

// Format holds the basic information of the WAV "fmt " chunk.
type Format struct {
	AudioFormat   uint16 // audio codec (1 is PCM)
	NumChannels   uint16 // number of audio channels
	SampleRate    uint32 // sampling frequency in Hz
	ByteRate      uint32 // bytes per second
	BlockAlign    uint16 // bytes per sample frame
	BitsPerSample uint16 // bits used per sample
}

// Sample ...
type Sample [2]int
