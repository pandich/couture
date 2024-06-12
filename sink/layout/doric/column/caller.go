package column

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/pandich/couture/event"
	"github.com/pandich/couture/sink/color"
	"github.com/pandich/couture/sink/layout"
	"github.com/i582/cfmt/cmd/cfmt"
)

const (
	callerPsuedoColumn       = "caller"
	entityStyleName          = "Entity"
	actionDelimiterStyleName = "ActionDelimiter"
	actionStyleName          = "Action"
	lineDelimiterStyleName   = "LineDelimiter"
	lineStyleName            = "Line"
)

type callerColumn struct {
	baseColumn
}

func newCallerColumn(
	entityStyle color.FgBgTuple,
	actionDelimiterStyle color.FgBgTuple,
	actionStyle color.FgBgTuple,
	lineDelimiterStyle color.FgBgTuple,
	lineStyle color.FgBgTuple,
	entityLayout layout.ColumnLayout,
) column {
	col := callerColumn{baseColumn{columnName: callerPsuedoColumn, colLayout: entityLayout}}

	linePadding := entityLayout.EffectivePadding()
	linePadding.Left = layout.NoPadding.Left
	lineLayout := layout.ColumnLayout{Padding: &linePadding}

	entityPadding := entityLayout.EffectivePadding()
	entityLayout.Padding = &entityPadding
	entityLayout.Padding.Right = layout.NoPadding.Right

	registerStyle(entityStyleName, entityStyle, entityLayout)
	registerStyle(actionDelimiterStyleName, actionDelimiterStyle, layout.NoPaddingLayout)
	registerStyle(actionStyleName, actionStyle, layout.NoPaddingLayout)
	registerStyle(lineDelimiterStyleName, lineDelimiterStyle, layout.NoPaddingLayout)
	registerStyle(lineStyleName, lineStyle, lineLayout)
	return col
}

func (col callerColumn) render(evt event.SinkEvent) string {
	entityName, actionName, lineNumber := col.callerParts(evt)

	var format = "{{%s}}::" + entityStyleName
	if evt.Action != "" {
		format += "{{∕}}::" + actionDelimiterStyleName
	}
	format += "{{%s}}::" + actionStyleName
	if evt.Line != event.NoLineNumber {
		format += "{{#}}::" + lineDelimiterStyleName
	}
	format += "{{%s}}::" + lineStyleName

	return cfmt.Sprintf(
		format,
		col.entityPartStyle(entityName, actionName, lineNumber).Render(entityName),
		actionName,
		lineNumber,
	)
}

func (col callerColumn) callerParts(event event.SinkEvent) (string, string, string) {
	var entityName = string(event.Entity.Abbreviate(int(col.colLayout.Width)))
	var actionName = string(event.Action)
	var lineNumber = ""
	if event.Line != 0 {
		lineNumber = fmt.Sprintf("%4d", event.Line)
	}
	return entityName, actionName, lineNumber
}

func (col callerColumn) entityPartStyle(entityName string, actionName string, lineNumber string) lipgloss.Style {
	const delimiterWidth = 1
	const sigilWidth = 2
	const minEntityWidth = 10

	totalWidth := sigilWidth + len(entityName) + delimiterWidth + len(actionName) + delimiterWidth + len(lineNumber)
	var entityWidth = int(col.colLayout.Width) - totalWidth + len(entityName)
	if entityWidth < minEntityWidth {
		entityWidth = minEntityWidth
	}
	return lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(entityWidth).
		MaxWidth(entityWidth)
}
