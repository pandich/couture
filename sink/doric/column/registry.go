package column

import (
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/pandich/couture/schema"
	"github.com/pandich/couture/sink"
	"github.com/pandich/couture/sink/layout"
	"github.com/pandich/couture/theme/color"
)

// DefaultColumns ...
var DefaultColumns = []schema.Column{
	sourcePseudoColumn,
	schema.Timestamp,
	schema.Application,
	schema.Context,
	callerPsuedoColumn,
	schema.Level,
	schema.Message,
}

type registry map[schema.Column]column

func newRegistry(config sink.Config) registry {
	return registry{
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
			levelMeter(config.LevelMeter != nil && *config.LevelMeter),
			useColor(config.Color != nil && *config.Color),
			config.Theme,
			config.Layout.Message,
		),
	}
}

func registerStyle(styleName string, style color.FgBgTuple, layout layout.ColumnLayout) {
	rawFormat := "{{%s%%s%s}}::" + style.Fg + "|bg" + style.Bg
	format := fmt.Sprintf(rawFormat, layout.Prefix(), layout.Suffix())
	cfmt.RegisterStyle(styleName, func(s string) string {
		return cfmt.Sprintf(format, s)
	})
}
