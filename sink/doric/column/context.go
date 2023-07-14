package column

import (
	"github.com/gagglepanda/couture/event"
	"github.com/gagglepanda/couture/schema"
	"github.com/gagglepanda/couture/sink/color"
	"github.com/gagglepanda/couture/sink/layout"
)

func newContextColumn(style color.FgBgTuple, layout layout.ColumnLayout) column {
	return newWeightedColumn(
		schema.Context,
		layout,
		style,
		func(event event.SinkEvent) []interface{} { return []interface{}{string(event.Context)} },
	)
}
