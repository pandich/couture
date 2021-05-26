package column

import (
	"couture/internal/pkg/sink"
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
		func(event sink.Event) []interface{} {
			return []interface{}{string(event.ApplicationNameOrBlank())}
		},
	)}
}
