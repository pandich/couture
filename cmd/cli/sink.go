package cli

import (
	"couture/internal/pkg/sink"
	"github.com/alecthomas/kong"
)

//init initializes all sink mappers
func init() {
	sinkMappers = append(sinkMappers, mapper(sink.GoString{}, sink.NewGoString)...)
	sinkMappers = append(sinkMappers, mapper(sink.Json{}, sink.NewJson)...)
	sinkMappers = append(sinkMappers, mapper(sink.Ansi{}, sink.NewAnsi)...)
}

var (
	//sinkCLI contains sink-specific cli args.
	sinkCLI struct {
		GoString *sink.GoString `group:"sink" help:"Dump string representation of event." name:"string" aliases:"go-string,str" xor:"console"`
		Json     *sink.Json     `group:"sink" help:"Dump JSON representation of event." name:"json" xor:"console"`
		Ansi     *sink.Ansi     `group:"sink" help:"ANSI output." name:"ansi" xor:"console"`
	}

	//sinkMappers contains sink-specific converters from string to a sink.Sink instance.
	sinkMappers []kong.Option

	//cliSinkOptions acts an adapter of sinkCLI to sink.Options.
	cliSinkOptions = sinkOptions{}
)

type (
	//sinkOptions is a trivial wrapper around coreCLI exposing its values as a sink.Options.
	sinkOptions struct{}
)

//Sinks returns all sink.Sink instances defined by the cli.
func Sinks() []interface{} {
	var sinks []interface{}
	if sinkCLI.GoString != nil {
		sinks = append(sinks, *sinkCLI.GoString)
	}
	if sinkCLI.Json != nil {
		sinks = append(sinks, *sinkCLI.Json)
	}
	if sinkCLI.Ansi != nil {
		sinks = append(sinks, *sinkCLI.Ansi)
	}
	return sinks
}

func (s sinkOptions) Wrap() uint {
	return coreCli.Wrap
}
