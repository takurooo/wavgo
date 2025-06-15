package wavgo

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/takurooo/wavgo/internal/riff"
)

func TestReader(t *testing.T) {
	var err error
	r := NewReader()
	err = r.Open("testdata/read_test.wav")
	require.Nil(t, err)
	defer func() {
		cerr := r.Close()
		require.Nil(t, cerr)
	}()

	err = r.Load()
	require.Nil(t, err)

	format := r.GetFormat()
	wantFormat := Format{
		AudioFormat:   1,
		NumChannels:   2,
		SampleRate:    44100,
		ByteRate:      176400,
		BlockAlign:    4,
		BitsPerSample: 16,
	}
	require.Equal(t, wantFormat, format)
	require.Equal(t, uint32(2), r.GetNumSamples())
	require.Equal(t, uint32(2), r.GetNumSamplesLeft())

	samples, err := r.GetSamples(2)
	require.Nil(t, err)
	require.Equal(t, 2, len(samples))
	require.Equal(t, 1, samples[0][0])
	require.Equal(t, 2, samples[0][1])
	require.Equal(t, 3, samples[1][0])
	require.Equal(t, 4, samples[1][1])
	require.Equal(t, uint32(0), r.GetNumSamplesLeft())
}

func TestReaderUnsupportedBitsPerSample(t *testing.T) {
	r := NewReader()
	err := r.Open("testdata/read_test.wav")
	require.NoError(t, err)
	defer r.Close()

	err = r.Load()
	require.NoError(t, err)

	r.format.BitsPerSample = 7
	samples, err := r.GetSamples(1)
	require.Nil(t, samples)
	require.True(t, errors.Is(err, ErrUnsupportedBitsPerSample))
}

func TestReaderGetSamplesValidation(t *testing.T) {
	r := NewReader()
	err := r.Open("testdata/read_test.wav")
	require.NoError(t, err)
	defer r.Close()

	err = r.Load()
	require.NoError(t, err)

	// Test negative numSamples
	samples, err := r.GetSamples(-1)
	require.Nil(t, samples)
	require.EqualError(t, err, "numSamples cannot be negative")

	// Test requesting more samples than available
	samples, err = r.GetSamples(int(r.GetNumSamples()) + 1)
	require.Nil(t, samples)
	require.EqualError(t, err, "requested samples exceed remaining samples")
}

func TestReaderOpenFileError(t *testing.T) {
	r := NewReader()
	err := r.Open("non/existent/file.wav")
	require.Error(t, err)
}

func TestReaderLoadBeforeOpen(t *testing.T) {
	r := NewReader()
	err := r.Load()
	require.Error(t, err)
}

func TestReaderCloseWithoutOpen(t *testing.T) {
	r := NewReader()
	err := r.Close()
	require.NoError(t, err) // Should not error when no file is open
}

func TestReaderGetSamplesBeforeLoad(t *testing.T) {
	r := NewReader()
	err := r.Open("testdata/read_test.wav")
	require.NoError(t, err)
	defer r.Close()

	// Try to get samples before loading
	samples, err := r.GetSamples(1)
	require.Nil(t, samples)
	require.Error(t, err) // Should error because br is nil
}

func TestReaderFormatValidation(t *testing.T) {
	t.Run("ZeroChannels", func(t *testing.T) {
		mockChunk := &riff.Chunk{
			ID:   "fmt ",
			Size: 16,
			Data: []byte{
				0x01, 0x00, // AudioFormat = 1
				0x00, 0x00, // NumChannels = 0
				0x44, 0xAC, 0x00, 0x00, // SampleRate = 44100
				0x88, 0x58, 0x01, 0x00, // ByteRate
				0x02, 0x00, // BlockAlign = 2
				0x10, 0x00, // BitsPerSample = 16
			},
		}

		_, err := parseFormatChunkData(mockChunk)
		require.Error(t, err)
		require.EqualError(t, err, "invalid NumChannels: must be greater than 0")
	})

	t.Run("ZeroSampleRate", func(t *testing.T) {
		mockChunk := &riff.Chunk{
			ID:   "fmt ",
			Size: 16,
			Data: []byte{
				0x01, 0x00, // AudioFormat = 1
				0x02, 0x00, // NumChannels = 2
				0x00, 0x00, 0x00, 0x00, // SampleRate = 0
				0x88, 0x58, 0x01, 0x00, // ByteRate
				0x02, 0x00, // BlockAlign = 2
				0x10, 0x00, // BitsPerSample = 16
			},
		}

		_, err := parseFormatChunkData(mockChunk)
		require.Error(t, err)
		require.EqualError(t, err, "invalid SampleRate: must be greater than 0")
	})

	t.Run("ZeroBlockAlign", func(t *testing.T) {
		mockChunk := &riff.Chunk{
			ID:   "fmt ",
			Size: 16,
			Data: []byte{
				0x01, 0x00, // AudioFormat = 1
				0x02, 0x00, // NumChannels = 2
				0x44, 0xAC, 0x00, 0x00, // SampleRate = 44100
				0x88, 0x58, 0x01, 0x00, // ByteRate
				0x00, 0x00, // BlockAlign = 0
				0x10, 0x00, // BitsPerSample = 16
			},
		}

		_, err := parseFormatChunkData(mockChunk)
		require.Error(t, err)
		require.EqualError(t, err, "invalid BlockAlign: must be greater than 0")
	})

	t.Run("ZeroBitsPerSample", func(t *testing.T) {
		mockChunk := &riff.Chunk{
			ID:   "fmt ",
			Size: 16,
			Data: []byte{
				0x01, 0x00, // AudioFormat = 1
				0x02, 0x00, // NumChannels = 2
				0x44, 0xAC, 0x00, 0x00, // SampleRate = 44100
				0x88, 0x58, 0x01, 0x00, // ByteRate
				0x04, 0x00, // BlockAlign = 4
				0x00, 0x00, // BitsPerSample = 0
			},
		}

		_, err := parseFormatChunkData(mockChunk)
		require.Error(t, err)
		require.EqualError(t, err, "invalid BitsPerSample: must be greater than 0")
	})
}

func TestReaderGetNumSamplesAndLeft(t *testing.T) {
	r := NewReader()
	err := r.Open("testdata/read_test.wav")
	require.NoError(t, err)
	defer r.Close()

	err = r.Load()
	require.NoError(t, err)

	// Check initial counts
	totalSamples := r.GetNumSamples()
	samplesLeft := r.GetNumSamplesLeft()
	require.Equal(t, totalSamples, samplesLeft)
	require.Equal(t, uint32(2), totalSamples)

	// Read one sample
	samples, err := r.GetSamples(1)
	require.NoError(t, err)
	require.Len(t, samples, 1)

	// Check counts after reading
	require.Equal(t, uint32(2), r.GetNumSamples())     // Total should not change
	require.Equal(t, uint32(1), r.GetNumSamplesLeft()) // Left should decrease

	// Read remaining sample
	samples, err = r.GetSamples(1)
	require.NoError(t, err)
	require.Len(t, samples, 1)

	// Check final counts
	require.Equal(t, uint32(2), r.GetNumSamples())
	require.Equal(t, uint32(0), r.GetNumSamplesLeft())
}
