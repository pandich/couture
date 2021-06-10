package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/schema"
	layout2 "couture/internal/pkg/sink/layout"
	theme2 "couture/internal/pkg/sink/theme"
)

func newApplicationColumn(style theme2.Style, layout layout2.ColumnLayout) column {
	return newWeightedColumn(
		schema.Application,
		layout,
		style,
		func(event model.SinkEvent) []interface{} { return []interface{}{string(event.Application)} },
	)
}
