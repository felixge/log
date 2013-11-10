package log

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

func NewFormat(layout string) *Format {
	f := &Format{layout: layout}
	f.compile()
	return f
}

type Format struct {
	layout    string
	isUTC     bool
	positions positions
}

type positions []position

func (p positions) Len() int {
	return len(p)
}
func (p positions) Less(i, j int) bool {
	return p[i].start < p[j].start
}
func (p positions) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type position struct {
	start int
	token string
}

var tokens = map[string]string{
	"line":     "%d",
	"file":     "%s",
	"level":    "%s",
	"function": "%s",
	"message":  "%s",
}

func (f *Format) Format(e Entry) string {
	layout := f.layout
	if f.isUTC {
		layout = e.Time.UTC().Format(layout)
	} else {
		layout = e.Time.Format(layout)
	}

	args := make([]interface{}, len(f.positions))
	for i, p := range f.positions {
		var val interface{}
		switch p.token {
		case "line":
			val = e.Line
		case "file":
			val = e.File
		case "level":
			val = e.Level
		case "function":
			val = e.Function
		case "message":
			val = e.Message
		}
		args[i] = val
	}

	return fmt.Sprintf(layout, args...)
}

// 2006/01/02 15:04:05.000 level message file/line/function
// layout: 2006/01/02 15:04:05.000 %s %s %s/%d/%s

func (f *Format) compile() {
	for token, _ := range tokens {
		r := regexp.MustCompile(token)
		matches := r.FindAllStringIndex(f.layout, -1)
		for _, match := range matches {
			f.positions = append(f.positions, position{match[0], token})
		}
	}

	for token, placeholder := range tokens {
		r := regexp.MustCompile(token)
		for {
			pos := r.FindStringIndex(f.layout)
			if pos == nil {
				break
			}

			s, e := pos[0], pos[1]
			f.layout = f.layout[0:s] + placeholder + f.layout[e:]
		}
	}

	f.isUTC = strings.Contains(f.layout, "UTC")
	sort.Sort(f.positions)
}
