package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"couture/internal/pkg/tty"
	"github.com/i582/cfmt/cmd/cfmt"
)

type levelColumn struct{}

// Name ...
func (l levelColumn) Name() string {
	return "level"
}

// Register ...
func (l levelColumn) Register(theme theme.Theme) {
	for _, lvl := range level.Levels {
		bgColor := theme.LevelColor(lvl)
		fgColor := tty.Contrast(bgColor)
		cfmt.RegisterStyle(l.Name()+string(lvl), func(s string) string {
			return cfmt.Sprintf("{{%s}}::bg"+bgColor+"|"+fgColor, s)
		})
	}
}

// Formatter ...
func (l levelColumn) Formatter(_ source.Source, event model.Event) string {
	return "{{ %1.1s }}::" + l.Name() + string(event.Level)
}

// Renderer ...
func (l levelColumn) Renderer(_ config.Config, _ source.Source, event model.Event) []interface{} {
	return []interface{}{string(event.Level)}
}
