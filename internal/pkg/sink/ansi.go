package sink

import (
	"couture/internal/pkg/model"
	"fmt"
	. "github.com/logrusorgru/aurora"
	"time"
)

//NewAnsi provides a configured Ansi sink.
func NewAnsi(_ string) interface{} {
	return Ansi{}
}

var (
	timeColor      = Cyan
	logLevelColors = map[model.LogLevel]func(interface{}) Value{
		model.LevelTrace: White,
		model.LevelDebug: BrightWhite,
		model.LevelInfo:  BrightGreen,
		model.LevelWarn:  BrightYellow,
		model.LevelError: exceptionColor,
	}
	threadColor    = Blue
	callerColor    = Yellow
	messageColor   = White
	exceptionColor = BrightRed
)

//Ansi provides colorized output.
type Ansi struct {
}

func (s Ansi) Accept(event model.Event) {
	var ok bool

	var levelColor func(arg interface{}) Value
	if levelColor, ok = logLevelColors[event.Level]; !ok {
		levelColor = logLevelColors[model.LevelInfo]
	}

	var exception = messageColor("")
	stackTrace := event.StackTrace()
	if stackTrace != nil {
		exception = exceptionColor("\n" + *stackTrace)
	}

	fmt.Println(Sprintf(
		"%s %-5s (%-20s) %-40s %s%s",
		timeColor(time.Time(event.Timestamp).Format(time.RFC3339)),
		levelColor(event.Level),
		threadColor(event.ThreadName),
		callerColor(fmt.Sprintf("%s:%s#%-4d", event.ClassName, event.MethodName, event.LineNumberAsInt())),
		messageColor(event.Message),
		exception,
	))
}
