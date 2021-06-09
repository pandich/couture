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

func newCallerColumn(cfg config.Config) column {
	layout := cfg.Layout.Caller
	return callerColumn{baseColumn{
		columnName: "caller",
		widthMode:  fixed,
		colLayout:  layout,
	}}
}

// Init ...
func (col callerColumn) Init(theme theme.Theme) {
	var prefix = ""
	if col.colLayout.Sigil != "" {
		prefix = " " + col.colLayout.Sigil + " "
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

// Render ...
func (col callerColumn) Render(_ config.Config, event model.SinkEvent) string {
	const delimiterCharacterCount = 4
	maxWidth := int(col.layout().Width) - delimiterCharacterCount

	var format = "{{%s}}::Entity"
	if event.Action != "" {
		format += "{{∕}}::ActionDelimiter"
	}
	format += "{{%s}}::Action"
	if event.Line != 0 {
		format += "{{#}}::LineNumberDelimiter"
	}
	format += "{{%s}}::Line"

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
		totalLength -= l
	}
	if l := len(actionName) - overage; overage > 0 && l >= 0 {
		actionName = actionName[l:]
		totalLength -= l
	}
	for i := maxWidth - totalLength; i > 0; i-- {
		entityName = " " + entityName
	}

	return cfmt.Sprintf(format, entityName, actionName, lineNumber)
}
