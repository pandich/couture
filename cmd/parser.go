package cmd

import (
	"couture/internal/pkg/couture"
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink/pretty"
	"couture/internal/pkg/sink/pretty/column"
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/alecthomas/kong"
	"reflect"
	"strings"
	"time"
)

const helpSummary = "Tails one or more event sources."

var maybeDie = parser.FatalIfErrorf

var parserVars = kong.Vars{
	"timeFormatNames":     strings.Join(timeFormatNames, ","),
	"columnNames":         strings.Join(column.Names(), ","),
	"themeNames":          strings.Join(theme.Names, ","),
	"defaultTheme":        theme.Default,
	"logLevels":           strings.Join(level.LowerNames(), ","),
	"defaultLogLevel":     level.Info.LowerNames(),
	"outputFormats":       strings.Join([]string{pretty.Name}, ","),
	"defaultOutputFormat": pretty.Name,
}

var parser = kong.Must(&cli,
	kong.Name(couture.Name),
	kong.Description(helpDescription()),
	kong.UsageOnError(),
	kong.ConfigureHelp(kong.HelpOptions{
		Summary:   true,
		FlagsLast: true,
	}),
	kong.TypeMapper(reflect.TypeOf(model.Filter{}), filterDecoder()),
	kong.TypeMapper(reflect.TypeOf(time.Time{}), timeLikeDecoder()),
	kong.Groups{
		"terminal": "Terminal Options",
		"display":  "Display Options",
		"content":  "Content Options",
		"filter":   "Filter Options",
	},
	kong.PostBuild(completionsHook),
	parserVars,
)

func helpDescription() string {
	var lines = []string{
		helpSummary,
		"",
		"Example Sources:",
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
