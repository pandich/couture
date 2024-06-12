package cloudwatch

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/pandich/couture/event"
	"github.com/pandich/couture/source"
	"github.com/pandich/couture/source/aws"
	errors2 "github.com/pkg/errors"
	"go.uber.org/ratelimit"
	"path"
	"reflect"
	"strings"
	"sync"
	"time"
)

// Metadata ...
func Metadata() source.Metadata {
	var exampleURLs []string
	for _, scheme := range []string{scheme, schemeAliasShort, schemeAliasFriendly} {
		exampleURLs = append(
			exampleURLs,
			fmt.Sprintf(
				"%s://<cloudwatch-log-path>?profile=<profile>&region=<region>&lookbackTime=<interval|date>",
				scheme,
			),
		)
	}
	exampleURLs = append(
		exampleURLs,
		"cloudwatch-lambda://<lambda-name>?profile=<profile>&region=<region>&lookbackTime=<interval|date>",
		"logs-appsync://<api-id>?profile=<profile>&region=<region>&lookbackTime=<interval|date>",
		"cw-rdsc://<cluster-name>?profile=<profile>&region=<region>&lookbackTime=<interval|date>",
		"cw-rdsi://<instance-name>?profile=<profile>&region=<region>&lookbackTime=<interval|date>",
	)
	return source.Metadata{
		Name: "AWS CloudWatch",
		Type: reflect.TypeOf(cloudwatchSource{}),
		CanHandle: func(url event.SourceURL) bool {
			_, ok := map[string]bool{
				scheme:                 true,
				schemeAliasShort:       true,
				schemeAliasFriendly:    true,
				schemeAliasAppSync:     true,
				schemeAliasCodeBuild:   true,
				schemeAliasLambda:      true,
				schemeAliasRDS:         true,
				schemeAliasRDSCluster:  true,
				schemeAliasRDSInstance: true,
			}[url.Scheme]
			return ok
		},
		Creator:     source.Single(newFromURL),
		ExampleURLs: exampleURLs,
	}
}

const (
	scheme                 = "cloudwatch"
	schemeAliasLambda      = "lambda"
	schemeAliasAppSync     = "appsync"
	schemeAliasCodeBuild   = "codebuild"
	schemeAliasRDS         = "rds"
	schemeAliasRDSInstance = "rdsi"
	schemeAliasRDSCluster  = "rdsc"
	schemeAliasShort       = "cw"
	schemeAliasFriendly    = "logs"
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
	nextToken             *string
	recentEvents          map[string]bool
	cloudWatchRateLimiter ratelimit.Limiter
}

// newFromURL Cloudwatch source.
func newFromURL(since *time.Time, sourceURL event.SourceURL) (*source.Source, error) {
	normalizeURL(&sourceURL)
	awsSource, err := aws.New('☂', &sourceURL)
	if err != nil {
		return nil, errors2.Wrapf(err, "bad CloudWatch URL: %+v\n", sourceURL)
	}
	return New(awsSource, since, sourceURL.Path), nil
}

// New makes a new AWS base source.
func New(
	awsSource aws.Source,
	lookbackTime *time.Time,
	logGroupName string,
) *source.Source {
	const maxRequestsPerMinute = 20
	src := cloudwatchSource{
		Source:       awsSource,
		lookbackTime: lookbackTime,
		logGroupName: logGroupName,
		logs:         cloudwatchlogs.NewFromConfig(awsSource.Config()),
		cloudWatchRateLimiter: ratelimit.New(
			maxRequestsPerMinute,
			ratelimit.Per(time.Minute),
			ratelimit.WithSlack(maxRequestsPerMinute),
		),
	}
	var p source.Source = &src
	return &p
}

// normalizeURL take the sourceURL and expands any syntactic sugar.
func normalizeURL(sourceURL *event.SourceURL) {
	if sourceURL.Path == "" && sourceURL.Host != "" {
		sourceURL.Path = sourceURL.Host
		sourceURL.Host = ""
	}
	switch sourceURL.Scheme {
	case scheme, schemeAliasShort, schemeAliasFriendly:
		sourceURL.Scheme = scheme
	case schemeAliasLambda:
		sourceURL.Scheme = scheme
		sourceURL.Path = path.Join("/aws/lambda", sourceURL.Path)
	case schemeAliasAppSync:
		sourceURL.Scheme = scheme
		sourceURL.Path = path.Join("/aws/appsync/apis", sourceURL.Path)
	case schemeAliasCodeBuild:
		sourceURL.Scheme = scheme
		sourceURL.Path = path.Join("/aws/codebuild/projects", sourceURL.Path)
	case schemeAliasRDS:
		sourceURL.Scheme = scheme
		sourceURL.Path = path.Join("/aws/rds", sourceURL.Path)
	case schemeAliasRDSCluster:
		sourceURL.Scheme = scheme
		sourceURL.Path = path.Join("/aws/rds/cluster", sourceURL.Path)
	case schemeAliasRDSInstance:
		sourceURL.Scheme = scheme
		sourceURL.Path = path.Join("/aws/rds/instance", sourceURL.Path)
	}
}

// Start ...
// Start ...
func (src *cloudwatchSource) Start(
	wg *sync.WaitGroup,
	running func() bool,
	srcChan chan source.Event,
	_ chan event.SinkEvent,
	errChan chan source.Error,
) error {
	var startTime *int64
	if src.lookbackTime != nil {
		i := src.lookbackTime.Unix()
		startTime = &i
	}
	_, err := src.logs.GetLogGroupFields(
		context.TODO(),
		&cloudwatchlogs.GetLogGroupFieldsInput{LogGroupName: &src.logGroupName},
	)
	if err != nil {
		return err
	}

	go func() {
		defer wg.Done()
		src.recentEvents = make(map[string]bool)
		for running() {
			src.cloudWatchRateLimiter.Take()
			logEvents, err := src.logs.FilterLogEvents(
				context.TODO(), &cloudwatchlogs.FilterLogEventsInput{
					LogGroupName: &src.logGroupName,
					NextToken:    src.nextToken,
					StartTime:    startTime,
				},
			)
			if err != nil {
				if strings.Contains(err.Error(), "ResourceNotFoundException") {
					errChan <- source.Error{SourceURL: src.URL(), Error: fmt.Errorf("log group not found: %s", src.logGroupName)}
				} else {
					errChan <- source.Error{SourceURL: src.URL(), Error: err}
				}
				continue
			}

			newEvents := false
			for _, logEvent := range logEvents.Events {
				if logEvent.Message != nil {
					_, found := src.recentEvents[*logEvent.EventId]
					if !found && (src.lookbackTime == nil || logEvent.Timestamp == nil || *logEvent.Timestamp > src.lookbackTime.UnixMilli()) {
						newEvents = true
						srcChan <- source.Event{Source: src, Event: *logEvent.Message}
					}
				}
			}

			if newEvents {
				src.recentEvents = make(map[string]bool)
				for _, logEvent := range logEvents.Events {
					id := *logEvent.EventId
					src.recentEvents[id] = true
				}
			}

			if logEvents.NextToken == nil || *logEvents.NextToken == "" {
				time.Sleep(1 * time.Second)
				continue
			}

			src.nextToken = logEvents.NextToken
		}
	}()
	return nil
}
