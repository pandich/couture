package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/i582/cfmt/cmd/cfmt"
	"time"
)

type timestampColumn struct {
	baseColumn
}

func newTimestampColumn() timestampColumn {
	return timestampColumn{baseColumn{
		columnName:  "timestamp",
		widthMode:   filling,
		widthWeight: 0,
	}}
}

// RegisterStyles ...
func (col timestampColumn) RegisterStyles(theme theme.Theme) {
	cfmt.RegisterStyle(col.name(), func(s string) string {
		return cfmt.Sprintf("{{ ☀︎ %s }}::bg"+theme.TimestampBg()+"|"+theme.TimestampFg(), s)
	})
}

// Format ...
func (col timestampColumn) Format(width uint, _ model.SinkEvent) string {
	return formatColumn(col, width)
}

// Render ...
func (col timestampColumn) Render(cfg config.Config, event model.SinkEvent) []interface{} {
	return []interface{}{orNoValue(time.Time(event.Timestamp).Format(cfg.TimeFormat))}
}
