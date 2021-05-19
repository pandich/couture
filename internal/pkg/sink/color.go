package sink

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	errors2 "github.com/pkg/errors"
	"math/rand"
	"time"
)

// NewColorCycle ...
func NewColorCycle(generator gamut.ColorGenerator) chan string {
	const cycleLength = 50
	rawColors, err := gamut.Generate(cycleLength, generator)
	if err != nil {
		panic(errors2.Wrap(err, "could not generate source color gamut"))
	}
	var colors []string
	for _, rawColor := range rawColors {
		const white = "#ffffff"
		color, ok := colorful.MakeColor(rawColor)
		if ok {
			colors = append(colors, color.Hex())
		} else {
			colors = append(colors, white)
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
