package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/model/theme"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/sink/pretty/config"
	"github.com/i582/cfmt/cmd/cfmt"
)

type levelColumn struct {
	baseColumn
}

func newLevelColumn(cfg config.Config) column {
	layout := cfg.Layout.Level
	return levelColumn{baseColumn{
		columnName: schema.Level,
		widthMode:  fixed,
		colLayout:  layout,
	}}
}

// Init ...
func (col levelColumn) Init(thm theme.Theme) {
	for _, lvl := range level.Levels {
		fgColor := thm.LevelColorFg(lvl)
		bgColor := thm.LevelColorBg(lvl)
		cfmt.RegisterStyle(col.name()+string(lvl), func(s string) string {
			return cfmt.Sprintf("{{"+col.format()+"}}::bg"+bgColor+"|"+fgColor, "", s, "")
		})
	}
}

// Render ...
func (col levelColumn) Render(_ config.Config, event model.SinkEvent) string {
	var lvl = event.Level
	if lvl == "" {
		lvl = level.Info
	}
	levelName := string(lvl)
	return cfmt.Sprintf(formatStyleOfWidth(col.name()+levelName, uint(col.weight())), levelName)
}
