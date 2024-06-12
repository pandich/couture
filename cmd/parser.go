package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/pandich/couture/couture"
	"github.com/pandich/couture/manager"
	"github.com/pandich/couture/mapping"
	theme2 "github.com/pandich/couture/sink/theme"
	"reflect"
	"strings"
	"time"
)

// parser for loading the cli struct.
var parser = kong.Must(
	&cli,

	kong.Name(couture.Name),
	kong.Description(helpDescription()),

	kong.UsageOnError(),

	kong.ConfigureHelp(
		kong.HelpOptions{
			Summary:   true,
			FlagsLast: true,
			Compact:   true,
		},
	),

	// more advanced time decoding than kong has built-in
	kong.TypeMapper(reflect.TypeOf(&time.Time{}), timeLikeDecoder()),

	kong.Groups{
		"diagnostic": "Diagnostic Options",
		"terminal":   "Terminal Options",
		"display":    "Display Options",
		"content":    "Content Options",
		"filter":     "Filter Options",
	},

	// here, if we are actually in completions mode (see completions.go)
	// we want to let kong generate the completions and exit on its own
	kong.PostBuild(completionsHook),

	// additional values available to kong at parse-time
	kong.Vars{
		"timeFormatNames": strings.Join(timeFormatNames, ","),
		"columnNames":     strings.Join(mapping.Names(), ","),
		"specialThemes":   strings.Join(theme2.Names(), ","),
	},
)

// helpDescription generates the description value for the help.
func helpDescription() string {
	// TODO flesh out the command's help description
	var lines = []string{
		"Tails one or more event sources.",
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
