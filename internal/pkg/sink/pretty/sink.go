package pretty

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"fmt"
	"io"
)

// prettySink provides render output.
type prettySink struct {
	*sink.Base
	// sourceStyles provides a (semi-)unique Color per source.
	styles      *styler
	paginate    bool
	shortPrefix bool
}

// New provides a configured prettySink sink.
func New(out io.Writer) *sink.Sink {
	pretty := &prettySink{
		Base:        sink.New(out),
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
	_, _ = fmt.Fprintln(snk.Out(), snk.styles.render(fields...))
}
