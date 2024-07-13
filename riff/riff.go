package riff

const (
	RIFF string = "RIFF"
	FMT  string = "fmt "
	DATA string = "data"
)

// Chunk ...
type Chunk struct {
	ID   string
	Size uint32
	Data []byte
}

// RIFFChunk ...
type riffChunk struct {
	ID        string
	Size      uint32
	Format    string
	SubChunks []*Chunk
}

func (r *riffChunk) AddSubChunk(id string, size uint32, data []byte) {
	r.SubChunks = append(r.SubChunks, &Chunk{id, size, data})
}

func (r *riffChunk) GetFMTChunk() (chunk *Chunk) {
	chunk = nil
	for _, c := range r.SubChunks {
		if c.ID == FMT {
			chunk = c
			break
		}
	}
	return chunk
}

func (r *riffChunk) GetDataChunk() (chunk *Chunk) {
	chunk = nil
	for _, c := range r.SubChunks {
		if c.ID == DATA {
			chunk = c
			break
		}
	}
	return chunk
}
