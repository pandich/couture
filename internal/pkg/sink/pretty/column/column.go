package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
)

func orNoValue(s string) string {
	const noValue = "—"

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
		RenderFormat(width uint, event model.SinkEvent) string
		RenderValue(cfg config.Config, event model.SinkEvent) []interface{}
		name() string
		mode() widthMode
		weight() widthWeight
	}

	baseColumn struct {
		columnName  string
		widthMode   widthMode
		widthWeight widthWeight
		sigil       *rune
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
	return col.widthWeight
}

type weightedColumn struct {
	baseColumn
	color func(theme.Theme) string
	value func(event model.SinkEvent) []interface{}
}

func newWeightedColumn(
	columnName string,
	sigil *rune,
	widthWeight widthWeight,
	color func(theme.Theme) string,
	value func(event model.SinkEvent) []interface{},
) weightedColumn {
	return weightedColumn{
		baseColumn: baseColumn{
			columnName:  columnName,
			widthMode:   weighted,
			widthWeight: widthWeight,
			sigil:       sigil,
		},
		color: color,
		value: value,
	}
}

// Init ...
func (col weightedColumn) Init(theme theme.Theme) {
	cfmt.RegisterStyle(col.name(), func(s string) string {
		var prefix = ""
		if col.sigil != nil {
			prefix = " " + string(*col.sigil) + " "
		}
		return cfmt.Sprintf("{{"+prefix+"%s }}::"+col.color(theme), s)
	})
}

// RenderFormat ...
func (col weightedColumn) RenderFormat(width uint, _ model.SinkEvent) string {
	return formatColumn(col, width)
}

// RenderValue ...
func (col weightedColumn) RenderValue(_ config.Config, event model.SinkEvent) []interface{} {
	return col.value(event)
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
