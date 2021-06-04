package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/theme"
)

type applicationColumn struct {
	weightedColumn
}

func newApplicationColumn() applicationColumn {
	const weight = 25
	sigil := '§'
	return applicationColumn{weightedColumn: newWeightedColumn(
		"application",
		&sigil,
		weight,
		func(th theme.Theme) string { return th.ApplicationFg() + "|bg" + th.ApplicationBg() },
		func(event model.SinkEvent) []interface{} {
			return []interface{}{orNoValue(string(event.Application))}
		},
	)}
}
