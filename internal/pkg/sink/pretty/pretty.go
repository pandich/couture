package pretty

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"fmt"
)

// Sink provides render output.
type Sink struct {
	sink.Base
	// sourceStyles provides a (semi-)unique Color per source.
	styles      *styler
	paginate    bool
	shortPrefix bool
}

// New provides a configured Sink sink.
func New(options sink.Options, _ string) interface{} {
	return Sink{
		Base:        sink.New(options),
		styles:      newStyler(),
		shortPrefix: true,
		paginate:    true,
	}
}

// Accept ...
func (snk *Sink) Accept(src source.Source, event model.Event) {
	var fields = []interface{}{
		src,
		event.ApplicationNameOrBlank(),
		event.Timestamp,
		event.Level,
		event.ThreadNameOrBlank(),
		event.ClassName.Abbreviate(classNameColumnWidth),
		methodNameDelimiter,
		event.MethodName,
		lineNumberDelimiter,
		event.LineNumber,
		event.Message,
	}

	stackTrace := event.StackTrace()
	if stackTrace != nil {
		fields = append(fields, "\n", *stackTrace)
	}

	fmt.Println(snk.styles.render(fields...))
}
