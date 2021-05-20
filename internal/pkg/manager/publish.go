package manager

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/source"
	"fmt"
	"net/url"
	"os"
	"time"
)

// eventTopic is the topic for all registry and sinks to communicate over.
const eventTopic = "topic:event"

// internalSource is the source used for all diagnostic messages.
var internalSource = source.New(model.SourceURL{})

func (mgr *publishingManager) publishError(
	methodName model.MethodName,
	level level.Level,
	err error,
	message string,
	args ...interface{},
) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	event := newDiagnosticEvent(level, methodName, message, args...)
	event.Exception = model.NewException(err)
	mgr.publishEvent(internalSource, event)
}

func (mgr *publishingManager) publishDiagnostic(level level.Level, methodName model.MethodName, message string) {
	event := newDiagnosticEvent(level, methodName, message)
	mgr.publishEvent(internalSource, event)
}

func (mgr *publishingManager) publishEvent(src source.Source, event model.Event) {
	if !event.Level.IsAtLeast(mgr.options.level) {
		return
	}
	if event.Matches(mgr.options.includeFilters, mgr.options.excludeFilters) {
		mgr.rateLimiter.Take()
		mgr.bus.Publish(eventTopic, src, event)
	}
}

func newDiagnosticEvent(
	level level.Level,
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
