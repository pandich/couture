package manager

import (
	"github.com/pandich/couture/internal/pkg/model"
	"github.com/pandich/couture/internal/pkg/source"
	"github.com/pandich/couture/internal/pkg/source/aws/cloudformation"
	"github.com/pandich/couture/internal/pkg/source/aws/cloudwatch"
	"github.com/pandich/couture/internal/pkg/source/aws/s3"
	"github.com/pandich/couture/internal/pkg/source/elasticsearch"
	"github.com/pandich/couture/internal/pkg/source/fake"
	"github.com/pandich/couture/internal/pkg/source/pipe/local"
	"github.com/pandich/couture/internal/pkg/source/pipe/ssh"
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
func GetSource(since *time.Time, sourceURL model.SourceURL) ([]source.Source, []error) {
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

func getSourceMetadata(sourceURL model.SourceURL) *source.Metadata {
	for _, metadata := range AvailableSources {
		if metadata.CanHandle(sourceURL) {
			return &metadata
		}
	}
	return nil
}

func (mgr *busManager) filter(evt *model.Event) model.FilterKind {
	if !evt.Level.IsAtLeast(mgr.config.Level) {
		return model.Exclude
	}
	return evt.Message.Matches(&mgr.config.Filters)
}
