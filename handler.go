package log

type Handler func(lvl Level, args ...interface{})

func (h Handler) Debug(args ...interface{}) {
	h(Debug, args...)
}
func (h Handler) Info(args ...interface{}) {
	h(Info, args...)
}
func (h Handler) Warn(args ...interface{}) error {
	return nil
}
func (h Handler) Error(args ...interface{}) error {
	return nil
}
func (h Handler) Fatal(args ...interface{}) {
}
