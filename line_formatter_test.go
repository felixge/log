package log

import (
	"fmt"
	"testing"
	"time"
)

func TestLineFormatterFormat_defaultLayout(t *testing.T) {
	message := "foo"
	e := Entry{
		Time:  time.Now(),
		Level: INFO,
		Args:  []interface{}{message},
		Stack: []StackFrame{{file: "bar.go", line: 23, function: "foo.bar"}},
	}

	f := NewLineFormatter(DefaultLayout, nil)
	str := f.Format(e)
	expected := fmt.Sprintf(
		"[%s UTC] [%s] %s (%s:%d)\n",
		e.Time.UTC().Format("2006-01-02 15:04:05.000"),
		INFO,
		message,
		e.Function(),
		e.Line(),
	)
	if str != expected {
		t.Errorf("Bad result: %q != %q", str, expected)
	}
}

func TestLineFormatterFormat_customFormat(t *testing.T) {
	message := "foo"
	e := Entry{
		Time:     time.Now(),
		Level:    INFO,
		Args:     []interface{}{message},
		Stack: []StackFrame{{file: "bar.go", line: 23, function: "foo.bar"}},
	}

	f := NewLineFormatter("2006/01/02 15:04:05.000 level message file/line/function", nil)
	str := f.Format(e)
	expected := fmt.Sprintf(
		"%s %s %s %s/%d/%s\n",
		e.Time.Format("2006/01/02 15:04:05.000"),
		INFO,
		message,
		e.File(),
		e.Line(),
		e.Function(),
	)
	if str != expected {
		t.Errorf("Bad result: %q != %q", str, expected)
	}
}
