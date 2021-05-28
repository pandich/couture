package model

import (
	"couture/internal/pkg/model/level"
)

// TimestampField ...
const TimestampField = "@timestamp"

// Event a log event
// TODO sources should be able to define one or more formats globally, or by source
//		perhaps with detection heuristics. Baically, remove JSON tags from this struct
//		and then have some separate ETL definition file. Standard types could be built-in.
//		Custom types in ~/.config/couture/mappings/*.yaml
//		see: https://github.com/tidwall/gjson - quick type detection?
//		see: https://github.com/jf-tech/omniparser - mapping?
//		see: https://github.com/Qntfy/kazaam - mapping?
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

// StackTrace ...
func (event Event) StackTrace() *StackTrace {
	if event.Exception != nil && event.Exception.StackTrace != "" {
		return &event.Exception.StackTrace
	}
	return nil
}
