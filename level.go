package log

import (
	"fmt"
	"strings"
)

// ParseLevel returns the Level value for the given string, or an error if
// such a level does not exist. e.G. "debug" will return DEBUG.
func ParseLevel(s string) (Level, error) {
	s = strings.ToLower(s)
	for lvl, lvlStr := range levels {
		if lvlStr == s {
			return lvl, nil
		}
	}
	return 0, fmt.Errorf("Unknown level: %s", s)
}

type Level int

// String returns the human readable name of the log level. e.G. Debug will
// return "debug"
func (l Level) String() string {
	return levels[l]
}

// The available log levels along with their recommended usage. Always log at
// the INFO level in production.
const (
	DEBUG Level = iota // Development details (e.g. raw input data)
	INFO               // Regular event (e.g. user login)
	WARN               // Undesireable event (e.g. invalid user input)
	ERROR              // E-mail somebody (e.g. could not save record)
	FATAL              // Call somebody (e.g. database down)
)

var levels = map[Level]string{
	DEBUG: "debug",
	INFO:  "info",
	WARN:  "warn",
	ERROR: "error",
	FATAL: "fatal",
}
