package column

import (
	"github.com/pandich/couture/internal/pkg/model"
	"github.com/pandich/couture/internal/pkg/schema"
	"github.com/pandich/couture/internal/pkg/sink"
	"github.com/pandich/couture/internal/pkg/sink/layout"
)

func newContextColumn(style sink.Style, layout layout.ColumnLayout) column {
	return newWeightedColumn(
		schema.Context,
		layout,
		style,
		func(event model.SinkEvent) []interface{} { return []interface{}{string(event.Context)} },
	)
}
