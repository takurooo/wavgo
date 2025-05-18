package wavgo

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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
	err = w.Open("testdata/write_test_output.wav")
	require.Nil(t, err)
	defer func() {
		if err := os.Remove("testdata/write_test_output.wav"); err != nil {
			t.Fatal(err)
		}
	}()

	samples := make([]Sample, 12)
	for i := 0; i < 12; i++ {
		for ch := 0; ch < int(format.NumChannels); ch++ {
			samples[i][ch] = i + ch
		}
	}
	err = w.WriteSamples(samples)
	require.Nil(t, err)

	err = w.Close()
	require.Nil(t, err)

	b, err := os.ReadFile("testdata/write_test_output.wav")
	require.Nil(t, err)
	want, err := os.ReadFile("testdata/write_test.wav.golden")
	require.Nil(t, err)
	require.Equal(t, want, b)
}

func TestWriterUnsupportedBitsPerSample(t *testing.T) {
	format := &Format{
		AudioFormat:   AudioFormatPCM,
		NumChannels:   2,
		SampleRate:    48000,
		ByteRate:      128000,
		BlockAlign:    4,
		BitsPerSample: 7,
	}

	w := NewWriter(format)
	err := w.Open("testdata/write_unsupported.wav")
	require.NoError(t, err)
	defer func() {
		w.Close()
		os.Remove("testdata/write_unsupported.wav")
	}()

	samples := make([]Sample, 1)
	err = w.WriteSamples(samples)
	require.True(t, errors.Is(err, ErrUnsupportedBitsPerSample))
}
