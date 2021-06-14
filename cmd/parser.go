package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/pandich/couture/couture"
	"github.com/pandich/couture/manager"
	"github.com/pandich/couture/schema"
	theme2 "github.com/pandich/couture/sink/theme"
	"reflect"
	"strings"
	"time"
)

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

var parserVars = kong.Vars{
	"timeFormatNames": strings.Join(timeFormatNames, ","),
	"columnNames":     strings.Join(schema.Names(), ","),
	"specialThemes":   strings.Join(theme2.Names(), ","),
}

const helpSummary = "Tails one or more event sources."

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
