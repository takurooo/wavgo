package riff

// Chunk ...
type Chunk struct {
	ID   string
	Size uint32
	Data []byte
}

// RIFFChunk ...
type RIFFChunk struct {
	ID        string
	Size      uint32
	Format    string
	SubChunks []*Chunk
}

const (
	RIFF string = "RIFF"
	FMT  string = "fmt "
	DATA string = "data"
)
