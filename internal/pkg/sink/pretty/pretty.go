package pretty

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"fmt"
	"github.com/muesli/gamut"
	"github.com/muesli/termenv"
	"github.com/olekukonko/ts"
	"io"
	"sync"
)

// noWrap ...
const noWrap = 0

// prettySink provides render output.
type prettySink struct {
	out              io.Writer
	terminalWidth    int
	sourceStyle      map[model.SourceURL]string
	sourceStyleMutex sync.RWMutex
	sourceColors     chan string
}

// New provides a configured prettySink sink.
func New(out io.Writer, wrap bool) *sink.Sink {
	var terminalWidth = noWrap
	if wrap {
		if size, err := ts.GetSize(); err == nil {
			terminalWidth = size.Col()
		}
	}
	var snk sink.Sink = &prettySink{
		out:              out,
		terminalWidth:    terminalWidth,
		sourceStyleMutex: sync.RWMutex{},
		sourceStyle:      map[model.SourceURL]string{},
		sourceColors:     sink.NewColorCycle(gamut.PastelGenerator{}),
	}
	return &snk
}

// Init ...
func (snk *prettySink) Init(_ []model.SourceURL) {
	termenv.Reset()
	termenv.ClearScreen()
}

// Accept ...
func (snk *prettySink) Accept(src source.Pushable, event model.Event) error {
	line, err := snk.renderEvent(src, event)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(snk.out, line)
	return err
}
