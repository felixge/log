package log

import (
	"io"
	"time"
)

func NewSlowWriter(w io.Writer, d time.Duration) *SlowWriter {
	return &SlowWriter{w: w, d: d}
}

type SlowWriter struct {
	w io.Writer
	d time.Duration
}

func (d *SlowWriter) Write(b []byte) (int, error) {
	time.Sleep(d.d)
	return d.w.Write(b)
}
