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
	// Action the action name. This field is optional.
	Action Action `json:"action" regroup:"action"`
	// Line the line number. This field is optional.
	Line Line `json:"line" regroup:"line"`
	// Context the context name. This field is optional.
	Context Context `json:"context" regroup:"context"`
	// Entity the entity name. This field is optional.
	Entity Entity `json:"entity" regroup:"entity"`
	// Error the error. This field is optional.
	Error Error `json:"error" regroup:"error"`
}

// SinkEvent ...
type SinkEvent struct {
	Event
	SourceURL SourceURL
	Filters   filters
	Schema    *schema.Schema
}
