package column

import (
	"github.com/dustin/go-humanize"
	"github.com/pandich/couture/internal/pkg/model"
	"github.com/pandich/couture/internal/pkg/schema"
	"github.com/pandich/couture/internal/pkg/sink"
	"github.com/pandich/couture/internal/pkg/sink/layout"
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
