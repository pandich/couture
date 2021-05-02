package sink

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/source"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-isatty"
	"os"
)

const (
	logLevelBackgroundColorIndex = 232
)

var (
	isTty          = isatty.IsTerminal(os.Stdout.Fd())
	au             = aurora.NewAurora(isTty) // create colorizer
	sourceColor    = func(name string) aurora.Value { return au.BrightBlue(name).BgBrightBlack() }
	exceptionColor = au.BrightRed
	timeColor      = au.BrightWhite
	logLevelColors = colorLevels{
		model.LevelTrace: au.BrightBlack,
		model.LevelDebug: au.BrightWhite,
		model.LevelInfo:  au.BrightGreen,
		model.LevelWarn:  au.BrightYellow,
		model.LevelError: au.BrightRed,
	}
	threadColor   = au.Blue
	callerColor   = au.Cyan
	messageColors = colorLevels{
		model.LevelTrace: au.BrightBlack,
		model.LevelDebug: au.White,
		model.LevelInfo:  au.White,
		model.LevelWarn:  au.Yellow,
		model.LevelError: au.Red,
	}
)

type (
	colorLevels map[model.LogLevel]func(interface{}) aurora.Value

	//Ansi provides colorized output.
	Ansi struct {
	}
)

//NewAnsi provides a configured Ansi sink.
func NewAnsi(_ string) interface{} {
	return Ansi{}
}

func (s Ansi) Accept(src source.Source, evt model.Event) {
	var ok bool

	var levelColor func(arg interface{}) aurora.Value
	if levelColor, ok = logLevelColors[evt.Level]; !ok {
		levelColor = logLevelColors[model.LevelInfo]
	}

	var messageColor func(arg interface{}) aurora.Value
	if messageColor, ok = messageColors[evt.Level]; !ok {
		messageColor = messageColors[model.LevelInfo]
	}

	stackTrace := exceptionColor(evt.StackTrace("\n"))

	var srcName string
	var timestamp string
	var threadName string
	var caller string
	var message string
	if isTty {
		timestamp = evt.Timestamp.GoString()
		threadName = evt.ThreadName.GoString()
		srcName = src.GoString()
		caller = evt.Caller().GoString()
		message = evt.Message.GoString()
	} else {
		timestamp = evt.Timestamp.String()
		threadName = string(evt.ThreadName)
		srcName = src.String()
		caller = evt.Caller().String()
		message = evt.Message.String()
	}

	fmt.Println(au.Sprintf(
		"%-32s %s %s %-20s %-40s %s%s",
		sourceColor(srcName),
		timeColor(timestamp),
		levelColor(evt.Level.Short()).BgIndex(logLevelBackgroundColorIndex),
		threadColor(threadName),
		callerColor(caller),
		messageColor(message),
		stackTrace,
	))
}
