package log

import (
	"os"
)

func NewTermWriter() *LineWriter {
	return NewLineWriter(os.Stdout)
}
