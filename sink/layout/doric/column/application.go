package column

import (
	"github.com/pandich/couture/event"
	"github.com/pandich/couture/mapping"
	"github.com/pandich/couture/sink/color"
	"github.com/pandich/couture/sink/layout"
)

func newApplicationColumn(style color.FgBgTuple, layout layout.ColumnLayout) column {
	return newWeightedColumn(
		mapping.Application,
		layout,
		style,
		func(event event.SinkEvent) []interface{} { return []interface{}{string(event.Application)} },
	)
}
