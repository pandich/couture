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
func (col threadColumn) name() string { return "thread" }

// weight ...
func (col threadColumn) weight() weight {
	const columnWeight = 20
	return columnWeight
}

// weightType ...
func (col threadColumn) weightType() weightType { return weighted }

// RegisterStyles ...
func (col threadColumn) RegisterStyles(theme theme.Theme) {
	cfmt.RegisterStyle(col.name(), func(s string) string {
		return cfmt.Sprintf("{{%s}}::"+theme.ThreadColor(), s)
	})
}

// Format ...
func (col threadColumn) Format(width uint, _ source.Source, _ model.Event) string {
	return "{{ ⇶ " + formatStringOfWidth(width) + " }}::" + col.name()
}

// Render ...
func (col threadColumn) Render(_ config.Config, _ source.Source, event model.Event) []interface{} {
	return []interface{}{string(event.ThreadNameOrBlank())}
}
