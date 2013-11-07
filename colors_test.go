package log

import (
	"fmt"
	"testing"
)

func ExampleTermStyle(b *testing.B) {
	fmt.Println((White | BgRed | Bold).Format("my text"))
}
