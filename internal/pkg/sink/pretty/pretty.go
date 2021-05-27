package pretty

import (
	"bufio"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/pretty/column"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/source"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/mattn/go-isatty"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/termenv"
	"io"
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
		printer:       newTTYWriter(os.Stdout),
	}
	return &snk
}

// Init ...
func (snk *prettySink) Init(sources []*source.Source) {
	for _, src := range sources {
		column.RegisterSource(snk.config.Theme, snk.config.ConsistentColors, *src)
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
	const resetSequence = termenv.CSI + termenv.ResetSeq + "m"

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
	snk.columnWidths = column.Widths(uint(config.TerminalWidth()), snk.columnOrder)
}

func (snk *prettySink) handleAutoResizeState() {
	resize := make(chan os.Signal)
	signal.Notify(resize, os.Interrupt, syscall.SIGWINCH)
	go func() {
		defer close(resize)
		for {
			<-resize
			snk.updateColumnWidths()
		}
	}()
}

func (snk *prettySink) handleColorState() {
	isTTY := isatty.IsTerminal(os.Stdout.Fd())
	isBlackOrWhite := snk.config.Theme.BaseColor == "#ffffff"
	if !isTTY || isBlackOrWhite {
		cfmt.DisableColors()
	}
}

func newTTYWriter(target io.Writer) chan string {
	delegate := make(chan string)
	go func() {
		defer close(delegate)
		writer := bufio.NewWriter(target)
		for {
			message := <-delegate
			_, err := writer.WriteString(message + "\n")
			if err != nil {
				panic(err)
			}
			err = writer.Flush()
			if err != nil {
				panic(err)
			}
		}
	}()
	return delegate
}
