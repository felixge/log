package log

import (
	"strings"
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

func (l *TestWriter) Contains(s string) bool {
	return l.ContainsLevel(s, -1)
}

func (l *TestWriter) ContainsLevel(s string, lvl Level) bool {
	for _, e := range l.Entries {
		if strings.Contains(e.Message, s) && e.Level == lvl || lvl == -1 {
			return true
		}
	}
	return false
}
