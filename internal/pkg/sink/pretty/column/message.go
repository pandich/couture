package column

import (
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"github.com/i582/cfmt/cmd/cfmt"
)

// TODO cleanup all the messy string handling

type messageColumn struct {
	baseColumn
}

func newMessageColumn() messageColumn {
	sigil := '¶'
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
		fgColor := theme.MessageColor()
		bgColor := theme.MessageBackgroundColor(lvl)
		cfmt.RegisterStyle(col.name()+string(lvl), func(s string) string {
			return cfmt.Sprintf("{{%s}}::"+fgColor+"|bg"+bgColor, s)
		})
	}
	for _, lvl := range level.Levels {
		fgColor := theme.MessageColor()
		bgColor := theme.HighlightBackgroundColor(lvl)
		cfmt.RegisterStyle(col.name()+"Highlight"+string(lvl), func(s string) string {
			return cfmt.Sprintf("{{%s}}::"+fgColor+"|bg"+bgColor, s)
		})
	}
	for _, lvl := range level.Levels {
		fgColor := theme.StackTraceColor()
		bgColor := theme.MessageBackgroundColor(lvl)
		cfmt.RegisterStyle(col.name()+"Error"+string(lvl), func(s string) string {
			return cfmt.Sprintf("{{%s}}::"+fgColor+"|bg"+bgColor, s)
		})
	}
}

// Format ...
func (col messageColumn) Format(_ uint, _ source.Source, _ sink.Event) string {
	return "%s"
}

// Render ...
func (col messageColumn) Render(config config.Config, _ source.Source, event sink.Event) []interface{} {
	var prefix string
	if config.Multiline {
		prefix += "\n"
	} else {
		prefix += " "
	}
	if col.sigil != nil {
		prefix += string(*col.sigil) + " "
	}

	lvl := event.Event.Level
	var formattedMessage = cfmt.Sprintf("{{"+prefix+"%s}}::"+col.name()+string(lvl), event.Event.Message)
	stackTrace := event.Event.StackTrace()
	if stackTrace != nil {
		formattedMessage += cfmt.Sprintf("\n{{"+prefix+"%s}}::"+col.name()+"Error"+string(lvl), *stackTrace)
	}

	for _, filter := range event.Filters {
		formattedMessage = filter.ReplaceAllStringFunc(formattedMessage, func(s string) string {
			return cfmt.Sprintf("{{%s}}::"+col.name()+"Highlight"+string(lvl), s)
		})
	}

	return []interface{}{formattedMessage}
}
