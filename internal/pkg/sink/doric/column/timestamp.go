package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/schema"
	layout2 "couture/internal/pkg/sink/layout"
	theme2 "couture/internal/pkg/sink/theme"
	"time"
)

func newTimestampColumn(timeFormat *string, style theme2.Style, layout layout2.ColumnLayout) column {
	return newWeightedColumn(
		schema.Timestamp,
		layout,
		style,
		func(event model.SinkEvent) []interface{} {
			return []interface{}{time.Time(event.Timestamp).Format(*timeFormat)}
		},
	)
}
