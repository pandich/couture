package manager

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/aws/cloudformation"
	"couture/internal/pkg/source/aws/cloudwatch"
	"couture/internal/pkg/source/elasticsearch"
	"couture/internal/pkg/source/fake"
	"couture/internal/pkg/source/pipe/local"
	"couture/internal/pkg/source/pipe/ssh"
	errors2 "github.com/pkg/errors"
)

// AvailableSources is a list of sourceMetadata sourceMetadata.
var AvailableSources = []source.Metadata{
	fake.Metadata(),
	cloudwatch.Metadata(),
	cloudformation.Metadata(),
	elasticsearch.Metadata(),
	local.Metadata(),
	ssh.Metadata(),
}

// GetSource ...
func GetSource(sourceURL model.SourceURL) ([]interface{}, []error) {
	var sources []interface{}
	var violations []error
	metadata := getSourceMetadata(sourceURL)
	if metadata != nil {
		configuredSource, err := metadata.Creator(sourceURL)
		if err != nil {
			violations = append(violations, err)
		} else {
			sources = append(sources, *configuredSource)
		}
	} else {
		violations = append(violations, errors2.Errorf("invalid source URL: %+v\n", sourceURL))
	}
	return sources, violations
}

// getSourceMetadata ...
func getSourceMetadata(sourceURL model.SourceURL) *source.Metadata {
	for _, metadata := range AvailableSources {
		if metadata.CanHandle(sourceURL) {
			return &metadata
		}
	}
	return nil
}

func (mgr publishingManager) shouldInclude(evt source.Event) bool {
	if !evt.Level.IsAtLeast(mgr.options.level) {
		return false
	}
	return evt.Message.Matches(mgr.options.includeFilters, mgr.options.excludeFilters)
}
