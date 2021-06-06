package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/i582/cfmt/cmd/cfmt"
)

type levelColumn struct {
	baseColumn
}

func newLevelColumn() column {
	const width = 4
	return levelColumn{baseColumn{
		columnName:  "level",
		widthMode:   fixed,
		widthWeight: width,
	}}
}

// Init ...
func (col levelColumn) Init(thm theme.Theme) {
	for _, lvl := range level.Levels {
		fgColor := thm.LevelColorFg(lvl)
		bgColor := thm.LevelColorBg(lvl)
		cfmt.RegisterStyle(col.name()+string(lvl), func(s string) string {
			return cfmt.Sprintf("{{ %1.1s }}::bg"+bgColor+"|"+fgColor, s)
		})
	}
}

// RenderFormat ...
func (col levelColumn) RenderFormat(_ uint, event model.SinkEvent) string {
	var lvl = event.Level
	if lvl == "" {
		lvl = level.Info
	}
	return formatStyleOfWidth(col.name()+string(lvl), uint(col.weight()))
}

// RenderValue ...
func (col levelColumn) RenderValue(_ config.Config, event model.SinkEvent) []interface{} {
	var lvl = event.Level
	if lvl == "" {
		lvl = level.Info
	}
	return []interface{}{string(lvl)}
}
