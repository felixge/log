**WORK IN PROGRESS:** Please come back later

# log

Package log is an attempt to provide the best logging library for Go.

## Simple and Beautiful

Logging should be simple:

```go
package main

import (
	"github.com/felixge/log"
)

func main() {
	log.Debug("A programming genius called Hank")
	log.Info("Wrote a system to %q his %q", "access", "bank")
	log.Warn("When his memory failed him")
	log.Error("They nailed him then jailed him")
	log.Fatal("Now his %q is %q and dank", "storage", "basic")
	// by W E Sword (http://goo.gl/R7Wjkv)
}
```

Produces:

![screenshot](http://felixge.github.io/log/screenshots/basic.png)
