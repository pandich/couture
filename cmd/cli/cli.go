package cli

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model/level"
	"github.com/alecthomas/kong"
	"github.com/posener/complete"
	"github.com/willabides/kongplete"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const (
	applicationName           = "couture"
	logCommand                = "log"
	installCompletionsCommand = "install-completions"
)

func newParser() *kong.Kong {
	parser := kong.Must(&cli,
		kong.Name(applicationName),
		// TODO description doesn't show up in standard help?
		kong.Description(description()),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{Summary: false, Tree: true}),
		kong.TypeMapper(reflect.TypeOf(regexp.Regexp{}), regexpDecoder()),
		kong.TypeMapper(reflect.TypeOf(time.Time{}), timeLikeDecoder()),
	)
	// FIXME kongplete doesn't do anything detectable
	kongplete.Complete(parser, kongplete.WithPredictor("file", complete.PredictFiles("*")))
	return parser
}

//nolint:lll
var cli struct {
	Log struct {
		OutputFormat string          `help:"The output format: ${enum}" enum:"pretty,json" default:"pretty" placeholder:"format" short:"f"`
		Paginator    string          `help:"Set the paginator for --paginate mode." default:"" placeholder:"command" env:"COUTURE_PAGINATOR"`
		Paginate     bool            `help:"Paginate the results using an external paginator" short:"p" negatable:"true"`
		Wrap         uint            `help:"Wrap the output. Use --no-wrap or --wrap=0 to disable." placeholder:"width" short:"w" xor:"wrap"`
		NoWrap       bool            `hide:"true" xor:"wrap"`
		Level        level.Level     `help:"The minimum log level to display: ${enum}" default:"info" placeholder:"level" enum:"trace,debug,info,warn,error"`
		Since        time.Time       `help:"How far back to look for events. May be a time or duration expression." default:"15m"`
		Include      []regexp.Regexp `help:"Include filter regular expressions. Performed before excludes." placeholder:"regex" short:"i"`
		Exclude      []regexp.Regexp `help:"Exclude filter regular expressions. Performed after includes." placeholder:"regex" short:"x"`
		Source       []url.URL       `arg:"true" help:"Log event sources." name:"url" required:"true"`
	} `cmd:""`

	InstallCompletions kongplete.InstallCompletions `cmd:""`
}

func description() string {
	const delimiter = "\n  "
	allExampleURLs := manager.SourceMetadata.ExampleURLs()
	return "Tails one or more event sources.\n\nExamples Source URLs:" +
		delimiter + strings.Join(allExampleURLs, delimiter)
}
