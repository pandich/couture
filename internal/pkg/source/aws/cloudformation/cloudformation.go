package cloudformation

import (
	"context"
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/aws"
	"couture/internal/pkg/source/aws/cloudwatch"
	"couture/pkg/model"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	errors2 "github.com/pkg/errors"
	"reflect"
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
		Type: reflect.TypeOf(cloudFormationSource{}),
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
	// includeStackEventsFlag is the url.URL query parameter indicating whether stack events should be included.
	includeStackEventsFlag = "events"
)

// logLevelByResourceStatus maps each possible resource status to a log level.
var logLevelByResourceStatus = map[types.ResourceStatus]model.Level{
	types.ResourceStatusCreateInProgress: model.InfoLevel,
	types.ResourceStatusCreateFailed:     model.ErrorLevel,
	types.ResourceStatusCreateComplete:   model.InfoLevel,

	types.ResourceStatusDeleteInProgress: model.InfoLevel,
	types.ResourceStatusDeleteFailed:     model.ErrorLevel,
	types.ResourceStatusDeleteComplete:   model.InfoLevel,
	types.ResourceStatusDeleteSkipped:    model.WarnLevel,

	types.ResourceStatusUpdateInProgress: model.InfoLevel,
	types.ResourceStatusUpdateFailed:     model.ErrorLevel,
	types.ResourceStatusUpdateComplete:   model.InfoLevel,

	types.ResourceStatusImportFailed:     model.ErrorLevel,
	types.ResourceStatusImportComplete:   model.InfoLevel,
	types.ResourceStatusImportInProgress: model.InfoLevel,

	types.ResourceStatusImportRollbackInProgress: model.WarnLevel,
	types.ResourceStatusImportRollbackFailed:     model.ErrorLevel,
	types.ResourceStatusImportRollbackComplete:   model.WarnLevel,
}

// URL schemes supported
const (
	scheme              = "cloudformation"
	schemeAliasShort    = "cf"
	schemeAliasFriendly = "stack"
)

// normalizeURL take the sourceURL and expands any syntactic sugar.
func normalizeURL(sourceURL model.SourceURL) {
	switch {
	case sourceURL.Scheme == schemeAliasShort:
	case sourceURL.Scheme == schemeAliasFriendly:
		sourceURL.Scheme = scheme
	}
}

// cloudFormationSource ...
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
		// For example: a lambda's log group's cloudwatch.Source.
		children []cloudwatch.Source
		// cf represents the clo
		cf        *cloudformation.Client
		stackName string
		// stackEventsNextToken is used to keep track of the last call to fetch stack events.
		stackEventsNextToken *string
	}
)

// create CloudFormation source casted to an *interface{}.
func create(sourceURL model.SourceURL) (*interface{}, error) {
	src, err := newSource(sourceURL)
	if err != nil {
		return nil, err
	}
	var i interface{} = src
	return &i, nil
}

// newSource CloudFormation source.
func newSource(sourceURL model.SourceURL) (*cloudFormationSource, error) {
	normalizeURL(sourceURL)
	stackName := sourceURL.Path
	awsSource, err := aws.New(&sourceURL)
	if err != nil {
		return nil, errors2.Wrapf(err, "bad CloudFormation URL: %+v", sourceURL)
	}

	lookbackTime, err := sourceURL.Since(lookbackTimeFlag)
	if err != nil {
		return nil, err
	}

	cf := cloudformation.NewFromConfig(awsSource.Config())

	var children []cloudwatch.Source
	// add lambda functions
	{
		lambdaResources, err := discoverLambdaResources(cf, stackName)
		if err != nil {
			return nil, err
		}
		for _, lambdaResource := range lambdaResources {
			// This needs work: it is really ugly to have to create the CW source this way
			// however, the constructor chain all expects model.URL instances
			// as the primary constructor element. It doesn't quite reach the level of a to-do task yet.
			var rawQuery = fmt.Sprintf("%s=%s&%s=%s", aws.RegionQueryFlag, awsSource.Region(), aws.ProfileQueryFlag, awsSource.Profile())
			if lookbackTime != nil {
				rawQuery += fmt.Sprintf("&%s=%s", lookbackTimeFlag, lookbackTime.Format(time.RFC3339))
			}
			child, err := cloudwatch.New(model.SourceURL{
				Scheme:   "lambda",
				Path:     *lambdaResource.PhysicalResourceId,
				RawQuery: rawQuery,
			})
			if err != nil {
				return nil, err
			}
			children = append(children, *child)
		}
	}

	return &cloudFormationSource{
		Source:               *awsSource,
		lookbackTime:         lookbackTime,
		includeStackEvents:   sourceURL.QueryFlag(includeStackEventsFlag),
		children:             children,
		cf:                   cf,
		stackName:            stackName,
		stackEventsNextToken: nil,
	}, nil
}

// Poll for more events.
func (source cloudFormationSource) Poll() ([]model.Event, error) {
	var events []model.Event

	childEvents, err := source.getChildEvents()
	if err != nil {
		return nil, err
	}
	events = append(events, childEvents...)

	if source.includeStackEvents {
		stackEvents, err := source.getStackEvents()
		if err != nil {
			return nil, err
		}
		events = append(events, stackEvents...)
	}

	return events, nil
}

// getChildEvents returns all CloudWatch log events for resources under this stack.
func (source cloudFormationSource) getChildEvents() ([]model.Event, error) {
	var events []model.Event

	for _, child := range source.children {
		results, err := child.Poll()
		if err != nil {
			return nil, err
		}
		events = append(events, results...)
	}
	return events, nil
}

// getChildEvents retrieves CloudFormation events for this stack.
func (source cloudFormationSource) getStackEvents() ([]model.Event, error) {
	stackEvents, err := source.cf.DescribeStackEvents(context.TODO(), &cloudformation.DescribeStackEventsInput{
		NextToken: source.stackEventsNextToken,
		StackName: &source.stackName,
	})
	if err != nil {
		return nil, err
	}
	source.stackEventsNextToken = stackEvents.NextToken

	var events []model.Event
	for _, stackEvent := range stackEvents.StackEvents {
		if source.lookbackTime == nil || source.lookbackTime.Before(*stackEvent.Timestamp) {
			level := logLevelByResourceStatus[stackEvent.ResourceStatus]

			var exception *model.Exception
			if level == model.ErrorLevel {
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
				Level:      level,
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
