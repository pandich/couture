package model

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

var (
	ErrNoMoreEvents = errors.New("no more events")
)

//goland:noinspection GoUnusedConst
const (
	//LevelTrace log level for tracing
	LevelTrace Level = "TRACE"
	//LevelDebug log level for debugging
	LevelDebug Level = "DEBUG"
	//LevelInfo log level for information
	LevelInfo Level = "INFO"
	//LevelWarn log level for warnings
	LevelWarn Level = "WARN"
	//LevelError log level for errors
	LevelError Level = "ERROR"
)

type (
	//Timestamp When the even occurred.
	Timestamp time.Time
	//MethodName a method name.
	MethodName string
	//Level a log level.
	Level string
	//LineNumber  a line number.
	LineNumber string
	//ThreadName a thread name.
	ThreadName string
	//ClassName a class name.
	ClassName string
	//Message a message.
	Message string
	//StackTrace a stack trace.
	StackTrace string
	//Caller represents <class>:<method>#<lime_number>.
	Caller string
	//Exception an exception.
	Exception struct {
		//StackTrace the full text of the stack trace.
		StackTrace StackTrace `json:"stackTrace"`
	}

	//Event a log event
	Event struct {
		//Timestamp the timestamp. This field is required, and should default to time.Now() if not present.
		Timestamp Timestamp
		//Level the level. This field is required, and should default to LevelInfo if not present.
		Level Level
		//Message the message. This field is required.
		Message Message
		//MethodName the method name. This field is optional.
		MethodName MethodName
		//LineNumber the line number. This field is optional.
		LineNumber LineNumber
		//ThreadName the thread name. This field is optional.
		ThreadName ThreadName
		//ClassName the class name. This field is optional.
		ClassName ClassName
		//Exception the exception. This field is optional.
		Exception *Exception
	}
)

func (event Event) LineNumberAsInt() uint64 {
	i, err := strconv.ParseUint(string(event.LineNumber), 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func (event Event) Caller() Caller {
	return Caller(fmt.Sprintf("%s:%s#%-4d", event.ClassName, event.MethodName, event.LineNumberAsInt()))
}

func (event Event) GoString() string {
	var ex = ""
	if event.Exception != nil {
		ex = "\nException: " + string((*event.Exception).StackTrace)
	}
	var ln = string(event.LineNumber)
	if i, err := strconv.ParseInt(string(event.LineNumber), 10, 64); err == nil {
		ln = fmt.Sprintf("%-4d", i)
	}
	return fmt.Sprintf(
		"%s [%-5s] (%s) %s#%s@%s - %s%s",
		event.Timestamp,
		event.Level,
		event.ThreadName,
		event.ClassName,
		event.MethodName,
		ln,
		event.Message,
		ex,
	)
}

func (event Event) StackTrace(prefix string) StackTrace {
	if event.Exception != nil {
		return StackTrace(prefix + string(event.Exception.StackTrace))
	}
	return ""
}

func (t Timestamp) String() string {
	return time.Time(t).Format(time.RFC3339)
}
func (t Timestamp) GoString() string {
	return t.String()
}

func (msg Message) String() string {
	return string(msg)
}
func (msg Message) GoString() string {
	return "❞ " + msg.String()
}

func (caller Caller) String() string {
	return string(caller)
}
func (caller Caller) GoString() string {
	return "➤ " + caller.String()
}

func (t ThreadName) GoString() string {
	return "⑂ " + string(t)
}

func (level Level) Short() string {
	return string(level[0])
}
