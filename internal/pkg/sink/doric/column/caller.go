package column

import (
	"couture/internal/pkg/model"
	layout2 "couture/internal/pkg/sink/layout"
	theme2 "couture/internal/pkg/sink/theme"
	"fmt"
	"github.com/charmbracelet/lipgloss"
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
	entityStyle theme2.Style,
	actionDelimiterStyle theme2.Style,
	actionStyle theme2.Style,
	lineDelimiterStyle theme2.Style,
	lineStyle theme2.Style,
	entityLayout layout2.ColumnLayout,
) column {
	col := callerColumn{baseColumn{columnName: callerPsuedoColumn, colLayout: entityLayout}}

	linePadding := entityLayout.EffectivePadding()
	linePadding.Left = layout2.NoPadding.Left
	lineLayout := layout2.ColumnLayout{Padding: &linePadding}

	entityPadding := entityLayout.EffectivePadding()
	entityLayout.Padding = &entityPadding
	entityLayout.Padding.Right = layout2.NoPadding.Right

	registerStyle(entityStyleName, entityStyle, entityLayout)
	registerStyle(actionDelimiterStyleName, actionDelimiterStyle, layout2.NoPaddingLayout)
	registerStyle(actionStyleName, actionStyle, layout2.NoPaddingLayout)
	registerStyle(lineDelimiterStyleName, lineDelimiterStyle, layout2.NoPaddingLayout)
	registerStyle(lineStyleName, lineStyle, lineLayout)
	return col
}

func (col callerColumn) render(event model.SinkEvent) string {
	entityName, actionName, lineNumber := col.callerParts(event)

	var format = "{{%s}}::" + entityStyleName
	if event.Action != "" {
		format += "{{∕}}::" + actionDelimiterStyleName
	}
	format += "{{%s}}::" + actionStyleName
	if event.Line != 0 {
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

func (col callerColumn) callerParts(event model.SinkEvent) (string, string, string) {
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
