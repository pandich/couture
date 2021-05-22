package pretty

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/pretty/column"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"couture/internal/pkg/tty"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/termenv"
	"go.uber.org/ratelimit"
	"log"
)

// Name ...
const Name = "pretty"

// TODO themes need to auto light/dark background adjust
// FIXME linebreaks messed up in highlighting process?

// prettySink provides render output.
type prettySink struct {
	terminalWidth int
	palette       chan string
	columnOrder   []string
	config        config.Config
	printer       chan string
}

// New provides a configured prettySink sink.
func New(cfg config.Config) *sink.Sink {
	if !tty.IsTTY() || cfg.Theme.BaseColor == theme.White {
		cfmt.DisableColors()
	}
	column.ByName.Init(cfg.Theme)
	cols := append([]string{"source"}, cfg.Columns...)
	var snk sink.Sink = &prettySink{
		terminalWidth: cfg.WrapWidth(),
		palette:       tty.NewColorCycle(cfg.Theme.SourceColors),
		columnOrder:   cols,
		config:        cfg,
		printer:       newPrinter(),
	}
	return &snk
}

// Init ...
func (snk *prettySink) Init(sources []source.Source) {
	for _, src := range sources {
		column.RegisterSourceStyle(src, <-snk.palette)
	}
	if snk.config.ClearScreen {
		termenv.ClearScreen()
	}
}

// Accept ...
func (snk *prettySink) Accept(src source.Source, event model.Event) error {
	snk.printer <- snk.render(src, event)
	return nil
}

func (snk *prettySink) render(src source.Source, event model.Event) string {
	// get format string and arguments
	var format = ""
	var values []interface{}
	for _, name := range snk.columnOrder {
		col := column.ByName[name]
		format += col.Formatter(src, event)
		format += tty.ResetSequence
		values = append(values, col.Renderer(snk.config, src, event)...)
	}

	// render
	line := cfmt.Sprintf(format, values...)
	wrapped := tty.Wrap(line, snk.terminalWidth)
	return wrapped
}

func newPrinter() chan string {
	// TODO do we want to rate limit?
	const ttyMaxEventsPerSecond = 200
	var throttle ratelimit.Limiter
	if tty.IsTTY() {
		throttle = ratelimit.New(ttyMaxEventsPerSecond)
	} else {
		throttle = ratelimit.NewUnlimited()
	}

	printer := make(chan string)
	go func() {
		for message := range printer {
			throttle.Take()
			log.Println(message)
		}
	}()
	return printer
}
