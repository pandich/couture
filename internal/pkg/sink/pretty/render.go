package pretty

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/source"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/wordwrap"
)

func (snk *prettySink) renderEvent(src source.Source, event model.Event) (string, error) {
	sourceStyleName := snk.palette.sourceStyle(src.URL())
	line := cfmt.Sprintf(
		"{{%s}}::"+sourceStyleName+
			"{{%s}}::Timestamp"+
			"{{%s}}::Application"+
			"{{%s}}::Thread"+
			"%s"+ // caller
			"{{%s}}::Level"+string(event.Level)+
			"%s"+ // message
			"%s", // stack trace
		src.URL().ShortForm(),
		event.Timestamp.Stamp(),
		event.ApplicationNameOrBlank(),
		event.ThreadNameOrBlank(),
		snk.renderCaller(event),
		event.Level,
		snk.renderHighlightedMessage(event),
		snk.renderHighlightedStackTrace(event),
	)
	return snk.wrapToTerminal(line)
}

func (snk *prettySink) renderCaller(event model.Event) string {
	const classNameWidth = 30
	const callerWidth = 55
	caller := padding.String(cfmt.Sprintf(
		"{{%s}}::Class{{/}}::MethodDelimiter{{%s}}::Method{{#}}::LineNumberDelimiter{{%d}}::LineNumber  ",
		event.ClassName.Abbreviate(classNameWidth),
		event.MethodName,
		event.LineNumber,
	), callerWidth)
	return caller
}

func (snk *prettySink) renderHighlightedMessage(event model.Event) string {
	var message = ""
	for _, chunk := range event.HighlightedMessage() {
		message += " "
		switch chunk.(type) {
		case model.HighlightedMessage:
			message += cfmt.Sprintf("{{%s}}::HighlightedMessage"+string(event.Level), chunk)
		case model.UnhighlightedMessage:
			message += cfmt.Sprintf("{{%s}}::Message"+string(event.Level), chunk)
		default:
			message += cfmt.Sprintf("{{%s}}::Message"+string(event.Level), chunk)
		}
	}
	return message
}

func (snk *prettySink) renderHighlightedStackTrace(event model.Event) string {
	var stackTrace = ""
	for _, chunk := range event.HighlightedStackTrace() {
		if stackTrace == "" {
			stackTrace += "\n"
		} else {
			stackTrace += " "
		}
		switch chunk.(type) {
		case model.HighlightedStackTrace:
			stackTrace += cfmt.Sprintf("{{%s}}::HighlightedStackTrace", chunk)
		case model.UnhighlightedStackTrace:
			stackTrace += cfmt.Sprintf("{{%s}}::StackTrace", chunk)
		default:
			stackTrace += cfmt.Sprintf("{{%s}}::StackTrace", chunk)
		}
	}
	return stackTrace
}

func (snk *prettySink) wrapToTerminal(s string) (string, error) {
	if snk.terminalWidth == noWrap {
		return s, nil
	}
	wrapper := wordwrap.NewWriter(snk.terminalWidth)
	wrapper.Breakpoints = []rune(" \t")
	wrapper.KeepNewlines = true
	if _, err := wrapper.Write([]byte(s)); err != nil {
		return "", err
	}
	wrapped := wrapper.String()
	return wrapped, nil
}
