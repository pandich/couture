package cloudformation

import (
	"context"
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/aws"
	"couture/internal/pkg/source/aws/cloudwatch"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	errors2 "github.com/pkg/errors"
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
		Creator: func(sourceURL model.SourceURL) (*interface{}, error) {
			src, err := newSource(sourceURL)
			if err != nil {
				return nil, err
			}
			var i interface{} = src
			return &i, nil
		},
		ExampleURLs: exampleURLs,
	}
}

const (
	// lookbackTimeFlag is the url.URL query parameter (optionally) defining how far to look back.
	lookbackTimeFlag = "since"
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
		// includeStackEvents specifies whether or not to include stack events in the log.
		includeStackEvents bool
		// children represents all child sources added during stack-resource discovery.
		// For example: a lambda's log group's cloudwatch.cloudwatchSource.
		children []*source.Source
		// cf represents the clo
		cf        *cloudformation.Client
		stackName string
		// stackEventsNextToken is used to keep track of the last call to fetch stack events.
		stackEventsNextToken *string
	}
)

// newSource CloudFormation source.
func newSource(sourceURL model.SourceURL) (*cloudFormationSource, error) {
	normalizeURL(&sourceURL)
	stackName := sourceURL.Path
	awsSource, err := aws.New('☁', &sourceURL)
	if err != nil {
		return nil, errors2.Wrapf(err, "bad CloudFormation URL: %+v", sourceURL)
	}

	lookbackTime, err := sourceURL.Since(lookbackTimeFlag)
	if err != nil {
		return nil, err
	}

	cf := cloudformation.NewFromConfig(awsSource.Config())

	var children []*source.Source
	// add lambda functions
	lambdaResources, err := discoverLambdaResources(cf, stackName)
	if err != nil {
		return nil, err
	}
	for _, lambdaResource := range lambdaResources {
		logGroupName := fmt.Sprintf("/log/lambda/%s", *lambdaResource.PhysicalResourceId)
		children = append(children, cloudwatch.New(awsSource, lookbackTime, logGroupName))
	}

	return &cloudFormationSource{
		Source:               awsSource,
		lookbackTime:         lookbackTime,
		includeStackEvents:   sourceURL.QueryFlag(includeStackEventsFlag),
		children:             children,
		cf:                   cf,
		stackName:            stackName,
		stackEventsNextToken: nil,
	}, nil
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
	errChan chan source.Error,
) error {
	for _, child := range src.children {
		err := (*child).Start(wg, running, srcChan, errChan)
		if err != nil {
			return err
		}
	}
	if src.includeStackEvents {
		go func() {
			for running() {
				stackEvents, err := src.getStackEvents()
				if err != nil {
					errChan <- source.Error{Source: src, Error: err}
				}
				for _, stackEvent := range stackEvents {
					srcChan <- source.Event{Source: src, Event: stackEvent}
				}
			}
		}()
	}
	return nil
}

// getChildEvents retrieves CloudFormation events for this stack.
func (src cloudFormationSource) getStackEvents() ([]model.Event, error) {
	stackEvents, err := src.cf.DescribeStackEvents(context.TODO(), &cloudformation.DescribeStackEventsInput{
		NextToken: src.stackEventsNextToken,
		StackName: &src.stackName,
	})
	if err != nil {
		return nil, err
	}
	src.stackEventsNextToken = stackEvents.NextToken

	var events []model.Event
	for _, stackEvent := range stackEvents.StackEvents {
		if src.lookbackTime == nil || src.lookbackTime.Before(*stackEvent.Timestamp) {
			lvl := logLevelByResourceStatus[stackEvent.ResourceStatus]

			var exception *model.Exception
			if lvl == level.Error {
				exception = &model.Exception{
					StackTrace: model.StackTrace(*stackEvent.ResourceStatusReason),
				}
			}

			threadName := model.ThreadName("cloudformation")
			events = append(events, model.Event{
				Timestamp:  model.Timestamp(*stackEvent.Timestamp),
				ThreadName: &threadName,
				ClassName:  model.ClassName(*stackEvent.StackName),
				MethodName: model.MethodName(*stackEvent.EventId),
				LineNumber: model.NoLineNumber,
				Level:      lvl,
				Message:    model.Message(stackEvent.ResourceStatus),
				Exception:  exception,
			})
		}
	}
	return events, nil
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
