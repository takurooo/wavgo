package wav

// BufferWriter ...
type BufferWriter struct {
	b []byte
	i int64
}

// NewBufferWriter ...
func NewBufferWriter(b []byte) *BufferWriter {
	return &BufferWriter{b, 0}
}

// Write ...
func (w *BufferWriter) Write(p []byte) (n int, err error) {
	w.b = append(w.b, p...)
	return len(p), nil
}

// WriteAt ...
func (w *BufferWriter) WriteAt(p []byte, off int64) (n int, err error) {

	bufLen := int64(len(w.b))
	if bufLen < off {
		w.b = append(w.b[:bufLen], p...)
	} else { //bufLen >= off
		w.b = append(w.b[:off], p...)
		if bufLen > int64(len(p))+off {
			w.b = append(w.b, w.b[off:]...)
		}
	}
	return len(p), nil
}
