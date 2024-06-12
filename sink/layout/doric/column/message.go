package column

import (
	"github.com/pandich/couture/event"
	"github.com/pandich/couture/event/level"
	"github.com/pandich/couture/mapping"
	"github.com/pandich/couture/sink/color"
	"github.com/pandich/couture/sink/layout"
	theme2 "github.com/pandich/couture/sink/theme"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/reflow/indent"
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
	levelMeter    bool
	expand        bool
	useColor      bool
	messageColumn struct {
		baseColumn
		highlight       highlight
		multiLine       multiLine
		levelMeter      levelMeter
		expand          expand
		prettyJSONStyle *pretty.Style
		useColor        useColor
	}
)

var levelMeterSigils = map[event.Bucket]string{
	event.Bucket1: " ",
	event.Bucket2: "▏",
	event.Bucket3: "▎",
	event.Bucket4: "▍",
	event.Bucket5: "▌",
	event.Bucket6: "▋",
	event.Bucket7: "▊",
	event.Bucket8: "▉",
}

func newMessageColumn(
	highlight highlight,
	expand expand,
	multiLine multiLine,
	levelMeter levelMeter,
	useColor useColor,
	th *theme2.Theme,
	layout layout.ColumnLayout,
) column {
	col := messageColumn{
		baseColumn:      baseColumn{columnName: mapping.Message, colLayout: layout},
		highlight:       highlight,
		levelMeter:      levelMeter,
		multiLine:       multiLine,
		expand:          expand,
		useColor:        useColor,
		prettyJSONStyle: th.AsPrettyJSONStyle(),
	}
	for _, lvl := range level.Levels {
		style := th.Message[lvl]
		errStyle := color.FgBgTuple{
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

func (col messageColumn) render(event event.SinkEvent) string {
	var message, expanded = col.renderMessage(event)
	if errorMessage := col.renderErrorMessage(event); errorMessage != "" {
		if message != "" {
			errorMessage = "\n" + errorMessage
		}
		message += errorMessage
	}

	if col.highlight {
		message = event.Filters.ReplaceAllStringFunc(
			message, func(s string) string {
				return col.levelSprintf(highlightSuffix, event.Level, s)
			},
		)
	}

	if col.levelMeter {
		message = col.renderLevelMeter(event) + message
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

func (col messageColumn) renderLevelMeter(evt event.SinkEvent) string {
	var bucketSigil, ok = levelMeterSigils[evt.LevelMeterBucket()]
	if !ok {
		bucketSigil = levelMeterSigils[event.BucketUpperBound]
	}
	return col.levelSprintf("", evt.Level, bucketSigil) + " "
}

func (col messageColumn) renderMessage(event event.SinkEvent) (string, bool) {
	var expanded = false
	var message = string(event.Message)
	if col.expand {
		if s, ok := col.expandMessage(event.Message); ok {
			expanded = true
			message = s
		}
	}
	message = col.levelSprintf("", event.Level, message)
	return message, expanded
}

func (col messageColumn) renderErrorMessage(evt event.SinkEvent) string {
	if evt.Error == "" {
		return ""
	}
	var errString = string(evt.Error)
	if col.expand {
		if s, ok := col.expandMessage(event.Message(evt.Error)); ok {
			errString = s
		}
	}
	errString = col.levelSprintf(errorSuffix, evt.Level, errString)
	errString = indent.String(errString, 4)
	return errString
}

func (col messageColumn) levelSprintf(suffix string, lvl level.Level, s interface{}) string {
	return cfmt.Sprintf("{{%s}}::"+col.levelStyleName(suffix, lvl), s)
}

func (col messageColumn) levelStyleName(suffix string, lvl level.Level) string {
	return string(col.name()) + suffix + string(lvl)
}

// expandMessage ...
func (col messageColumn) expandMessage(msg event.Message) (string, bool) {
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
