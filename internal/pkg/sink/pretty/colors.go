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

var (
	colorProfile         = termenv.EnvColorProfile()
	errorColor           = color("#EB1313")
	warnColor            = color("#DCEB3B")
	infoColor            = color("#33A654")
	debugColor           = color("#b37312")
	traceColor           = color("#877150")
	levelForegroundColor = color("#ffffff")
	timestampColor       = color("#117df0")
	applicationNameColor = color("#9357ff")
	threadNameColor      = color("#56A8F8")
	classNameColor       = color("#EBD700")
	methodNameColor      = color("#17C0EB")
	lineNumberColor      = color("#9E9857")
	messageColor         = color("#f2f1da")
)

func color(hex string) lipgloss.Color {
	return lipgloss.Color(fmt.Sprint(colorProfile.Color(hex)))
}

func newColorCycle(count uint8, generator gamut.ColorGenerator) chan lipgloss.Color {
	rawColors, err := gamut.Generate(int(count), generator)
	if err != nil {
		panic(errors2.Wrap(err, "could not generate source color gamut"))
	}
	var colors []lipgloss.Color
	for _, c := range rawColors {
		hex := fmt.Sprint(colorProfile.FromColor(c))
		colors = append(colors, lipgloss.Color(hex))
	}

	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(colors), func(i, j int) { colors[i], colors[j] = colors[j], colors[i] })

	cycle := make(chan lipgloss.Color)
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
