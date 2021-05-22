package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"github.com/i582/cfmt/cmd/cfmt"
)

type threadColumn struct{}

// Name ...
func (t threadColumn) Name() string {
	return "thread"
}

// Register ...
func (t threadColumn) Register(theme theme.Theme) {
	cfmt.RegisterStyle(t.Name(), func(s string) string {
		return cfmt.Sprintf("{{%s}}::"+theme.ThreadColor(), s)
	})
}

// Formatter ...
func (t threadColumn) Formatter(_ source.Source, _ model.Event) string {
	return "{{ ⇶ %-15.15s }}::" + t.Name()
}

// Renderer ...
func (t threadColumn) Renderer(_ config.Config, _ source.Source, event model.Event) []interface{} {
	return []interface{}{string(event.ThreadNameOrBlank())}
}
