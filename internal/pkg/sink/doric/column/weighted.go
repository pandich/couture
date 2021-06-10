package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/layout"
	"github.com/i582/cfmt/cmd/cfmt"
)

type extractor func(event model.SinkEvent) []interface{}

type extractorColumn struct {
	baseColumn
	extractor extractor
}

func newWeightedColumn(
	columnName string,
	layout layout.ColumnLayout,
	style sink.Style,
	value func(event model.SinkEvent) []interface{},
) extractorColumn {
	col := extractorColumn{
		baseColumn: baseColumn{
			columnName: columnName,
			colLayout:  layout,
		},
		extractor: value,
	}
	registerStyle(col.columnName, style, layout)
	return col
}

func (col extractorColumn) render(event model.SinkEvent) string {
	value := col.extractor(event)
	return cfmt.Sprintf(col.format(), value...)
}
