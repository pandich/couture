package column

import (
	"bytes"
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/termenv"
)

const (
	highlightSuffix = "Highlight"
	errorSuffix     = "Error"
)

type messageColumn struct {
	baseColumn
	lexer      chroma.Lexer
	formatter  chroma.Formatter
	bgColorSeq map[level.Level]string
}

func newMessageColumn() messageColumn {
	sigil := '▸'

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
		formatter = formatters.Get("noop")
	}

	return messageColumn{
		baseColumn: baseColumn{
			columnName:  "message",
			widthMode:   filling,
			widthWeight: 0,
			sigil:       &sigil,
		},
		lexer:      chroma.Coalesce(lexers.Get("json")),
		formatter:  formatter,
		bgColorSeq: map[level.Level]string{},
	}
}

// RegisterStyles ...
func (col messageColumn) RegisterStyles(theme theme.Theme) {
	for _, lvl := range level.Levels {
		fg := theme.MessageFg()
		bg := theme.MessageBg(lvl)
		col.bgColorSeq[lvl] = termenv.CSI + termenv.EnvColorProfile().Color(bg).Sequence(true) + "m"
		cfmt.RegisterStyle(
			col.levelStyleName("", lvl),
			colorFormat(fg, bg),
		)
		cfmt.RegisterStyle(
			col.levelStyleName(highlightSuffix, lvl),
			colorFormat(theme.HighlightFg(), bg),
		)
		cfmt.RegisterStyle(
			col.levelStyleName(errorSuffix, lvl),
			colorFormat(theme.StackTraceFg(), bg),
		)
	}
}

// Format ...
func (col messageColumn) Format(_ uint, _ sink.Event) string {
	return "%s"
}

// Render ...
func (col messageColumn) Render(cfg config.Config, event sink.Event) []interface{} {
	lvl := event.Level
	stackTrace := event.StackTrace()
	var message = string(event.Message)
	if cfg.ExpandJSON {
		iterator, err := col.lexer.Tokenise(nil, message)
		if err == nil {
			var buf bytes.Buffer
			err := col.formatter.Format(&buf, cfg.Theme.JSONColorTheme, iterator)
			if err == nil {
				message = buf.String()
				if !cfg.Multiline {
					message = "\n" + message
				}
			}
		}
	}

	var msg = col.levelSprintf(col.prefix(cfg), "", lvl, message)
	msg += col.stackTrace(lvl, stackTrace)

	if cfg.Highlight {
		for _, filter := range event.Filters {
			msg = filter.ReplaceAllStringFunc(msg, func(s string) string {
				return col.levelSprintf("", highlightSuffix, lvl, s) + col.bgColorSeq[lvl]
			})
		}
	}

	return []interface{}{msg}
}

//
// Helpers
//

func (col messageColumn) stackTrace(lvl level.Level, stackTrace *model.StackTrace) string {
	if stackTrace == nil {
		return ""
	}
	return "\n" + indent.String(col.levelSprintf("", errorSuffix, lvl, *stackTrace), 4)
}

func (col messageColumn) prefix(config config.Config) string {
	var prefix string
	if config.Multiline {
		prefix += "\n"
	} else {
		prefix += " "
	}
	if col.sigil != nil {
		prefix += string(*col.sigil) + " "
	}
	return prefix
}

func (col messageColumn) levelSprintf(prefix string, suffix string, lvl level.Level, s interface{}) string {
	return cfmt.Sprintf("{{"+prefix+"%s}}::"+col.levelStyleName(suffix, lvl), s)
}

func (col messageColumn) levelStyleName(suffix string, lvl level.Level) string {
	return col.name() + suffix + string(lvl)
}

func colorFormat(fgColor string, bgColor string) func(s string) string {
	return func(s string) string {
		return cfmt.Sprintf("{{%s}}::"+fgColor+"|bg"+bgColor, s)
	}
}
