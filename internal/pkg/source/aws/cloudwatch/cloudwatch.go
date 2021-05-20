package cloudwatch

// TODO CloudWatch source

import (
	"context"
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/aws"
	"encoding/json"
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
		Creator: func(sourceURL model.SourceURL) (*interface{}, error) {
			src, err := newFromURL(sourceURL)
			if err != nil {
				return nil, err
			}
			var i interface{} = src
			return &i, err
		},
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
	*aws.Source
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
func newFromURL(sourceURL model.SourceURL) (*source.Pollable, error) {
	normalizeURL(&sourceURL)
	awsSource, err := aws.New(&sourceURL)
	if err != nil {
		return nil, errors2.Wrapf(err, "bad CloudWatch URL: %+v\n", sourceURL)
	}
	lookbackTime, err := sourceURL.Since(lookbackTimeFlag)
	if err != nil {
		return nil, err
	}
	return New(awsSource, lookbackTime, sourceURL.Path), nil
}

// New ...
func New(
	awsSource *aws.Source,
	lookbackTime *time.Time,
	logGroupName string,
) *source.Pollable {
	src := cloudwatchSource{
		Source:       awsSource,
		lookbackTime: lookbackTime,
		logGroupName: logGroupName,
		logs:         cloudwatchlogs.NewFromConfig(awsSource.Config()),
	}
	var p source.Pollable = &src
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

// Poll for more events.
func (src *cloudwatchSource) Poll() ([]model.Event, error) {
	var events []model.Event
	var startTime *int64
	if src.lookbackTime != nil {
		i := src.lookbackTime.Unix()
		startTime = &i
	}
	logEvents, err := src.logs.FilterLogEvents(context.TODO(), &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName: &src.logGroupName,
		NextToken:    src.nextToken,
		StartTime:    startTime,
	})
	if err != nil {
		return nil, errors2.Unwrap(err)
	}
	src.nextToken = logEvents.NextToken
	for _, logEvent := range logEvents.Events {
		if logEvent.Message != nil {
			var event = model.Event{}
			err := json.Unmarshal([]byte(*logEvent.Message), &event)
			if err != nil {
				var message model.Message
				if logEvent.Message != nil {
					message = model.Message(*logEvent.Message)
				}
				threadName := model.ThreadName(*logEvent.LogStreamName)
				events = append(events, model.Event{
					Timestamp:  model.Timestamp(time.Unix(*logEvent.Timestamp, 0)),
					Level:      level.Info,
					Message:    message,
					ThreadName: &threadName,
					ClassName:  model.ClassName(*logEvent.EventId),
					Exception:  nil,
				})
			} else {
				events = append(events, event)
			}
		}
	}
	return events, nil
}
