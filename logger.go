package log

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func ParseLevel(s string) (Level, error) {
	s = strings.ToLower(s)
	for lvl, lvlStr := range levels {
		if lvlStr == s {
			return lvl, nil
		}
	}
	return 0, fmt.Errorf("Unknown level: %s", s)
}

const DefaultFormat = "2006-01-02T15:04:05.000Z [level] message (file:line)"

type Level int

func (l Level) String() string {
	return levels[l]
}

const (
	Debug Level = iota
	Info
	Warn
	Error
	Fatal
)

var levels = map[Level]string{
	Debug: "debug",
	Info:  "info",
	Warn:  "warn",
	Error: "error",
	Fatal: "fatal",
}

type Interface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{}) error
	Fatal(args ...interface{})
}

type Handler interface {
	HandleLog(Entry)
}

type Entry struct {
	Time    time.Time
	Level   Level
	Message string
	File    string
	Line    int
}

func (e Entry) Format(layout string) string {
	layout = e.Time.Format(layout)
	layout = strings.Replace(layout, "level", e.Level.String(), -1)
	layout = strings.Replace(layout, "message", e.Message, -1)
	layout = strings.Replace(layout, "file", e.File, -1)
	layout = strings.Replace(layout, "line", strconv.FormatInt(int64(e.Line), 10), -1)
	return layout
}

func NewLogger() *Logger {
	return &Logger{}
}

type Logger struct {
	handlers []*logHandler
}

type logHandler struct {
	lvl     Level
	handler Handler
}

func (l *Logger) Debug(args ...interface{}) {
	l.log(Debug, args)
}
func (l *Logger) Info(args ...interface{}) {
	l.log(Info, args)
}
func (l *Logger) Warn(args ...interface{}) {
	l.log(Warn, args)
}
func (l *Logger) Error(args ...interface{}) error {
	return entryToError(l.log(Error, args))
}
func (l *Logger) Fatal(args ...interface{}) {
	l.log(Fatal, args)
	os.Exit(1)
}

func (l *Logger) Handle(lvl Level, handler Handler) {
	l.handlers = append(l.handlers, &logHandler{lvl, handler})
}

func (l *Logger) log(lvl Level, args []interface{}) Entry {
	msg := formatMessage(args)
	_, file, lines, _ := runtime.Caller(2)

	e := Entry{
		Time:    time.Now().UTC(),
		Level:   lvl,
		Message: msg,
		File:    file,
		Line:    lines,
	}

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
