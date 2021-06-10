package column

import (
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/sink/doric/config"
	layout2 "couture/internal/pkg/sink/layout"
	theme2 "couture/internal/pkg/sink/theme"
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

func newRegistry(config config.Config) registry {
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
			config.Highlight,
			config.Expand,
			config.Multiline,
			config.Theme.Level[level.Error].Bg,
			config.Theme.Message,
			config.Layout.Message,
		),
	}
}

func registerStyle(styleName string, style theme2.Style, layout layout2.ColumnLayout) {
	rawFormat := "{{%s%%s%s}}::" + style.Fg + "|bg" + style.Bg
	format := fmt.Sprintf(rawFormat, layout.Prefix(), layout.Suffix())
	cfmt.RegisterStyle(styleName, func(s string) string {
		return cfmt.Sprintf(format, s)
	})
}
