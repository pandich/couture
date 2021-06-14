package color

import (
	"github.com/muesli/gamut"
	imgcolor "image/color"
)

type adaptorPalette []AdaptorColor

// PleasingPalette ...
func (rgb rgbAdaptorColor) PleasingPalette(colorCount uint) adaptorPalette {
	const distanceCutoff = 0.7

	colorDistance := rgb.DistancesRgb()
	meetsDistanceCutoff := colorDistance.min() > distanceCutoff

	var generator gamut.ColorGenerator
	switch {
	// cool or blue-dominant
	case rgb.IsCool(),
		meetsDistanceCutoff && colorDistance.closestToBlue():
		generator = gamut.PastelGenerator{}

	// pastel or green-dominant
	case rgb.IsPastel(),
		meetsDistanceCutoff && colorDistance.closestToGreen():
		generator = gamut.WarmGenerator{}

	// warm
	case rgb.IsWarm():
		switch {
		// red dominant
		case meetsDistanceCutoff && colorDistance.closestToRed():
			generator = gamut.HappyGenerator{}
		// other warm
		default:
			generator = gamut.WarmGenerator{}
		}

	// default to pastel
	default:
		generator = gamut.PastelGenerator{}
	}

	paletteColors, _ := gamut.Generate(int(colorCount), generator)

	var out adaptorPalette
	for _, pc := range paletteColors {
		out = append(out, byImageColor(pc).Blend(rgb.Complementary(), 20))
	}
	return out
}

// Clamped ...
func (ap adaptorPalette) Clamped(pal gamut.Palette) adaptorPalette {
	var colors []imgcolor.Color
	for _, c := range ap {
		colors = append(colors, c.AsImageColor())
	}
	var out []AdaptorColor
	for _, c := range pal.Clamped(colors) {
		out = append(out, byGamutColor(c))
	}
	return out
}
