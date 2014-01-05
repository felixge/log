package log

import (
	"runtime"
	"time"
)

// @TODO remove?
func NewEntry(lvl Level, args ...interface{}) Entry {
	return Entry{
		Time:  time.Now(),
		Level: lvl,
		Args:  args,
	}
}

func NewEntryWithStack(lvl Level, skip int, count int, args ...interface{}) Entry {
	return Entry{
		Time:  time.Now(),
		Level: lvl,
		Args:  args,
		Stack: CaptureStack(skip, count),
	}
}

type Entry struct {
	Time  time.Time
	Level Level
	Args  []interface{}
	Stack []StackFrame
}

func (e Entry) File() (file string) {
	if len(e.Stack) > 0 {
		file = e.Stack[0].File()
	}
	return
}

func (e Entry) Function() (function string) {
	if len(e.Stack) > 0 {
		function = e.Stack[0].Function()
	}
	return
}

func (e Entry) Line() (line int) {
	if len(e.Stack) > 0 {
		line = e.Stack[0].Line()
	}
	return
}

func CaptureStack(skip int, count int) (stack []StackFrame) {
	for ; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}

		function := runtime.FuncForPC(pc).Name()

		frame := StackFrame{file: file, line: line, function: function}
		stack = append(stack, frame)

		count--
		if count == 0 {
			break
		}
	}
	return stack
}

type StackFrame struct {
	file     string
	line     int
	function string
}

func (s StackFrame) File() string {
	return s.file
}

func (s StackFrame) Line() int {
	return s.line
}

func (s StackFrame) Function() string {
	return s.function
}
