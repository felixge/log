package log

import (
	"io"
)

func NewLineWriter(w io.Writer) *LineWriter {
	return &LineWriter{w: w, format: DefaultFormat}
}

type LineWriter struct {
	w io.Writer
	format string
}

func (l *LineWriter) HandleLog(e Entry) {
	io.WriteString(l.w, e.Format(l.format)+"\n")
}
