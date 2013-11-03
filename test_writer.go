package log

import (
	"regexp"
)

func NewTestWriter() *TestWriter {
	return &TestWriter{}
}

type TestWriter struct {
	Entries []Entry
}

func (l *TestWriter) HandleLog(e Entry) {
	l.Entries = append(l.Entries, e)
}

func (l *TestWriter) Match(expr string) bool {
	return l.MatchLevel(expr, -1)
}

func (l *TestWriter) MatchLevel(expr string, lvl Level) bool {
	r := regexp.MustCompile(expr)
	for _, e := range l.Entries {
		if r.MatchString(e.Message) && e.Level == lvl || lvl == -1 {
			return true
		}
	}
	return false
}
