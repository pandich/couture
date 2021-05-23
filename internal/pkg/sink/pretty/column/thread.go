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
		func(th theme.Theme) string { return th.ThreadColor() },
		func(event model.Event) []interface{} { return []interface{}{string(event.ThreadNameOrBlank())} },
	)}
}
