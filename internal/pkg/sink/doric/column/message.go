package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/sink"
	layout2 "couture/internal/pkg/sink/layout"
	theme2 "couture/internal/pkg/sink/theme"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/reflow/indent"
)

const (
	errorSuffix     = "Error"
	highlightSuffix = "Highlight"
	sigilSuffix     = "Sigil"
)

type messageColumn struct {
	baseColumn
	highlight bool
	multiLine bool
	expand    bool
}

func newMessageColumn(
	highlight *bool,
	expand *bool,
	multiLine *bool,
	errorFg string,
	messageStyles map[level.Level]theme2.Style,
	layout layout2.ColumnLayout,
) column {
	col := messageColumn{
		baseColumn: baseColumn{
			columnName: schema.Message,
			colLayout:  layout,
		},
		highlight: highlight != nil && *highlight,
		multiLine: multiLine != nil && *multiLine,
		expand:    expand != nil && *expand,
	}
	for _, lvl := range level.Levels {
		style := messageStyles[lvl]
		errStyle := theme2.Style{
			Fg: errorFg,
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
	message = col.highlightMessage(event, message)

	if message != "" {
		if col.multiLine || expanded {
			message = "\n" + message
		} else {
			message = " " + message
		}
	}

	return cfmt.Sprint(message)
}

func (col messageColumn) renderErrorMessage(event model.SinkEvent) string {
	var errString = string(event.Error)
	if errString != "" {
		if col.expand {
			if s, ok := sink.ExpandText(errString); ok {
				errString = s
			}
		}
		errString = col.levelSprintf("", errorSuffix, event.Level, errString)
		errString = indent.String(errString, 4)
	}
	return errString
}

func (col messageColumn) renderMessage(event model.SinkEvent) (string, bool) {
	var expanded = false
	var message = string(event.Message)
	if col.expand {
		if s, ok := sink.ExpandText(message); ok {
			expanded = true
			message = s
		}
	}
	message = col.levelSprintf("", "", event.Level, message)
	return message, expanded
}

func (col messageColumn) highlightMessage(event model.SinkEvent, message string) string {
	if col.highlight {
		for _, filter := range event.Filters {
			if filter.Kind.IsHighlighted() {
				message = filter.Pattern.ReplaceAllStringFunc(message, func(s string) string {
					return col.levelSprintf("", highlightSuffix, event.Level, s)
				})
			}
		}
	}
	return message
}

func (col messageColumn) levelSprintf(prefix string, suffix string, lvl level.Level, s interface{}) string {
	return cfmt.Sprintf("{{"+prefix+"%s}}::"+col.levelStyleName(suffix, lvl), s)
}

func (col messageColumn) levelStyleName(suffix string, lvl level.Level) string {
	return col.name() + suffix + string(lvl)
}
