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

func TestWriterEmptyFile(t *testing.T) {
	format := &Format{
		AudioFormat:   AudioFormatPCM,
		NumChannels:   1,
		SampleRate:    44100,
		ByteRate:      88200,
		BlockAlign:    2,
		BitsPerSample: 16,
	}

	w := NewWriter(format)
	err := w.Open("testdata/empty_test.wav")
	require.NoError(t, err)
	defer os.Remove("testdata/empty_test.wav")

	// Close without writing any samples
	err = w.Close()
	require.NoError(t, err)

	// Verify file exists and has correct header size
	info, err := os.Stat("testdata/empty_test.wav")
	require.NoError(t, err)
	require.Greater(t, info.Size(), int64(0)) // File should exist and have some content
}

func TestWriterMultipleFormats(t *testing.T) {
	testCases := []struct {
		name    string
		format  *Format
		samples []Sample
	}{
		{
			name: "8-bit mono",
			format: &Format{
				AudioFormat:   AudioFormatPCM,
				NumChannels:   1,
				SampleRate:    22050,
				ByteRate:      22050,
				BlockAlign:    1,
				BitsPerSample: 8,
			},
			samples: []Sample{{100, 0}},
		},
		{
			name: "16-bit stereo",
			format: &Format{
				AudioFormat:   AudioFormatPCM,
				NumChannels:   2,
				SampleRate:    44100,
				ByteRate:      176400,
				BlockAlign:    4,
				BitsPerSample: 16,
			},
			samples: []Sample{{-1000, 1000}},
		},
		{
			name: "24-bit stereo",
			format: &Format{
				AudioFormat:   AudioFormatPCM,
				NumChannels:   2,
				SampleRate:    48000,
				ByteRate:      288000,
				BlockAlign:    6,
				BitsPerSample: 24,
			},
			samples: []Sample{{-100000, 100000}},
		},
		{
			name: "32-bit stereo",
			format: &Format{
				AudioFormat:   AudioFormatPCM,
				NumChannels:   2,
				SampleRate:    96000,
				ByteRate:      768000,
				BlockAlign:    8,
				BitsPerSample: 32,
			},
			samples: []Sample{{-1000000, 1000000}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filename := "testdata/TestWriterMultipleFormats.wav"
			w := NewWriter(tc.format)
			err := w.Open(filename)
			require.NoError(t, err)
			defer os.Remove(filename)

			err = w.WriteSamples(tc.samples)
			require.NoError(t, err)
			err = w.Close()
			require.NoError(t, err)
		})
	}
}

func TestWriterIncrementalWrites(t *testing.T) {
	format := &Format{
		AudioFormat:   AudioFormatPCM,
		NumChannels:   2,
		SampleRate:    44100,
		ByteRate:      176400,
		BlockAlign:    4,
		BitsPerSample: 16,
	}

	w := NewWriter(format)
	err := w.Open("testdata/TestWriterIncrementalWrites.wav")
	require.NoError(t, err)
	defer os.Remove("testdata/TestWriterIncrementalWrites.wav")

	// Write samples incrementally
	for i := 0; i < 5; i++ {
		samples := []Sample{{i, i + 1}}
		err = w.WriteSamples(samples)
		require.NoError(t, err)
	}

	err = w.Close()
	require.NoError(t, err)

	// Verify file was created correctly
	info, err := os.Stat("testdata/TestWriterIncrementalWrites.wav")
	require.NoError(t, err)
	require.Greater(t, info.Size(), int64(44)) // Should be larger than header
}

func TestWriterZeroSamples(t *testing.T) {
	format := &Format{
		AudioFormat:   AudioFormatPCM,
		NumChannels:   1,
		SampleRate:    44100,
		ByteRate:      88200,
		BlockAlign:    2,
		BitsPerSample: 16,
	}

	w := NewWriter(format)
	err := w.Open("testdata/TestWriterZeroSamples.wav")
	require.NoError(t, err)
	defer os.Remove("testdata/TestWriterZeroSamples.wav")

	// Write empty slice of samples
	err = w.WriteSamples([]Sample{})
	require.NoError(t, err)

	err = w.Close()
	require.NoError(t, err)
}

func TestWriterOpenFileError(t *testing.T) {
	format := &Format{
		AudioFormat:   AudioFormatPCM,
		NumChannels:   1,
		SampleRate:    44100,
		ByteRate:      88200,
		BlockAlign:    2,
		BitsPerSample: 16,
	}

	w := NewWriter(format)
	// Try to open file in non-existent directory
	err := w.Open("/non/existent/directory/test.wav")
	require.Error(t, err)
}

func TestWriterLargeFile(t *testing.T) {
	format := &Format{
		AudioFormat:   AudioFormatPCM,
		NumChannels:   1,
		SampleRate:    44100,
		ByteRate:      88200,
		BlockAlign:    2,
		BitsPerSample: 16,
	}

	w := NewWriter(format)
	err := w.Open("testdata/TestWriterLargeFile.wav")
	require.NoError(t, err)
	defer os.Remove("testdata/TestWriterLargeFile.wav")

	// Write a large number of samples
	samples := make([]Sample, 1000)
	for i := range samples {
		samples[i] = Sample{i % 32767, 0}
	}

	err = w.WriteSamples(samples)
	require.NoError(t, err)

	err = w.Close()
	require.NoError(t, err)

	// Verify file size is correct
	info, err := os.Stat("testdata/TestWriterLargeFile.wav")
	require.NoError(t, err)
	expectedSize := int64(44 + 1000*2) // Header + samples * bytes per sample
	require.Equal(t, expectedSize, info.Size())
}
