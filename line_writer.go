package log

import (
	"io"
)

func NewLineWriter(w io.Writer) *LineWriter {
	return NewLineWriterFormat(w, DefaultFormat)
}

func NewLineWriterFormat(w io.Writer, format string) *LineWriter {
	return &LineWriter{w: w, format: format}
}

type LineWriter struct {
	w      io.Writer
	format string
}

func (l *LineWriter) HandleLog(e Entry) {
	io.WriteString(l.w, e.Format(l.format)+"\n")
}
