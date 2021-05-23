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
func (col levelColumn) name() string { return "level" }

// weight ...
func (col levelColumn) weight() weight {
	const columnWidth = 4
	return columnWidth
}

// weightType ...
func (col levelColumn) weightType() weightType { return fixed }

// RegisterStyles ...
func (col levelColumn) RegisterStyles(theme theme.Theme) {
	for _, lvl := range level.Levels {
		bgColor := theme.LevelColor(lvl)
		fgColor := tty.Contrast(bgColor)
		cfmt.RegisterStyle(col.name()+string(lvl), func(s string) string {
			return cfmt.Sprintf("{{%s}}::bg"+bgColor+"|"+fgColor, s)
		})
	}
}

// Format ...
func (col levelColumn) Format(_ uint, _ source.Source, event model.Event) string {
	return "{{ %1.1s }}::" + col.name() + string(event.Level)
}

// Render ...
func (col levelColumn) Render(_ config.Config, _ source.Source, event model.Event) []interface{} {
	return []interface{}{string(event.Level)}
}
