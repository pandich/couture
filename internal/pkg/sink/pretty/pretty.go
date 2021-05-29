package pretty

import (
	"bufio"
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/pretty/column"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/source"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/mattn/go-isatty"
	"os"
)

// Name ...
const Name = "pretty"

// prettySink provides render output.
type prettySink struct {
	terminalWidth uint
	table         *column.Table
	config        config.Config
	out           chan string
}

// New provides a configured prettySink sink.
func New(cfg config.Config) *sink.Sink {
	isTTY := isatty.IsTerminal(os.Stdout.Fd())
	isBlackOrWhite := cfg.Theme.BaseColor == "#ffffff"
	if !isTTY || isBlackOrWhite {
		cfmt.DisableColors()
	}

	if len(cfg.Columns) == 0 {
		cfg.Columns = column.DefaultColumns
	}

	var snk sink.Sink = &prettySink{
		terminalWidth: cfg.EffectiveTerminalWidth(),
		table:         column.NewTable(cfg),
		config:        cfg,
		out:           newOut(),
	}
	return &snk
}

// Init ...
func (snk *prettySink) Init(sources []*source.Source) {
	for _, src := range sources {
		column.RegisterSource(snk.config.Theme, snk.config.ConsistentColors, *src)
	}
}

// Accept ...
func (snk *prettySink) Accept(event model.SinkEvent) error {
	snk.out <- snk.table.RenderEvent(event)
	return nil
}

func newOut() chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		writer := bufio.NewWriter(os.Stdout)
		for {
			message := <-out
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
	return out
}
