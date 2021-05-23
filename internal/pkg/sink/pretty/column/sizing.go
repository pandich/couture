package column

import (
	"fmt"
)

type widthMode int

type widthWeight uint

const (
	weighted widthMode = iota
	fixed
	filling
)

// Widths ...
func Widths(terminalWidth uint, columnNames []string) map[string]uint {
	const fixedWidth = 0
	const maxGrowthWidthPercent = 5.0 / 4.0
	const nonMessageAreaWidthPercent = 1.0 / 3.0

	var remainingWidth = terminalWidth
	var totalWeight widthWeight
	for _, name := range columnNames {
		col := ByName[name]
		switch col.mode() {
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
		widths[name] = width
	}

	return widths
}

func formatStringOfWidth(width uint) string {
	if width <= 0 {
		return "%s"
	}
	return fmt.Sprintf("%%-%[1]d.%[1]ds", width)
}

func formatStyleOfWidth(style string, width uint) string {
	return "{{" + formatStringOfWidth(width) + "}}::" + style
}
func formatColumn(col column, width uint) string {
	return formatStyleOfWidth(col.name(), width)
}
