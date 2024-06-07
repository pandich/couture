package column

import (
	"github.com/dustin/go-humanize"
	"github.com/gagglepanda/couture/event"
	"github.com/gagglepanda/couture/mapping"
	"github.com/gagglepanda/couture/sink/color"
	"github.com/gagglepanda/couture/sink/layout"
	"strings"
	"time"
)

func newTimestampColumn(timeFormat *string, style color.FgBgTuple, layout layout.ColumnLayout) column {
	var count int
	switch {
	case timeFormat != nil && *timeFormat == event.HumanTimeFormat:
		count = 1
	case timeFormat != nil:
		count = len(*timeFormat)
	default:
		count = 1
	}
	empty := []interface{}{"-" + strings.Repeat(" ", count-1)}

	return newWeightedColumn(
		mapping.Timestamp,
		layout,
		style,
		func(evt event.SinkEvent) []interface{} {
			then := time.Time(evt.Timestamp)

			if *timeFormat == event.HumanTimeFormat {
				if then == (time.Time{}) {
					return empty
				}
				humanized := humanize.Time(then)
				if humanized != "now" {
					return []interface{}{humanized}
				}
				return []interface{}{then.Format(time.Stamp)}
			}
			if then == (time.Time{}) {
				return empty
			}
			return []interface{}{then.Format(*timeFormat)}
		},
	)
}
