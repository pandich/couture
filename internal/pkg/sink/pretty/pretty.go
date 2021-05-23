package pretty

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/pretty/column"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"couture/internal/pkg/tty"
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
	"go.uber.org/ratelimit"
)

// Name ...
const Name = "pretty"

// TODO themes need to auto light/dark background adjust
// FIXME linebreaks messed up in highlighting process?

// prettySink provides render output.
type prettySink struct {
	terminalWidth uint
	columnOrder   []string
	config        config.Config
	printer       chan string
	columnWidths  map[string]uint
}

// New provides a configured prettySink sink.
func New(cfg config.Config) *sink.Sink {
	if !tty.IsTTY() || cfg.Theme.BaseColor == theme.White {
		cfmt.DisableColors()
	}
	column.ByName.Init(cfg.Theme)
	cols := append([]string{"source"}, cfg.Columns...)
	pretty := &prettySink{
		terminalWidth: cfg.EffectiveTerminalWidth(),
		columnOrder:   cols,
		config:        cfg,
		printer:       newPrinter(),
	}
	pretty.updateColumnWidths()
	var snk sink.Sink = pretty
	return &snk
}

func (snk *prettySink) updateColumnWidths() {
	snk.columnWidths = column.Widths(uint(tty.TerminalWidth()), snk.columnOrder)
}

// Init ...
func (snk *prettySink) Init(sources []source.Source) {
	for _, src := range sources {
		column.RegisterSource(snk.config.Theme, src)
	}
	snk.updateColumnWidths()
	if snk.config.ClearScreen {
		snk.printer <- tty.ClearScreenSequence
	}
	snk.printer <- tty.HomeCursorSequence
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
		format += col.Format(snk.columnWidths[name], src, event)
		values = append(values, col.Render(snk.config, src, event)...)
	}
	format += tty.ResetSequence

	// render
	var line = cfmt.Sprintf(format, values...)
	if snk.config.Wrap {
		line = tty.Wrap(line, snk.config.EffectiveTerminalWidth())
	}
	return line
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
			fmt.Println(message)
		}
	}()
	return printer
}
