package cloudformation

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/gagglepanda/couture/model"
	"github.com/gagglepanda/couture/model/level"
	"github.com/gagglepanda/couture/source"
	"github.com/gagglepanda/couture/source/aws"
	"github.com/gagglepanda/couture/source/aws/cloudwatch"
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
		exampleURLs = append(exampleURLs,
			fmt.Sprintf("%s://<stack-name>?profile=<profile>&region=<region>&lookbackTime=<interval|date>&events(=<true|false>)", scheme),
		)
	}
	return source.Metadata{
		Name: "AWS CloudFormation",
		Type: reflect.TypeOf(cloudFormationSource{}),
		CanHandle: func(url model.SourceURL) bool {
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

const (
	// includeStackEventsFlag is the url.URL query parameter indicating whether stack events should be included.
	includeStackEventsFlag = "events"
)

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
		// children represents all child sources added during stack-resource discovery.
		// For example: a lambda's log group's cloudwatch.cloudwatchSource.
		children []*source.Source
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
func newSource(since *time.Time, sourceURL model.SourceURL) (*source.Source, error) {
	const maxRequestsPerMinute = 20

	normalizeURL(&sourceURL)
	stackName := sourceURL.Path[1:]
	awsSource, err := aws.New('☁', &sourceURL)
	if err != nil {
		return nil, errors2.Wrapf(err, "bad CloudFormation URL: %+v", sourceURL)
	}

	cf := cloudformation.NewFromConfig(awsSource.Config())

	var children []*source.Source
	// add lambda functions
	lambdaResources, err := discoverLambdaResources(cf, stackName)
	if err != nil {
		return nil, err
	}
	for _, lambdaResource := range lambdaResources {
		logGroupName := path.Join(aws.LambdaLogGroupPrefix, *lambdaResource.PhysicalResourceId)
		children = append(children, cloudwatch.New(awsSource, since, logGroupName))
	}

	var src source.Source = &cloudFormationSource{
		Source:               awsSource,
		lookbackTime:         since,
		includeStackEvents:   sourceURL.QueryFlag(includeStackEventsFlag),
		children:             children,
		cf:                   cf,
		stackName:            stackName,
		stackEventsNextToken: nil,
		stackEventTimes:      map[int64]bool{},
		stackEventRateLimiter: ratelimit.New(
			maxRequestsPerMinute,
			ratelimit.Per(time.Minute),
			ratelimit.WithSlack(maxRequestsPerMinute),
		),
	}
	return &src, nil
}

// normalizeURL take the sourceURL and expands any syntactic sugar.
func normalizeURL(sourceURL *model.SourceURL) {
	sourceURL.Normalize()
	switch {
	case sourceURL.Scheme == schemeAliasShort:
	case sourceURL.Scheme == schemeAliasFriendly:
		sourceURL.Scheme = scheme
	}
}

// Start ...
func (src *cloudFormationSource) Start(
	wg *sync.WaitGroup,
	running func() bool,
	srcChan chan source.Event,
	snkChan chan model.SinkEvent,
	errChan chan source.Error,
) error {
	for _, child := range src.children {
		err := (*child).Start(wg, running, srcChan, snkChan, errChan)
		if err != nil {
			return err
		}
	}
	if src.includeStackEvents {
		go func() {
			for running() {
				stackEvents, err := src.getStackEvents()
				if err != nil {
					errChan <- source.Error{SourceURL: src.URL(), Error: err}
				}
				for _, evt := range stackEvents {
					snkChan <- evt
				}
			}
		}()
	}
	return nil
}

// getChildEvents retrieves CloudFormation events for this stack.
func (src *cloudFormationSource) getStackEvents() ([]model.SinkEvent, error) {
	src.stackEventRateLimiter.Take()
	stackEvents, err := src.cf.DescribeStackEvents(context.TODO(), &cloudformation.DescribeStackEventsInput{
		NextToken: src.stackEventsNextToken,
		StackName: &src.stackName,
	})
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
) ([]model.SinkEvent, error) {
	var events []model.SinkEvent
	for _, stackEvent := range stackEvents.StackEvents {
		event := src.stackEventToModelEvent(stackEvent)
		events = append(events, model.SinkEvent{Event: event, SourceURL: src.URL()})
	}
	return events, nil
}

func (src *cloudFormationSource) stackEventToModelEvent(stackEvent types.StackEvent) model.Event {
	var message model.Message
	var evtError model.Error

	lvl := logLevelByResourceStatus[stackEvent.ResourceStatus]
	if lvl == level.Error {
		evtError = model.Error(stackEvent.ResourceStatus)
	} else {
		message = model.Message(stackEvent.ResourceStatus)
	}
	var entity = model.Entity(*stackEvent.StackName)
	if s := stackEvent.PhysicalResourceId; s != nil && *s != "" {
		entity = model.Entity(*s)
	}
	event := model.Event{
		Timestamp:   model.Timestamp(*stackEvent.Timestamp),
		Application: model.Application(*stackEvent.ResourceType),
		Context:     model.Context(*stackEvent.EventId),
		Entity:      entity,
		Action:      model.Action(""),
		Line:        model.NoLineNumber,
		Level:       lvl,
		Message:     message,
		Error:       evtError,
	}
	return event
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
			lambdaFunctions = append(lambdaFunctions, resource)
		case "AWS::CloudFormation::Stack":
			subResources, err := discoverLambdaResources(cf, *resource.LogicalResourceId)
			if err != nil {
				return nil, err
			}
			lambdaFunctions = append(lambdaFunctions, subResources...)
		}
	}
	return lambdaFunctions, nil
}
