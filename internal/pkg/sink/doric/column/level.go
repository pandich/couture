package column

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/pandich/couture/internal/pkg/model"
	"github.com/pandich/couture/internal/pkg/model/level"
	"github.com/pandich/couture/internal/pkg/schema"
	"github.com/pandich/couture/internal/pkg/sink"
	"github.com/pandich/couture/internal/pkg/sink/layout"
)

type levelColumn struct {
	extractorColumn
}

func newLevelColumn(styles map[level.Level]sink.Style, layout layout.ColumnLayout) column {
	for _, lvl := range level.Levels {
		formatLevel := schema.Level + string(lvl)
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
