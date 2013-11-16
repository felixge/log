package main

import (
	"github.com/felixge/log"
)

func main() {
	log.Debug("A programming genius called Hank")
	log.Info("Wrote a system to '%s' his '%s'", "access", "bank")
	log.Warn("When his memory failed him")
	log.Error("They nailed him then jailed him")
	log.Fatal("Now his '%s' is '%s' and dank", "storage", "basic")
}
