package sink

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/source"
	"github.com/mattn/go-isatty"
	"os"
)

// IsTTY ...
func IsTTY() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}

// Sink of events. Responsible for consuming an event.
type Sink interface {
	// Init called prior to the beginning of logging.
	Init(sources []model.SourceURL)
	// Accept consumes an event, typically for display.
	Accept(src source.Pushable, event model.Event) error
}
