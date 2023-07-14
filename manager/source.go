package manager

import (
	"github.com/gagglepanda/couture/event"
	"github.com/gagglepanda/couture/model"
	"github.com/gagglepanda/couture/source"
	"github.com/gagglepanda/couture/source/aws/cloudformation"
	"github.com/gagglepanda/couture/source/aws/cloudwatch"
	"github.com/gagglepanda/couture/source/aws/s3"
	"github.com/gagglepanda/couture/source/elasticsearch"
	"github.com/gagglepanda/couture/source/fake"
	"github.com/gagglepanda/couture/source/pipe/local"
	"github.com/gagglepanda/couture/source/pipe/ssh"
	errors2 "github.com/pkg/errors"
	"time"
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
func GetSource(since *time.Time, sourceURL event.SourceURL) ([]source.Source, []error) {
	if sourceURL.Scheme == "complete" {
		return nil, nil
	}
	var sources []source.Source
	var violations []error
	metadata := getSourceMetadata(sourceURL)
	if metadata != nil {
		configuredSource, err := metadata.Creator(since, sourceURL)
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

func getSourceMetadata(sourceURL event.SourceURL) *source.Metadata {
	for _, metadata := range AvailableSources {
		if metadata.CanHandle(sourceURL) {
			return &metadata
		}
	}
	return nil
}

func (mgr *busManager) filter(evt *event.Event) model.FilterKind {
	if !evt.Level.IsAtLeast(mgr.config.Level) {
		return model.Exclude
	}
	return evt.Message.Matches(&mgr.config.Filters)
}
