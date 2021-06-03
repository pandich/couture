package manager

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/aws/cloudformation"
	"couture/internal/pkg/source/aws/cloudwatch"
	"couture/internal/pkg/source/aws/s3"
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
	s3.Metadata(),
	local.Metadata(),
	ssh.Metadata(),
}

// GetSource gets a source, if possible, for the specified sourceURL.
func GetSource(sourceURL model.SourceURL) ([]source.Source, []error) {
	if sourceURL.Scheme == "complete" {
		return nil, nil
	}
	var sources []source.Source
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

func getSourceMetadata(sourceURL model.SourceURL) *source.Metadata {
	for _, metadata := range AvailableSources {
		if metadata.CanHandle(sourceURL) {
			return &metadata
		}
	}
	return nil
}

func (mgr busManager) shouldInclude(evt *model.Event) bool {
	if !evt.Level.IsAtLeast(mgr.config.Level) {
		return false
	}
	return evt.Message.Matches(mgr.config.Filters)
}
