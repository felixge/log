package log

import (
	"errors"
	"strings"
	"sync"
	"time"
)

var (
	ErrFlushTimeout = errors.New("Flush timed out.")
)

type Config struct {
	FlushTimeout time.Duration
}

func NewLogger(config Config, handlers ...Handler) *Logger {
	l := &Logger{config: config}
	for _, h := range handlers {
		l.Handle(DEBUG, h)
	}
	return l
}

type Logger struct {
	config   Config
	handlers []*logHandler
}

type logHandler struct {
	lvl     Level
	handler Handler
}

func NewError(e Entry) error {
	message := DefaultMessageFormatter.Format(e)
	message = strings.TrimRight(message, "\n")
	return errors.New(message)
}

// Debug logs at the Debug level.
func (l *Logger) Debug(args ...interface{}) {
	l.Log(NewEntryWithStack(DEBUG, 3, 1, args...))
}

// Debug logs at the Info level.
func (l *Logger) Info(args ...interface{}) {
	l.Log(NewEntryWithStack(INFO, 3, 1, args...))
}

// Warn logs at the Warn level.
func (l *Logger) Warn(args ...interface{}) {
	l.Log(NewEntryWithStack(WARN, 3, 1, args...))
}

// Error logs at the Error level and returns the formatted error message as
// an error for convenience.
func (l *Logger) Error(args ...interface{}) error {
	e := NewEntryWithStack(ERROR, 3, 1, args...)
	l.Log(e)
	return NewError(e)
}

// Panic logs at the Panic level, calls Flush() and then os.Exit(1).
func (l *Logger) Panic(args ...interface{}) {
	e := NewEntryWithStack(PANIC, 3, 1, args...)
	l.Log(e)
	panic(NewError(e))
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

func (l *Logger) Log(e Entry) {
	for _, h := range l.handlers {
		if e.Level >= h.lvl {
			h.handler.Log(e)
		}
	}
}
