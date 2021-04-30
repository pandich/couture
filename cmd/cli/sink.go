package cli

import (
	"couture/internal/pkg/sink"
	"errors"
	"github.com/alecthomas/kong"
)

func init() {
	sinkMappers = append(sinkMappers, mapper(sink.GoString{}, sink.NewGoString)...)
	sinkMappers = append(sinkMappers, mapper(sink.Json{}, sink.NewJson)...)
}

var (
	//sinkCLI contains sink-specific cli args.
	sinkCLI struct {
		GoString *sink.GoString `group:"Output" help:"Dump string representation of event." name:"string" aliases:"go-string,str" xor:"console"`
		Json     *sink.Json     `group:"Output" help:"Dump JSON representation of event." name:"json" xor:"console"`
	}

	//sinkMappers contains sink-specific converters from string to a sink.Sink instance.
	sinkMappers []kong.Option
)

var (
	ErrNoSinks = errors.New("at least one sinks must be specified")
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

type sinksValidator struct{}

func (v sinksValidator) Validate() error {
	if len(Sinks()) == 0 {
		return ErrNoSinks
	}
	return nil
}
