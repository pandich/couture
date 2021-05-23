package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
)

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
