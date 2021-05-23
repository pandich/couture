package tty

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"github.com/muesli/termenv"
	errors2 "github.com/pkg/errors"
	"math/rand"
	"time"
)

// NewColorCycle ...
func NewColorCycle(generator gamut.ColorGenerator) chan string {
	colors := newColorPool(generator)
	cycle := make(chan string)
	go func() {
		defer close(cycle)
		var i = 0
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
//goland:noinspection GoUnusedExportedFunction,GoUnnecessarilyExportedIdentifiers
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
//goland:noinspection GoUnusedExportedFunction,GoUnnecessarilyExportedIdentifiers
func IsDark(hex string) bool {
	const gray = 0.5
	col := hexColorful(hex)
	_, _, l := col.Hcl()
	return l < gray
}

// SimilarBg ...
func SimilarBg(hex string) string {
	return hex + "|bg" + Fainter(hex, 0.96)
}

// Tinted ...
func Tinted(baseHex string, hex string) string {
	complementary, _ := colorful.MakeColor(gamut.Complementary(gamut.Hex(baseHex)))
	return complementary.BlendHcl(hexColorful(hex), 0.85).Hex()
}

func hexColorful(hex string) colorful.Color {
	cl, _ := colorful.MakeColor(gamut.Hex(hex))
	return cl
}

func newColorPool(generator gamut.ColorGenerator) []string {
	const cycleLength = 50
	rawColors, err := gamut.Generate(cycleLength, generator)
	if err != nil {
		panic(errors2.Wrap(err, "could not generate source color gamut"))
	}
	var colors []string
	for _, rawColor := range rawColors {
		c, _ := colorful.MakeColor(rawColor)
		colors = append(colors, c.Hex())
	}

	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(colors), func(i, j int) { colors[i], colors[j] = colors[j], colors[i] })
	return colors
}
