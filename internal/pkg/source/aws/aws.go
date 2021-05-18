package aws

import (
	"context"
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"time"
)

const (
	// regionQueryFlag is the url.URL query parameter used to indicate the AWS region.
	regionQueryFlag = "region"
	// profileQueryFlag is the url.URL query parameter used to indicate the AWS profile.
	profileQueryFlag = "profile"
)

// Source ...
// Source represents an AWS entity in a specific region and profile.
type Source struct {
	*source.Polling
	// entity an arbitrary name whose meaning is implementation specific.
	entity string
	// config is the config for AWS clients.
	config aws.Config
	// region is the AWS region.
	region string
	// profile is the AWS profile.
	profile string
}

// New parses the baseSource.sourceURL into region, profile, and entity.
func New(sourceURL *model.SourceURL) (*Source, error) {
	sourceURL.Normalize()
	cfg, err := config.LoadDefaultConfig(context.TODO(), configOptions(sourceURL)...)
	if err != nil {
		return nil, err
	}
	return &Source{
		Polling: source.NewPollable(*sourceURL, time.Second),
		entity:  sourceURL.Path,
		config:  cfg,
	}, nil
}

func configOptions(sourceURL *model.SourceURL) []func(*config.LoadOptions) error {
	var options []func(*config.LoadOptions) error
	if region, ok := sourceURL.QueryKey(regionQueryFlag); ok {
		options = append(options, config.WithRegion(region))
	}
	if profile, ok := sourceURL.QueryKey(profileQueryFlag); ok {
		options = append(options, config.WithSharedConfigProfile(profile))
	}
	return options
}

// Config is the AWS configuration.
func (source Source) Config() aws.Config {
	return source.config
}
