package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/layout"
	"couture/internal/pkg/model/theme"
	"couture/internal/pkg/sink/pretty/config"
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
)

func orNoValue(s string) string {
	const noValue = ""

	if s == "" {
		return noValue
	}
	return s
}

// DefaultColumns ...
var DefaultColumns = []string{
	"source",
	"timestamp",
	"application",
	"context",
	"caller",
	"level",
	"message",
}

type (
	column interface {
		Init(th theme.Theme)
		Render(cfg config.Config, event model.SinkEvent) string
		name() string
		mode() widthMode
		layout() layout.ColumnLayout
		format() string
	}

	baseColumn struct {
		columnName string
		widthMode  widthMode
		colLayout  layout.ColumnLayout
	}
)

// GetName ...
func (col baseColumn) name() string {
	return col.columnName
}

// WeightType ...
func (col baseColumn) mode() widthMode {
	return col.widthMode
}

// Weight ...
func (col baseColumn) weight() widthWeight {
	return widthWeight(col.colLayout.Width)
}

func (col baseColumn) format() string {
	return fmt.Sprintf(
		"%%%[1]d.%[1]ds%%%[2]d.%[2]ds%%%[3]d.%[3]ds",
		col.colLayout.Padding.Left,
		col.colLayout.Width,
		col.colLayout.Padding.Right,
	)
}

func (col baseColumn) layout() layout.ColumnLayout {
	return col.colLayout
}

type weightedColumn struct {
	baseColumn
	color func(theme.Theme) string
	value func(event model.SinkEvent) []interface{}
}

func newWeightedColumn(
	columnName string,
	layout layout.ColumnLayout,
	color func(theme.Theme) string,
	value func(event model.SinkEvent) []interface{},
) weightedColumn {
	return weightedColumn{
		baseColumn: baseColumn{
			columnName: columnName,
			widthMode:  weighted,
			colLayout:  layout,
		},
		color: color,
		value: value,
	}
}

// Render ...
func (col weightedColumn) Render(_ config.Config, event model.SinkEvent) string {
	return cfmt.Sprintf(formatColumn(col, col.layout().Width), col.value(event)...)
}

// Init ...
func (col weightedColumn) Init(theme theme.Theme) {
	cfmt.RegisterStyle(col.name(), func(s string) string {
		var prefix = ""
		if col.colLayout.Sigil != "" {
			prefix = " " + col.colLayout.Sigil + " "
		}
		return cfmt.Sprintf("{{"+prefix+"%s }}::"+col.color(theme), s)
	})
}

func formatStringOfWidth(width uint) string {
	if width <= 0 {
		return "%s"
	}
	return fmt.Sprintf("%%-%[1]d.%[1]ds", width)
}

func formatStyleOfWidth(style string, width uint) string {
	return "{{" + formatStringOfWidth(width) + "}}::" + style
}

func formatColumn(col column, width uint) string {
	return formatStyleOfWidth(col.name(), width)
}
