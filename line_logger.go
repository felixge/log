package log

import (
	"fmt"
	"io"
)

func NewLineLogger(w io.Writer) *LineLogger {
	l := &LineLogger{w: w}
	l.Handler = l.handle
	return l
}

type LineLogger struct {
	Handler
	w io.Writer
}

func (l *LineLogger) handle(lvl Level, args ...interface{}) {
	str := format(args...)
	io.WriteString(l.w, str)
}

func format(args ...interface{}) string {
	if len(args) > 0 {
		if format, ok := args[0].(string); ok {
			return fmt.Sprintf(format, args[1:]...)
		}
	}
	return fmt.Sprint(args...)
}
