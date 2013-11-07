package log

import (
	"io"
	"os"
	"regexp"
	"strings"
)

// DefaultTermStyle defines the default colors/style used by NewTermLogger
var DefaultTermStyle = map[Level]TermStyle{
	Debug: DarkGrey,
	Info:  Black,
	Warn:  Yellow,
	Error: Red,
	Fatal: White | BgRed,
}

// DefaultFormat defines the default log format used by NewTermLogger.
const DefaultFormat = "2006-01-02T15:04:05.000Z [level] message (file:line)"

// NewLineWriter returns a Handler that writes newline separated log entries
// to the given io.Writer w using the provided format and style.
func NewLineWriter(w io.Writer, format string, style map[Level]TermStyle) *LineWriter {
	return &LineWriter{w: w, format: format, style: style}
}

// NewTermWriter returns a *Logger that writes to os.Stdout using the
// DefaultFormat and DefaultTermStyle.
func NewTermLogger() *Logger {
	return NewLogger(NewLineWriter(os.Stdout, DefaultFormat, DefaultTermStyle))
}

// LineWriter is a Handler that provides newline separated logging.
type LineWriter struct {
	w      io.Writer
	format string
	style  map[Level]TermStyle
}

// HandleLog writes the given log entry to a new line.
// @TODO Process entries in another goroutine.
func (l *LineWriter) HandleLog(e Entry) {
	line := strings.Replace(e.Format(l.format), "\n", "", -1) + "\n"
	if style, ok := l.style[e.Level]; ok {
		line = style.Format(line)
	}
	io.WriteString(l.w, line)
}

// Flush waits for any buffered log Entries to be written out.
func (l *LineWriter) Flush() {
	// do nothing
}

// NewTestWriter returns a new *TestWriter.
func NewTestWriter() *TestWriter {
	return &TestWriter{}
}

// TestWriter is a Handler that simplifies writing unit tests for logging.
type TestWriter struct {
	Entries []Entry
}

// HandleLog attaches the given Entry to the Entries slice.
func (w *TestWriter) HandleLog(e Entry) {
	w.Entries = append(w.Entries, e)
}

// Flush satifies the Handler interface, but is a no-op for *TestWriter.
func (w *TestWriter) Flush() {}

// Match returns true if the regular expr matches the Message of a log
// Entry in the Entries slices.
func (w *TestWriter) Match(expr string) bool {
	return w.MatchLevel(expr, -1)
}

// Match returns true if the regular expr and lvl matches a log Entry in the
// Entries slices.
func (w *TestWriter) MatchLevel(expr string, lvl Level) bool {
	r := regexp.MustCompile(expr)
	for _, e := range w.Entries {
		if r.MatchString(e.Message) && e.Level == lvl || lvl == -1 {
			return true
		}
	}
	return false
}
