package pretty

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/pretty/column"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"couture/internal/pkg/tty"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/reflow/wordwrap"
	"os"
	"os/signal"
	"syscall"
)

// Name ...
const Name = "pretty"

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
	if len(cfg.Columns) == 0 {
		cfg.Columns = column.DefaultColumns
	}
	column.ByName.Init(cfg.Theme)
	var snk sink.Sink = &prettySink{
		terminalWidth: cfg.EffectiveTerminalWidth(),
		columnOrder:   cfg.Columns,
		config:        cfg,
		printer:       tty.NewTTYWriter(os.Stdout),
	}
	return &snk
}

// Init ...
func (snk *prettySink) Init(sources []*source.Source) {
	for _, src := range sources {
		column.RegisterSource(snk.config.Theme, *src)
	}
	snk.handleColorState()
	snk.updateColumnWidths()
	if snk.config.AutoResize {
		snk.handleAutoResizeState()
	}
}

// Accept ...
func (snk *prettySink) Accept(event sink.Event) error {
	format, values := snk.columnFormat(event)
	var line = cfmt.Sprintf(format, values...)
	if snk.config.Wrap {
		line = wordwrap.String(line, int(snk.config.EffectiveTerminalWidth()))
	}
	snk.printer <- line
	return nil
}

func (snk *prettySink) columnFormat(event sink.Event) (string, []interface{}) {
	const resetSequence = "\x1b[2m"

	// get format string and arguments
	var format = ""
	var values []interface{}
	for _, name := range snk.columnOrder {
		col := column.ByName[name]
		format += col.Format(snk.columnWidths[name], event)
		values = append(values, col.Render(snk.config, event)...)
	}
	format += resetSequence
	return format, values
}

func (snk *prettySink) updateColumnWidths() {
	snk.columnWidths = column.Widths(uint(tty.TerminalWidth()), snk.columnOrder)
}

func (snk *prettySink) handleAutoResizeState() {
	resize := make(chan os.Signal)
	signal.Notify(resize, os.Interrupt, syscall.SIGWINCH)
	go func() {
		defer close(resize)
		for range resize {
			snk.updateColumnWidths()
		}
	}()
}

func (snk *prettySink) handleColorState() {
	isBlackOrWhite := snk.config.Theme.BaseColor == theme.White || snk.config.Theme.BaseColor == theme.Black
	if !tty.IsTTY() || isBlackOrWhite {
		cfmt.DisableColors()
	}
}
