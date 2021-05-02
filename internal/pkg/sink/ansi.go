package sink

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/source"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mitchellh/go-wordwrap"
	"strings"
)

var (
	//au aurora colorizer instance.
	au = aurora.NewAurora(isTty) // create colorizer
	//sourceColor is the color in with which to display the name of the source.Source of the event.
	sourceColor = func(name string) aurora.Value { return au.BrightBlue(name).BgBrightBlack() }
	//exceptionColor is the color with which to display model.Exception.
	exceptionColor = au.BrightRed
	//timeColor is the color with which to display model.Timestamp.
	timeColor = au.BrightWhite
	//levelColors are the colors associated with each model.Level.
	levelColors = colorLevels{
		model.LevelTrace: au.BrightBlack,
		model.LevelDebug: au.BrightWhite,
		model.LevelInfo:  au.BrightGreen,
		model.LevelWarn:  au.BrightYellow,
		model.LevelError: au.BrightRed,
	}
	//threadColor specifies the color in which to display model.ThreadName.
	threadColor = au.Blue
	//threadColor specifies the color in which to display model.ClassName, model.MethodName, and model.LineNumber
	//via model.Caller().
	callerColor = au.Cyan

	//messageColors defines the color of the message for each model.Level.
	//TODO make this an option where messages can always be display as a single color
	messageColors = colorLevels{
		model.LevelTrace: func(i interface{}) aurora.Value { return au.Gray(8, i) },
		model.LevelDebug: func(i interface{}) aurora.Value { return au.Gray(12, i) },
		model.LevelInfo:  func(i interface{}) aurora.Value { return au.Gray(18, i) },
		model.LevelWarn:  func(i interface{}) aurora.Value { return au.Gray(20, i) },
		model.LevelError: func(i interface{}) aurora.Value { return au.Gray(24, i) },
	}
)

type (
	colorCreator func(arg interface{}) aurora.Value
	//colorLevels is a map of a log level to a color function.
	colorLevels map[model.Level]colorCreator

	//Ansi provides colorized output.
	Ansi struct {
		baseSink
	}
)

//NewAnsi provides a configured Ansi sink.
func NewAnsi(options Options, _ string) interface{} {
	return Ansi{baseSink{options: options}}
}

func (sink Ansi) Accept(src source.Source, event model.Event) {
	var srcName string
	var timestamp string
	var threadName string
	var caller string
	var message string
	if isTty {
		timestamp = event.Timestamp.GoString()
		threadName = event.ThreadName.GoString()
		srcName = src.GoString()
		caller = event.Caller().GoString()
		message = event.Message.GoString()
	} else {
		timestamp = event.Timestamp.String()
		threadName = string(event.ThreadName)
		srcName = src.String()
		caller = event.Caller().String()
		message = event.Message.String()
	}
	message = strings.TrimRight(message, "\n ")

	line := au.Sprintf(
		"%-32s %s [%s] %-20s %-40s %s%s",
		sourceColor(srcName),
		timeColor(timestamp),
		sink.levelColor(event.Level)(event.Level.Short()),
		threadColor(threadName),
		callerColor(caller),
		sink.messageColor(event.Level)(message),
		sink.stackTrace(event),
	)
	if sink.options.Wrap() > 0 {
		line = wordwrap.WrapString(line, sink.options.Wrap())
	}
	fmt.Println(line)
}

//messageColor gets the color to display a model.Message for a given model.Level.
func (sink Ansi) messageColor(level model.Level) colorCreator {
	var ok bool
	var messageColor func(arg interface{}) aurora.Value
	if messageColor, ok = messageColors[level]; !ok {
		messageColor = messageColors[model.LevelInfo]
	}
	return messageColor
}

//levelColor gets the color to display a model.Level.
func (sink Ansi) levelColor(level model.Level) colorCreator {
	var ok bool
	var levelColor func(arg interface{}) aurora.Value
	if levelColor, ok = levelColors[level]; !ok {
		levelColor = levelColors[model.LevelInfo]
	}
	return levelColor
}

//stackTrace returns a nicely formatted model.StackTrace.
func (sink Ansi) stackTrace(event model.Event) string {
	var stackTrace = "\n"
	for _, line := range strings.Split(string(event.StackTrace("\n")), "\n") {
		if line != "" {
			stackTrace += aurora.BgBrightRed(" ").Black().String() + " " + exceptionColor(line).String() + "\n"
		}
	}
	stackTrace = strings.TrimRight(stackTrace, "\n")
	return stackTrace
}
