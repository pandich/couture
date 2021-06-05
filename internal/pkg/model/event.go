package model

import (
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/schema"
)

// Event a log event
type Event struct {
	// Timestamp the timestamp. This field is required, and should default to time.Now() if not present.
	Timestamp Timestamp `json:"timestamp" regroup:"timestamp"`
	// Level the level. This field is required, and should default to Info if not present.
	Level level.Level `json:"level" regroup:"level"`
	// Message the message. This field is required.
	Message Message `json:"message" regroup:"message"`
	// Application is the name of the application that generated this event. This field is optional.
	Application Application `json:"application" regroup:"application"`
	// Method the method name. This field is optional.
	Method Method `json:"method" regroup:"method"`
	// Line the line number. This field is optional.
	Line Line `json:"line" regroup:"line"`
	// Thread the thread name. This field is optional.
	Thread Thread `json:"thread" regroup:"thread"`
	// Class the class name. This field is optional.
	Class Class `json:"class" regroup:"class"`
	// Exception the exception. This field is optional.
	Exception Exception `json:"exception" regroup:"exception"`
}

// SinkEvent ...
type SinkEvent struct {
	Event
	SourceURL SourceURL
	Filters   []Filter
	Schema    *schema.Schema
}
