package cmd

import (
	"net/url"
	"time"
)

// cli contains all allowed arguments.
var cli struct {
	DumpMetrics dumpMetrics `hidden:"true" group:"diagnostic"`
	DumpUnknown dumpUnknown `hidden:"true" group:"diagnostic"`
	RateLimit   rateLimit   `hidden:"true" group:"diagnostic" env:"COUTURE_RATE_LIMIT" default:"0"`
	ShowMapping showMapping `hidden:"true" group:"diagnostic" env:"COUTURE_SHOW_MAPPING"`

	AutoResize autoResize `group:"terminal" help:"Automatically resize columns when the terminal resizes." negatable:"true" env:"COUTURE_AUTO_RESIZE"`
	NoColor    noColor    `group:"terminal" help:"Disable color mode." short:"C"`
	TTY        tty        `group:"terminal" help:"Enable TTY mode." short:"T"`
	Width      width      `group:"terminal" help:"Set the wrap width. Defaults to the current terminal width." placeholder:"width" short:"W" env:"COUTURE_WIDTH"`
	Wrap       wrap       `group:"terminal" help:"Wrap the output to the terminal width, or to the width specified by --width." short:"w" negatable:"true" env:"COUTURE_WRAP"`

	ColorMode        colorMode        `group:"display" optional:"true" help:"Set the color mode: ${enum}." enum:"auto,dark,light" default:"auto" env:"COUTURE_COLOR_MODE"`
	ConsistentColors consistentColors `group:"display" help:"Maintain consistent source URL colors across sessions." negatable:"true" env:"COUTURE_CONSISTENT_COLORS"`
	Expand           expand           `group:"display" help:"Expand structured message bodies, like JSON." negatable:"true" env:"COUTURE_EXPAND"`
	LevelMeter       levelMeter       `group:"display" help:"Display a frequency meter before each message." negatable:"true" env:"COUTURE_LEVEL_METER"`
	MultiLine        multiLine        `group:"display" help:"Use a multi-line format for log events. Enabled by --expand-json." negatable:"true"`
	SourceStyle      sourceStyle      `group:"display" help:"Choose a theme for source field background colors: ${enum}." enum:"happy,pastel,similar,warm" short:"S" default:"pastel"`
	Theme            themeName        `group:"display" help:"Specify the theme by hex code or name, including custom names: ${specialThemes}." default:"prince" env:"COUTURE_THEME"`

	Filter    filterLike `group:"filter" help:"Apply filter regular expressions." placeholder:"[+|-|!]regex" short:"f" sep:"|"`
	Highlight highlight  `group:"filter" help:"Highlight matches from specified patterns in --include." negatable:"true" env:"COUTURE_HIGHLIGHT"`
	Level     levelLike  `group:"filter" help:"Set the minimum log level to display: ${enum}." default:"trace" short:"l" enum:"trace,debug,info,warn,error" env:"COUTURE_LEVEL"`
	Since     *time.Time `group:"filter" help:"Look for events since a specific time." placeholder:"(time|duration)" short:"s"`

	Column     columns    `group:"content" help:"Specify columns to display: ${enum}." placeholder:"column" enum:"${columnNames}" env:"COUTURE_COLUMN_NAMES"`
	TimeFormat timeFormat `group:"content" help:"Use a Go-standard or named time format: ${timeFormatNames}. See Go time format constants." short:"t" default:"stamp" env:"COUTURE_TIME_FORMAT"`
	Source     []url.URL  `arg:"true" help:"Specify log event source URL, alias, or alias group." name:"source-url"`
}
