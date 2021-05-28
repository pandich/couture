package column

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/pretty/config"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/termenv"
	"os"
	"os/signal"
	"syscall"
)

type widthMode int

type widthWeight uint

const (
	weighted widthMode = iota
	fixed
	filling
)

// Table ...
type Table struct {
	config config.Config
	widths map[string]uint
}

// NewTable ...
func NewTable(config config.Config) *Table {
	for _, name := range config.Columns {
		col := registry[name]
		col.RegisterStyles(config.Theme)
	}
	table := Table{
		config: config,
		widths: map[string]uint{},
	}
	table.updateColumnWidths()
	if config.AutoResize {
		table.autoUpdateColumnWidths()
	}
	return &table
}

// RenderEvent ...
func (table *Table) RenderEvent(event sink.Event) string {
	const resetSequence = termenv.CSI + termenv.ResetSeq + "m"

	// get format string and arguments
	var format = ""
	var values []interface{}
	for _, name := range table.config.Columns {
		col := registry[name]
		format += col.Format(table.widths[name], event)
		values = append(values, col.Render(table.config, event)...)
	}
	format += resetSequence
	var line = cfmt.Sprintf(format, values...)
	if table.config.Wrap {
		line = wordwrap.String(line, int(table.config.EffectiveTerminalWidth()))
	}
	return line
}

func (table *Table) updateColumnWidths() {
	const fixedWidth = 0
	const maxGrowthWidthPercent = 5.0 / 4.0
	const nonMessageAreaWidthPercent = 1.0 / 3.0

	var remainingWidth = table.config.EffectiveTerminalWidth()
	var totalWeight widthWeight
	for _, name := range table.config.Columns {
		col := registry[name]
		switch col.mode() {
		case weighted:
			totalWeight += col.weight()
		case fixed:
			remainingWidth -= uint(col.weight())
		}
	}

	// reserve room for the message
	remainingWidth = uint(float64(remainingWidth) * nonMessageAreaWidthPercent)

	for _, name := range table.config.Columns {
		col := registry[name]
		var width uint = fixedWidth
		switch col.mode() {
		case weighted:
			weighting := float64(col.weight()) / float64(totalWeight)
			width = uint(weighting * float64(remainingWidth))
			max := uint(float64(col.weight()) * maxGrowthWidthPercent)
			if width > max {
				width = max
			}
		case fixed:
			width = uint(col.weight())
		}
		table.widths[name] = width
	}
}

func (table *Table) autoUpdateColumnWidths() {
	resize := make(chan os.Signal)
	signal.Notify(resize, os.Interrupt, syscall.SIGWINCH)
	go func() {
		defer close(resize)
		for {
			<-resize
			table.updateColumnWidths()
		}
	}()
}
