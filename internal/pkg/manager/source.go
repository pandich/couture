package manager

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/aws/cloudformation"
	"couture/internal/pkg/source/aws/cloudwatch"
	"couture/internal/pkg/source/elasticsearch"
	"couture/internal/pkg/source/fake"
	local2 "couture/internal/pkg/source/pipe/local"
	ssh2 "couture/internal/pkg/source/pipe/ssh"
	errors2 "github.com/pkg/errors"
	"sync"
)

// AvailableSources is a list of sourceMetadata sourceMetadata.
var AvailableSources = []source.Metadata{
	fake.Metadata(),
	cloudwatch.Metadata(),
	cloudformation.Metadata(),
	elasticsearch.Metadata(),
	local2.Metadata(),
	ssh2.Metadata(),
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

type internalSource struct{}

// ID ...
func (i internalSource) ID() string {
	return "Manager"
}

// URL ...
func (i internalSource) URL() model.SourceURL {
	return model.SourceURL{}
}

// Start ...
func (i internalSource) Start(_ *sync.WaitGroup, _ func() bool, _ func(event model.Event)) error {
	return nil
}

var managerSource = internalSource{}
