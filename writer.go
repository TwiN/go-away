package goaway

import "io"

func NewWriter(base io.Writer, detector *ProfanityDetector) *Writer {
	return &Writer{
		base:     base,
		detector: detector,
	}
}

type Writer struct {
	base     io.Writer
	detector *ProfanityDetector
	buf      []byte
}

func (w *Writer) Write(payload []byte) (int, error) {
	last := 0
	for i, char := range payload {
		if char != byte('\n') {
			continue
		}

		result := append(w.buf, payload[last:i+1]...)
		_, err := w.base.Write([]byte(w.detector.Censor(string(result))))
		if err != nil {
			return 0, err
		}

		w.buf = w.buf[:0]
		last = i + 1
	}
	w.buf = payload[last:]
	return len(payload), nil
}

func (w *Writer) Flush() error {
	if len(w.buf) == 0 {
		return nil
	}
	_, err := w.base.Write([]byte(w.detector.Censor(string(w.buf))))
	w.buf = w.buf[:0]
	return err
}
