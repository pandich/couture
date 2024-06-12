package column

import (
	"github.com/pandich/couture/event"
	"github.com/pandich/couture/mapping"
	"github.com/pandich/couture/sink/color"
	"github.com/pandich/couture/sink/layout"
	"github.com/i582/cfmt/cmd/cfmt"
)

type extractor func(event event.SinkEvent) []interface{}

type extractorColumn struct {
	baseColumn
	extractor extractor
}

func newWeightedColumn(
	columnName mapping.Column,
	layout layout.ColumnLayout,
	style color.FgBgTuple,
	value func(event event.SinkEvent) []interface{},
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

func (col extractorColumn) render(event event.SinkEvent) string {
	value := col.extractor(event)
	return cfmt.Sprintf(col.format(), value...)
}
