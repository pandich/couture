package color

import (
	"github.com/gookit/color"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"github.com/muesli/termenv"
	imgcolor "image/color"
)

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

// AsHexPair ...
func (rgb rgbColor) AsHexPair() HexPair {
	return HexPair{
		Bg: rgb.AsHexColor(),
		Fg: rgb.Contrast().AsHexColor(),
	}
}
