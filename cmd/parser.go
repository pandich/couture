package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/pandich/couture/couture"
	"github.com/pandich/couture/manager"
	"github.com/pandich/couture/model/level"
	"github.com/pandich/couture/schema"
	"github.com/pandich/couture/theme"
	"reflect"
	"strings"
	"time"
)

const helpSummary = "Tails one or more event sources."

var maybeDie = parser.FatalIfErrorf

var parserVars = kong.Vars{
	"timeFormatNames": strings.Join(timeFormatNames, ","),
	"columnNames":     strings.Join(schema.Names(), ","),
	"themeNames":      strings.Join(theme.Names(), ","),
	"defaultTheme":    theme.Prince,
	"logLevels":       strings.Join(level.LowerNames(), ","),
	"defaultLogLevel": level.Info.LowerName(),
}

var parser = kong.Must(&cli,
	kong.Name(couture.Name),
	kong.Description(helpDescription()),
	kong.UsageOnError(),
	kong.ConfigureHelp(kong.HelpOptions{
		Summary:   true,
		FlagsLast: true,
	}),
	kong.TypeMapper(reflect.TypeOf(&time.Time{}), timeLikeDecoder()),
	kong.Groups{
		"diagnostic": "Diagnostic Options",
		"terminal":   "Terminal Options",
		"display":    "Display Options",
		"content":    "Content Options",
		"filter":     "Filter Options",
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
