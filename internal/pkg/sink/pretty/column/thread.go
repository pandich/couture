package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/theme"
)

type contextColumn struct {
	weightedColumn
}

func newContextColumn() column {
	const weight = 20
	sigil := '⇶'
	return contextColumn{newWeightedColumn(
		"context",
		&sigil,
		weight,
		func(th theme.Theme) string { return th.ContextFg() },
		func(event model.SinkEvent) []interface{} {
			return []interface{}{orNoValue(string(event.Context))}
		},
	)}
}
