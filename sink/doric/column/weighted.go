package column

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/pandich/couture/model"
	"github.com/pandich/couture/schema"
	"github.com/pandich/couture/sink/layout"
	"github.com/pandich/couture/theme"
)

type extractor func(event model.SinkEvent) []interface{}

type extractorColumn struct {
	baseColumn
	extractor extractor
}

func newWeightedColumn(
	columnName schema.Column,
	layout layout.ColumnLayout,
	style theme.Style,
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
