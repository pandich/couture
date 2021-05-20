package pretty

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"fmt"
	"github.com/muesli/termenv"
	"github.com/olekukonko/ts"
	"image/color"
	"io"
)

// TODO color schemes

// TODO configuration column widths
// TODO adaptive column widths
// TODO configurable column order

// FIXME source URLs don't display well
// FIXME column widths are bad
// FIXME linebreaks messed up in highlighting process?

// noWrap ...
const noWrap = 0

// prettySink provides render output.
type prettySink struct {
	out           io.Writer
	terminalWidth int
	palette       palette
}

// New provides a configured prettySink sink.
func New(out io.Writer, wrap bool, baseColor color.Color) *sink.Sink {
	var snk sink.Sink = &prettySink{
		out:           out,
		terminalWidth: terminalWidth(wrap),
		palette:       newPalette(baseColor),
	}
	return &snk
}

// Init ...
func (snk *prettySink) Init(sources []source.Source) {
	for _, src := range sources {
		snk.palette.registerSource(src.URL())
	}
	termenv.Reset()
	termenv.ClearScreen()
}

func terminalWidth(wrap bool) int {
	var terminalWidth = noWrap
	if wrap {
		if size, err := ts.GetSize(); err == nil {
			terminalWidth = size.Col()
		}
	}
	return terminalWidth
}

// Accept ...
func (snk *prettySink) Accept(src source.Source, event model.Event) error {
	line, err := snk.renderEvent(src, event)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(snk.out, line)
	return err
}
