package cloudwatch

import (
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/aws"
	"couture/pkg/model"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	errors2 "github.com/pkg/errors"
	"reflect"
	"time"
)

// Metadata ...
func Metadata() source.Metadata {
	var exampleURLs []string
	for _, scheme := range []string{scheme, schemeAliasShort, schemeAliasFriendly} {
		exampleURLs = append(
			exampleURLs,
			fmt.Sprintf("%s://<cloudwatch-log-path>?profile=<profile>&region=<region>&lookbackTime=<interval|date>", scheme),
		)
	}
	exampleURLs = append(exampleURLs, "lambda://<lambda-name>?profile=<profile>&region=<region>&lookbackTime=<interval|date>")
	return source.Metadata{
		Type: reflect.TypeOf(Source{}),
		CanHandle: func(url model.SourceURL) bool {
			_, ok := map[string]bool{
				scheme:              true,
				schemeAliasShort:    true,
				schemeAliasFriendly: true,
			}[url.Scheme]
			return ok
		},
		Creator:     create,
		ExampleURLs: exampleURLs,
	}
}

const (
	// lookbackTimeFlag is the url.URL query parameter (optionally) defining how far to look back.
	lookbackTimeFlag = "since"
)

const (
	scheme              = "cloudwatch"
	schemeAlias         = "lambda"
	schemeAliasShort    = "cw"
	schemeAliasFriendly = "logs"
)

// Source a Cloudwatch log poller.
type Source struct {
	aws.Source
	// lookbackTime is how far back to look for log events.
	lookbackTime *time.Time
	// logs is the CloudWatch logs client.
	logs *cloudwatchlogs.Client
	// logGroupName is the name of the log group.
	logGroupName string
	// nextToken for calls to get log events.
	nextToken *string
}

// create Cloudwatch source casted to an *interface{}.
func create(sourceURL model.SourceURL) (*interface{}, error) {
	src, err := New(sourceURL)
	if err != nil {
		return nil, err
	}
	var i interface{} = src
	return &i, err
}

// New Cloudwatch source.
func New(sourceURL model.SourceURL) (*Source, error) {
	normalizeURL(sourceURL)
	awsSource, err := aws.New(&sourceURL)
	if err != nil {
		return nil, errors2.Wrapf(err, "bad CloudWatch URL: %+v", sourceURL)
	}
	lookbackTime, err := sourceURL.Since(lookbackTimeFlag)
	if err != nil {
		return nil, err
	}
	return &Source{
		Source:       *awsSource,
		lookbackTime: lookbackTime,
		logGroupName: sourceURL.Path,
		logs:         cloudwatchlogs.NewFromConfig(awsSource.Config()),
	}, nil
}

// normalizeURL take the sourceURL and expands any syntactic sugar.
func normalizeURL(sourceURL model.SourceURL) {
	switch {
	case sourceURL.Scheme == schemeAlias:
		sourceURL.Scheme = scheme
		sourceURL.Path = "/aws/lambda" + sourceURL.Path
	case sourceURL.Scheme == schemeAliasShort:
	case sourceURL.Scheme == schemeAliasFriendly:
		sourceURL.Scheme = scheme
	}
}

// Poll for more events.
func (source Source) Poll() ([]model.Event, error) {
	var events []model.Event

	// filter with source.lookbackTime
	source.nextToken = nil // TODO implement polling of CloudWatch events.

	return events, nil
}
