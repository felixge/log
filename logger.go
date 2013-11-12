package log

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	ErrFlushTimeout = errors.New("Flush timed out.")
)

// ParseLevel returns the Level value for the given string, or an error if
// such a level does not exist. e.G. "debug" will return DEBUG.
func ParseLevel(s string) (Level, error) {
	s = strings.ToLower(s)
	for lvl, lvlStr := range levels {
		if lvlStr == s {
			return lvl, nil
		}
	}
	return 0, fmt.Errorf("Unknown level: %s", s)
}

type Level int

// String returns the human readable name of the log level. e.G. Debug will
// return "debug"
func (l Level) String() string {
	return levels[l]
}

// The available log levels along with their recommended usage. Always log at
// the INFO level in production.
const (
	DEBUG Level = iota // Development details (e.g. raw input data)
	INFO               // Regular event (e.g. user login)
	WARN               // Undesireable event (e.g. invalid user input)
	ERROR              // E-mail somebody (e.g. could not save record)
	FATAL              // Call somebody (e.g. database down)
)

var levels = map[Level]string{
	DEBUG: "debug",
	INFO:  "info",
	WARN:  "warn",
	ERROR: "error",
	FATAL: "fatal",
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
	// HandleLog processes the given Entry (e.g. writes it to a file, sends it to
	// a log service)
	HandleLog(Entry)
	// Flush waits for any buffered data to be flushed and blocks new calls
	// to HandleLog until it returns.
	Flush()
}

type Config struct {
	FlushTimeout time.Duration
	FatalExit    bool
}

var DefaultConfig = Config{
	FlushTimeout: 30 * time.Second,
	FatalExit:    true,
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
	flushLock sync.Mutex
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
	l.flushLock.Lock()
	defer l.flushLock.Unlock()

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

	l.flushLock.Lock()
	defer l.flushLock.Unlock()

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
