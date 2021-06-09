package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/theme"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/sink/pretty/config"
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
	config   config.Config
	widths   map[string]uint
	registry map[string]column
}

// NewTable ...
func NewTable(config config.Config) *Table {
	registry := map[string]column{
		"source":           newSourceColumn(config),
		schema.Timestamp:   newTimestampColumn(config),
		schema.Application: newApplicationColumn(config),
		schema.Context:     newContextColumn(config),
		"caller":           newCallerColumn(config),
		schema.Level:       newLevelColumn(config),
		schema.Message:     newMessageColumn(config),
	}
	for _, name := range config.Columns {
		col := registry[name]
		if config.Theme == nil {
			th := theme.Registry[theme.Prince]
			config.Theme = &th
		}
		col.Init(*config.Theme)
	}
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

// RenderEvent ...
func (table *Table) RenderEvent(event model.SinkEvent) string {
	const resetSequence = termenv.CSI + termenv.ResetSeq + "m"

	// get format string and arguments
	var line string
	for _, name := range table.config.Columns {
		col := table.registry[name]
		line += col.Render(table.config, event) + resetSequence
	}
	if table.config.Wrap != nil && *table.config.Wrap {
		line = wordwrap.String(line, int(table.config.EffectiveTerminalWidth()))
	}
	return line
}

func (table *Table) updateColumnWidths() {
	const fixedWidth = 0
	const maxGrowthWidthPercent = 5.0 / 4.0
	const nonMessageAreaWidthPercent = 1.0 / 3.0

	var remainingWidth = widthWeight(table.config.EffectiveTerminalWidth())
	var totalWeight widthWeight
	for _, name := range table.config.Columns {
		col := table.registry[name]
		weight := widthWeight(col.layout().Width)
		switch col.mode() {
		case weighted:
			totalWeight += weight
		case fixed:
			remainingWidth -= weight
		}
	}

	// reserve room for the message
	remainingWidth = widthWeight(float64(remainingWidth) * nonMessageAreaWidthPercent)

	for _, name := range table.config.Columns {
		col := table.registry[name]
		var width uint = fixedWidth
		switch col.mode() {
		case weighted:
			weigth := col.layout().Width
			weighting := float64(weigth) / float64(totalWeight)
			width = uint(weighting * float64(remainingWidth))
			max := uint(float64(weigth) * maxGrowthWidthPercent)
			if width > max {
				width = max
			}
		case fixed:
			width = col.layout().Width
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
