package wavgo

import (
	"os"
	"testing"
)

func TestWriter(t *testing.T) {
	var err error
	format := &Format{
		AudioFormat:   AudioFormatPCM,
		NumChannels:   2,
		SampleRate:    48000,
		ByteRate:      128000,
		BlockAlign:    4,
		BitsPerSample: 16,
	}

	w := NewWriter(format)
	err = w.Open("testdata/output.wav")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		w.Close()
		if err := os.Remove("testdata/output.wav"); err != nil {
			t.Fatal(err)
		}
	}()

	samples := make([]Sample, 12)
	for i := 0; i < 12; i++ {
		for ch := 0; i < int(format.NumChannels); i++ {
			samples[i][ch] = i + ch
		}
	}

	err = w.WriteSamples(samples)
	if err != nil {
		t.Fatal(err)
	}
}
