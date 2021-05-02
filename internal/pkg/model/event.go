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
	LevelTrace LogLevel = "TRACE"
	//LevelDebug log level for debugging
	LevelDebug LogLevel = "DEBUG"
	//LevelInfo log level for information
	LevelInfo LogLevel = "INFO"
	//LevelWarn log level for warnings
	LevelWarn LogLevel = "WARN"
	//LevelError log level for errors
	LevelError LogLevel = "ERROR"
)

type (
	//Timestamp When the even occurred.
	Timestamp time.Time
	//MethodName a method name.
	MethodName string
	//LogLevel a log level.
	LogLevel string
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
		Level LogLevel
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

func (t ThreadName) GoString() string {
	return "🧵" + string(t)
}

func (e Event) LineNumberAsInt() uint64 {
	i, err := strconv.ParseUint(string(e.LineNumber), 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func (e Event) Caller() Caller {
	return Caller(fmt.Sprintf("%s:%s#%-4d", e.ClassName, e.MethodName, e.LineNumberAsInt()))
}

func (e Event) GoString() string {
	var ex = ""
	if e.Exception != nil {
		ex = "\nException: " + string((*e.Exception).StackTrace)
	}
	var ln = string(e.LineNumber)
	if i, err := strconv.ParseInt(string(e.LineNumber), 10, 64); err == nil {
		ln = fmt.Sprintf("%-4d", i)
	}
	return fmt.Sprintf(
		"%s [%-5s] (%s) %s#%s@%s - %s%s",
		e.Timestamp,
		e.Level,
		e.ThreadName,
		e.ClassName,
		e.MethodName,
		ln,
		e.Message,
		ex,
	)
}

func (e Event) StackTrace(prefix string) StackTrace {
	if e.Exception != nil {
		return StackTrace(prefix + string(e.Exception.StackTrace))
	}
	return ""
}

func (t Timestamp) String() string {
	return time.Time(t).Format(time.RFC3339)
}
func (t Timestamp) GoString() string {
	return "⌚︎" + t.String()
}

func (m Message) String() string {
	return string(m)
}
func (m Message) GoString() string {
	return "📕" + m.String()
}

func (c Caller) String() string {
	return string(c)
}
func (c Caller) GoString() string {
	return "📞" + c.String()
}

func (l LogLevel) Short() string {
	return string(l[0])
}
