package log

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
	// Log processes the given Entry (e.g. writes it to a file, sends it to
	// a log service)
	Log(Entry)
	// Flush waits for any buffered data to be flushed and blocks new calls
	// to Log until it returns.
	Flush()
}

type Formatter interface {
	Format(e Entry) string
}

type ErrorHandler func(error)

type ErrEntryDropped struct {
	Entry Entry
}

func (e *ErrEntryDropped) Error() string {
	return "Dropped log entry."
}
