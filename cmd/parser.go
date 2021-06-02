package cmd

import (
	"couture/internal/pkg/couture"
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink/pretty"
	"couture/internal/pkg/sink/pretty/column"
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/alecthomas/kong"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const helpSummary = "Tails one or more event sources."

var maybeDie = parser.FatalIfErrorf

var parser = kong.Must(&cli,
	kong.Name(couture.Name),
	kong.Description(helpDescription()),
	kong.UsageOnError(),
	kong.ConfigureHelp(kong.HelpOptions{
		Summary:   true,
		FlagsLast: true,
	}),
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
		"themeNames":          strings.Join(theme.Names, ","),
		"defaultTheme":        theme.Default,
		"logLevels":           strings.Join(level.Lower(), ","),
		"defaultLogLevel":     level.Info.Lower(),
		"outputFormats":       strings.Join([]string{pretty.Name}, ","),
		"defaultOutputFormat": pretty.Name,
	},
	kong.Groups{
		"terminal": "Terminal Options",
		"display":  "Display Options",
		"content":  "Content Options",
		"filter":   "Filter Options",
	},
	kong.PostBuild(completionsHook()),
)

func helpDescription() string {
	var lines = []string{
		helpSummary,
		"",
		"Examples Sources:",
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
