package sink

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"github.com/muesli/termenv"
	errors2 "github.com/pkg/errors"
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

// Contrast ...
func Contrast(hex string) string {
	cf, _ := colorful.MakeColor(gamut.Contrast(gamut.Hex(hex)))
	return cf.Hex()
}

// Lighter ...
func Lighter(hex string, percent float64) string {
	cf, _ := colorful.MakeColor(gamut.Lighter(gamut.Hex(hex), percent))
	return cf.Hex()
}

// Darker ...
func Darker(hex string, percent float64) string {
	cf, _ := colorful.MakeColor(gamut.Darker(gamut.Hex(hex), percent))
	return cf.Hex()
}

// Fainter ...
func Fainter(hex string, percent float64) string {
	const count = 1000
	i := int(count * percent)
	bg := termenv.ConvertToRGB(termenv.BackgroundColor())
	fainter := gamut.Blends(bg, gamut.Hex(hex), count)[count-i]
	col, _ := colorful.MakeColor(fainter)
	return col.Hex()
}

// IsDark ...
func IsDark(hex string) bool {
	const gray = 0.5
	col, _ := colorful.MakeColor(gamut.Hex(hex))
	_, _, l := col.Hcl()
	return l < gray
}

// WithFaintBg ...
func WithFaintBg(hex string) string {
	return hex + "|bg" + Fainter(hex, 0.96)
}
