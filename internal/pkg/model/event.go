package model

import (
	"couture/internal/pkg/model/level"
	"regexp"
)

// Event a log event
type Event struct {
	// Timestamp the timestamp. This field is required, and should default to time.Now() if not present.
	Timestamp Timestamp
	// Level the level. This field is required, and should default to Info if not present.
	Level level.Level
	// Message the message. This field is required.
	Message Message
	// Application is the name of the application that generated this event. This field is optional.
	Application Application
	// Method the method name. This field is optional.
	Method Method
	// Line the line number. This field is optional.
	Line Line
	// Thread the thread name. This field is optional.
	Thread Thread
	// Class the class name. This field is optional.
	Class Class
	// Exception the exception. This field is optional.
	Exception Exception
}

// SinkEvent ...
type SinkEvent struct {
	Event
	SourceURL SourceURL
	Filters   []regexp.Regexp
}
