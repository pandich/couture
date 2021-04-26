package model

import (
	"fmt"
	"time"
)

//goland:noinspection GoUnusedConst
const (
	// LevelMissing log level when no valid level is present
	LevelMissing Level = ""
	// LevelTrace log level for tracing
	LevelTrace Level = "TRACE"
	// LevelDebug log level for debugging
	LevelDebug Level = "DEBUG"
	// LevelInfo log level for information
	LevelInfo Level = "INFO"
	// LevelWarn log level for warnings
	LevelWarn Level = "WARN"
	// LevelError log level for errors
	LevelError Level = "ERROR"

	// MessageMissing for when an event lacks a message
	MessageMissing Message = "⁉"
)

type (
	// Timestamp ISO-8601 timestamp of when an event occurred.
	Timestamp string
	// MethodName a method name.
	MethodName string
	// Level a log level.
	Level string
	// LineNumber  a line number.
	LineNumber uint32
	// ThreadName a thread name.
	ThreadName string
	// ClassName a class name.
	ClassName string
	// Message a message.
	Message string
	// StackTrace a stack trace.
	StackTrace string
	// Exception an exception.
	Exception struct {
		// StackTrace the full text of the stack trace.
		StackTrace StackTrace `json:"stackTrace"`
	}

	// Event a log event
	Event struct {
		// Timestamp the timestamp. This field is required, and should default to time.Now() if not present.
		Timestamp Timestamp
		// Level the level. This field is required, and should default to LevelMissing if not present.
		Level Level
		// Message the message. This field is required, and should default to MessageMissing if not present.
		Message Message
		// MethodName the method name. This field is optional.
		MethodName *MethodName
		// LineNumber the line number. This field is optional.
		LineNumber *LineNumber
		// ThreadName the thread name. This field is optional.
		ThreadName *ThreadName
		// ClassName the class name. This field is optional.
		ClassName *ClassName
		// Exception the exception. This field is optional.
		Exception *Exception
	}
)

// NewEvent creates an event with sensible defaults for nil required values.
func NewEvent(
	timestamp time.Time,
	level Level,
	message Message,
	methodName *MethodName,
	lineNumber *LineNumber,
	threadName *ThreadName,
	className *ClassName,
	stackTrace *StackTrace,
) *Event {
	var exception *Exception
	if stackTrace != nil {
		exception = &Exception{
			StackTrace: *stackTrace,
		}
	}

	return &Event{
		Timestamp:  Timestamp(timestamp.Format(time.RFC3339)),
		Level:      level,
		Message:    message,
		MethodName: methodName,
		LineNumber: lineNumber,
		ThreadName: threadName,
		ClassName:  className,
		Exception:  exception,
	}
}

func (e Event) GoString() string {
	var ex = ""
	if e.Exception != nil {
		ex = "\nException: " + string((*e.Exception).StackTrace)
	}
	return fmt.Sprintf(
		"%s [%-5s] (%s) %s#%s@%d - %s%s",
		e.Timestamp,
		e.Level,
		*e.ThreadName,
		*e.ClassName,
		*e.MethodName,
		*e.LineNumber,
		e.Message,
		ex,
	)
}
