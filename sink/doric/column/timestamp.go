package column

import (
	"github.com/dustin/go-humanize"
	"github.com/gagglepanda/couture/model"
	"github.com/gagglepanda/couture/schema"
	"github.com/gagglepanda/couture/sink/color"
	"github.com/gagglepanda/couture/sink/layout"
	"time"
)

func newTimestampColumn(timeFormat *string, style color.FgBgTuple, layout layout.ColumnLayout) column {
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
