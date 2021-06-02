package cmd

import (
	"couture/internal/pkg/model/level"
	"net/url"
	"regexp"
	"time"
)

//nolint:lll
var cli struct {
	Metrics bool `group:"debug" hidden:"true" default:"false"`

	TTY        tty        `group:"terminal" help:"Force TTY mode." short:"T" default:"false"`
	Color      color      `group:"terminal" help:"Force Color mode." short:"T" default:"true" negatable:"true"`
	Wrap       wrap       `group:"terminal" help:"Wrap the output tp the terminal width, or that specified by --width." short:"w" default:"false" negatable:"true"`
	Width      width      `group:"terminal" help:"Wrap width. Default is the current terminal width." placeholder:"width" short:"W" predictor:"width"`
	AutoResize autoResize `group:"terminal" help:"Auto-resize columns when the terminal resizes." negatable:"true" default:"true"`

	Theme            themeName        `group:"display" help:"Specify the core Theme color: ${enum}." placeholder:"Theme" default:"${defaultTheme}" enum:"${themeNames}" env:"COUTURE_THEME" predictor:"themes"`
	ConsistentColors consistentColors `group:"display" help:"Maintain consistent source URL colors between runs." negatable:"true" default:"true"`
	Multiline        multiline        `group:"display" help:"Display each log event in multiline format. (Enabled by --expand-json)" negatable:"true" default:"false"`
	ExpandJSON       expandJSON       `group:"display" help:"Example JSON message bodies. Warning: has a significant performance impact." negatable:"true" default:"false"`

	Level     level.Level     `group:"filter" help:"The minimum log level to display: ${enum}." default:"${defaultLogLevel}" placeholder:"level" short:"l" enum:"${logLevels}" env:"COUTURE_DEFAULT_LEVEL" predictor:"level"`
	Since     time.Time       `group:"filter" help:"How far back to look for events. Parses most time and duration formats including human friendly." placeholder:"(time|duration)" short:"s" default:"15m" env:"COUTURE_DEFAULT_SINCE"`
	Highlight highlight       `group:"filter" help:"Highlight matches from the patterns specified in --include." negatable:"true" default:"true"`
	Include   []regexp.Regexp `group:"filter" help:"Include filter regular expressions; they are performed before excludes." placeholder:"regex" short:"i" sep:"|"`
	Exclude   []regexp.Regexp `group:"filter" help:"Exclude filter regular expressions; they are performed after includes." placeholder:"regex" short:"x" sep:"|"`

	TimeFormat timeFormat `group:"content" help:"Go-standard time format string or a named format: ${timeFormatNames}." short:"t" default:"stamp" env:"COUTURE_DEFAULT_TIME_FORMAT" predictor:"time_format"`
	Column     columns    `group:"content" help:"Specify one or more columns to display: ${enum}." placeholder:"column" enum:"${columnNames}" env:"COUTURE_DEFAULT_COLUMN_NAMES" predictor:"column_names"`

	Source []url.URL `arg:"true" help:"Log event source URLs." name:"source_url" predictor:"sources"`
}
