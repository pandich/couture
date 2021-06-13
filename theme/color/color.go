package color

// TODO use this to deal with color in a clean central way

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"github.com/muesli/gamut/palette"
	"github.com/muesli/termenv"
	imgcolor "image/color"
)

var (
	colorfulWhite, _ = colorful.Hex("#ffffff")
	colorfulBlack, _ = colorful.Hex("#000000")
)

// Hex ...
func Hex(hex string) AdaptorColor {
	values := color.Hex(hex).Values()
	r, g, b := uint8(values[0]), uint8(values[1]), uint8(values[2])
	return rgbColor{r, g, b}
}

func byName(name string) (AdaptorColor, bool) {
	if c, ok := palette.AllPalettes().Color(name); ok {
		cf, _ := colorful.MakeColor(c)
		return Hex(cf.Hex()), true
	}
	return nil, false
}

// FromImageColor ...
func FromImageColor(imgColor imgcolor.Color) AdaptorColor {
	cf, _ := colorful.MakeColor(imgColor)
	return Hex(cf.Hex())
}

// MustByName ...
func MustByName(name string) AdaptorColor {
	if c, ok := byName(name); ok {
		return c
	}
	panic(name)
}

// AdaptorColor ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
type AdaptorColor interface {
	fmt.Stringer
	fmt.GoStringer
	AsColorfulColor() colorful.Color
	AsGoColor() color.Color
	AsHexColor() string
	AsImageColor() imgcolor.Color
	AsPrettyJSONColor() [2]string
	AsRGBColor() color.RGBColor
	AsTermenvColor() termenv.Color
	AdjustConstrast(mode ContrastPolarity, polarity ContrastPolarity, amount float64) AdaptorColor
	Blend(other AdaptorColor, blendPercent int) AdaptorColor
	Contrast() AdaptorColor
}

type rgbColor [3]uint8

// AsRGBColor ...
func (rgb rgbColor) AsRGBColor() color.RGBColor {
	return color.RGB(rgb[0], rgb[1], rgb[2])
}

// AsHexColor ...
func (rgb rgbColor) AsHexColor() string {
	return "#" + rgb.AsRGBColor().Hex()
}

// AsGoColor ...
func (rgb rgbColor) AsGoColor() color.Color {
	return rgb.AsRGBColor().Color()
}

// AsImageColor ...
func (rgb rgbColor) AsImageColor() imgcolor.Color {
	return gamut.Hex(rgb.AsHexColor())
}

// AsColorfulColor ...
func (rgb rgbColor) AsColorfulColor() colorful.Color {
	cf, _ := colorful.Hex(rgb.AsHexColor())
	return cf
}

// AsTermenvColor ...
func (rgb rgbColor) AsTermenvColor() termenv.Color {
	return termenv.ANSI256.FromColor(rgb.AsImageColor())
}

// AsPrettyJSONColor ...
func (rgb rgbColor) AsPrettyJSONColor() [2]string {
	start := termenv.CSI + rgb.AsTermenvColor().Sequence(false) + "m"
	end := termenv.CSI + "39m"
	return [2]string{start, end}
}

// String ...
func (rgb rgbColor) String() string {
	return rgb.AsColorfulColor().Hex()
}

// GoString ...
func (rgb rgbColor) GoString() string {
	if names, _ := palette.AllPalettes().Name(rgb.AsColorfulColor()); len(names) > 0 {
		return names[0].Name
	}
	return rgb.String()
}

// AdjustConstrast ...
func (rgb rgbColor) AdjustConstrast(mode ContrastPolarity, polarity ContrastPolarity, amount float64) AdaptorColor {
	var base colorful.Color
	switch mode {
	case DarkMode:
		switch polarity {
		case MoreContrast:
			base = colorfulBlack
		case LessContrast:
			base = colorfulWhite
		}
	case LightMode:
		switch polarity {
		case MoreContrast:
			base = colorfulWhite
		case LessContrast:
			base = colorfulBlack
		}
	}
	blended := rgb.AsColorfulColor().BlendHsv(base, amount)
	return Hex(blended.Hex())
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
		return FromImageColor(blends[blendPercent])
	}
}

// Contrast ...
func (rgb rgbColor) Contrast() AdaptorColor {
	return FromImageColor(gamut.Contrast(rgb.AsImageColor()))
}
