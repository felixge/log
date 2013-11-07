package log

import (
	"fmt"
	"testing"
)

func ExampleNewTermLogger(t *testing.T) {
	l := NewTermLogger()
	l.Debug("this is debugging")
	l.Info("this is an info")
	l.Warn("this is a warning")
	l.Error("this is an error")
	l.Fatal("this is fatal")
}

func ExampleTermStyle(b *testing.B) {
	fmt.Println((White | BgRed | Bold).Format("my text"))
}
