package pretty

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"fmt"
	"github.com/muesli/gamut"
	"github.com/olekukonko/ts"
	"io"
	"sync"
)

// AutoWrap ...
const AutoWrap = -1

// NoWrap ...
const NoWrap = 0

// prettySink provides render output.
type prettySink struct {
	out              io.Writer
	terminalWidth    int
	sourceStyle      map[model.SourceURL]string
	sourceStyleMutex sync.RWMutex
	sourceColors     chan string
}

// New provides a configured prettySink sink.
func New(out io.Writer, wrap int) *sink.Sink {
	var terminalWidth = 72
	switch wrap {
	case AutoWrap:
		if size, err := ts.GetSize(); err == nil {
			terminalWidth = size.Col()
		}
	case NoWrap:
		terminalWidth = NoWrap
	default:
		terminalWidth = wrap
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

// Accept ...
func (snk *prettySink) Accept(src source.Pushable, event model.Event) error {
	line, err := snk.renderEvent(src, event)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(snk.out, line)
	return err
}
