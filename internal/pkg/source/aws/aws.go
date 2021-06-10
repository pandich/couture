package aws

import (
	"context"
	"couture/internal/pkg/model"
	"couture/internal/pkg/source"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// LambdaLogGroupPrefix ...
const LambdaLogGroupPrefix = "/aws/lambda"

const (
	// regionQueryFlag is the url.URL query parameter used to indicate the AWS region.
	regionQueryFlag = "region"
	// profileQueryFlag is the url.URL query parameter used to indicate the AWS profile.
	profileQueryFlag = "profile"
)

// Source represents an AWS entity in a specific region and profile.
type Source struct {
	source.BaseSource
	// entity an arbitrary name whose meaning is implementation specific.
	entity string
	// config is the config for AWS clients.
	config aws.Config
}

// New parses the baseSource.sourceURL into region, profile, and entity.
func New(sigil rune, sourceURL *model.SourceURL) (Source, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), configOptions(sourceURL)...)
	if err != nil {
		return Source{}, err
	}
	return Source{
		BaseSource: source.New(sigil, *sourceURL),
		entity:     sourceURL.Path,
		config:     cfg,
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
