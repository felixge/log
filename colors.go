package log

import (
	"fmt"
	"strconv"
	"strings"
)

// @TODO rename file to termstyle.go

// TermStyle provides a simple abstraction for ANSI/VT100 color/style escape
// sequences.
type TermStyle int

const (
	// Foreground colors
	Black TermStyle = 1 << iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	LightGrey
	DarkGrey
	LightRed
	LightGreen
	LightYellow
	LightBlue
	LightMagenta
	LightCyan
	White

	// Background colors
	BgBlack
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgLightGrey
	BgDarkGrey
	BgLightRed
	BgLightGreen
	BgLightYellow
	BgLightBlue
	BgLightMagenta
	BgLightCyan
	BgWhite

	// Special styles
	Bold
	Dim
	Underlined
	Blink // does not work with most terminal emulators (e.g. Terminal/iTerm2 on OSX)
	Reverse
	Hidden
)

// from http://misc.flogisoft.com/bash/tip_colors_and_formatting
var styleMap = map[TermStyle]uint8{
	Bold:       1,
	Dim:        2,
	Underlined: 4,
	Blink:      5,
	Reverse:    7,
	Hidden:     8,

	Black:        30,
	Red:          31,
	Green:        32,
	Yellow:       33,
	Blue:         34,
	Magenta:      35,
	Cyan:         36,
	LightGrey:    37,
	DarkGrey:     90,
	LightRed:     91,
	LightGreen:   92,
	LightYellow:  93,
	LightBlue:    94,
	LightMagenta: 95,
	LightCyan:    96,
	White:        97,

	BgBlack:        40,
	BgRed:          41,
	BgGreen:        42,
	BgYellow:       43,
	BgBlue:         44,
	BgMagenta:      45,
	BgCyan:         46,
	BgLightGrey:    47,
	BgDarkGrey:     100,
	BgLightRed:     101,
	BgLightGreen:   102,
	BgLightYellow:  103,
	BgLightBlue:    104,
	BgLightMagenta: 105,
	BgLightCyan:    106,
	BgWhite:        107,
}

// Format wraps the given str with the right terminal escape sequences.
func (s TermStyle) Format(str string) string {
	codes := []string{}
	for style, code := range styleMap {
		if s&style > 0 {
			codes = append(codes, strconv.FormatInt(int64(code), 10))
		}
	}
	str = fmt.Sprintf("\033[%sm%s\033[0m", strings.Join(codes, ";"), str)
	return str
}
