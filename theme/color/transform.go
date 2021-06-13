package color

import (
	"github.com/muesli/gamut"
)

// Complementary ...
func (rgb rgbColor) Complementary() AdaptorColor {
	return ByImageColor(gamut.Complementary(rgb.AsImageColor()))
}

// Analogous ...
func (rgb rgbColor) Analogous() analogous {
	raw := gamut.Analogous(rgb.AsImageColor())
	return analogous{
		ByImageColor(raw[0]),
		ByImageColor(raw[1]),
	}
}

// Triadic ...
func (rgb rgbColor) Triadic() triadic {
	raw := gamut.Triadic(rgb.AsImageColor())
	return triadic{
		ByImageColor(raw[0]),
		ByImageColor(raw[1]),
	}
}

// Lighter ...
func (rgb rgbColor) Lighter(percent float64) AdaptorColor {
	return ByImageColor(gamut.Lighter(rgb.AsImageColor(), percent))
}

// Darker ...
func (rgb rgbColor) Darker(percent float64) AdaptorColor {
	return ByImageColor(gamut.Darker(rgb.AsImageColor(), percent))
}

// SplitComplementary ...
func (rgb rgbColor) SplitComplementary() splitComplementary {
	raw := gamut.SplitComplementary(rgb.AsImageColor())
	return splitComplementary{
		ByImageColor(raw[0]),
		ByImageColor(raw[1]),
	}
}

// Monochromatic ...
func (rgb rgbColor) Monochromatic() shades {
	const count = 256
	imageColors := gamut.Monochromatic(rgb.AsImageColor(), count)
	colors := shades{}
	for i, imageColor := range imageColors {
		colors[i] = ByImageColor(imageColor)
	}
	return colors
}

// AdjustConstrast ...
func (rgb rgbColor) AdjustConstrast(polarity contrastPolarity, amount float64) AdaptorColor {
	var base AdaptorColor
	switch polarity {
	case LessNoticable:
		switch Mode {
		case DarkMode:
			base = White
		default:
			base = Black
		}
	case MoreNoticable:
		fallthrough
	default:
		switch Mode {
		case DarkMode:
			base = Black
		default:
			base = White
		}
	}
	return ByHex(rgb.
		AsColorfulColor().
		BlendHsv(base.AsColorfulColor(), amount).
		Hex())
}

// Blend ...
func (rgb rgbColor) Blend(other AdaptorColor, blendPercent int) AdaptorColor {
	const minPercent = 0
	const maxPercent = 100
	switch {
	case blendPercent <= minPercent:
		return rgb
	case blendPercent >= maxPercent:
		return other
	default:

		const blendCount = 100
		blends := gamut.Blends(rgb.AsImageColor(), other.AsImageColor(), blendCount)
		return ByImageColor(blends[blendPercent])
	}
}

// Contrast ...
func (rgb rgbColor) Contrast() AdaptorColor {
	return ByImageColor(gamut.Contrast(rgb.AsImageColor()))
}

// HueOffset ...
func (rgb rgbColor) HueOffset(degrees int) AdaptorColor {
	return ByImageColor(gamut.HueOffset(rgb.AsImageColor(), degrees))
}

// SimiliarHues ...
func (rgb rgbColor) SimiliarHues(count int) []AdaptorColor {
	return rgb.generate(count, rgb.similarHueGenerator())
}

func (rgb rgbColor) generate(count int, generator gamut.ColorGenerator) []AdaptorColor {
	colors, err := gamut.Generate(count, generator)
	if err != nil {
		panic(err)
	}
	return byImageColors(colors)
}

func (rgb rgbColor) similarHueGenerator() gamut.ColorGenerator {
	return gamut.SimilarHueGenerator{
		FineGranularity: gamut.FineGranularity{},
		Color:           rgb.AsImageColor(),
	}
}
