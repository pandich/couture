package cli

//See https://github.com/alecthomas/kong for documentation on the command-line parsing.

import (
	"couture/internal/pkg/model"
	"github.com/alecthomas/kong"
	"regexp"
	"time"
)

func init() {
	coreCli.Plugins = kong.Plugins{&sinkCLI, &sourceCLI}
}

//coreCli contains all global arguments.
//TODO continue porting https://github.com/gaggle-net/ldt-scripts/blob/main/cloudwatch/aws-log-tail/src/main/kotlin/com/gaggle/awstail/config.kt
var coreCli struct {
	AwsRegion  string `group:"aws" help:"AWS region" default:"us-west-2" name:"aws-region" aliases:"region" env:"AWS_REGION"`
	AwsProfile string `group:"aws" help:"AWS profile" default:"integration" name:"aws-profile" aliases:"profile" env:"AWS_PROFILE"`

	Quiet      bool `group:"display" help:"Only log lines are displayed. Header and diagnostics are suppressed." name:"quiet" aliases:"silent" short:"q" xor:"verbosity"`
	Verbose    bool `group:"display" help:"Display additional diagnostic data." name:"verbose" short:"v" xor:"verbosity"`
	ShowPrefix bool `group:"display" help:"Display a prefix before each log line indicting its source." name:"prefix" default:"true" negatable:"true"`
	ShortNames bool `group:"display" help:"Display a abbreviated source names." name:"short-names" aliases:"short" default:"true" negatable:"true"`
	Wrap       uint `group:"display" help:"Wrap output to the specified width." name:"wrap" short:"w" default:"0"`

	Follow       bool          `group:"behavior" help:"Follow the logs." default:"false" name:"follow" short:"f" negatable:"true"`
	PollInterval time.Duration `group:"behavior" help:"How long to sleep between polls. (Applies only to some sources.)" default:"2s" name:"interval" aliases:"sleep" short:"i"`
	LineCount    uint32        `group:"behavior" help:"How many lines of history to include. (Applies only to some sources.)" default:"20" name:"lines" aliases:"count" short:"c"`
	Since        time.Duration `group:"behavior" help:"How far back to search for events." default:"5m" name:"since" aliases:"back,lookback" short:"b"`

	Patterns []*regexp.Regexp `group:"filter" help:"Filter patterns." name:"filter" short:"f" sep:","`
	LogLevel model.Level      `group:"filter" help:"Minimum log level to display (${enum})." default:"DEBUG" name:"log-level" aliases:"level" short:"l" enum:"ERROR,WARN,INFO,DEBUG,TRACE"`

	coreValidator
	kong.Plugins
}

//MustLoad parses all arguments, including all sources and sinks.
func MustLoad() *kong.Context {
	var opts []kong.Option
	opts = append(opts, sourceMappers...)
	opts = append(opts, sinkMappers...)
	opts = append(opts, coreMappers...)
	opts = append(opts, coreOptions...)
	return kong.Parse(&coreCli, opts...)
}
