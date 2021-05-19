package sink

import (
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"github.com/mattn/go-isatty"
	"os"
)

// IsTTY ...
func IsTTY() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}

// Sink of events. Responsible for consuming an event.
type Sink interface {
	// Accept consumes an event, typically for display.
	Accept(src source.Pushable, event model.Event) error
}
