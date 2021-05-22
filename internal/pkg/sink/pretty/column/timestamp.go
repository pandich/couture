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
func (t timestampColumn) Name() string {
	return "timestamp"
}

// Register ...
func (t timestampColumn) Register(theme theme.Theme) {
	cfmt.RegisterStyle(t.Name(), func(s string) string {
		return cfmt.Sprintf("{{%s}}::"+theme.TimestampColor(), s)
	})
}

// Formatter ...
func (t timestampColumn) Formatter(_ source.Source, _ model.Event) string {
	return "{{ ⌚ %s }}::" + t.Name()
}

// Renderer ...
func (t timestampColumn) Renderer(config config.Config, _ source.Source, event model.Event) []interface{} {
	return []interface{}{time.Time(event.Timestamp).Format(config.TimeFormat)}
}
