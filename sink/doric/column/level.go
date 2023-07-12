package column

import (
	"github.com/gagglepanda/couture/model"
	"github.com/gagglepanda/couture/model/level"
	"github.com/gagglepanda/couture/schema"
	"github.com/gagglepanda/couture/sink/color"
	"github.com/gagglepanda/couture/sink/layout"
	"github.com/i582/cfmt/cmd/cfmt"
)

type levelColumn struct {
	extractorColumn
}

func newLevelColumn(styles map[level.Level]color.FgBgTuple, layout layout.ColumnLayout) column {
	for _, lvl := range level.Levels {
		formatLevel := string(schema.Level) + string(lvl)
		cfmt.RegisterStyle(formatLevel, styles[lvl].Format())
	}
	return levelColumn{
		extractorColumn: extractorColumn{
			baseColumn: baseColumn{columnName: schema.Level, colLayout: layout},
			extractor: func(event model.SinkEvent) []interface{} {
				return []interface{}{string(event.Level)}
			},
		},
	}
}

func (col levelColumn) render(event model.SinkEvent) string {
	format := col.formatWithSuffix(string(event.Level))
	value := col.extractor(event)
	return cfmt.Sprintf(format, value...)
}
