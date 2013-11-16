package log

import (
	"regexp"
)

func NewTestLogger() (*Logger, *TestHandler) {
	t := NewTestHandler()
	l := NewLogger(DefaultConfig, t)
	return l, t
}

// NewTestHandler returns a new *TestHandler.
func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

// TestHandler is a Handler that simplifies writing unit tests for logging.
type TestHandler struct {
	Entries []Entry
}

// Log attaches the given Entry to the Entries slice.
func (w *TestHandler) Log(e Entry) {
	w.Entries = append(w.Entries, e)
}

// Flush satifies the Handler interface, but is a no-op for *TestHandler.
func (w *TestHandler) Flush() {}

// Match returns true if the regular expr matches the Message of a log
// Entry in the Entries slices.
func (w *TestHandler) Match(expr string) bool {
	return w.MatchLevel(expr, -1)
}

// Match returns true if the regular expr and lvl matches a log Entry in the
// Entries slices.
func (w *TestHandler) MatchLevel(expr string, lvl Level) bool {
	r := regexp.MustCompile(expr)
	for _, e := range w.Entries {
		if r.MatchString(e.Message) && e.Level == lvl || lvl == -1 {
			return true
		}
	}
	return false
}
