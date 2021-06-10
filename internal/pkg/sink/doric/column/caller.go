package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/layout"
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
	entityStyle sink.Style,
	actionDelimiterStyle sink.Style,
	actionStyle sink.Style,
	lineDelimiterStyle sink.Style,
	lineStyle sink.Style,
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

func (col callerColumn) render(event model.SinkEvent) string {
	entityName, actionName, lineNumber := col.callerParts(event)

	var format = "{{%s}}::" + entityStyleName
	if event.Action != "" {
		format += "{{∕}}::" + actionDelimiterStyleName
	}
	format += "{{%s}}::" + actionStyleName
	if event.Line != model.NoLineNumber {
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