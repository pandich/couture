package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
)

type sourceColumn struct{}

// Name ...
func (s sourceColumn) Name() string {
	return "source"
}

// Register ...
func (s sourceColumn) Register(_ theme.Theme) {
}

// Formatter ...
func (s sourceColumn) Formatter(src source.Source, _ model.Event) string {
	return "{{ " + string(src.Sigil()) + " %-30.30s }}::" + src.ID()
}

// Renderer ...
func (s sourceColumn) Renderer(_ config.Config, src source.Source, _ model.Event) []interface{} {
	return []interface{}{src.URL().ShortForm()}
}
