package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"github.com/i582/cfmt/cmd/cfmt"
	"time"
)

type timestampColumn struct {
	baseColumn
}

func newTimestampColumn() timestampColumn {
	return timestampColumn{baseColumn{
		columnName:  "timestamp",
		weightType:  filling,
		widthWeight: 0,
	}}
}

// RegisterStyles ...
func (col timestampColumn) RegisterStyles(theme theme.Theme) {
	cfmt.RegisterStyle(col.name(), func(s string) string {
		return cfmt.Sprintf("{{ ⌚ %s }}::"+theme.TimestampColor(), s)
	})
}

// Format ...
func (col timestampColumn) Format(width uint, _ source.Source, _ model.Event) string {
	return formatColumn(col, width)
}

// Render ...
func (col timestampColumn) Render(config config.Config, _ source.Source, event model.Event) []interface{} {
	return []interface{}{time.Time(event.Timestamp).Format(config.TimeFormat)}
}
