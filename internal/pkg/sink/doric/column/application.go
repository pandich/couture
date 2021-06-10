package column

import (
	"github.com/pandich/couture/internal/pkg/model"
	"github.com/pandich/couture/internal/pkg/schema"
	"github.com/pandich/couture/internal/pkg/sink"
	"github.com/pandich/couture/internal/pkg/sink/layout"
)

func newApplicationColumn(style sink.Style, layout layout.ColumnLayout) column {
	return newWeightedColumn(
		schema.Application,
		layout,
		style,
		func(event model.SinkEvent) []interface{} { return []interface{}{string(event.Application)} },
	)
}
