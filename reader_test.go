package wav

import (
	"testing"
)

func TestReader(t *testing.T) {
	var err error

	r := NewReader()
	if err = r.Open("files/test.wav"); err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	err = r.ReadOnMemory()
	if err != nil {
		t.Fatal(err)
	}

	var format Format
	r.GetFormat(&format)

	if format.NumChannels > 2 {
		t.Fatalf("Invalid NumChannels: %d", format.NumChannels)
	}

	if format.BitsPerSample != 8 &&
		format.BitsPerSample != 16 &&
		format.BitsPerSample != 24 &&
		format.BitsPerSample != 32 {
		t.Fatalf("Invalid BitsPerSample: %d", format.BitsPerSample)
	}

	samples, err := r.GetSamples(2)

	if err != nil {
		t.Fatal(err)
	}

	if len(samples) != 2 {
		t.Fatalf("Invalid NumSamples: %d", len(samples))
	}

	if samples[0][0] != 1 ||
		samples[0][1] != 2 ||
		samples[1][0] != 3 ||
		samples[1][1] != 4 {
		t.Fatalf("Invalid Samples: %d", samples)
	}
}
