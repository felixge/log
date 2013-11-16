package log

import (
	"fmt"
	"testing"
	"time"
)

func TestLineFormatterFormat_defaultFormat(t *testing.T) {
	e := Entry{
		Time:     time.Now(),
		Level:    INFO,
		Message:  "foo",
		File:     "bar.go",
		Line:     23,
		Function: "foo.bar",
	}

	f := NewLineFormatter(DefaultLayout, nil)
	str := f.Format(e)
	expected := fmt.Sprintf(
		"[%s UTC] [%s] %s (%s:%d)",
		e.Time.UTC().Format("2006-01-02 15:04:05.000"),
		INFO,
		e.Message,
		e.Function,
		e.Line,
	)
	if str != expected {
		t.Errorf("Bad result: %q != %q", str, expected)
	}
}

func TestLineFormatterFormat_customFormat(t *testing.T) {
	e := Entry{
		Time:     time.Now(),
		Level:    INFO,
		Message:  "foo",
		File:     "bar.go",
		Line:     23,
		Function: "foo.bar",
	}

	f := NewLineFormatter("2006/01/02 15:04:05.000 level message file/line/function", nil)
	str := f.Format(e)
	expected := fmt.Sprintf(
		"%s %s %s %s/%d/%s",
		e.Time.Format("2006/01/02 15:04:05.000"),
		INFO,
		e.Message,
		e.File,
		e.Line,
		e.Function,
	)
	if str != expected {
		t.Errorf("Bad result: %q != %q", str, expected)
	}
}
