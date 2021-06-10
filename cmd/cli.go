package cmd

import (
	"net/url"
	"time"
)

//nolint:lll
var cli struct {
	DumpMetrics dumpMetrics `group:"diagnostic" hidden:"true"`
	DumpUnknown dumpUnknown `group:"diagnostic" hidden:"true"`
	RateLimit   rateLimit   `group:"diagnostic" hidden:"true" default:"0" env:"COUTURE_RATE_LIMIT"`
	ShowSchema  showSchema  `group:"diagnostic" hidden:"true" env:"COUTURE_SHOW_SCHEMA"`

	TTY        tty        `group:"terminal" help:"Force TTY mode." short:"T"`
	Color      color      `group:"terminal" help:"Force Color mode." short:"T" negatable:"true"`
	Wrap       wrap       `group:"terminal" help:"Wrap the output tp the terminal width, or that specified by --width." short:"w" negatable:"true" env:"COUTURE_WIDTH"`
	Width      width      `group:"terminal" help:"Wrap width. Default is the current terminal width." placeholder:"width" short:"W"`
	AutoResize autoResize `group:"terminal" help:"Auto-resize columns when the terminal resizes." negatable:"true"`

	Theme            themeName        `group:"display" help:"Specify the the color theme: ${enum}." placeholder:"theme" default:"${defaultTheme}" enum:"${themeNames}" env:"COUTURE_THEME"`
	ConsistentColors consistentColors `group:"display" help:"Maintain consistent source URL colors between runs." negatable:"true"`
	MultiLine        multiLine        `group:"display" help:"Display each log event in multi-line format. (Enabled by --expand-json)" negatable:"true"`
	Expand           expand           `group:"display" help:"Example structured message bodies (e.g. JSON)." negatable:"true"`

	Level     levelLike  `group:"filter" help:"The minimum log level to display: ${enum}." default:"${defaultLogLevel}" placeholder:"level" short:"l" enum:"${logLevels}" env:"COUTURE_LEVEL"`
	Since     *time.Time `group:"filter" help:"How far back to look for events. Parses most time and duration formats including human friendly." placeholder:"(time|duration)" short:"s" env:"COUTURE_SINCE"`
	Highlight highlight  `group:"filter" help:"Highlight matches from the patterns specified in --include." negatable:"true"`
	Filter    filterLike `group:"filter" help:"Filter regular expressions. Format." placeholder:"[+|-|!]regex" short:"f" sep:"|"`

	TimeFormat timeFormat `group:"content" help:"Go-standard time format string or a named format: ${timeFormatNames}. (See https://golang.org/pkg/time/#pkg-constants)" short:"t" default:"stamp" env:"COUTURE_TIME_FORMAT"`
	Column     columns    `group:"content" help:"Specify one or more columns to display: ${enum}." placeholder:"column" enum:"${columnNames}" env:"COUTURE_COLUMN_NAMES"`

	Source []url.URL `arg:"true" help:"Log event source URL, alias, or alias group." name:"source-url"`
}
