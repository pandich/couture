package color

import (
	"github.com/gookit/color"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"github.com/muesli/termenv"
	imgcolor "image/color"
)

// AsRGBColor ...
func (rgb rgbAdaptorColor) AsRGBColor() color.RGBColor {
	return color.RGB(rgb[0], rgb[1], rgb[2])
}

// AsHexColor ...
func (rgb rgbAdaptorColor) AsHexColor() string {
	return "#" + rgb.AsRGBColor().Hex()
}

// AsGamutColor ...
func (rgb rgbAdaptorColor) AsGamutColor() gamut.Color {
	return gamut.Color{Color: rgb.AsImageColor()}
}

// AsGooKitColor ...
func (rgb rgbAdaptorColor) AsGooKitColor() color.Color {
	return rgb.AsRGBColor().Color()
}

// AsImageColor ...
func (rgb rgbAdaptorColor) AsImageColor() imgcolor.Color {
	return gamut.Hex(rgb.AsHexColor())
}

// AsColorfulColor ...
func (rgb rgbAdaptorColor) AsColorfulColor() colorful.Color {
	cf, _ := colorful.Hex(rgb.AsHexColor())
	return cf
}

// AsTermenvColor ...
func (rgb rgbAdaptorColor) AsTermenvColor() termenv.Color {
	return termenv.ANSI256.FromColor(rgb.AsImageColor())
}

// AsPrettyJSONColor ...
func (rgb rgbAdaptorColor) AsPrettyJSONColor() [2]string {
	start := termenv.CSI + rgb.AsTermenvColor().Sequence(false) + "m"
	end := termenv.CSI + "39m"
	return [2]string{start, end}
}

// AsHexPair ...
func (rgb rgbAdaptorColor) AsHexPair() FgBgTuple {
	return FgBgTuple{
		Bg: rgb.AsHexColor(),
		Fg: rgb.Contrast().AsHexColor(),
	}
}
