package log

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
	Fatal
)

type Interface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{}) error
	Error(args ...interface{}) error
	Fatal(args ...interface{})
}
