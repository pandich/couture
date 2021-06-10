package column

import (
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/termenv"
	"github.com/pandich/couture/internal/pkg/model"
	"github.com/pandich/couture/internal/pkg/model/level"
	"github.com/pandich/couture/internal/pkg/sink"
	"os"
	"os/signal"
	"syscall"
)

// Table ...
type Table struct {
	config   sink.Config
	widths   map[string]uint
	registry registry
}

// NewTable ...
func NewTable(config sink.Config) *Table {
	registry := newRegistry(config)
	table := Table{
		config:   config,
		widths:   map[string]uint{},
		registry: registry,
	}
	table.updateColumnWidths()
	if config.AutoResize != nil && *config.AutoResize {
		table.autoUpdateColumnWidths()
	}
	return &table
}

// Render ...
func (table *Table) Render(event model.SinkEvent) string {
	const resetSequence = termenv.CSI + termenv.ResetSeq + "m"
	if event.Level == "" {
		event.Level = level.Default
	}
	// get format string and arguments
	var line string
	for _, name := range table.config.Columns {
		col := table.registry[name]
		line += col.render(event) + resetSequence
	}
	if table.config.Wrap != nil && *table.config.Wrap {
		line = wordwrap.String(line, int(table.config.EffectiveTerminalWidth()))
	}
	return line
}

func (table *Table) updateColumnWidths() {
	const maxGrowthWidthPercent = 5.0 / 4.0
	const nonMessageAreaWidthPercent = 1.0 / 3.0

	var remainingWidth = table.config.EffectiveTerminalWidth()
	var totalWeight uint
	for _, name := range table.config.Columns {
		col := table.registry[name]
		totalWeight += col.layout().Width
	}

	// reserve room for the message
	remainingWidth = uint(float64(remainingWidth) * nonMessageAreaWidthPercent)

	for _, name := range table.config.Columns {
		col := table.registry[name]
		weigth := col.layout().Width
		weighting := float64(weigth) / float64(totalWeight)
		var width = uint(weighting * float64(remainingWidth))
		max := uint(float64(weigth) * maxGrowthWidthPercent)
		if width > max {
			width = max
		}
		table.widths[name] = width
	}
}

func (table *Table) autoUpdateColumnWidths() {
	resize := make(chan os.Signal, 1)
	signal.Notify(resize, os.Interrupt, syscall.SIGWINCH)
	go func() {
		defer close(resize)
		for {
			<-resize
			table.updateColumnWidths()
		}
	}()
}
