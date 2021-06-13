package color

import (
	"github.com/muesli/gamut"
)

// Complementary ...
func (rgb rgbAdaptorColor) Complementary() AdaptorColor {
	return ByImageColor(gamut.Complementary(rgb.AsImageColor()))
}

// Analogous ...
func (rgb rgbAdaptorColor) Analogous() analogous {
	raw := gamut.Analogous(rgb.AsImageColor())
	return analogous{
		ByImageColor(raw[0]),
		ByImageColor(raw[1]),
	}
}

// Triadic ...
func (rgb rgbAdaptorColor) Triadic() triadic {
	raw := gamut.Triadic(rgb.AsImageColor())
	return triadic{
		ByImageColor(raw[0]),
		ByImageColor(raw[1]),
	}
}

// Lighter ...
func (rgb rgbAdaptorColor) Lighter(percent percent) AdaptorColor {
	return ByImageColor(gamut.Lighter(rgb.AsImageColor(), percent.asFloat64()))
}

// Darker ...
func (rgb rgbAdaptorColor) Darker(percent percent) AdaptorColor {
	return ByImageColor(gamut.Darker(rgb.AsImageColor(), percent.asFloat64()))
}

// SplitComplementary ...
func (rgb rgbAdaptorColor) SplitComplementary() splitComplementary {
	raw := gamut.SplitComplementary(rgb.AsImageColor())
	return splitComplementary{
		ByImageColor(raw[0]),
		ByImageColor(raw[1]),
	}
}

// Monochromatic ...
func (rgb rgbAdaptorColor) Monochromatic() shades {
	const count = 256
	imageColors := gamut.Monochromatic(rgb.AsImageColor(), count)
	colors := shades{}
	for i, imageColor := range imageColors {
		colors[i] = ByImageColor(imageColor)
	}
	return colors
}

// AdjustConstrast ...
func (rgb rgbAdaptorColor) AdjustConstrast(polarity contrastPolarity, percent percent) AdaptorColor {
	var base AdaptorColor
	switch polarity {
	case MoreContrast:
		switch Mode {
		case DarkMode:
			base = White
		default:
			base = Black
		}
	case LessContrast:
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
		BlendHsv(base.AsColorfulColor(), percent.asFloat64()).
		Hex())
}

// Blend ...
func (rgb rgbAdaptorColor) Blend(other AdaptorColor, blendPercent percent) AdaptorColor {
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
func (rgb rgbAdaptorColor) Contrast() AdaptorColor {
	return ByImageColor(gamut.Contrast(rgb.AsImageColor()))
}

// HueOffset ...
func (rgb rgbAdaptorColor) HueOffset(degrees º) AdaptorColor {
	return ByImageColor(gamut.HueOffset(rgb.AsImageColor(), degrees.asInt()))
}

// Similar ...
func (rgb rgbAdaptorColor) Similar(count int) []AdaptorColor {
	return rgb.generate(count, rgb.similarHueGenerator())
}

func (rgb rgbAdaptorColor) generate(count int, generator gamut.ColorGenerator) []AdaptorColor {
	colors, err := gamut.Generate(count, generator)
	if err != nil {
		panic(err)
	}
	return byImageColors(colors)
}

func (rgb rgbAdaptorColor) similarHueGenerator() gamut.ColorGenerator {
	return gamut.SimilarHueGenerator{
		FineGranularity: gamut.FineGranularity{},
		Color:           rgb.AsImageColor(),
	}
}
