package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/layout"
	"couture/internal/pkg/model/theme"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/sink/pretty/config"
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
	"time"
)

type timestampColumn struct {
	baseColumn
}

func newTimestampColumn(style theme.Style, layout layout.ColumnLayout) column {
	col := timestampColumn{
		baseColumn: baseColumn{columnName: schema.Timestamp, widthMode: filling, colLayout: layout},
	}
	cfmt.RegisterStyle(col.name(), func(s string) string {
		return cfmt.Sprintf("{{%s %s}}::bg"+style.Fg+"|"+style.Bg, col.colLayout.Sigil, s)
	})
	return col
}

// Render ...
func (col timestampColumn) Render(cfg config.Config, event model.SinkEvent) string {
	if cfg.TimeFormat == nil {
		return fmt.Sprint(event.Timestamp)
	}
	t := time.Time(event.Timestamp)
	txt := t.Format(*cfg.TimeFormat)
	return cfmt.Sprintf(formatColumn(col, col.colLayout.Width), orNoValue(txt))
}
