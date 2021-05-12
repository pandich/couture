package model

import (
	"regexp"
)

// TimestampField ...
const TimestampField = "@timestamp"

// Event a log event
type Event struct {
	// Timestamp the timestamp. This field is required, and should default to time.Now() if not present.
	Timestamp Timestamp `json:"@timestamp"`
	// Level the level. This field is required, and should default to LevelInfo if not present.
	Level Level `json:"level"`
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

// Matches determines if an event matches the filters criteria.
func (event Event) Matches(level Level, include []*regexp.Regexp, exclude []*regexp.Regexp) bool {
	// return false if the log level is too low
	if !event.Level.isAtLeast(level) {
		return false
	}

	// process the includes returning true on the first match
	for _, filter := range include {
		if filter.MatchString(string(event.Message)) {
			return true
		}
	}
	// if we made it this far and have include filters, none of them matched, so we return false
	if len(include) > 0 {
		return false
	}

	// process the excludes returning false on the first match
	for _, filter := range exclude {
		if filter.MatchString(string(event.Message)) {
			return false
		}
	}

	// return true
	return true
}

// StackTrace ...
func (event Event) StackTrace() *StackTrace {
	if event.Exception != nil && event.Exception.StackTrace != "" {
		return &event.Exception.StackTrace
	}
	return nil
}
