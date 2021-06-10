package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/layout"
	"github.com/dustin/go-humanize"
	"time"
)

func newTimestampColumn(timeFormat *string, style sink.Style, layout layout.ColumnLayout) column {
	return newWeightedColumn(
		schema.Timestamp,
		layout,
		style,
		func(event model.SinkEvent) []interface{} {
			then := time.Time(event.Timestamp)
			if *timeFormat == model.HumanTimeFormat {
				humanized := humanize.Time(then)
				if humanized != "now" {
					return []interface{}{humanized}
				}
				return []interface{}{then.Format(time.Stamp)}
			}
			return []interface{}{then.Format(*timeFormat)}
		},
	)
}
