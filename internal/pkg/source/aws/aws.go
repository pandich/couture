package aws

import (
	"context"
	"couture/internal/pkg/source/polling"
	"couture/pkg/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"os"
	"time"
)

const (
	// RegionQueryFlag is the url.URL query parameter used to indicate the AWS region.
	RegionQueryFlag = "region"
	// ProfileQueryFlag is the url.URL query parameter used to indicate the AWS profile.
	ProfileQueryFlag = "profile"
)

// Source ...
type (
	// Source represents an AWS entity in a specific region and profile.
	Source struct {
		polling.Source
		// entity an arbitrary name whose meaning is implementation specific.
		entity string
		// config is the config for AWS clients.
		config aws.Config
		// region is the AWS region.
		region string
		// profile is the AWS profile.
		profile string
	}
)

// New parses the baseSource.sourceURL into region, profile, and entity.
func New(sourceURL *model.SourceURL) (*Source, error) {
	sourceURL.Normalize()

	region := extractRegion(sourceURL)
	profile := extractProfile(sourceURL)
	var loadOptions []func(*config.LoadOptions) error
	if region != "" {
		loadOptions = append(loadOptions, config.WithRegion(region))
	}
	if profile != "" {
		loadOptions = append(loadOptions, config.WithSharedConfigProfile(profile))
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), loadOptions...)

	if err != nil {
		return nil, err
	}
	return &Source{
		Source:  polling.New(*sourceURL, time.Second),
		entity:  sourceURL.Path,
		config:  cfg,
		region:  region,
		profile: profile,
	}, nil
}

// Config is the AWS configuration.
func (source Source) Config() aws.Config {
	return source.config
}

// Profile is the AWS profile being used.
func (source Source) Profile() string {
	return source.profile
}

// Region is the AWS region being used.
func (source Source) Region() string {
	return source.region
}

// extractRegion tries to get the region
func extractRegion(sourceURL *model.SourceURL) string {
	var region, ok = sourceURL.QueryKey(RegionQueryFlag)
	if ok {
		return region
	}
	if envRegion, ok := os.LookupEnv("AWS_REGION"); ok {
		return envRegion
	}

	if envRegion, ok := os.LookupEnv("AWS_DEFAULT_REGION"); ok {
		return envRegion
	}

	return ""
}

func extractProfile(sourceURL *model.SourceURL) string {
	var profile, ok = sourceURL.QueryKey(ProfileQueryFlag)
	if ok {
		return profile
	}
	if envProfile, ok := os.LookupEnv("AWS_PROFILE"); ok {
		return envProfile
	}

	if envProfile, ok := os.LookupEnv("AWS_DEFAULT_PROFILE"); ok {
		return envProfile
	}

	return ""
}
