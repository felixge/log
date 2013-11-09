package log

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ParseLevel returns the Level value for the given string, or an error if
// such a level does not exist. e.G. "debug" will return Debug.
func ParseLevel(s string) (Level, error) {
	s = strings.ToLower(s)
	for lvl, lvlStr := range levels {
		if lvlStr == s {
			return lvl, nil
		}
	}
	return 0, fmt.Errorf("Unknown level: %s", s)
}

// Level defines a log level.
type Level int

// String returns the human readable name of the log level. e.G. Debug will
// return "debug"
func (l Level) String() string {
	return levels[l]
}

// The available log levels along with their recommended usage. Always log at
// the info Level in production.
const (
	Debug Level = iota // Development details (e.g. raw input data)
	Info               // Regular event (e.g. user login)
	Warn               // Undesireable event (e.g. invalid user input)
	Error              // E-mail somebody (e.g. could not save record)
	Fatal              // Call somebody (e.g. database down)
)

var levels = map[Level]string{
	Debug: "debug",
	Info:  "info",
	Warn:  "warn",
	Error: "error",
	Fatal: "fatal",
}

// Interface defines the log interface provided by this package. Use this when
// passing *Logger instances around.
type Interface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{}) error
	Fatal(args ...interface{})
	Panic()
}

// Handler is used to implement log handlers.
type Handler interface {
	// HandleLog processes the given Entry. Can be async if needed.
	HandleLog(Entry)
	// Flush waits for any buffered data to be processed.
	Flush()
}

func NewEntry(lvl Level, args ...interface{}) Entry {
	fn, file, line := getCaller()
	return Entry{
		Time:     time.Now(),
		Level:    lvl,
		Message:  formatMessage(args),
		File:     file,
		Line:     line,
		Function: fn,
	}
}

func getCaller() (fn, file string, line int) {
	var (
		skip     = 0
		thisFile string
		ok       bool
		pc       uintptr
	)

	for ; ; skip++ {
		pc, file, line, ok = runtime.Caller(skip)
		if !ok {
			break
		} else if skip == 0 {
			thisFile = file
			continue
		} else if file != thisFile {
			fn = runtime.FuncForPC(pc).Name()
			if fn != "runtime.panic" {
				break
			}
		}
	}
	return fn, file, line
}

type Entry struct {
	Time     time.Time
	Level    Level
	Message  string
	File     string
	Function string
	Line     int
}

func (e Entry) Format(layout string) string {
	if strings.Contains(layout, "UTC") {
		layout = e.Time.UTC().Format(layout)
	} else {
		layout = e.Time.Format(layout)
	}
	layout = strings.Replace(layout, "level", e.Level.String(), -1)
	layout = strings.Replace(layout, "message", e.Message, -1)
	layout = strings.Replace(layout, "file", e.File, -1)
	layout = strings.Replace(layout, "function", e.Function, -1)
	layout = strings.Replace(layout, "line", strconv.FormatInt(int64(e.Line), 10), -1)
	return layout
}

func NewLogger(handlers ...Handler) *Logger {
	l := &Logger{flushTimeout: DefaultFlushTimeout, exit: DefaultExit}
	for _, h := range handlers {
		l.Handle(Debug, h)
	}
	return l
}

type Logger struct {
	handlers     []*logHandler
	flushTimeout time.Duration
	exit         bool
}

type logHandler struct {
	lvl     Level
	handler Handler
}

// Debug logs at the Debug level.
func (l *Logger) Debug(args ...interface{}) {
	l.log(Debug, args)
}

// Debug logs at the Info level.
func (l *Logger) Info(args ...interface{}) {
	l.log(Info, args)
}

// Warn logs at the Warn level.
func (l *Logger) Warn(args ...interface{}) {
	l.log(Warn, args)
}

// Error logs at the Error level and returns the formatted error message as
// an error for convenience.
func (l *Logger) Error(args ...interface{}) error {
	return entryToError(l.log(Error, args))
}

var DefaultExit = true

// Fatal logs at the Fatal level, calls Flush() and then os.Exit(1).
func (l *Logger) Fatal(args ...interface{}) {
	l.log(Fatal, args)
	l.Flush()
	if l.exit {
		os.Exit(1)
	}
}

func (l *Logger) Panic() {
	if p := recover(); p != nil {
		switch p.(type) {
		case string:
			l.Fatal("panic: %s", p)
		default:
			l.Fatal("panic: %#v", p)
		}
	}
}

var (
	DefaultFlushTimeout = 30 * time.Second
	ErrFlushTimeout     = errors.New("Flush timed out.")
)

func (l *Logger) Flush() error {
	var wg sync.WaitGroup
	for _, h := range l.handlers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			h.handler.Flush()
		}()
	}

	err := make(chan error, 1)
	go func() {
		wg.Wait()
		err <- nil
	}()
	go func() {
		time.Sleep(l.flushTimeout)
		err <- ErrFlushTimeout
	}()
	return <-err
}

func (l *Logger) SetFlushTimeout(d time.Duration) {
	l.flushTimeout = d
}

func (l *Logger) SetExit(exit bool) {
	l.exit = false
}

func (l *Logger) Handle(lvl Level, handler Handler) {
	l.handlers = append(l.handlers, &logHandler{lvl, handler})
}

func (l *Logger) log(lvl Level, args []interface{}) Entry {
	e := NewEntry(lvl, args...)
	for _, h := range l.handlers {
		if e.Level >= h.lvl {
			h.handler.HandleLog(e)
		}
	}
	return e
}

func formatMessage(args []interface{}) string {
	if len(args) > 0 {
		if formatMessage, ok := args[0].(string); ok {
			return fmt.Sprintf(formatMessage, args[1:]...)
		}
	}
	return fmt.Sprint(args...)
}

func entryToError(e Entry) error {
	return errors.New(e.Message)
}
