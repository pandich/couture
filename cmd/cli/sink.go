package cli

import (
	"couture/internal/pkg/sink"
	"github.com/alecthomas/kong"
)

//sinkCLI contains sink-specific cli args.
var sinkCLI struct {
	GoString *sink.GoString `group:"sink" help:"Dump string representation of event." name:"string" aliases:"go-string,str" xor:"console"`
	Json     *sink.Json     `group:"sink" help:"Dump JSON representation of event." name:"json" xor:"console"`
	Logrus   *sink.Logrus   `group:"sink" help:"Formatting log output." name:"log" xor:"console"`
}

func init() {
	sinkMappers = append(sinkMappers, mapper(sink.GoString{}, sink.NewGoString)...)
	sinkMappers = append(sinkMappers, mapper(sink.Json{}, sink.NewJson)...)
	sinkMappers = append(sinkMappers, mapper(sink.Logrus{}, sink.NewLogrus)...)
}

//Sinks returns all sink.Sink instances defined by the cli.
func Sinks() []interface{} {
	var sinks []interface{}
	if sinkCLI.GoString != nil {
		sinks = append(sinks, *sinkCLI.GoString)
	}
	if sinkCLI.Json != nil {
		sinks = append(sinks, *sinkCLI.Json)
	}
	if sinkCLI.Logrus != nil {
		sinks = append(sinks, *sinkCLI.Logrus)
	}
	return sinks
}

//sinkMappers contains sink-specific converters from string to a sink.Sink instance.
var sinkMappers []kong.Option
