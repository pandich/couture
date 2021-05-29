package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/theme"
)

type threadColumn struct {
	weightedColumn
}

func newThreadColumn() threadColumn {
	const weight = 20
	sigil := '⇶'
	return threadColumn{newWeightedColumn(
		"thread",
		&sigil,
		weight,
		func(th theme.Theme) string { return th.ThreadFg() },
		func(event model.SinkEvent) []interface{} {
			return []interface{}{string(event.Thread)}
		},
	)}
}
