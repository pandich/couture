package manager

import (
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"fmt"
	"net/url"
	"time"
)

const (
	// eventTopic is the topic for all registry and sinks to communicate over.
	eventTopic = "topic:event"
)

// internalSource is the source used for all diagnostic messages.
var internalSource = source.New(model.SourceURL{})

func (mgr *publishingManager) publishError(
	methodName model.MethodName,
	level model.Level,
	err error,
	message string,
	args ...interface{},
) {
	event := newDiagnosticEvent(level, methodName, message, args...)
	event.Exception = model.NewException(err)
	mgr.publishEvent(internalSource, event)
}

func (mgr *publishingManager) publishDiagnostic(level model.Level, methodName model.MethodName, message string) {
	event := newDiagnosticEvent(level, methodName, message)
	mgr.publishEvent(internalSource, event)
}

func (mgr *publishingManager) publishEvent(src source.Source, event model.Event) {
	if event.Matches(mgr.options.level, mgr.options.includeFilters, mgr.options.excludeFilters) {
		mgr.bus.Publish(eventTopic, src, event)
	}
}

func newDiagnosticEvent(
	level model.Level,
	methodName model.MethodName,
	message string,
	messageArgs ...interface{},
) model.Event {
	u := url.URL(internalSource.URL())
	threadName := model.ThreadName("[-]")
	return model.Event{
		Timestamp:  model.Timestamp(time.Now()),
		Level:      level,
		Message:    model.Message(fmt.Sprintf(message, messageArgs...)),
		MethodName: methodName,
		LineNumber: model.NoLineNumber,
		ThreadName: &threadName,
		ClassName:  model.ClassName(u.String()),
	}
}