package column

import (
	"github.com/gagglepanda/couture/event"
	"github.com/gagglepanda/couture/mapping"
	"github.com/gagglepanda/couture/sink/color"
	"github.com/gagglepanda/couture/sink/layout"
)

func newApplicationColumn(style color.FgBgTuple, layout layout.ColumnLayout) column {
	return newWeightedColumn(
		mapping.Application,
		layout,
		style,
		func(event event.SinkEvent) []interface{} { return []interface{}{string(event.Application)} },
	)
}
