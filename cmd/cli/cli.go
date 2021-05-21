package cli

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink/pretty"
	"github.com/alecthomas/kong"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const (
	applicationName = "couture"
	helpSummary     = "Tails one or more event sources."
)

//nolint:lll
var cli struct {
	OutputFormat string          `group:"Display Options" help:"The output format: ${enum}." enum:"pretty,json" default:"pretty" placeholder:"format" short:"f" required:"true" env:"COUTURE_DEFAULT_FORMAT"`
	Paginator    string          `group:"Display Options" help:"Set the paginator for --paginate mode." default:"more" placeholder:"command" env:"PAGER"`
	Paginate     bool            `group:"Display Options" help:"Paginate the results using an external paginator.  (default=${default})" short:"p" default:"false" negatable:"true"`
	Wrap         bool            `group:"Display Options" help:"Wrap the output. (default=${default})" placeholder:"width" short:"w" default:"true" negatable:"true"`
	Theme        string          `group:"Display Options" help:"Specify the core Theme color: ${enum}." placeholder:"Theme" default:"prince" enum:"none,prince"`
	Column       []string        `group:"Display Options" help:"Specify one or more columns to display: ${enum}." placeholder:"column" enum:"application,caller,level,message,stackTrace,thread,timestamp"`
	Level        level.Level     `group:"Filter Options" help:"The minimum log level to display: ${enum}." default:"info" placeholder:"level" short:"l" enum:"trace,debug,info,warn,error" env:"COUTURE_DEFAULT_LEVEL"`
	Since        time.Time       `group:"Filter Options" help:"How far back to look for events. May be a time or duration expression." placeholder:"(time|duration)" short:"s" default:"15m" env:"COUTURE_DEFAULT_SINCE"`
	Include      []regexp.Regexp `group:"Filter Options" help:"Include filter regular expressions; they are performed before excludes." placeholder:"regex" short:"i" sep:"|"`
	Exclude      []regexp.Regexp `group:"Filter Options" help:"Exclude filter regular expressions; they are performed after includes." placeholder:"regex" short:"x" sep:"|"`

	Source []url.URL `arg:"true" help:"Log event sources." name:"url" required:"true"`
}

// Run ...
func Run() {
	// new cli parser
	parser := kong.Must(&cli,
		kong.Name(applicationName),
		kong.Description(description()),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{Summary: false, Tree: true}),
		kong.TypeMapper(reflect.TypeOf(regexp.Regexp{}), regexpDecoder()),
		kong.TypeMapper(reflect.TypeOf(time.Time{}), timeLikeDecoder()),
		kong.TypeMapper(reflect.TypeOf(pretty.ColumnName("")), timeLikeDecoder()),
	)
	if runtime.GOOS == "windows" {
		parser.Fatalf("unsupported operating system: %s", runtime.GOOS)
	}

	// load config
	err := loadConfig()
	parser.FatalIfErrorf(err)

	// parse args
	_, err = parser.Parse(evaluateArgs())
	parser.FatalIfErrorf(err)

	// get cli flags
	mgrOptions, err := getFlags()
	parser.FatalIfErrorf(err)

	// get cli args
	sources, err := getArgs()
	parser.FatalIfErrorf(err)

	// create the manager and start it
	mgr, err := manager.New(append(mgrOptions, sources...)...)
	parser.FatalIfErrorf(err)
	err = (*mgr).Start()
	parser.FatalIfErrorf(err)
}

func description() string {
	var lines = []string{
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
	examples := strings.Join(lines, "\n")
	return helpSummary + "\n\n" + examples
}
