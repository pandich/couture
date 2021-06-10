package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/layout"
	"couture/internal/pkg/model/theme"
	"couture/internal/pkg/schema"
	"time"
)

func newTimestampColumn(timeFormat *string, style theme.Style, layout layout.ColumnLayout) column {
	return newWeightedColumn(
		schema.Timestamp,
		layout,
		style,
		func(event model.SinkEvent) []interface{} {
			return []interface{}{time.Time(event.Timestamp).Format(*timeFormat)}
		},
	)
}
