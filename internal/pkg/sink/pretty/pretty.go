package pretty

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"fmt"
	"strings"
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
		event.AsCaller(),
		event.Message,
	}

	stackTrace := event.StackTrace()
	if stackTrace != nil {
		fields = append(fields, "\n", *stackTrace)
	}

	fmt.Println(snk.styles.render(fields...))
}

func abbreviateClassName(className model.ClassName, aspirationalWidth int) string {
	var s = string(className)
	pieces := strings.Split(s, ".")
	var l = len(s) - (len(pieces) - 1)
	var changed = true
	for l > aspirationalWidth && changed {
		changed = false
		for i := 0; i < len(pieces)-1; i++ {
			if len(pieces[i]) > 1 {
				l -= len(pieces[i]) - 1
				pieces[i] = string(pieces[i][0])
				changed = true
			}
			if l <= aspirationalWidth {
				break
			}
		}
	}
	return strings.Join(pieces, ".")
}
