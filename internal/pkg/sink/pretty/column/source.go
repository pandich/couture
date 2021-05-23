package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
)

type sourceColumn struct{}

// Name ...
func (col sourceColumn) name() string { return "source" }

// weight ...
func (col sourceColumn) weight() weight {
	const columnWeight = 40
	return columnWeight
}

// weightType ...
func (col sourceColumn) weightType() weightType { return weighted }

// RegisterStyles ...
func (col sourceColumn) RegisterStyles(_ theme.Theme) {}

// Format ...
func (col sourceColumn) Format(width uint, src source.Source, _ model.Event) string {
	return "{{ " + string(src.Sigil()) + " " + formatStringOfWidth(width) + " }}::" + src.ID()
}

// Render ...
func (col sourceColumn) Render(_ config.Config, src source.Source, _ model.Event) []interface{} {
	return []interface{}{src.URL().ShortForm()}
}
