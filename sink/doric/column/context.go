package column

import (
	"github.com/pandich/couture/model"
	"github.com/pandich/couture/schema"
	"github.com/pandich/couture/sink/layout"
	"github.com/pandich/couture/theme"
)

func newContextColumn(style theme.Style, layout layout.ColumnLayout) column {
	return newWeightedColumn(
		schema.Context,
		layout,
		style,
		func(event model.SinkEvent) []interface{} { return []interface{}{string(event.Context)} },
	)
}
