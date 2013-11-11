package log

import (
	"io"
	"os"
	"regexp"
	"time"
)

// NewTermWriter returns a *Logger that writes to os.Stdout using the
// DefaultFormat and DefaultTermStyle.
func NewTermLogger() *Logger {
	return NewLogger(NewLineWriter(os.Stdout, DefaultFormat, DefaultTermStyle))
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
