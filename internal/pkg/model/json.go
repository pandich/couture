package model

import (
	"bytes"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/muesli/termenv"
)

// NewChromaLexer ...
func NewChromaLexer(name string) chroma.Lexer {
	return chroma.Coalesce(lexers.Get(name))
}

// NewChromaFormatter ...
func NewChromaFormatter() chroma.Formatter {
	var formatter chroma.Formatter
	switch termenv.EnvColorProfile() {
	case termenv.ANSI:
		formatter = formatters.Get("terminal8")
	case termenv.ANSI256:
		formatter = formatters.Get("terminal256")
	case termenv.TrueColor:
		formatter = formatters.Get("terminal16m")
	case termenv.Ascii:
		fallthrough
	default:
		formatter = formatters.Fallback
	}
	return formatter
}

var jsonLexer = NewChromaLexer("json")

// PrettyJSON ...
func PrettyJSON(s string) string {
	iterator, err := jsonLexer.Tokenise(nil, s)
	if err == nil {
		var buf bytes.Buffer
		err := formatters.TTY.Format(&buf, styles.BlackWhite, iterator)
		if err == nil {
			s = buf.String()
		}
	}
	return s
}
