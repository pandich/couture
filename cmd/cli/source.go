package cli

import (
	"couture/internal/pkg/source"
	"errors"
	"github.com/alecthomas/kong"
)

//sourceCLI contains source-specific cli args.
var sourceCLI struct {
	Fake *source.Fake `group:"Input" name:"fake" hidden:"true"`
}

func init() {
	sourceMappers = append(sourceMappers, mapper(source.Fake{}, source.NewFake)...)
}

var (
	ErrNoSources = errors.New("at least one source must be specified")
)

//sourceMappers contains source-specific converters from string to a source.Source instance.
var sourceMappers []kong.Option

//Sources returns all source.Source instances defined by the cli.
func Sources() []interface{} {
	var sources []interface{}
	if sourceCLI.Fake != nil {
		sources = append(sources, *sourceCLI.Fake)
	}
	return sources
}

type sourcesValidator struct{}

func (v sourcesValidator) Validate() error {
	if len(Sources()) == 0 {
		return ErrNoSources
	}
	return nil
}
