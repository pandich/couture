package cli

import (
	"couture/pkg/couture/model"
	"github.com/alecthomas/kong"
	"time"
)

var CommandLineArguments struct {
	Silent      bool `help:"Only log lines are displayed. Header and diagnostics are suppressed." name:"silent" short:"s" default:"false" group:"display"`
	ClearScreen bool `help:"Clear screen prior to start." name:"clear" default:"true" group:"display" negatable:"true"`
	ShowPrefix  bool `help:"Display a prefix before each log line indicting its source." name:"prefix" default:"true" group:"display" negatable:"true"`
	ShortNames  bool `help:"Display a abbreviated source names." name:"short-names" default:"true" group:"display" negatable:"true"`

	Follow           bool          `help:"Follow the logs." default:"true" name:"follow" short:"f" group:"main" negatable:"true"`
	FollowInterval   time.Duration `help:"How long to sleep between polls." default:"2s" short:"i" name:"follow-interval" group:"main"`
	LookbackInterval time.Duration `help:"How far back to search for events." default:"5m" short:"b" name:"lookback-interval" group:"main"`

	AwsRegion  string `help:"AWS region" default:"us-west-2" name:"aws-region" env:"AWS_REGION" group:"aws"`
	AwsProfile string `help:"AWS profile" default:"integration" name:"aws-profile" env:"AWS_PROFILE" group:"aws"`

	Patterns []string    `help:"Filter patterns." name:"filter" short:"f" sep:"," group:"filter"`
	LogLevel model.Level `help:"Minimum log level to display." default:"DEBUG" name:"log-level" short:"l" group:"filter"`
}

func ParseCommandLine() *kong.Context {
	return kong.Parse(
		&CommandLineArguments,
		kong.ShortUsageOnError(),
		kong.Name("couture"),
		kong.Description("Tail multiple log sources."),
		kong.ConfigureHelp(kong.HelpOptions{
			Tree:     true,
			Indenter: kong.TreeIndenter,
		}),
	)
}
