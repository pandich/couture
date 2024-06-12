package cloudformation

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/pandich/couture/event"
	"github.com/pandich/couture/event/level"
	"github.com/pandich/couture/global"
	"github.com/pandich/couture/source"
	"github.com/pandich/couture/source/aws"
	"github.com/pandich/couture/source/aws/cloudwatch"
	errors2 "github.com/pkg/errors"
	"go.uber.org/ratelimit"
	"path"
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
			fmt.Sprintf(
				"%s://<stack-name>?profile=<profile>&region=<region>&lookbackTime=<interval|date>&events(=<true|false>)",
				scheme,
			),
		)
	}
	return source.Metadata{
		Name: "AWS CloudFormation",
		Type: reflect.TypeOf(cloudFormationSource{}),
		CanHandle: func(url event.SourceURL) bool {
			_, ok := map[string]bool{
				scheme:              true,
				schemeAliasShort:    true,
				schemeAliasFriendly: true,
			}[url.Scheme]
			return ok
		},
		Creator:     newSource,
		ExampleURLs: exampleURLs,
	}
}

// logLevelByResourceStatus maps each possible resource status to a log level.
var logLevelByResourceStatus = map[types.ResourceStatus]level.Level{
	types.ResourceStatusCreateInProgress: level.Info,
	types.ResourceStatusCreateFailed:     level.Error,
	types.ResourceStatusCreateComplete:   level.Info,

	types.ResourceStatusDeleteInProgress: level.Info,
	types.ResourceStatusDeleteFailed:     level.Error,
	types.ResourceStatusDeleteComplete:   level.Info,
	types.ResourceStatusDeleteSkipped:    level.Warn,

	types.ResourceStatusUpdateInProgress: level.Info,
	types.ResourceStatusUpdateFailed:     level.Error,
	types.ResourceStatusUpdateComplete:   level.Info,

	types.ResourceStatusImportFailed:     level.Error,
	types.ResourceStatusImportComplete:   level.Info,
	types.ResourceStatusImportInProgress: level.Info,

	types.ResourceStatusImportRollbackInProgress: level.Warn,
	types.ResourceStatusImportRollbackFailed:     level.Error,
	types.ResourceStatusImportRollbackComplete:   level.Warn,
}

// URL schemes supported
const (
	scheme              = "cloudformation"
	schemeAliasShort    = "cf"
	schemeAliasFriendly = "stack"
)

type (
	// cloudFormationSource a CloudFormation stack event, and stack resource log watcher.
	// Stack resources are recursively searched to discover Cloudwatch log groups related to the stack.
	// Currently Supported Resources:
	//		Lambda Functions
	cloudFormationSource struct {
		aws.Source
		// lookbackTime is how far back to look for log events.
		lookbackTime *time.Time
		// includeStackEvents specifies whether to include stack events in the log.
		includeStackEvents bool
		// cf represents the clo
		cf        *cloudformation.Client
		stackName string
		// stackEventsNextToken is used to keep track of the last call to fetch stack events.
		stackEventsNextToken  *string
		stackEventTimes       map[int64]bool
		stackEventRateLimiter ratelimit.Limiter
	}
)

// newSource CloudFormation source.
func newSource(since *time.Time, sourceURL event.SourceURL) ([]source.Source, error) {
	normalizeURL(&sourceURL)
	stackName := sourceURL.Path[1:]
	awsSource, err := aws.New('☁', &sourceURL)
	if err != nil {
		return nil, errors2.Wrapf(err, "bad CloudFormation URL: %+v", sourceURL)
	}

	cf := cloudformation.NewFromConfig(awsSource.Config())

	var sources []source.Source
	// add lambda functions
	lambdaResources, err := discoverLambdaResources(cf, stackName)
	if err != nil {
		return nil, err
	}
	for _, lambdaResource := range lambdaResources {
		logGroupName := path.Join("/aws/lambda", *lambdaResource.PhysicalResourceId)
		src, err := aws.New(
			'𝞴', &event.SourceURL{
				Scheme:      "cloudwatch",
				Host:        sourceURL.Host,
				Path:        logGroupName,
				RawFragment: sourceURL.RawFragment,
				RawQuery:    sourceURL.RawQuery,
			},
		)
		if err != nil {
			return nil, err
		}
		url := src.URL()
		if val, found := url.QueryKey("lookbackTime"); found {
			dur, err := time.ParseDuration(val)
			if err == nil {
				local := time.Now().Add(-dur)
				since = &local
			}
		}
		child := cloudwatch.New(src, since, logGroupName)
		if child != nil {
			sources = append(sources, *child)
		}
	}

	return sources, nil
}

