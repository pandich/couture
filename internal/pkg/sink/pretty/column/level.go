package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/layout"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/model/theme"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/sink/pretty/config"
	"github.com/i582/cfmt/cmd/cfmt"
)

type levelColumn struct {
	baseColumn
}

func newLevelColumn(styles map[level.Level]theme.Style, layout layout.ColumnLayout) column {
	col := levelColumn{
		baseColumn: baseColumn{columnName: schema.Level, widthMode: fixed, colLayout: layout},
	}
	for _, lvl := range level.Levels {
		style := styles[lvl]
		cfmt.RegisterStyle(col.name()+string(lvl), func(s string) string {
			return cfmt.Sprintf("{{"+col.format()+"}}::bg"+style.Bg+"|"+style.Fg, "", s, "")
		})
	}
	return col
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
