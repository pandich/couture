package pretty

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"couture/internal/pkg/tty"
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/termenv"
	"sync"
)

// TODO themes need to auto light/dark background adjust
// FIXME linebreaks messed up in highlighting process?

// prettySink provides render output.
type prettySink struct {
	terminalWidth int
	palette       chan string
	columnOrder   []ColumnName
	printLock     sync.Mutex
	config        Config
}

// New provides a configured prettySink sink.
func New(config Config) *sink.Sink {
	if !tty.IsTTY() || config.Theme.BaseColor == "" {
		cfmt.DisableColors()
	}
	columns.init(config.Theme)
	var snk sink.Sink = &prettySink{
		terminalWidth: config.wrapWidth(),
		palette:       tty.NewColorCycle(config.Theme.SourceColors),
		columnOrder:   config.effectiveColumns(),
		printLock:     sync.Mutex{},
		config:        config,
	}
	return &snk
}

// Init ...
func (snk *prettySink) Init(sources []source.Source) {
	for _, src := range sources {
		registerSourceStyle(src, <-snk.palette)
	}
	termenv.Reset()
	termenv.ClearScreen()
}

// Accept ...
func (snk *prettySink) Accept(src source.Source, event model.Event) error {
	snk.printLock.Lock()
	termenv.Reset()
	defer snk.printLock.Unlock()

	// get format string and arguments
	var format = ""
	var values []interface{}
	for _, name := range snk.columnOrder {
		col := columns[name]
		format += col.formatter(src, event)
		values = append(values, col.renderer(snk.config, src, event))
	}

	// render
	line := cfmt.Sprintf(format, values...)

	// print
	fmt.Println(tty.Wrap(line, snk.terminalWidth))
	return nil
}
