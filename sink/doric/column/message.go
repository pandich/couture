package column

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/reflow/indent"
	"github.com/pandich/couture/model"
	"github.com/pandich/couture/model/level"
	"github.com/pandich/couture/schema"
	"github.com/pandich/couture/sink/layout"
	"github.com/pandich/couture/theme"
	"github.com/pandich/couture/theme/color"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"strconv"
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
	useColor      bool
	messageColumn struct {
		baseColumn
		highlight       highlight
		multiLine       multiLine
		expand          expand
		prettyJSONStyle *pretty.Style
		useColor        useColor
	}
)

func newMessageColumn(
	highlight highlight,
	expand expand,
	multiLine multiLine,
	useColor useColor,
	th *theme.Theme,
	layout layout.ColumnLayout,
) column {
	col := messageColumn{
		baseColumn:      baseColumn{columnName: schema.Message, colLayout: layout},
		highlight:       highlight,
		multiLine:       multiLine,
		expand:          expand,
		useColor:        useColor,
		prettyJSONStyle: th.AsPrettyJSONStyle(),
	}
	for _, lvl := range level.Levels {
		style := th.Message[lvl]
		errStyle := color.HexPair{
			Fg: th.Level[level.Error].Bg,
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
		if s, ok := col.expandMessage(event.Message); ok {
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
		if s, ok := col.expandMessage(model.Message(event.Error)); ok {
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
	return string(col.name()) + suffix + string(lvl)
}

// expandMessage ...
func (col messageColumn) expandMessage(msg model.Message) (string, bool) {
	var in = string(msg)
	if in == "" {
		return in, false
	}
	if in[0] == '"' {
		s, err := strconv.Unquote(in)
		if err != nil {
			return in, false
		}
		in = s
	}
	if !gjson.Valid(in) {
		return in, false
	}
	var out = pretty.Pretty([]byte(in))
	if col.useColor {
		out = pretty.Color(out, col.prettyJSONStyle)
	}
	return string(out), true
}
