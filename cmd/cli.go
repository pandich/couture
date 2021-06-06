package cmd

import (
	"net/url"
	"time"
)

//nolint:lll
var cli struct {
	Metrics    dumpMetrics `group:"diagnostic" hidden:"true" default:"false" negatable:"true"`
	RateLimit  rateLimit   `group:"diagnostic" hidden:"true" default:"0" env:"COUTURE_RATE_LIMIT"`
	ShowSchema showSchema  `group:"diagnostic" hidden:"true" default:"false" negatable:"true" env:"COUTURE_SHOW_SCHEMA"`

	TTY        tty        `group:"terminal" help:"Force TTY mode." short:"T" default:"false"`
	Color      color      `group:"terminal" help:"Force Color mode." short:"T" default:"true" negatable:"true"`
	Wrap       wrap       `group:"terminal" help:"Wrap the output tp the terminal width, or that specified by --width." short:"w" default:"false" negatable:"true" env:"COUTURE_WIDTH"`
	Width      width      `group:"terminal" help:"Wrap width. Default is the current terminal width." placeholder:"width" short:"W"`
	AutoResize autoResize `group:"terminal" help:"Auto-resize columns when the terminal resizes." negatable:"true" default:"true"`

	Theme            themeName        `group:"display" help:"Specify the core Theme color: ${enum}." placeholder:"Theme" default:"${defaultTheme}" enum:"${themeNames}" env:"COUTURE_THEME"`
	ConsistentColors consistentColors `group:"display" help:"Maintain consistent source URL colors between runs." negatable:"true" default:"true"`
	Multiline        multiline        `group:"display" help:"Display each log event in multiline format. (Enabled by --expand-json)" negatable:"true" default:"false"`
	ExpandJSON       expandJSON       `group:"display" help:"Example JSON message bodies. Warning: has a significant performance impact." negatable:"true" default:"false"`

	Level     levelLike  `group:"filter" help:"The minimum log level to display: ${enum}." default:"${defaultLogLevel}" placeholder:"level" short:"l" enum:"${logLevels}" env:"COUTURE_LEVEL"`
	Since     *time.Time `group:"filter" help:"How far back to look for events. Parses most time and duration formats including human friendly." placeholder:"(time|duration)" short:"s" env:"COUTURE_SINCE"`
	Highlight highlight  `group:"filter" help:"Highlight matches from the patterns specified in --include." negatable:"true" default:"true"`
	Filter    filterLike `group:"filter" help:"Filter regular expressions. Format." placeholder:"[+|-|!]regex" short:"f" sep:"|"`

	TimeFormat timeFormat `group:"content" help:"Go-standard time format string or a named format: ${timeFormatNames}. (See https://golang.org/pkg/time/#pkg-constants)" short:"t" default:"stamp" env:"COUTURE_TIME_FORMAT"`
	Column     columns    `group:"content" help:"Specify one or more columns to display: ${enum}." placeholder:"column" enum:"${columnNames}" env:"COUTURE_COLUMN_NAMES"`

	Source []url.URL `arg:"true" help:"Log event source URL, alias, or alias group." name:"source-url"`
}
