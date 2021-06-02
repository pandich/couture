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

func newLevelColumn() levelColumn {
	const width = 4
	return levelColumn{baseColumn{
		columnName:  "level",
		widthMode:   fixed,
		widthWeight: width,
	}}
}

// RegisterStyles ...
func (col levelColumn) RegisterStyles(thm theme.Theme) {
	for _, lvl := range level.Levels {
		fgColor := thm.LevelColorFg(lvl)
		bgColor := thm.LevelColorBg(lvl)
		cfmt.RegisterStyle(col.name()+string(lvl), func(s string) string {
			return cfmt.Sprintf("{{ %1.1s }}::bg"+bgColor+"|"+fgColor, s)
		})
	}
}

// Format ...
func (col levelColumn) Format(_ uint, event model.SinkEvent) string {
	var lvl = event.Level
	if lvl == "" {
		lvl = level.Info
	}
	return formatStyleOfWidth(col.name()+string(lvl), uint(col.weight()))
}

// Render ...
func (col levelColumn) Render(_ config.Config, event model.SinkEvent) []interface{} {
	var lvl = event.Level
	if lvl == "" {
		lvl = level.Info
	}
	return []interface{}{string(lvl)}
}
