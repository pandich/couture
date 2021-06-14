package cmd

import (
	"net/url"
	"time"
)

//nolint: lll
var cli struct {
	DumpMetrics dumpMetrics `hidden:"true" group:"diagnostic"`
	DumpUnknown dumpUnknown `hidden:"true" group:"diagnostic"`
	RateLimit   rateLimit   `hidden:"true" group:"diagnostic" env:"COUTURE_RATE_LIMIT" default:"0"`
	ShowSchema  showSchema  `hidden:"true" group:"diagnostic" env:"COUTURE_SHOW_SCHEMA"`

	AutoResize autoResize `group:"terminal" help:"Auto-resize columns when the terminal resizes." negatable:"true" env:"COUTURE_AUTO_RESIZE"`
	NoColor    noColor    `group:"terminal" help:"Force no color mode." short:"C"`
	TTY        tty        `group:"terminal" help:"Force TTY mode." short:"T"`
	Width      width      `group:"terminal" help:"Wrap width. Default is the current terminal width." placeholder:"width" short:"W" env:"COUTURE_WIDTH"`
	Wrap       wrap       `group:"terminal" help:"Wrap the output tp the terminal width, or that specified by --width." short:"w" negatable:"true" env:"COUTURE_WRAP"`

	ColorMode        string           `group:"display" optional:"true" help:"Specify the the color mode: ${enum}." enum:"auto,dark,light" default:"auto" env:"COUTURE_COLOR_MODE"`
	ConsistentColors consistentColors `group:"display" help:"Maintain consistent source URL colors between runs." negatable:"true" env:"COUTURE_CONSISTENT_COLORS"`
	Expand           expand           `group:"display" help:"Example structured message bodies (e.g. JSON)." negatable:"true" env:"COUTURE_EXPAND"`
	LevelMeter       levelMeter       `group:"display" help:"Display a meter before each message indicating how often it has been called." negatable:"true" env:"COUTURE_LEVEL_METER"`
	MultiLine        multiLine        `group:"display" help:"Display each log event in multi-line format. (Enabled by --expand-json)" negatable:"true"`
	SourceStyle      sourceStyle      `group:"display" help:"Select the theme for generating distinct source field background colors: ${enum}." enum:"happy,pastel,similar,warm" short:"S" default:"pastel"`
	Theme            themeName        `group:"display" help:"Specify the the hex code or name of the theme base noColor. Any hex code or name. Custom names: ${specialThemes}." default:"prince" env:"COUTURE_THEME"`

	Filter    filterLike `group:"filter" help:"Filter regular expressions. Format." placeholder:"[+|-|!]regex" short:"f" sep:"|"`
	Highlight highlight  `group:"filter" help:"Highlight matches from the patterns specified in --include." negatable:"true" env:"COUTURE_HIGHLIGHT"`
	Level     levelLike  `group:"filter" help:"The minimum log level to display: ${enum}." default:"trace" short:"l" enum:"trace,debug,info,warn,error" env:"COUTURE_LEVEL"`
	Since     *time.Time `group:"filter" help:"How far back to look for events. Parses most time and duration formats including human friendly." placeholder:"(time|duration)" short:"s"`

	Column     columns    `group:"content" help:"Specify one or more columns to display: ${enum}." placeholder:"column" enum:"${columnNames}" env:"COUTURE_COLUMN_NAMES"`
	TimeFormat timeFormat `group:"content" help:"Go-standard time format string or a named format: ${timeFormatNames}. (See https://golang.org/pkg/time/#pkg-constants)" short:"t" default:"stamp" env:"COUTURE_TIME_FORMAT"`

	Source []url.URL `arg:"true" help:"Log event source URL, alias, or alias group." name:"source-url"`
}
