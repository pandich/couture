package cli

import (
	"couture/internal/pkg/sink"
	"github.com/alecthomas/kong"
)

var (
	//sinkCLI contains sink-specific cli args.
	sinkCLI struct {
		GoString *sink.GoString `group:"Output Options" help:"Dump string representation of event." name:"string" aliases:"go-string,str" xor:"console"`
		Json     *sink.Json     `group:"Output Options" help:"Dump JSON representation of event." name:"json" xor:"console"`
	}

	//sinkMappers contains sink-specific converters from string to a sink.Sink instance.
	sinkMappers = []kong.Option{
		mapper(one(sink.GoString{}), sink.NewGoString),
		mapper(one(sink.Json{}), sink.NewJson),
	}
)

//Sinks returns all sink.Sink instances defined by the cli.
func Sinks() []interface{} {
	var sinks []interface{}
	if sinkCLI.GoString != nil {
		sinks = append(sinks, *sinkCLI.GoString)
	} else if sinkCLI.Json != nil {
		sinks = append(sinks, *sinkCLI.Json)
	}
	return sinks
}
