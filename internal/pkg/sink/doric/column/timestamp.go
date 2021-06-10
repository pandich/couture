package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/layout"
	"time"
)

func newTimestampColumn(timeFormat *string, style sink.Style, layout layout.ColumnLayout) column {
	return newWeightedColumn(
		schema.Timestamp,
		layout,
		style,
		func(event model.SinkEvent) []interface{} {
			return []interface{}{time.Time(event.Timestamp).Format(*timeFormat)}
		},
	)
}
