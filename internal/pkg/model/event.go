package model

import (
	"couture/internal/pkg/model/level"
	"regexp"
)

// TimestampField ...
const TimestampField = "@timestamp"

// Event a log event
type Event struct {
	// Timestamp the timestamp. This field is required, and should default to time.Now() if not present.
	Timestamp Timestamp `json:"@timestamp"`
	// Level the level. This field is required, and should default to Info if not present.
	Level level.Level `json:"level"`
	// Message the message. This field is required.
	Message Message `json:"message"`
	// ApplicationName is the name of the application that generated this event. This field is optional.
	ApplicationName *ApplicationName `json:"application,omitempty"`
	// MethodName the method name. This field is optional.
	MethodName MethodName `json:"method"`
	// LineNumber the line number. This field is optional.
	LineNumber LineNumber `json:"line_number"`
	// ThreadName the thread name. This field is optional.
	ThreadName *ThreadName `json:"thread_name"`
	// ClassName the class name. This field is optional.
	ClassName ClassName `json:"class"`
	// Exception the exception. This field is optional.
	Exception *Exception `json:"exception,omitempty"`
	// highlightMarks all the matches found in this event
	highlightMarks highlightMarks
	// highlightMarks all the matches found in this event
	stackTraceHighlightMarks highlightMarks
}

// ApplicationNameOrBlank ...
func (event Event) ApplicationNameOrBlank() ApplicationName {
	if event.ApplicationName != nil {
		return *event.ApplicationName
	}
	return ""
}

// ThreadNameOrBlank ...
func (event Event) ThreadNameOrBlank() ThreadName {
	if event.ThreadName != nil {
		return *event.ThreadName
	}
	return ""
}

func (event Event) stackTrace() *StackTrace {
	if event.Exception != nil && event.Exception.StackTrace != "" {
		return &event.Exception.StackTrace
	}
	return nil
}

// Matches ...
func (event *Event) Matches(includes []regexp.Regexp, excludes []regexp.Regexp) bool {
	if len(includes) == 0 && len(excludes) == 0 {
		return true
	}

	// setting state inside of here is pretty ugly
	event.highlightMarks = []highlightMark{}
	if highlightMarks, matches := event.Message.matches(includes, excludes); matches {
		event.highlightMarks = highlightMarks
	}

	event.stackTraceHighlightMarks = []highlightMark{}
	trace := event.stackTrace()
	if trace != nil {
		if highlightMarks, matches := Message(*trace).matches(includes, excludes); matches {
			event.stackTraceHighlightMarks = highlightMarks
		}
	}
	return len(event.highlightMarks) > 0 || len(event.stackTraceHighlightMarks) > 0
}

// HighlightedMessage ...
func (event Event) HighlightedMessage() []interface{} {
	return event.Message.highlighted(
		event.highlightMarks.merged(),
		func(msg Message) interface{} { var i interface{} = HighlightedMessage(msg); return i },
		func(msg Message) interface{} { var i interface{} = UnhighlightedMessage(msg); return i },
	)
}

// HighlightedStackTrace ...
func (event Event) HighlightedStackTrace() []interface{} {
	trace := event.stackTrace()
	if trace == nil {
		return []interface{}{}
	}
	return event.Message.highlighted(
		event.highlightMarks.merged(),
		func(msg Message) interface{} { var i interface{} = HighlightedStackTrace(msg); return i },
		func(msg Message) interface{} { var i interface{} = UnhighlightedStackTrace(msg); return i },
	)
}
