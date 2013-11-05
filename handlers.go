package log

import (
	"io"
	"os"
	"regexp"
)

var DefaultStyle = map[Level]Style{
	Debug: DarkGrey,
	Info:  Black,
	Warn:  Yellow,
	Error: Red,
	Fatal: White | BgRed,
}

func NewLineWriter(w io.Writer, format string, style map[Level]Style) *LineWriter {
	return &LineWriter{w: w, format: format, style: style}
}

func NewTermWriter() *LineWriter {
	return &LineWriter{w: os.Stdout, format: DefaultFormat, style: DefaultStyle}
}

func NewTermLogger() *Logger {
	return NewLogger(NewTermWriter())
}

type LineWriter struct {
	w      io.Writer
	format string
	style map[Level]Style
}

func (l *LineWriter) HandleLog(e Entry) {
	line := e.Format(l.format) + "\n"
	if style, ok := l.style[e.Level]; ok {
		line = style.Apply(line)
	}
	io.WriteString(l.w, line)
}

func NewTestWriter() *TestWriter {
	return &TestWriter{}
}

type TestWriter struct {
	Entries []Entry
}

func (w *TestWriter) HandleLog(e Entry) {
	w.Entries = append(w.Entries, e)
}

func (w *TestWriter) Match(expr string) bool {
	return w.MatchLevel(expr, -1)
}

func (w *TestWriter) MatchLevel(expr string, lvl Level) bool {
	r := regexp.MustCompile(expr)
	for _, e := range w.Entries {
		if r.MatchString(e.Message) && e.Level == lvl || lvl == -1 {
			return true
		}
	}
	return false
}
