package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/layout"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/model/theme"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/sink/pretty/config"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/reflow/indent"
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
}

func newMessageColumn(errorFg string, messageStyles map[level.Level]theme.Style, layout layout.ColumnLayout) messageColumn {
	col := messageColumn{
		baseColumn: baseColumn{
			columnName: schema.Message,
			colLayout:  layout,
			widthMode:  filling,
		},
	}
	for _, lvl := range level.Levels {
		style := messageStyles[lvl]
		cfmt.RegisterStyle(
			col.levelStyleName("", lvl),
			colorFormat(style.Fg, style.Bg),
		)
		cfmt.RegisterStyle(
			col.levelStyleName(highlightSuffix, lvl),
			colorFormat(style.Bg, style.Fg),
		)
		cfmt.RegisterStyle(
			col.levelStyleName(errorSuffix, lvl),
			colorFormat(errorFg, style.Bg),
		)
		cfmt.RegisterStyle(
			col.levelStyleName(sigilSuffix, lvl),
			colorFormat(style.Bg, style.Fg),
		)
	}
	return col
}

// Render ...
func (col messageColumn) Render(cfg config.Config, event model.SinkEvent) string {
	if event.Level == "" {
		event.Level = level.Info
	}
	var expanded = false
	var message = string(event.Message)
	if cfg.Expand != nil && *cfg.Expand {
		if s, ok := expand(message); ok {
			expanded = true
			message = s
		}
	}
	message = col.levelSprintf("", "", event.Level, message)

	var errString = string(event.Error)
	if errString != "" {
		if cfg.Expand != nil && *cfg.Expand {
			if s, ok := expand(errString); ok {
				errString = s
			}
		}
		if message != "" {
			message += "\n"
		}
		message += indent.String(col.levelSprintf("", errorSuffix, event.Level, errString), 4)
	}

	if cfg.Highlight != nil && *cfg.Highlight {
		for _, filter := range event.Filters {
			if filter.Kind.IsHighlighted() {
				message = filter.Pattern.ReplaceAllStringFunc(message, func(s string) string {
					return col.levelSprintf("", highlightSuffix, event.Level, s)
				})
			}
		}
	}

	if (cfg.Multiline != nil && *cfg.Multiline) || expanded {
		message = "\n" + message
	} else {
		message = " " + message
	}

	return cfmt.Sprint(message)
}

//
// Helpers
//

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
