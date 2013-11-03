package log

import (
	"bytes"
	"testing"
)

func TestLineLogger(t *testing.T) {
	b := bytes.NewBuffer(nil)
	l := Interface(NewLineLogger(b))
	l.Debug("Hello %s", "World")
	t.Log(b.String())
}
