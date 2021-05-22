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
func (a applicationColumn) Name() string {
	return "application"
}

// Register ...
func (a applicationColumn) Register(theme theme.Theme) {
	cfmt.RegisterStyle(a.Name(), func(s string) string {
		return cfmt.Sprintf("{{%s}}::"+theme.ApplicationColor(), s)
	})
}

// Formatter ...
func (a applicationColumn) Formatter(_ source.Source, _ model.Event) string {
	return "{{ § %-20.20s }}::" + a.Name()
}

// Renderer ...
func (a applicationColumn) Renderer(_ config.Config, _ source.Source, event model.Event) []interface{} {
	return []interface{}{string(event.ApplicationNameOrBlank())}
}
