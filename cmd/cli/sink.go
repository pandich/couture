package cli

import (
	"couture/internal/pkg/sink/json"
	"couture/internal/pkg/sink/pretty"
	string2 "couture/internal/pkg/sink/simple"
	"github.com/alecthomas/kong"
)

// init initializes sources sink mappers
func init() {
	sinkMappers = append(sinkMappers, mapper(string2.Sink{}, string2.New)...)
	sinkMappers = append(sinkMappers, mapper(json.Sink{}, json.New)...)
	sinkMappers = append(sinkMappers, mapper(pretty.Sink{}, pretty.New)...)
}

// configuredSink returns sources sink.Sink instances defined by the cli.
func configuredSink() *interface{} {
	var i interface{}
	switch {
	case cli.Log.Simple != nil:
		i = cli.Log.Simple
	case cli.Log.JSON != nil:
		i = cli.Log.JSON
	case cli.Log.Pretty != nil:
		i = cli.Log.Pretty
	default:
		return nil
	}
	return &i
}

var (
	// sinkMappers contains sink-specific converters from string to a sink.Sink instance.
	sinkMappers []kong.Option

	// cliSinkOptions acts an adapter of sinkCLI to sink.Options.
	cliSinkOptions = sinkOptionsDecorator{}
)

type (
	// sinkOptionsDecorator is a trivial wrapper around cli exposing its values as a sink.Options.
	sinkOptionsDecorator struct{}
)

// Wrap ...
func (s sinkOptionsDecorator) Wrap() uint {
	return cli.Log.Wrap
}

// Emphasis ...
func (s sinkOptionsDecorator) Emphasis() bool {
	return cli.Log.Emphasis
}
