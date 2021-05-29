package cli

import (
	"couture/internal/pkg/couture"
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink/pretty"
	"couture/internal/pkg/sink/pretty/column"
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/alecthomas/kong"
	"github.com/posener/complete/cmd/install"
	"github.com/willabides/kongplete"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// FEATURE shell completion
// FEATURE use the kong config integration for all the non-alias config

const helpSummary = "Tails one or more event sources."

//nolint:lll
var cli struct {
	OutputFormat     string `group:"display" hidden:"true" help:"The output format: ${enum}." enum:"${outputFormats}" default:"${defaultOutputFormat}" placeholder:"format" short:"f" required:"true" env:"COUTURE_DEFAULT_FORMAT"`
	Wrap             bool   `group:"display" help:"Wrap the output tp the terminal width, or that specified by --width." short:"w" default:"true" negatable:"true"`
	Width            uint   `group:"display" help:"Wrap width." placeholder:"width" short:"W"`
	Theme            string `group:"display" help:"Specify the core Theme color: ${enum}." placeholder:"Theme" default:"${defaultTheme}" enum:"${themeNames}" env:"COUTURE_THEME"`
	Multiline        bool   `group:"display" help:"Display each log event in multiline format." negatable:"true" default:"false"`
	Highlight        bool   `group:"display" help:"Highlight matches from the patterns specified in --include." negatable:"true" default:"true"`
	AutoResize       bool   `group:"display" help:"Auto-resize columns when the terminal resizes." negatable:"true" default:"true"`
	ConsistentColors bool   `group:"display" help:"Maintain consistent source URL colors between runs." negatable:"true" default:"true"`

	// FIXME format has to be by source
	Format     string     `group:"content" help:"Specify the JSON format used for the log events: ${enum}" short:"f" enum:"logstash" default:"logstash"`
	Column     []string   `group:"content" help:"Specify one or more columns to display: ${enum}." placeholder:"column" enum:"${columnNames}" env:"COUTURE_DEFAULT_COLUMN_NAMES"`
	TimeFormat timeFormat `group:"content" help:"Go-standard time format string or a named format: ${timeFormatNames}." short:"t" default:"stamp" env:"COUTURE_DEFAULT_TIME_FORMAT"`
	ExpandJSON bool       `group:"content" help:"Example JSON message bodies. Warning: has a significant performance impact." negatable:"true" default:"false"`

	Level   level.Level     `group:"filter" help:"The minimum log level to display: ${enum}." default:"${defaultLogLevel}" placeholder:"level" short:"l" enum:"${logLevels}" env:"COUTURE_DEFAULT_LEVEL"`
	Since   time.Time       `group:"filter" help:"How far back to look for events. Parses most time and duration formats including human friendly." placeholder:"(time|duration)" short:"s" default:"15m" env:"COUTURE_DEFAULT_SINCE"`
	Include []regexp.Regexp `group:"filter" help:"Include filter regular expressions; they are performed before excludes." placeholder:"regex" short:"i" sep:"|"`
	Exclude []regexp.Regexp `group:"filter" help:"Exclude filter regular expressions; they are performed after includes." placeholder:"regex" short:"x" sep:"|"`

	Source []url.URL `arg:"true" help:"Log event source URLs." name:"source_url" required:"true"`
}

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
		"themeNames":          strings.Join(theme.Names(), ","),
		"defaultTheme":        theme.Default,
		"logLevels":           strings.Join(level.Lower(), ","),
		"defaultLogLevel":     level.Info.Lower(),
		"outputFormats":       strings.Join([]string{pretty.Name}, ","),
		"defaultOutputFormat": pretty.Name,
	},
	kong.Groups{
		"display": "Display Config",
		"content": "Content Config",
		"filter":  "Filter Config",
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

// FIXME completions are not working properly
func completionsHook() func(k *kong.Kong) error {
	return func(k *kong.Kong) error {
		commandName := k.Model.Name
		doInstall := os.Getenv("COMP_INSTALL") == "1"
		doUninstall := os.Getenv("COMP_UNINSTALL") == "1"
		if doInstall || doUninstall {
			kongplete.Complete(k)
			var err error
			if doInstall {
				if install.IsInstalled(commandName) {
					_ = install.Uninstall(commandName)
				}
				err = install.Install(commandName)
			} else {
				err = install.Uninstall(commandName)
			}
			if err != nil {
				return err
			}
			os.Exit(0)
		}
		return nil
	}
}
