package log

func ExampleTermWriter() {
	l := NewLogger()
	l.Handle(Debug, NewTermWriter())
	l.Debug("Hello %s", "World")
}
