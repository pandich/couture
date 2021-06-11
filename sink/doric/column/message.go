package column

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/reflow/indent"
	"github.com/pandich/couture/model"
	"github.com/pandich/couture/model/level"
	"github.com/pandich/couture/schema"
	"github.com/pandich/couture/sink"
	"github.com/pandich/couture/sink/layout"
)

const (
	errorSuffix     = "Error"
	highlightSuffix = "Highlight"
	sigilSuffix     = "Sigil"
)

type (
	highlight     bool
	multiLine     bool
	expand        bool
	messageColumn struct {
		baseColumn
		highlight highlight
		multiLine multiLine
		expand    expand
	}
)

func newMessageColumn(
	highlight highlight,
	expand expand,
	multiLine multiLine,
	errorStyle sink.Style,
	messageStyles map[level.Level]sink.Style,
	layout layout.ColumnLayout,
) column {
	col := messageColumn{
		baseColumn: baseColumn{
			columnName: schema.Message,
			colLayout:  layout,
		},
		highlight: highlight,
		multiLine: multiLine,
		expand:    expand,
	}
	for _, lvl := range level.Levels {
		style := messageStyles[lvl]
		errStyle := sink.Style{
			Fg: errorStyle.Bg,
			Bg: style.Bg,
		}
		cfmt.RegisterStyle(
			col.levelStyleName("", lvl),
			style.Format(),
		)
		cfmt.RegisterStyle(
			col.levelStyleName(highlightSuffix, lvl),
			style.Reverse().Format(),
		)
		cfmt.RegisterStyle(
			col.levelStyleName(errorSuffix, lvl),
			errStyle.Format(),
		)
		cfmt.RegisterStyle(
			col.levelStyleName(sigilSuffix, lvl),
			style.Reverse().Format(),
		)
	}
	return col
}

func (col messageColumn) render(event model.SinkEvent) string {
	var message, expanded = col.renderMessage(event)
	if errorMessage := col.renderErrorMessage(event); errorMessage != "" {
		if message != "" {
			errorMessage = "\n" + errorMessage
		}
		message += errorMessage
	}

	if col.highlight {
		message = event.Filters.ReplaceAllStringFunc(message, func(s string) string {
			return col.levelSprintf("", highlightSuffix, event.Level, s)
		})
	}

	if message != "" {
		if bool(col.multiLine) || expanded {
			message = "\n" + message
		} else {
			message = " " + message
		}
	}

	return cfmt.Sprint(message)
}

func (col messageColumn) renderMessage(event model.SinkEvent) (string, bool) {
	var expanded = false
	var message = string(event.Message)
	if col.expand {
		if s, ok := event.Message.Expand(); ok {
			expanded = true
			message = s
		}
	}
	message = col.levelSprintf("", "", event.Level, message)
	return message, expanded
}

func (col messageColumn) renderErrorMessage(event model.SinkEvent) string {
	if event.Error == "" {
		return ""
	}
	var errString = string(event.Error)
	if col.expand {
		if s, ok := model.Message(event.Error).Expand(); ok {
			errString = s
		}
	}
	errString = col.levelSprintf("", errorSuffix, event.Level, errString)
	errString = indent.String(errString, 4)
	return errString
}

func (col messageColumn) levelSprintf(prefix string, suffix string, lvl level.Level, s interface{}) string {
	return cfmt.Sprintf("{{"+prefix+"%s}}::"+col.levelStyleName(suffix, lvl), s)
}

func (col messageColumn) levelStyleName(suffix string, lvl level.Level) string {
	return col.name() + suffix + string(lvl)
}
