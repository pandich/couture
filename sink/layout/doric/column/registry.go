package column

import (
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/pandich/couture/mapping"
	"github.com/pandich/couture/sink"
	"github.com/pandich/couture/sink/color"
	"github.com/pandich/couture/sink/layout"
)

// DefaultColumns ...
var DefaultColumns = []mapping.Column{
	sourcePseudoColumn,
	mapping.Timestamp,
	mapping.Application,
	mapping.Context,
	callerPsuedoColumn,
	mapping.Level,
	mapping.Message,
}

type registry map[mapping.Column]column

func newRegistry(config sink.Config) registry {
	return registry{
		sourcePseudoColumn: newSourceColumn(
			config.Layout.Source,
		),
		mapping.Timestamp: newTimestampColumn(
			config.TimeFormat,
			config.Theme.Timestamp,
			config.Layout.Timestamp,
		),
		mapping.Application: newApplicationColumn(
			config.Theme.Application,
			config.Layout.Application,
		),
		mapping.Context: newContextColumn(
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
		mapping.Level: newLevelColumn(
			config.Theme.Level,
			config.Layout.Level,
		),
		mapping.Message: newMessageColumn(
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
	cfmt.RegisterStyle(
		styleName, func(s string) string {
			return cfmt.Sprintf(format, s)
		},
	)
}
