package column

import (
	"github.com/gagglepanda/couture/model"
	"github.com/gagglepanda/couture/schema"
	"github.com/gagglepanda/couture/sink/color"
	"github.com/gagglepanda/couture/sink/layout"
)

func newApplicationColumn(style color.FgBgTuple, layout layout.ColumnLayout) column {
	return newWeightedColumn(
		schema.Application,
		layout,
		style,
		func(event model.SinkEvent) []interface{} { return []interface{}{string(event.Application)} },
	)
}
