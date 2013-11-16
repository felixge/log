package log

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	ErrFlushTimeout = errors.New("Flush timed out.")
)

type Config struct {
	FlushTimeout time.Duration
	FatalExit    bool
}

func NewLogger(config Config, handlers ...Handler) *Logger {
	l := &Logger{config: config}
	for _, h := range handlers {
		l.Handle(DEBUG, h)
	}
	return l
}

type Logger struct {
	config    Config
	handlers  []*logHandler
}

type logHandler struct {
	lvl     Level
	handler Handler
}

// Debug logs at the Debug level.
func (l *Logger) Debug(args ...interface{}) {
	l.log(DEBUG, args)
}

// Debug logs at the Info level.
func (l *Logger) Info(args ...interface{}) {
	l.log(INFO, args)
}

// Warn logs at the Warn level.
func (l *Logger) Warn(args ...interface{}) {
	l.log(WARN, args)
}

// Error logs at the Error level and returns the formatted error message as
// an error for convenience.
func (l *Logger) Error(args ...interface{}) error {
	return entryToError(l.log(ERROR, args))
}

// Fatal logs at the Fatal level, calls Flush() and then os.Exit(1).
func (l *Logger) Fatal(args ...interface{}) {
	l.log(FATAL, args)
	l.Flush()
	if l.config.FatalExit {
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
		time.Sleep(l.config.FlushTimeout)
		err <- ErrFlushTimeout
	}()
	return <-err
}

func (l *Logger) Handle(lvl Level, handler Handler) {
	l.handlers = append(l.handlers, &logHandler{lvl, handler})
}

func (l *Logger) log(lvl Level, args []interface{}) Entry {
	e := NewEntry(lvl, args...)

	for _, h := range l.handlers {
		if e.Level >= h.lvl {
			h.handler.Log(e)
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
