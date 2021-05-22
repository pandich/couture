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
	"github.com/muesli/termenv"
	"sync"
)

// TODO themes need to auto light/dark background adjust
// FIXME linebreaks messed up in highlighting process?

// prettySink provides render output.
type prettySink struct {
	terminalWidth int
	palette       chan string
	columnOrder   []string
	printLock     sync.Mutex
	config        config.Config
}

// New provides a configured prettySink sink.
func New(cfg config.Config) *sink.Sink {
	if !tty.IsTTY() || cfg.Theme.BaseColor == theme.White {
		cfmt.DisableColors()
	}
	column.ByName.Init(cfg.Theme)
	var effectiveColumns = column.DefaultOrder
	if len(cfg.Columns) > 0 {
		effectiveColumns = cfg.Columns
	}
	effectiveColumns = append([]string{column.SourceColumn}, effectiveColumns...)
	var snk sink.Sink = &prettySink{
		terminalWidth: cfg.WrapWidth(),
		palette:       tty.NewColorCycle(cfg.Theme.SourceColors),
		columnOrder:   effectiveColumns,
		printLock:     sync.Mutex{},
		config:        cfg,
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
	snk.println(snk.render(src, event))
	return nil
}

func (snk *prettySink) render(src source.Source, event model.Event) string {
	// get format string and arguments
	var format = ""
	var values []interface{}
	for _, name := range snk.columnOrder {
		col := column.ByName[name]
		format += col.Formatter(src, event)
		values = append(values, col.Renderer(snk.config, src, event)...)
	}

	// render
	line := cfmt.Sprintf(format, values...)
	wrapped := tty.Wrap(line, snk.terminalWidth)
	return wrapped
}

func (snk *prettySink) println(s string) {
	snk.printLock.Lock()
	termenv.Reset()
	defer snk.printLock.Unlock()
	fmt.Println(s)
}
