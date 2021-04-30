package cli

import (
	"couture/internal/pkg/source"
	"github.com/alecthomas/kong"
)

func init() {
	sourceMappers = append(sourceMappers, mapper(source.Fake{}, source.NewFake)...)
}

var (
	//sourceCLI contains source-specific cli args.
	sourceCLI struct {
		Fake *source.Fake `group:"Input" name:"fake" hidden:"true"`
	}
	//sourceMappers contains source-specific converters from string to a source.Source instance.
	sourceMappers []kong.Option
)

//Sources returns all source.Source instances defined by the cli.
func Sources() []interface{} {
	var sources []interface{}
	if sourceCLI.Fake != nil {
		sources = append(sources, *sourceCLI.Fake)
	}
	return sources
}
