package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
)

type callerColumn struct {
	baseColumn
}

func newCallerColumn() column {
	const width = 50
	sigil := '☎'
	return callerColumn{baseColumn{
		columnName:  "caller",
		widthMode:   fixed,
		widthWeight: width,
		sigil:       &sigil,
	}}
}

// RegisterStyles ...
func (col callerColumn) RegisterStyles(theme theme.Theme) {
	var prefix = ""
	if col.sigil != nil {
		prefix = " " + string(*col.sigil) + " "
	}

	bgColor := theme.CallerBg()

	cfmt.RegisterStyle("Entity", func(s string) string {
		return cfmt.Sprintf("{{"+prefix+"︎%s}}::bg"+bgColor+"|"+theme.EntityFg(), s)
	})
	cfmt.RegisterStyle("ActionDelimiter", func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+bgColor+"|"+theme.ActionDelimiterFg(), s)
	})
	cfmt.RegisterStyle("Action", func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+bgColor+"|"+theme.ActionFg(), s)
	})
	cfmt.RegisterStyle("LineNumberDelimiter", func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+bgColor+"|"+theme.LineNumberDelimiterFg(), s)
	})
	cfmt.RegisterStyle("Line", func(s string) string {
		return cfmt.Sprintf("{{%s }}::bg"+bgColor+"|"+theme.LineNumberFg(), s)
	})
}

// Format ...
func (col callerColumn) Format(_ uint, evt model.SinkEvent) string {
	var s = "{{%s}}::Entity"
	if evt.Action != "" {
		s += "{{∕}}::ActionDelimiter"
	}
	s += "{{%s}}::Action"
	if evt.Line != 0 {
		s += "{{#}}::LineNumberDelimiter"
	}
	s += "{{%s}}::Line"
	return s
}

// Render ...
func (col callerColumn) Render(_ config.Config, event model.SinkEvent) []interface{} {
	const maxEntityNameWidth = 30
	const maxWidth = 60

	var padding = ""
	var entityName = orNoValue(string(event.Entity.Abbreviate(maxEntityNameWidth)))
	var actionName = string(event.Action)
	var lineNumber = ""
	if event.Line != 0 {
		lineNumber = fmt.Sprintf("%4d", event.Line)
	}
	totalLength := len(entityName) + len(actionName) + len(lineNumber)
	for i := totalLength; i < maxWidth; i++ {
		padding += " "
	}
	extraChars := totalLength - maxWidth
	if extraChars > 0 {
		actionName = actionName[:len(actionName)-extraChars-1]
	}

	return []interface{}{
		padding + entityName,
		actionName,
		lineNumber,
	}
}
