package column

import (
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/layout"
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
)

// DefaultColumns ...
var DefaultColumns = []string{
	sourcePseudoColumn,
	schema.Timestamp,
	schema.Application,
	schema.Context,
	callerPsuedoColumn,
	schema.Level,
	schema.Message,
}

type registry map[string]column

func newRegistry(config sink.Config) registry {
	errorStyle := config.Theme.Level[level.Error]
	return map[string]column{
		sourcePseudoColumn: newSourceColumn(
			config.Layout.Source,
		),
		schema.Timestamp: newTimestampColumn(
			config.TimeFormat,
			config.Theme.Timestamp,
			config.Layout.Timestamp,
		),
		schema.Application: newApplicationColumn(
			config.Theme.Application,
			config.Layout.Application,
		),
		schema.Context: newContextColumn(
			config.Theme.Context,
			config.Layout.Context,
		),
		callerPsuedoColumn: newCallerColumn(
			config.Theme.Entity,
			config.Theme.ActionDelimiter,
			config.Theme.Action,
			config.Theme.LineDelimiter,
			config.Theme.Line,
			config.Layout.Caller,
		),
		schema.Level: newLevelColumn(
			config.Theme.Level,
			config.Layout.Level,
		),
		schema.Message: newMessageColumn(
			highlight(config.Highlight != nil && *config.Highlight),
			expand(config.Expand != nil && *config.Expand),
			multiLine(config.MultiLine != nil && *config.MultiLine),
			errorStyle,
			config.Theme.Message,
			config.Layout.Message,
		),
	}
}

func registerStyle(styleName string, style sink.Style, layout layout.ColumnLayout) {
	rawFormat := "{{%s%%s%s}}::" + style.Fg + "|bg" + style.Bg
	format := fmt.Sprintf(rawFormat, layout.Prefix(), layout.Suffix())
	cfmt.RegisterStyle(styleName, func(s string) string {
		return cfmt.Sprintf(format, s)
	})
}
