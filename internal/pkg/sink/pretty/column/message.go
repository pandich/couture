package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/reflow/indent"
)

const (
	highlightSuffix = "Highlight"
	errorSuffix     = "Error"
)

type messageColumn struct {
	baseColumn
}

func newMessageColumn() messageColumn {
	sigil := '▸'
	return messageColumn{baseColumn{
		columnName:  "message",
		widthMode:   filling,
		widthWeight: 0,
		sigil:       &sigil,
	}}
}

// RegisterStyles ...
func (col messageColumn) RegisterStyles(theme theme.Theme) {
	for _, lvl := range level.Levels {
		cfmt.RegisterStyle(
			col.levelStyleName("", lvl),
			colorFormat(theme.MessageFg(), theme.MessageBg(lvl)),
		)
	}
	for _, lvl := range level.Levels {
		cfmt.RegisterStyle(
			col.levelStyleName(highlightSuffix, lvl),
			colorFormat(theme.MessageFg(), theme.HighlightBg(lvl)),
		)
	}
	for _, lvl := range level.Levels {
		cfmt.RegisterStyle(
			col.levelStyleName(errorSuffix, lvl),
			colorFormat(theme.StackTraceFg(), theme.MessageBg(lvl)),
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

	var msg = col.levelSprintf(col.prefix(cfg), "", lvl, event.Message)
	msg += col.stackTrace(lvl, stackTrace)

	if cfg.Highlight {
		for _, filter := range event.Filters {
			msg = filter.ReplaceAllStringFunc(msg, func(s string) string {
				return col.levelSprintf("", highlightSuffix, lvl, s)
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
