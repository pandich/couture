package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

// DefaultColumns ...
var DefaultColumns = []string{
	"source",
	"timestamp",
	"application",
	"thread",
	"caller",
	"level",
	"message",
}

type (
	column interface {
		RegisterStyles(th theme.Theme)
		Format(width uint, event model.SinkEvent) string
		Render(cfg config.Config, event model.SinkEvent) []interface{}
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

// Name ...
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

// RegisterStyles ...
func (col weightedColumn) RegisterStyles(theme theme.Theme) {
	cfmt.RegisterStyle(col.name(), func(s string) string {
		var prefix = ""
		if col.sigil != nil {
			prefix = " " + string(*col.sigil) + " "
		}
		return cfmt.Sprintf("{{"+prefix+"%s }}::"+col.color(theme), s)
	})
}

// Format ...
func (col weightedColumn) Format(width uint, _ model.SinkEvent) string {
	return formatColumn(col, width)
}

// Render ...
func (col weightedColumn) Render(_ config.Config, event model.SinkEvent) []interface{} {
	return col.value(event)
}

func contrast(hex string) string {
	cf, _ := colorful.MakeColor(gamut.Contrast(gamut.Hex(hex)))
	return cf.Hex()
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
