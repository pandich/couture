package cli

import (
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/aws/cloudformation"
	"couture/internal/pkg/source/aws/cloudwatch"
	"couture/internal/pkg/source/elasticsearch"
	"couture/internal/pkg/source/fake"
	"couture/pkg/model"
	"errors"
	errors2 "github.com/pkg/errors"
	"gopkg.in/multierror.v1"
)

var (
	errNoHandlerForURL = errors.New("unhandled src URL")

	// sources is a list of sources sources.
	sources = []source.Metadata{
		fake.Metadata(),
		cloudwatch.Metadata(),
		cloudformation.Metadata(),
		elasticsearch.Metadata(),
	}
)

// configuredSources returns sources source.Source instances defined by the cli.
func configuredSources() ([]interface{}, error) {
	var violations []error
	var configuredSources []interface{}
	for _, sourceArgs := range cli.Sources {
		sourceURL := model.SourceURL(sourceArgs)
		var handled bool
		for _, metadata := range sources {
			if !metadata.CanHandle(sourceURL) {
				continue
			}
			handled = true
			configuredSource, err := metadata.Creator(sourceURL)
			if err != nil {
				violations = append(violations, err)
			} else {
				configuredSources = append(configuredSources, *configuredSource)
			}
		}
		if !handled {
			violations = append(violations, errors2.WithMessagef(errNoHandlerForURL, "%+v", sourceURL))
		}
	}
	if len(violations) > 0 {
		return nil, multierror.New(violations)
	}
	return configuredSources, nil
}
