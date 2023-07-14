package column

import (
	"github.com/dustin/go-humanize"
	"github.com/gagglepanda/couture/event"
	"github.com/gagglepanda/couture/mapping"
	"github.com/gagglepanda/couture/sink/color"
	"github.com/gagglepanda/couture/sink/layout"
	"time"
)

func newTimestampColumn(timeFormat *string, style color.FgBgTuple, layout layout.ColumnLayout) column {
	return newWeightedColumn(
		mapping.Timestamp,
		layout,
		style,
		func(evt event.SinkEvent) []interface{} {
			then := time.Time(evt.Timestamp)
			if *timeFormat == event.HumanTimeFormat {
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
