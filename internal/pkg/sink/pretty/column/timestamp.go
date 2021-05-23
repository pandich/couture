package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"github.com/i582/cfmt/cmd/cfmt"
	"time"
)

type timestampColumn struct{}

// Name ...
func (col timestampColumn) name() string { return "timestamp" }

// weight ...
func (col timestampColumn) weight() weight { return 0 }

// weightType ...
func (col timestampColumn) weightType() weightType { return filling }

// RegisterStyles ...
func (col timestampColumn) RegisterStyles(theme theme.Theme) {
	cfmt.RegisterStyle(col.name(), func(s string) string {
		return cfmt.Sprintf("{{%s}}::"+theme.TimestampColor(), s)
	})
}

// Format ...
func (col timestampColumn) Format(_ uint, _ source.Source, _ model.Event) string {
	return "{{ ⌚ %s }}::" + col.name()
}

// Render ...
func (col timestampColumn) Render(config config.Config, _ source.Source, event model.Event) []interface{} {
	return []interface{}{time.Time(event.Timestamp).Format(config.TimeFormat)}
}
