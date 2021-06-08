package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/theme"
	"couture/internal/pkg/sink/pretty/config"
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
	"time"
)

type timestampColumn struct {
	baseColumn
}

func newTimestampColumn() column {
	return timestampColumn{baseColumn{
		columnName:  "timestamp",
		widthMode:   filling,
		widthWeight: 0,
	}}
}

// Init ...
func (col timestampColumn) Init(theme theme.Theme) {
	cfmt.RegisterStyle(col.name(), func(s string) string {
		return cfmt.Sprintf("{{ ☀︎ %s }}::bg"+theme.TimestampBg()+"|"+theme.TimestampFg(), s)
	})
}

// RenderFormat ...
func (col timestampColumn) RenderFormat(width uint, _ model.SinkEvent) string {
	return formatColumn(col, width)
}

// RenderValue ...
func (col timestampColumn) RenderValue(cfg config.Config, event model.SinkEvent) []interface{} {
	if cfg.TimeFormat == nil {
		return []interface{}{fmt.Sprint(event.Timestamp)}
	}
	t := time.Time(event.Timestamp)
	txt := t.Format(*cfg.TimeFormat)
	return []interface{}{orNoValue(txt)}
}
