package log

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

func NewLineFormatter(layout string, style map[Level]TermStyle) *LineFormatter {
	f := &LineFormatter{layout: layout, style: style}
	f.compile()
	return f
}

type LineFormatter struct {
	layout    string
	isUTC     bool
	positions positions
	style     map[Level]TermStyle
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
	"function": "%s",
	"level":    "%s",
	"message":  "%s",
}

func (f *LineFormatter) Format(e Entry) string {
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
		case "level":
			val = e.Level
		case "line":
			val = e.Line()
		case "file":
			val = e.File()
		case "function":
			val = e.Function()
		case "message":
			val = f.formatMessage(e.Args)
		}
		args[i] = val
	}

	if style, ok := f.style[e.Level]; ok {
		layout = style.Format(layout)
	}

	return fmt.Sprintf(layout, args...) + "\n"
}

func (f *LineFormatter) formatMessage(args []interface{}) string {
	fullContext := Context{}
	for i, arg := range args {
		if context, ok := arg.(Context); ok {
			for key, val := range context {
				fullContext[key] = val
			}
			args = append(args[0:i], args[i+1:]...)
		}
	}
	contextString := ""
	for key, val := range fullContext {
		contextString += " " + key + "=" + fmt.Sprint(val)
	}

	if len(args) > 0 {
		if formatMessage, ok := args[0].(string); ok {
			return fmt.Sprintf(formatMessage, args[1:]...)+contextString
		}
	}
	return fmt.Sprint(args...) + contextString
}

// 2006/01/02 15:04:05.000 level message file/line/function
// layout: 2006/01/02 15:04:05.000 %s %s %s/%d/%s

func (f *LineFormatter) compile() {
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
