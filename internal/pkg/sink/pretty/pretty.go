package pretty

import (
	"bufio"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/pretty/column"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"couture/internal/pkg/tty"
	"github.com/i582/cfmt/cmd/cfmt"
	"go.uber.org/ratelimit"
	"os"
	"os/signal"
	"syscall"
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
	if len(cfg.Columns) == 0 {
		cfg.Columns = column.DefaultColumns
	}
	column.ByName.Init(cfg.Theme)
	pretty := &prettySink{
		terminalWidth: cfg.EffectiveTerminalWidth(),
		columnOrder:   cfg.Columns,
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
	resizeChan := make(chan os.Signal)
	signal.Notify(resizeChan, os.Interrupt, syscall.SIGWINCH)
	go func() {
		for range resizeChan {
			snk.updateColumnWidths()
		}
	}()
}

// Accept ...
func (snk *prettySink) Accept(src source.Source, event sink.Event) error {
	snk.printer <- snk.render(src, event)
	return nil
}

func (snk *prettySink) render(src source.Source, event sink.Event) string {
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
	const maxTTYLinesPerSecond = 250

	var limiter = ratelimit.NewUnlimited()
	if tty.IsTTY() {
		limiter = ratelimit.New(maxTTYLinesPerSecond)
	}

	writer := bufio.NewWriter(os.Stdout)

	printer := make(chan string)
	go func() {
		var i = 0
		limiter.Take()
		for {
			message := <-printer
			_, err := writer.WriteString(message + "\n")
			if err != nil {
				panic(err)
			}
			i++
			if i%10 == 0 {
				writer.Flush()
				i = 0
			}
		}
	}()
	return printer
}
