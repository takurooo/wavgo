package wav

import (
	"os"
	"testing"
)

func TestWriter(t *testing.T) {
	var err error

	format := &Format{}
	format.AudioFormat = AudioFormatPCM
	format.NumChannels = 2
	format.SampleRate = 48000
	format.ByteRate = 128000
	format.BlockAlign = 4
	format.BitsPerSample = 16

	w := NewWriter(format)
	err = w.Open("test.wav")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		w.Close()
		if err := os.Remove("test.wav"); err != nil {
			t.Fatal(err)
		}
	}()

	samples := make([]Sample, 12)
	for i := 0; i < 12; i++ {
		for ch := 0; i < int(format.NumChannels); i++ {
			samples[i][ch] = i + ch
		}
	}

	w.WriteSamples(samples)
}
