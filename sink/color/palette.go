package color

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	imgcolor "image/color"
)

type adaptorPalette []AdaptorColor

func byColorfulColors(colors []colorful.Color) adaptorPalette {
	var p adaptorPalette
	for _, c := range colors {
		p = append(p, byColorfulColor(c))
	}
	return p
}

// ToPleasingPalette ...
func (rgb rgbAdaptorColor) ToPleasingPalette(count uint) adaptorPalette {
	switch {
	case rgb.IsCool():
		colors, _ := colorful.HappyPalette(int(count))
		return byColorfulColors(colors)

	case rgb.IsWarm():
		colors, _ := colorful.WarmPalette(int(count))
		return byColorfulColors(colors)

	case rgb.IsPastel():
		fallthrough
	default:
		colors, _ := colorful.SoftPalette(int(count))
		return byColorfulColors(colors)
	}
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
