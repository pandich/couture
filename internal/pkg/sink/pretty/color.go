package pretty

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/gamut"
	"github.com/muesli/termenv"
	errors2 "github.com/pkg/errors"
	"math/rand"
	"time"
)

func pastels() chan lipgloss.TerminalColor {
	const colorCount = 50
	return newColorCycle(colorCount, gamut.PastelGenerator{})
}

var (
	colorProfile         = termenv.EnvColorProfile()
	errorColor           = color("#EB1313")
	warnColor            = color("#DCEB3B")
	infoColor            = color("#33A654")
	debugColor           = color("#b37312")
	traceColor           = color("#877150")
	levelForegroundColor = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}
	timestampColor       = color("#117df0")
	applicationNameColor = color("#9357ff")
	threadNameColor      = color("#909090")
	classNameColor       = color("#EBD700")
	methodNameColor      = color("#17C0EB")
	lineNumberColor      = color("#9E9857")
	messageColor         = color("#f2f1da")
)

func color(hex string) lipgloss.TerminalColor {
	return lipgloss.Color(fmt.Sprint(colorProfile.Color(hex)))
}

func newColorCycle(count uint8, generator gamut.ColorGenerator) chan lipgloss.TerminalColor {
	rawColors, err := gamut.Generate(int(count), generator)
	if err != nil {
		panic(errors2.Wrap(err, "could not generate source color gamut"))
	}
	var colors []lipgloss.TerminalColor
	for _, c := range rawColors {
		hex := fmt.Sprint(colorProfile.FromColor(c))
		colors = append(colors, lipgloss.Color(hex))
	}

	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(colors), func(i, j int) { colors[i], colors[j] = colors[j], colors[i] })

	cycle := make(chan lipgloss.TerminalColor)
	go func() {
		var i = 0
		defer close(cycle)
		for {
			cycle <- colors[i]
			i++
			if i >= len(colors) {
				i = 0
			}
		}
	}()
	return cycle
}
