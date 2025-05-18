package wavgo

import "errors"

// ErrUnsupportedBitsPerSample is returned when the number of bits per sample is not supported.
var ErrUnsupportedBitsPerSample = errors.New("unsupported BitsPerSample")