// normalizeURL take the sourceURL and expands any syntactic sugar.
func normalizeURL(sourceURL *event.SourceURL) {
	sourceURL.Normalize()
	switch {
	case sourceURL.Scheme == schemeAliasShort:
	case sourceURL.Scheme == schemeAliasFriendly:
		sourceURL.Scheme = scheme
	}
}

// Start ...
func (src *cloudFormationSource) Start(
	_ *sync.WaitGroup,
	running func() bool,
	_ chan source.Event,
	snkChan chan event.SinkEvent,
	errChan chan source.Error,
) error {
	if src.includeStackEvents {
		go func() {
			for running() {
				stackEvents, err := src.getStackEvents()
				if err != nil {
					errChan <- source.Error{SourceURL: src.URL(), Error: err}
				}
				for _, evt := range stackEvents {
					if src.lookbackTime == nil || time.Time(evt.Timestamp).After(*src.lookbackTime) {
						snkChan <- evt
					}
				}
			}
		}()
	}
	return nil
}

// getChildEvents retrieves CloudFormation events for this stack.
func (src *cloudFormationSource) getStackEvents() ([]event.SinkEvent, error) {
	src.stackEventRateLimiter.Take()
	stackEvents, err := src.cf.DescribeStackEvents(
		context.TODO(), &cloudformation.DescribeStackEventsInput{
			NextToken: src.stackEventsNextToken,
			StackName: &src.stackName,
		},
	)
	if err != nil {
		return nil, err
	}
	src.stackEventsNextToken = stackEvents.NextToken
	if len(stackEvents.StackEvents) == 0 {
		return nil, nil
	}
	timestamp := stackEvents.StackEvents[0].Timestamp.Unix()
	if _, ok := src.stackEventTimes[timestamp]; ok {
		return nil, nil
	}
	src.stackEventTimes[timestamp] = true
	return src.stackEventsToModelEvents(stackEvents)
}

func (src *cloudFormationSource) stackEventsToModelEvents(
	stackEvents *cloudformation.DescribeStackEventsOutput,
) ([]event.SinkEvent, error) {
	var events []event.SinkEvent
	for _, stackEvent := range stackEvents.StackEvents {
		evt := src.stackEventToModelEvent(stackEvent)
		events = append(events, event.SinkEvent{Event: evt, SourceURL: src.URL()})
	}
	return events, nil
}

func (src *cloudFormationSource) stackEventToModelEvent(stackEvent types.StackEvent) event.Event {
	var message event.Message
	var evtError event.Error

	lvl := logLevelByResourceStatus[stackEvent.ResourceStatus]
	if lvl == level.Error {
		evtError = event.Error(stackEvent.ResourceStatus)
	} else {
		message = event.Message(stackEvent.ResourceStatus)
	}
	var entity = event.Entity(*stackEvent.StackName)
	if s := stackEvent.PhysicalResourceId; s != nil && *s != "" {
		entity = event.Entity(*s)
	}
	evt := event.Event{
		Timestamp:   event.Timestamp(*stackEvent.Timestamp),
		Application: event.Application(*stackEvent.ResourceType),
		Context:     event.Context(*stackEvent.EventId),
		Entity:      entity,
		Action:      event.Action(""),
		Line:        event.NoLineNumber,
		Level:       lvl,
		Message:     message,
		Error:       evtError,
	}
	return evt
}

// discoverLambdaResources discovers all lambdas under the stack or its child stacks.
func discoverLambdaResources(cf *cloudformation.Client, stackName string) ([]types.StackResource, error) {
	resources, err := cf.DescribeStackResources(
		context.TODO(),
		&cloudformation.DescribeStackResourcesInput{StackName: &stackName},
	)
	if err != nil {
		return nil, err
	}
	var lambdaFunctions []types.StackResource
	for _, resource := range resources.StackResources {
		switch *resource.ResourceType {
		case "AWS::Lambda::Function":
			global.DiscoveryBus <- global.Resource{
				Kind: "aws::lambda::function",
				Name: *resource.PhysicalResourceId,
			}
			lambdaFunctions = append(lambdaFunctions, resource)
		case "AWS::CloudFormation::Stack":
			subResources, err := discoverLambdaResources(cf, *resource.LogicalResourceId)
			if err != nil {
				return nil, err
			}

			for _, subResource := range subResources {
				global.DiscoveryBus <- global.Resource{
					Kind: "aws::cloudformation=>aws::lambda::function",
					Name: *subResource.PhysicalResourceId,
				}
			}

			lambdaFunctions = append(lambdaFunctions, subResources...)
		}
	}
	return lambdaFunctions, nil
}
