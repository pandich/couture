package sink

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	errors2 "github.com/pkg/errors"
	"image/color"
	"math/rand"
	"time"
)

// NewColorCycle ...
func NewColorCycle(generator gamut.ColorGenerator, defaultColor string) chan string {
	const cycleLength = 50
	rawColors, err := gamut.Generate(cycleLength, generator)
	if err != nil {
		panic(errors2.Wrap(err, "could not generate source color gamut"))
	}
	var colors []string
	for _, rawColor := range rawColors {
		c, ok := colorful.MakeColor(rawColor)
		if ok {
			colors = append(colors, c.Hex())
		} else {
			colors = append(colors, defaultColor)
		}
	}

	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(colors), func(i, j int) { colors[i], colors[j] = colors[j], colors[i] })

	cycle := make(chan string)
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

// Triple ...
func Triple(center color.Color) (string, string, string) {
	surrounding := gamut.Analogous(center)
	leftColor, _ := colorful.MakeColor(surrounding[0])
	centerColor, _ := colorful.MakeColor(center)
	rightColor, _ := colorful.MakeColor(surrounding[1])
	return leftColor.Hex(), centerColor.Hex(), rightColor.Hex()
}

// HexContrast ...
func HexContrast(hex string) string {
	cf, _ := colorful.MakeColor(gamut.Contrast(gamut.Hex(hex)))
	return cf.Hex()
}
