package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/theme"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/sink/pretty/config"
)

type applicationColumn struct {
	weightedColumn
}

func newApplicationColumn(cfg config.Config) applicationColumn {
	layout := cfg.Layout.Application
	color := func(th theme.Theme) string {
		return th.ApplicationFg() + "|bg" + th.ApplicationBg()
	}
	value := func(event model.SinkEvent) []interface{} {
		return []interface{}{orNoValue(string(event.Application))}
	}
	return applicationColumn{weightedColumn: newWeightedColumn(schema.Application, layout, color, value)}
}
