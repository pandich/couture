package cli

import (
	"couture/internal/pkg/couture"
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model/level"
	"github.com/alecthomas/kong"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const helpSummary = "Tails one or more event sourceArgs."

//nolint:lll
var cli struct {
	OutputFormat string `group:"Display Options" help:"The output format: ${enum}." enum:"pretty,json" default:"pretty" placeholder:"format" short:"f" required:"true" env:"COUTURE_DEFAULT_FORMAT"`
	Wrap         bool   `group:"Display Options" help:"Wrap the output tp the terminal width, or that specified by --width." short:"w" default:"true" negatable:"true"`
	Width        uint   `group:"Display Options" help:"Wrap width." placeholder:"width" short:"W" default:"0"`
	Theme        string `group:"Display Options" help:"Specify the core Theme color: ${enum}." placeholder:"Theme" default:"${defaultTheme}" enum:"${themeNames}"`
	MultiLine    bool   `group:"Display Options" help:"Display each log event in multi-line format." negatable:"true" default:"false"`
	Sigil        bool   `group:"Display Options" help:"Display column prefix sigils to help denote them." negatable:"true" default:"true"`
	ClearScreen  bool   `group:"Display Options" help:"Clear the screen prior to displaying events." negatable:"true" default:"true"`

	Column     []string `group:"Content Options" help:"Specify one or more columns to display: ${enum}." placeholder:"column" enum:"${columnNames}"`
	TimeFormat string   `group:"Content Options" help:"Go-standard time format string or a named format: ${timeFormatNames}." short:"t" default:"stamp"`

	Level   level.Level     `group:"Filter Options" help:"The minimum log level to display: ${enum}." default:"${defaultLogLevel}" placeholder:"level" short:"l" enum:"${logLevels}" env:"COUTURE_DEFAULT_LEVEL"`
	Since   time.Time       `group:"Filter Options" help:"How far back to look for events. Parses most time and duration formats including human friendly." placeholder:"(time|duration)" short:"s" default:"15m" env:"COUTURE_DEFAULT_SINCE"`
	Include []regexp.Regexp `group:"Filter Options" help:"Include filter regular expressions; they are performed before excludes." placeholder:"regex" short:"i" sep:"|"`
	Exclude []regexp.Regexp `group:"Filter Options" help:"Exclude filter regular expressions; they are performed after includes." placeholder:"regex" short:"x" sep:"|"`

	Source []url.URL `arg:"true" help:"Log event sourceArgs." name:"url" required:"true"`
}

// Run ...
func Run() {
	parser := kong.Must(&cli,
		kong.Name(couture.Name),
		kong.Description(description()),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{Summary: false, Tree: true}),
		kong.TypeMapper(reflect.TypeOf(regexp.Regexp{}), regexpDecoder()),
		kong.TypeMapper(reflect.TypeOf(time.Time{}), timeLikeDecoder()),
		parserVars(),
	)

	if runtime.GOOS == "windows" {
		parser.Fatalf("unsupported operating system: %s", runtime.GOOS)
	}

	// load config
	err := loadConfig()
	parser.FatalIfErrorf(err)

	// expand aliases, etc.
	_, err = parser.Parse(evaluatedOsArgs())
	parser.FatalIfErrorf(err)

	// get cli managerOptionFlags and args
	mgrOptions, err := managerOptionFlags()
	parser.FatalIfErrorf(err)
	sources, err := sourceArgs()
	parser.FatalIfErrorf(err)

	// create the manager and start it
	mgr, err := manager.New(append(mgrOptions, sources...)...)
	parser.FatalIfErrorf(err)

	err = (*mgr).Start()
	parser.FatalIfErrorf(err)
}

func description() string {
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
