package pretty

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"fmt"
)

// prettySink provides render output.
type prettySink struct {
	// sourceStyles provides a (semi-)unique Color per source.
	styles      *styler
	paginate    bool
	shortPrefix bool
}

// New provides a configured prettySink sink.
func New() *sink.Sink {
	pretty := &prettySink{
		styles:      newStyler(),
		shortPrefix: true,
		paginate:    true,
	}
	var snk sink.Sink = pretty
	return &snk
}

type caller string

// Accept ...
func (snk *prettySink) Accept(src source.Source, event model.Event) {
	var fields = []interface{}{
		src,
		event.ApplicationNameOrBlank(),
		event.Timestamp.Stamp(),
		event.Level,
		caller(snk.styles.render(
			event.ClassName.Abbreviate(callerWidth),
			methodNameDelimiter,
			event.MethodName,
			lineNumberDelimiter,
		)),
		event.LineNumber,
		event.ThreadNameOrBlank(),
		newLine,
		model.Message(snk.styles.render(event.HighlightedMessage()...)),
	}
	if stackTrace := event.StackTrace(); stackTrace != nil {
		fields = append(fields, "\n", model.StackTrace(snk.styles.render(event.HighlightedStackTrace()...)))
	}
	fmt.Println(snk.styles.render(fields...))
}
