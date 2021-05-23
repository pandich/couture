package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"github.com/i582/cfmt/cmd/cfmt"
)

type threadColumn struct {
	baseColumn
}

func newThreadColumn() threadColumn {
	const weight = 20
	return threadColumn{baseColumn{
		columnName:  "thread",
		weightType:  weighted,
		widthWeight: weight,
	}}
}

// RegisterStyles ...
func (col threadColumn) RegisterStyles(theme theme.Theme) {
	cfmt.RegisterStyle(col.name(), func(s string) string {
		return cfmt.Sprintf("{{ ⇶ %s }}::"+theme.ThreadColor(), s)
	})
}

// Format ...
func (col threadColumn) Format(width uint, _ source.Source, _ model.Event) string {
	return formatColumn(col, width)
}

// Render ...
func (col threadColumn) Render(_ config.Config, _ source.Source, event model.Event) []interface{} {
	return []interface{}{string(event.ThreadNameOrBlank())}
}
