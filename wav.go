package wav

// AudioFormatPCM ...
const (
	AudioFormatPCM = 0x0001
)

// Format ...
type Format struct {
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
}

// Sample ...
type Sample [2]int
