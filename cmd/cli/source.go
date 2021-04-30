package cli

import (
	"couture/internal/pkg/source"
	"github.com/alecthomas/kong"
)

//sourceCLI contains source-specific cli args.
var sourceCLI struct {
	Fake *source.Fake `group:"Input Options" name:"fake" hidden:"true"`
	// Files []source.File `group:"Input Options" name:"file" short:"F"`
}

func init() {
	sourceMappers = append(sourceMappers, mapper(source.Fake{}, source.NewFake)...)
	// sourceMappers = append(sourceMappers, mapper(source.File{}, source.NewFile)...)
}

//Sources returns all source.Source instances defined by the cli.
func Sources() []interface{} {
	var sources []interface{}
	if sourceCLI.Fake != nil {
		sources = append(sources, *sourceCLI.Fake)
	}
	// sources = append(sources, sourceCLI.Files)
	return sources
}

//sourceMappers contains source-specific converters from string to a source.Source instance.
var sourceMappers []kong.Option
