package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"github.com/i582/cfmt/cmd/cfmt"
)

type applicationColumn struct {
	baseColumn
}

func newApplicationColumn() applicationColumn {
	const weight = 25
	return applicationColumn{baseColumn{
		columnName:  "application",
		weightType:  weighted,
		widthWeight: weight,
	}}
}

// RegisterStyles ...
func (col applicationColumn) RegisterStyles(theme theme.Theme) {
	cfmt.RegisterStyle(col.name(), func(s string) string {
		return cfmt.Sprintf("{{ § %s }}::"+theme.ApplicationColor(), s)
	})
}

// Format ...
func (col applicationColumn) Format(width uint, _ source.Source, _ model.Event) string {
	return formatColumn(col, width)
}

// Render ...
func (col applicationColumn) Render(_ config.Config, _ source.Source, event model.Event) []interface{} {
	return []interface{}{string(event.ApplicationNameOrBlank())}
}
