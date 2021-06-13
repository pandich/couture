package color

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut/palette"
	"github.com/muesli/termenv"
	imgcolor "image/color"
)

// White ...
var White = ByHex("#ffffff")

// Black ...
var Black = ByHex("#000000")

var specialNames = map[string]string{
	"prince":    "Logan",
	"halloween": "Burnt Orange",
}

type shades [256]AdaptorColor
type splitComplementary [2]AdaptorColor
type analogous [2]AdaptorColor
type triadic [2]AdaptorColor

type (
	// AdaptorColor ...
	AdaptorColor interface {
		fmt.Stringer
		fmt.GoStringer
		AdjustConstrast(polarity contrastPolarity, amount float64) AdaptorColor
		Analogous() analogous
		AsColorfulColor() colorful.Color
		AsGoColor() color.Color
		AsHexColor() string
		AsHexPair() HexPair
		AsImageColor() imgcolor.Color
		AsPrettyJSONColor() [2]string
		AsRGBColor() color.RGBColor
		AsTermenvColor() termenv.Color
		Blend(other AdaptorColor, blendPercent int) AdaptorColor
		Complementary() AdaptorColor
		Contrast() AdaptorColor
		Darker(percent float64) AdaptorColor
		HueOffset(degrees int) AdaptorColor
		Lighter(percent float64) AdaptorColor
		Monochromatic() shades
		SimiliarHues(count int) []AdaptorColor
		SplitComplementary() splitComplementary
		Triadic() triadic
	}

	rgbColor [3]uint8

	// HexPair ...
	HexPair struct {
		Fg string `yaml:"fg"`
		Bg string `yaml:"bg"`
	}
)

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

// Reverse ...
func (s HexPair) Reverse() HexPair {
	return HexPair{
		Fg: s.Bg,
		Bg: s.Fg,
	}
}

// Format ...
func (s HexPair) Format() func(value string) string {
	return func(value string) string {
		return cfmt.Sprintf("{{%s}}::"+s.Fg+"|bg"+s.Bg, value)
	}
}
