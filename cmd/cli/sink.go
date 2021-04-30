package cli

import (
	"couture/internal/pkg/sink"
	"github.com/alecthomas/kong"
)

//sinkCLI contains sink-specific cli args.
var sinkCLI struct {
	GoString *sink.GoString `group:"Output Options" help:"Dump string representation of event." name:"string" aliases:"go-string,str" xor:"console"`
	Json     *sink.Json     `group:"Output Options" help:"Dump JSON representation of event." name:"json" xor:"console"`
}

func init() {
	sinkMappers = append(sinkMappers, mapper(sink.GoString{}, sink.NewGoString)...)
	sinkMappers = append(sinkMappers, mapper(sink.Json{}, sink.NewJson)...)
}

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

//sinkMappers contains sink-specific converters from string to a sink.Sink instance.
var sinkMappers []kong.Option
