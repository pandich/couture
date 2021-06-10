package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/layout"
)

func newApplicationColumn(style sink.Style, layout layout.ColumnLayout) column {
	return newWeightedColumn(
		schema.Application,
		layout,
		style,
		func(event model.SinkEvent) []interface{} { return []interface{}{string(event.Application)} },
	)
}
