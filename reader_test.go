package wavgo

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	var err error
	r := NewReader()
	err = r.Open("testdata/read_test.wav")
	require.Nil(t, err)
	defer r.Close()

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
