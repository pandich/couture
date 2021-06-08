package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/theme"
	"couture/internal/pkg/sink/pretty/config"
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

// Init ...
func (col callerColumn) Init(theme theme.Theme) {
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

// RenderFormat ...
func (col callerColumn) RenderFormat(_ uint, evt model.SinkEvent) string {
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

// RenderValue ...
func (col callerColumn) RenderValue(_ config.Config, event model.SinkEvent) []interface{} {
	const maxWidth = 54

	var entityName = orNoValue(string(event.Entity.Abbreviate(maxWidth)))
	var actionName = string(event.Action)
	var lineNumber = ""
	if event.Line != 0 {
		lineNumber = fmt.Sprintf("%4d", event.Line)
	}
	var totalLength = len(entityName) + len(actionName) + len(lineNumber)

	// pad
	for i := totalLength; i < maxWidth; i++ {
		entityName = " " + entityName
		totalLength++
	}

	// trim
	var overage = totalLength - maxWidth
	if l := len(entityName) - overage; overage > 0 && l >= 0 {
		entityName = entityName[len(entityName)-l:]
		overage -= l
	}
	if l := len(actionName) - overage; overage > 0 && l >= 0 {
		actionName = actionName[l:]
	}

	// return
	return []interface{}{entityName, actionName, lineNumber}
}
