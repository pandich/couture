package column

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/i582/cfmt/cmd/cfmt"
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

var columns = []column{
	newSourceColumn(),
	newTimestampColumn(),
	newApplicationColumn(),
	newThreadColumn(),
	newCallerColumn(),
	newLevelColumn(),
	newMessageColumn(),
}

type (
	column interface {
		RegisterStyles(th theme.Theme)
		Format(width uint, event sink.Event) string
		Render(cfg config.Config, event sink.Event) []interface{}
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
	value func(event sink.Event) []interface{}
}

func newWeightedColumn(
	columnName string,
	sigil *rune,
	widthWeight widthWeight,
	color func(theme.Theme) string,
	value func(event sink.Event) []interface{},
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
func (col weightedColumn) Format(width uint, _ sink.Event) string {
	return formatColumn(col, width)
}

// Render ...
func (col weightedColumn) Render(_ config.Config, event sink.Event) []interface{} {
	return col.value(event)
}
