package cli

import (
	"couture/internal/pkg/sink"
	"github.com/alecthomas/kong"
)

var (
	//sinkCLI contains sink-specific cli args.
	sinkCLI struct {
		String *sink.GoString `group:"Output Options" help:"Dump string representation of event." name:"string" aliases:"go-string,str"`
	}

	//sinkMappers contains sink-specific converters from string to a sink.Sink instance.
	sinkMappers = []kong.Option{
		mapper(one(sink.GoString{}), sink.NewGoString),
	}
)

//Sinks returns all sink.Sink instances defined by the cli.
func Sinks() []interface{} {
	var sinks []interface{}
	if sinkCLI.String != nil {
		sinks = append(sinks, *sinkCLI.String)
	}
	return sinks
}
