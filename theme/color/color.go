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

type shades [256]AdaptorColor
type splitComplementary [2]AdaptorColor
type analogous [2]AdaptorColor
type triadic [2]AdaptorColor

type (
	percent uint8
	º       int

	// AdaptorColor ...
	AdaptorColor interface {
		fmt.Stringer
		fmt.GoStringer
		AdjustConstrast(polarity contrastPolarity, percent percent) AdaptorColor
		Analogous() analogous
		AsColorfulColor() colorful.Color
		AsGoColor() color.Color
		AsHexColor() string
		AsHexPair() FgBgTuple
		AsImageColor() imgcolor.Color
		AsPrettyJSONColor() [2]string
		AsRGBColor() color.RGBColor
		AsTermenvColor() termenv.Color
		Blend(other AdaptorColor, blendPercent percent) AdaptorColor
		Complementary() AdaptorColor
		Contrast() AdaptorColor
		Darker(percent percent) AdaptorColor
		HueOffset(degrees º) AdaptorColor
		Lighter(percent percent) AdaptorColor
		Monochromatic() shades
		Similar(count int) []AdaptorColor
		SplitComplementary() splitComplementary
		Triadic() triadic
	}

	rgbAdaptorColor [3]uint8

	// FgBgTuple ...
	FgBgTuple struct {
		Fg string `yaml:"fg"`
		Bg string `yaml:"bg"`
	}
)

// String ...
func (rgb rgbAdaptorColor) String() string {
	return rgb.AsColorfulColor().Hex()
}

// GoString ...
func (rgb rgbAdaptorColor) GoString() string {
	if names, _ := palette.AllPalettes().Name(rgb.AsColorfulColor()); len(names) > 0 {
		return names[0].Name
	}
	return rgb.String()
}

// Reverse ...
func (s FgBgTuple) Reverse() FgBgTuple {
	return FgBgTuple{
		Fg: s.Bg,
		Bg: s.Fg,
	}
}

// Format ...
func (s FgBgTuple) Format() func(value string) string {
	return func(value string) string {
		return cfmt.Sprintf("{{%s}}::"+s.Fg+"|bg"+s.Bg, value)
	}
}

func (p percent) asFloat64() float64 {
	const oneHundredPercent = 100.0
	return float64(p) / oneHundredPercent
}

func (d º) asInt() int {
	const circleDegrees = 360
	i := d % circleDegrees
	if i < 0 {
		i += circleDegrees
	}
	return int(i)
}
