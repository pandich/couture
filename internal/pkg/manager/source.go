package manager

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/aws/cloudformation"
	"couture/internal/pkg/source/aws/cloudwatch"
	"couture/internal/pkg/source/elasticsearch"
	"couture/internal/pkg/source/fake"
	"couture/internal/pkg/source/ssh"
	"couture/internal/pkg/source/tail"
	errors2 "github.com/pkg/errors"
)

// SourceMetadata is a list of sourceMetadata sourceMetadata.
var SourceMetadata = []source.Metadata{
	fake.Metadata(),
	cloudwatch.Metadata(),
	cloudformation.Metadata(),
	elasticsearch.Metadata(),
	tail.Metadata(),
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
	for _, metadata := range SourceMetadata {
		if metadata.CanHandle(sourceURL) {
			return &metadata
		}
	}
	return nil
}
