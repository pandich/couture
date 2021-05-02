package cli

import (
	"couture/internal/pkg/source"
	"fmt"
	"github.com/alecthomas/kong"
	"gopkg.in/multierror.v1"
	"net/url"
)

var (
	//sourceCLI contains source-specific cli args.
	sourceCLI struct {
		Sources []url.URL `arg:"true" group:"source" name:"sources"`
	}
)

//Sources returns all source.Source instances defined by the cli.
func Sources() []interface{} {
	var errors []error
	var sources []interface{}
	for _, srcUrl := range sourceCLI.Sources {
		var handled bool
		for _, src := range source.Available() {
			var err error
			var creator source.Creator
			creator, err = source.CreatorFor(src)
			if err != nil {
				errors = append(errors, err)
			} else {
				if src.CanHandle(srcUrl) {
					handled = true
					sources = append(sources, creator(srcUrl))
					break
				}
			}
		}
		if !handled {
			errors = append(errors, fmt.Errorf("unhandled sourcr URL: %v", srcUrl))
		}
	}
	if len(errors) > 0 {
		panic(multierror.New(errors))
	}
	return sources
}

//sourceMappers contains source-specific converters from string to a source.Source instance.
var sourceMappers []kong.Option
