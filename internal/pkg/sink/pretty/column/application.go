package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"github.com/i582/cfmt/cmd/cfmt"
)

type applicationColumn struct{}

// Name ...
func (col applicationColumn) name() string { return "application" }

// weight ...
func (col applicationColumn) weight() weight {
	const columnWeight = 25
	return columnWeight
}

// weightType ...
func (col applicationColumn) weightType() weightType { return weighted }

// RegisterStyles ...
func (col applicationColumn) RegisterStyles(theme theme.Theme) {
	cfmt.RegisterStyle(col.name(), func(s string) string {
		return cfmt.Sprintf("{{%s}}::"+theme.ApplicationColor(), s)
	})
}

// Format ...
func (col applicationColumn) Format(width uint, _ source.Source, _ model.Event) string {
	return "{{ § " + formatStringOfWidth(width) + " }}::" + col.name()
}

// Render ...
func (col applicationColumn) Render(_ config.Config, _ source.Source, event model.Event) []interface{} {
	return []interface{}{string(event.ApplicationNameOrBlank())}
}
