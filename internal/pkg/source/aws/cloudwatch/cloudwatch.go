package cloudwatch

import (
	"context"
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/schema"
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/aws"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	errors2 "github.com/pkg/errors"
	"reflect"
	"sync"
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
		Name: "AWS CloudWatch",
		Type: reflect.TypeOf(cloudwatchSource{}),
		CanHandle: func(url model.SourceURL) bool {
			_, ok := map[string]bool{
				scheme:              true,
				schemeAliasShort:    true,
				schemeAliasFriendly: true,
				schemeAliasLambda:   true,
			}[url.Scheme]
			return ok
		},
		Creator:     newFromURL,
		ExampleURLs: exampleURLs,
	}
}

const (
	// lookbackTimeFlag is the url.URL query parameter (optionally) defining how far to look back.
	lookbackTimeFlag = "since"
)

const (
	scheme              = "cloudwatch"
	schemeAliasLambda   = "lambda"
	schemeAliasShort    = "cw"
	schemeAliasFriendly = "logs"
)

// cloudwatchSource a Cloudwatch log poller.
type cloudwatchSource struct {
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

// newFromURL Cloudwatch source.
func newFromURL(sourceURL model.SourceURL) (*source.Source, error) {
	normalizeURL(&sourceURL)
	awsSource, err := aws.New('☂', &sourceURL)
	if err != nil {
		return nil, errors2.Wrapf(err, "bad CloudWatch URL: %+v\n", sourceURL)
	}
	lookbackTime, err := sourceURL.Since(lookbackTimeFlag)
	if err != nil {
		return nil, err
	}
	return New(awsSource, lookbackTime, sourceURL.Path), nil
}

// New makes a new AWS base source.
func New(
	awsSource aws.Source,
	lookbackTime *time.Time,
	logGroupName string,
) *source.Source {
	src := cloudwatchSource{
		Source:       awsSource,
		lookbackTime: lookbackTime,
		logGroupName: logGroupName,
		logs:         cloudwatchlogs.NewFromConfig(awsSource.Config()),
	}
	var p source.Source = &src
	return &p
}

// normalizeURL take the sourceURL and expands any syntactic sugar.
func normalizeURL(sourceURL *model.SourceURL) {
	sourceURL.Normalize()
	switch {
	case sourceURL.Scheme == schemeAliasLambda:
		sourceURL.Scheme = scheme
		sourceURL.Path = "/aws/lambda" + sourceURL.Path
	case sourceURL.Scheme == schemeAliasShort:
	case sourceURL.Scheme == schemeAliasFriendly:
		sourceURL.Scheme = scheme
	}
}

// Start ...
func (src *cloudwatchSource) Start(
	wg *sync.WaitGroup,
	running func() bool,
	srcChan chan source.Event,
	_ chan model.SinkEvent,
	errChan chan source.Error,
) error {
	var startTime *int64
	if src.lookbackTime != nil {
		i := src.lookbackTime.Unix()
		startTime = &i
	}

	go func() {
		defer wg.Done()
		for running() {
			logEvents, err := src.logs.FilterLogEvents(context.TODO(), &cloudwatchlogs.FilterLogEventsInput{
				LogGroupName: &src.logGroupName,
				NextToken:    src.nextToken,
				StartTime:    startTime,
			})
			if err != nil {
				errChan <- source.Error{SourceURL: src.URL(), Error: err}
				continue
			}
			src.nextToken = logEvents.NextToken
			for _, logEvent := range logEvents.Events {
				if logEvent.Message != nil {
					if err != nil {
						errChan <- source.Error{SourceURL: src.URL(), Error: err}
					} else {
						srcChan <- source.Event{Source: src, Event: *logEvent.Message, Schema: schema.Logstash}
					}
				}
			}
		}
	}()
	return nil
}
