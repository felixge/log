**WORK IN PROGRESS:** Please come back later

# log

Package log is an attempt to provide the best logging library for Go.

## Good by Default

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

Log entries are written to stdout, using UTC, millisecond precision, ANSI
colors, and include the call site they were created from:

![screenshot](http://felixge.github.io/log/screenshots/basic.png)

So if you're looking for a logging library that allows you to get started in no
time, package log is for you.

## Simple Interface

```go
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{}) error
	Panic(args ...interface{})
}
```

The Logger interface makes it easy to pick appropiate log levels and allows you
to decouple your app code from the underlaying logging implementation.

So if you're looking for a simple and well defined logging interface, package
log is for you.

However, if you're looking for a log package with 10+ log levels, package log
is not for you.

## Modular Design

@TODO insert code snippet

Everybody likes his logging just a little bit different, so package log fully
exposes its modular design, allowing you to put together the logger of your
dreams.

@TODO insert screenshot

So if you're looking for a logging library that won't get into your way,
package log is for you.

## Decent Performance

Package log comes with good benchmarks and has been observed to handle 1+
million / entries per second.

So if you'd like to sleep well, knowing that the CPU overhead of your logger is
negligible, package log is for you.

## Missing Features
