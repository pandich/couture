package sink

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/source"
	"github.com/mattn/go-isatty"
	"os"
)

var (
	//isTty specifies whether or not we are writing to a TTY.
	isTty = isatty.IsTerminal(os.Stdout.Fd())
)

type (
	//Implementations go here. Each implementation struct should be unexported and exposed with a var.
	//For each implementation, update cmd/couture/cli/sink.

	//Sink of events. Responsible for consuming an event.
	Sink interface {
		//Accept consumes an event, typically for display.
		Accept(src source.Source, event model.Event)
	}

	//Options for displaying output. Each Sink may use or ignore these values as is appropriate to their type
	//the state of isTty, and other considerations.
	Options interface {
		Wrap() uint
	}

	//baseSink is meant to be included in all Sink implementations.
	baseSink struct {
		//options contains the options for this sink.
		options Options
	}
)
