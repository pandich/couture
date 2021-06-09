package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/theme"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/sink/pretty/config"
)

type contextColumn struct {
	weightedColumn
}

func newContextColumn(cfg config.Config) column {
	layout := cfg.Layout.Context
	color := func(th theme.Theme) string { return th.ContextFg() }
	value := func(event model.SinkEvent) []interface{} {
		return []interface{}{orNoValue(string(event.Context))}
	}
	return contextColumn{newWeightedColumn(schema.Context, layout, color, value)}
}
