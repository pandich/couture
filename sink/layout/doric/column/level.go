package column

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/pandich/couture/event"
	"github.com/pandich/couture/event/level"
	"github.com/pandich/couture/mapping"
	"github.com/pandich/couture/sink/color"
	"github.com/pandich/couture/sink/layout"
)

type levelColumn struct {
	extractorColumn
}

func newLevelColumn(styles map[level.Level]color.FgBgTuple, layout layout.ColumnLayout) column {
	for _, lvl := range level.Levels {
		formatLevel := string(mapping.Level) + string(lvl)
		cfmt.RegisterStyle(formatLevel, styles[lvl].Format())
	}
	return levelColumn{
		extractorColumn: extractorColumn{
			baseColumn: baseColumn{columnName: mapping.Level, colLayout: layout},
			extractor: func(event event.SinkEvent) []interface{} {
				return []interface{}{string(event.Level)}
			},
		},
	}
}

func (col levelColumn) render(event event.SinkEvent) string {
	format := col.formatWithSuffix(string(event.Level))
	value := col.extractor(event)
	return cfmt.Sprintf(format, value...)
}
