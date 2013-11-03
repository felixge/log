package log

import (
	"io"
	"os"
)

func NewTermWriter() *TermWriter {
	return &TermWriter{format: DefaultFormat}
}

type TermWriter struct {
	format string
}

func (l *TermWriter) HandleLog(e Entry) {
	io.WriteString(os.Stdout, Format(l.format, e)+"\n")
}
