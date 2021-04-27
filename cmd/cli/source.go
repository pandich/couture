package cli

import (
	"couture/internal/pkg/source"
	"github.com/alecthomas/kong"
)

var (
	//sourceCLI contains source-specific cli args.
	sourceCLI struct {
		Fakes []source.Fake `group:"Input Options" help:"A filename, URI, or pattern." name:"file" aliases:"files"`
	}

	//sourceMappers contains source-specific converters from string to a source.Source instance.
	sourceMappers = []kong.Option{
		mapper(many(source.Fake{}), source.NewFake),
	}
)

//Sources returns all source.Source instances defined by the cli.
func Sources() []interface{} {
	var sources []interface{}
	for _, src := range sourceCLI.Fakes {
		sources = append(sources, src)
	}
	return sources
}
