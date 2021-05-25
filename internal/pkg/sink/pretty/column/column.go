package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"github.com/i582/cfmt/cmd/cfmt"
)

// DefaultColumns ...
var DefaultColumns = []string{
	"timestamp",
	"application",
	"thread",
	"caller",
	"level",
	"message",
	"error",
}

var columns = []column{
	newSourceColumn(),
	newTimestampColumn(),
	newApplicationColumn(),
	newThreadColumn(),
	newCallerColumn(),
	newLevelColumn(),
	newMessageColumn(),
	newStackTraceColumn(),
}

type (
	column interface {
		RegisterStyles(th theme.Theme)
		Format(width uint, src source.Source, event model.Event) string
		Render(cfg config.Config, src source.Source, event model.Event) []interface{}
		name() string
		mode() widthMode
		weight() widthWeight
	}

	baseColumn struct {
		columnName  string
		weightType  widthMode
		widthWeight widthWeight
	}
)

// Name ...
func (col baseColumn) name() string {
	return col.columnName
}

// WeightType ...
func (col baseColumn) mode() widthMode {
	return col.weightType
}

// Weight ...
func (col baseColumn) weight() widthWeight {
	return col.widthWeight
}

type weightedColumn struct {
	baseColumn
	sigil *rune
	color func(theme.Theme) string
	value func(model.Event) []interface{}
}

func newWeightedColumn(
	columnName string,
	sigil *rune,
	widthWeight widthWeight,
	color func(theme.Theme) string,
	value func(model.Event) []interface{},
) weightedColumn {
	return weightedColumn{
		baseColumn: baseColumn{
			columnName:  columnName,
			weightType:  weighted,
			widthWeight: widthWeight,
		},
		sigil: sigil,
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
func (col weightedColumn) Format(width uint, _ source.Source, _ model.Event) string {
	return formatColumn(col, width)
}

// Render ...
func (col weightedColumn) Render(_ config.Config, _ source.Source, event model.Event) []interface{} {
	return col.value(event)
}
