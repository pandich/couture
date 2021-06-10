package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/layout"
)

func newContextColumn(style sink.Style, layout layout.ColumnLayout) column {
	return newWeightedColumn(
		schema.Context,
		layout,
		style,
		func(event model.SinkEvent) []interface{} { return []interface{}{string(event.Context)} },
	)
}
