package column

// TODO get rid of cfmt for this section - this is too complex for it

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
	highlightSuffix = "Highlight"
	errorSuffix     = "Error"
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
func (col messageColumn) Format(_ uint, _ model.SinkEvent) string {
	return "%s"
}

// Render ...
func (col messageColumn) Render(cfg config.Config, event model.SinkEvent) []interface{} {
	var lvl = event.Level
	if lvl == "" {
		lvl = level.Info
	}
	var message = string(event.Message)
	if cfg.ExpandJSON {
		if s, ok := expandJSON(message); ok {
			message = "\n" + s
		}
	}
	message = col.levelSprintf(col.prefix(cfg), "", lvl, message)
	var exception = string(event.Exception)
	if exception != "" && cfg.ExpandJSON {
		if s, ok := expandJSON(exception); ok {
			exception = s
		}
	}
	if exception != "" {
		exception = "\n" + indent.String(
			col.levelSprintf("", errorSuffix, lvl, exception),
			4,
		)
	}
	message += exception

	if cfg.Highlight {
		for _, filter := range event.Filters {
			message = filter.ReplaceAllStringFunc(message, func(s string) string {
				return col.levelSprintf("", highlightSuffix, lvl, s) + col.bgColorSeq[lvl]
			})
		}
	}

	return []interface{}{message}
}

//
// Helpers
//

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

func expandJSON(in string) (string, bool) {
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
