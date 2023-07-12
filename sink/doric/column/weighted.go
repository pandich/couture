package column

import (
	"github.com/gagglepanda/couture/model"
	"github.com/gagglepanda/couture/schema"
	"github.com/gagglepanda/couture/sink/color"
	"github.com/gagglepanda/couture/sink/layout"
	"github.com/i582/cfmt/cmd/cfmt"
)

type extractor func(event model.SinkEvent) []interface{}

type extractorColumn struct {
	baseColumn
	extractor extractor
}

func newWeightedColumn(
	columnName schema.Column,
	layout layout.ColumnLayout,
	style color.FgBgTuple,
	value func(event model.SinkEvent) []interface{},
) extractorColumn {
	col := extractorColumn{
		baseColumn: baseColumn{
			columnName: columnName,
			colLayout:  layout,
		},
		extractor: value,
	}
	registerStyle(string(col.columnName), style, layout)
	return col
}

func (col extractorColumn) render(event model.SinkEvent) string {
	value := col.extractor(event)
	return cfmt.Sprintf(col.format(), value...)
}
