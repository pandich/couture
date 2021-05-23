package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"fmt"
)

type weightType int

type weight uint

const (
	weighted weightType = iota
	fixed
	filling
)

type column interface {
	name() string
	RegisterStyles(th theme.Theme)
	Format(width uint, src source.Source, event model.Event) string
	Render(cfg config.Config, src source.Source, event model.Event) []interface{}
	weightType() weightType
	weight() weight
}

var columns = []column{
	sourceColumn{},
	timestampColumn{},
	applicationColumn{},
	threadColumn{},
	callerColumn{},
	levelColumn{},
	messageColumn{},
	stackTraceColumn{},
}

// Widths ...
func Widths(terminalWidth uint, columnNames []string) map[string]uint {
	const maxGrowthWidthPercent = 4 / 3
	const nonMessageAreaWidthPercent = 0.4
	const fixedWidth = 0

	var remainingWidth = terminalWidth
	var totalWeight weight
	for _, name := range columnNames {
		col := ByName[name]
		switch col.weightType() {
		case weighted:
			totalWeight += col.weight()
		case fixed:
			remainingWidth -= uint(col.weight())
		}
	}

	// reserve room for the message
	remainingWidth = uint(float64(remainingWidth) * nonMessageAreaWidthPercent)

	var widths = map[string]uint{}
	for _, name := range columnNames {
		col := ByName[name]
		var width uint = fixedWidth
		switch col.weightType() {
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
		widths[name] = width
	}

	return widths
}

func formatStringOfWidth(width uint) string {
	return fmt.Sprintf("%%-%[1]d.%[1]ds", width)
}
