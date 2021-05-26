package cli

import (
	"couture/internal/pkg/couture"
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink/pretty"
	"couture/internal/pkg/sink/pretty/column"
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/alecthomas/kong"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const helpSummary = "Tails one or more event sources."

//nolint:lll
var cli struct {
	OutputFormat string `hidden:"true" group:"Display Options" help:"The output format: ${enum}." enum:"${outputFormats}" default:"${defaultOutputFormat}" placeholder:"format" short:"f" required:"true" env:"COUTURE_DEFAULT_FORMAT"`
	Wrap         bool   `group:"Display Options" help:"Wrap the output tp the terminal width, or that specified by --width." short:"w" default:"true" negatable:"true"`
	Width        uint   `group:"Display Options" help:"Wrap width." placeholder:"width" short:"W"`
	Theme        string `group:"Display Options" help:"Specify the core Theme color: ${enum}." placeholder:"Theme" default:"${defaultTheme}" enum:"${themeNames}"`
	Multiline    bool   `group:"Display Options" help:"Display each log event in multiline format." negatable:"true" default:"false"`
	Highlight    bool   `group:"Display Options" help:"Highlight matches from the patterns specified in --include." negatable:"true" default:"true"`

	Column     []string   `group:"Content Options" help:"Specify one or more columns to display: ${enum}." placeholder:"column" enum:"${columnNames}"`
	TimeFormat timeFormat `group:"Content Options" help:"Go-standard time format string or a named format: ${timeFormatNames}." short:"t" default:"stamp"`

	Level   level.Level     `group:"Filter Options" help:"The minimum log level to display: ${enum}." default:"${defaultLogLevel}" placeholder:"level" short:"l" enum:"${logLevels}" env:"COUTURE_DEFAULT_LEVEL"`
	Since   time.Time       `group:"Filter Options" help:"How far back to look for events. Parses most time and duration formats including human friendly." placeholder:"(time|duration)" short:"s" default:"15m" env:"COUTURE_DEFAULT_SINCE"`
	Include []regexp.Regexp `group:"Filter Options" help:"Include filter regular expressions; they are performed before excludes." placeholder:"regex" short:"i" sep:"|"`
	Exclude []regexp.Regexp `group:"Filter Options" help:"Exclude filter regular expressions; they are performed after includes." placeholder:"regex" short:"x" sep:"|"`

	Source []url.URL `arg:"true" help:"Log event source URLs." name:"source_url" required:"true"`
}

var parser = kong.Must(&cli,
	kong.Name(couture.Name),
	kong.Description(helpDescription()),
	kong.UsageOnError(),
	kong.ConfigureHelp(kong.HelpOptions{Summary: false, Tree: true}),
	kong.TypeMapper(reflect.TypeOf(regexp.Regexp{}), regexpDecoder()),
	kong.TypeMapper(reflect.TypeOf(time.Time{}), timeLikeDecoder()),
	kong.Vars{
		"timeFormatNames": strings.Join([]string{
			"c",
			"iso8601",
			"iso8601-nanos",
			"kitchen",
			"rfc1123",
			"rfc1123-utc",
			"rfc3339",
			"rfc3339-nanos",
			"rfc822",
			"rfc822-utc",
			"rfc850",
			"ruby",
			"stamp",
			"stamp-micros",
			"stamp-millis",
			"stamp-nanos",
			"unix",
		}, ","),
		"columnNames":         strings.Join(column.Names(), ","),
		"themeNames":          strings.Join(theme.Names(), ","),
		"defaultTheme":        theme.Prince,
		"logLevels":           strings.Join(level.SimpleNames(), ","),
		"defaultLogLevel":     level.Info.SimpleName(),
		"outputFormats":       strings.Join([]string{pretty.Name}, ","),
		"defaultOutputFormat": pretty.Name,
	},
)

func helpDescription() string {
	var lines = []string{
		helpSummary,
		"",
		"Examples Source URLs:",
		"",
	}
	for _, src := range manager.AvailableSources {
		if len(src.ExampleURLs) > 0 {
			lines = append(lines, "  "+src.Name+":")
			for _, u := range src.ExampleURLs {
				lines = append(lines, "    "+u)
			}
			lines = append(lines, "")
		}
	}
	return strings.Join(lines, "\n")
}
