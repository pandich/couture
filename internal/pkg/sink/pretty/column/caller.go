package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/layout"
	"couture/internal/pkg/model/theme"
	"couture/internal/pkg/sink/pretty/config"
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
)

type callerColumn struct {
	baseColumn
}

func newCallerColumn(
	entityStyle theme.Style,
	actionDelimiterStyle theme.Style,
	actionStyle theme.Style,
	lineDelimiterStyle theme.Style,
	lineStyle theme.Style,
	layout layout.ColumnLayout,
) column {
	col := callerColumn{
		baseColumn: baseColumn{
			columnName: "caller",
			widthMode:  fixed,
			colLayout:  layout,
		},
	}

	var prefix = ""
	if col.colLayout.Sigil != "" {
		prefix = " " + col.colLayout.Sigil + " "
	}

	cfmt.RegisterStyle("Entity", func(s string) string {
		return cfmt.Sprintf("{{"+prefix+"︎%s}}::bg"+entityStyle.Bg+"|"+entityStyle.Fg, s)
	})
	cfmt.RegisterStyle("ActionDelimiter", func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+actionDelimiterStyle.Bg+"|"+actionDelimiterStyle.Fg, s)
	})
	cfmt.RegisterStyle("Action", func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+actionStyle.Bg+"|"+actionStyle.Fg, s)
	})
	cfmt.RegisterStyle("LineDelimiter", func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+lineDelimiterStyle.Bg+"|"+lineDelimiterStyle.Fg, s)
	})
	cfmt.RegisterStyle("Line", func(s string) string {
		return cfmt.Sprintf("{{%s }}::bg"+lineStyle.Bg+"|"+lineStyle.Fg, s)
	})

	return col
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
		format += "{{#}}::LineDelimiter"
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
	underage := int(col.layout().Width) - totalLength
	for i := underage; i > 0; i-- {
		switch col.layout().Align {
		case layout.AlignRight:
			entityName = " " + entityName
		case layout.AlignLeft:
			fallthrough
		default:
			entityName += " "
		}
	}

	return cfmt.Sprintf(format, entityName, actionName, lineNumber)
}
