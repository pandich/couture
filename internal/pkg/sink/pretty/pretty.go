package pretty

import (
	"bufio"
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/gamut"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/termenv"
	"github.com/olekukonko/ts"
	"io"
	"sync"
)

// TODO configuration column widths
// TODO adaptive column widths
// FIXME column widths are bad
// FIXME linebreaks messed up in highlighting process?

// noWrap ...
const noWrap = 0

// Config ...
type (
	// Config ...
	Config struct {
		Wrap    bool
		Theme   Theme
		Columns []ColumnName
	}
	// Theme ...
	Theme struct {
		BaseColor        string
		ApplicationColor string
		DefaultColor     string
		TimestampColor   string
		ErrorColor       string
		TraceColor       string
		DebugColor       string
		InfoColor        string
		WarnColor        string
		MessageColor     string
		StackTraceColor  string
		SourceColors     gamut.ColorGenerator
	}
	// prettySink provides render output.
	prettySink struct {
		out           io.Writer
		terminalWidth int
		palette       chan string
		columnOrder   []ColumnName
		printLock     sync.Mutex
		flusher       func()
	}
)

// New provides a configured prettySink sink.
func New(out io.Writer, config Config) *sink.Sink {
	theme := config.Theme
	if !sink.IsTTY() || theme.BaseColor == "" {
		cfmt.DisableColors()
	}
	for _, col := range columns {
		col.register(theme)
	}

	var columnOrder = defaultColumnOrder
	if len(config.Columns) > 0 {
		columnOrder = config.Columns
	}
	columnOrder = append([]ColumnName{sourceColumn}, columnOrder...)
	var snk sink.Sink = &prettySink{
		out:           out,
		terminalWidth: terminalWidth(config.Wrap),
		palette:       sink.NewColorCycle(theme.SourceColors, theme.DefaultColor),
		columnOrder:   columnOrder,
		printLock:     sync.Mutex{},
		flusher:       func() { _ = bufio.NewWriter(out).Flush() },
	}
	return &snk
}

// Init ...
func (snk *prettySink) Init(sources []source.Source) {
	for _, src := range sources {
		styleColor := <-snk.palette
		styleName := "source" + src.ID()
		cfmt.RegisterStyle(styleName, func(s string) string { return cfmt.Sprintf("{{/%-30.30s }}::"+styleColor, s) })
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
	snk.printLock.Lock()
	defer fmt.Fprintln(snk.out, termenv.CSI+termenv.ResetSeq+"m")
	defer snk.flusher()
	defer snk.printLock.Unlock()

	line, err := snk.renderEvent(src, event)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprint(snk.out, line)
	return nil
}

func (snk *prettySink) renderEvent(src source.Source, event model.Event) (string, error) {
	var format = ""
	var values []interface{}

	for _, colName := range snk.columnOrder {
		col := columns[colName]
		format += col.format(src, event)
		values = append(values, col.value(src, event))
	}

	line := cfmt.Sprintf(format, values...)
	return snk.wrapToTerminal(line)
}

func (snk *prettySink) wrapToTerminal(s string) (string, error) {
	if snk.terminalWidth == noWrap {
		return s, nil
	}
	wrapper := wordwrap.NewWriter(snk.terminalWidth)
	wrapper.Breakpoints = []rune(" \t")
	wrapper.KeepNewlines = true
	if _, err := wrapper.Write([]byte(s)); err != nil {
		return "", err
	}
	wrapped := wrapper.String()
	return wrapped, nil
}
