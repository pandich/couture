package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/termenv"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"strconv"
	"strings"
)

const (
	errorSuffix     = "Error"
	highlightSuffix = "Highlight"
	sigilSuffix     = "Sigil"
)

type messageColumn struct {
	baseColumn
	bgColorSeq map[level.Level]string
}

func newMessageColumn() messageColumn {
	sigil := '▸'

	return messageColumn{
		baseColumn: baseColumn{
			columnName:  "message",
			widthMode:   filling,
			widthWeight: 0,
			sigil:       &sigil,
		},
		bgColorSeq: map[level.Level]string{},
	}
}

// Init ...
func (col messageColumn) Init(theme theme.Theme) {
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
			colorFormat(theme.HighlightFg(lvl), theme.HighlightBg()),
		)
		cfmt.RegisterStyle(
			col.levelStyleName(errorSuffix, lvl),
			colorFormat(theme.StackTraceFg(), bg),
		)
		cfmt.RegisterStyle(
			col.levelStyleName(sigilSuffix, lvl),
			colorFormat(theme.HighlightFg(lvl), theme.HighlightBg()),
		)
	}
}

// RenderFormat ...
func (col messageColumn) RenderFormat(_ uint, _ model.SinkEvent) string {
	return "%s"
}

// RenderValue ...
func (col messageColumn) RenderValue(cfg config.Config, event model.SinkEvent) []interface{} {
	if event.Level == "" {
		event.Level = level.Info
	}
	var expanded = false
	var message = string(event.Message)
	if cfg.Expand {
		if s, ok := expand(message); ok {
			expanded = true
			message = s
		}
	}
	message = col.levelSprintf("", "", event.Level, message)

	var errString = string(event.Error)
	if errString != "" {
		if cfg.Expand {
			if s, ok := expand(errString); ok {
				errString = s
			}
		}
		message += "\n" + indent.String(
			col.levelSprintf("", errorSuffix, event.Level, errString),
			4)
	}

	if cfg.Highlight {
		for _, filter := range event.Filters {
			if filter.Kind.IsHighlighted() {
				message = filter.Pattern.ReplaceAllStringFunc(message, func(s string) string {
					return col.levelSprintf("", highlightSuffix, event.Level, s) + col.bgColorSeq[event.Level]
				})
			}
		}
	}

	message = col.prefix(cfg, event) + message

	if cfg.Multiline || expanded {
		message = "\n" + message
	} else {
		message = " " + message
	}

	return []interface{}{message}
}

//
// Helpers
//

func (col messageColumn) prefix(cfg config.Config, evt model.SinkEvent) string {
	var s = " "
	if col.sigil != nil {
		s += string(*col.sigil) + " "
	}
	if cfg.ShowSchema {
		s += "[" + evt.Schema.Name + "] "
	}
	return " " + col.levelSprintf(s, sigilSuffix, evt.Level, "") + " "
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

func expand(in string) (string, bool) {
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
	in = string(pretty.Pretty([]byte(in)))
	in = strings.TrimRight(in, "\n")
	return in, true
}
