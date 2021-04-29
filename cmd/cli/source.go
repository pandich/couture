package cli

import (
	"couture/internal/pkg/source"
	"github.com/alecthomas/kong"
)

var (
	//sourceCLI contains source-specific cli args.
	sourceCLI struct {
		Fake *source.Fake `group:"Input Options" name:"fake" hidden:"true"`
	}

	//sourceMappers contains source-specific converters from string to a source.Source instance.
	sourceMappers = []kong.Option{
		mapper(one(source.Fake{}), source.NewFake),
	}
)

//Sources returns all source.Source instances defined by the cli.
func Sources() []interface{} {
	var sources []interface{}
	if sourceCLI.Fake != nil {
		sources = append(sources, *sourceCLI.Fake)
	}
	return sources
}
