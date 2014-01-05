package main

import (
	"github.com/felixge/log"
)

func main() {
	log.Debug("Hello world.", log.Context{"a": 1, "foo": "bar"})
}
