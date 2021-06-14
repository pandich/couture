package color

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"github.com/muesli/gamut/palette"
	"github.com/muesli/termenv"
	imgcolor "image/color"
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
var (
	// White ...
	White = ByHex("#ffffff")
	// Black ...
	Black = ByHex("#000000")
	// Red ...
	Red = ByHex("#FF0000")
	// Green ...
	Green = ByHex("#00FF00")
	// Blue ...
	Blue = ByHex("#0000FF")
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type (
	// Distance ...
	Distance struct{ R, G, B float64 }

	shades             [256]AdaptorColor
	splitComplementary [2]AdaptorColor
	analogous          [2]AdaptorColor
	triadic            [2]AdaptorColor
	percent            uint8
	º                  int

	// AdaptorColor ...
	AdaptorColor interface {
		fmt.Stringer
		fmt.GoStringer
		AdjustConstrast(polarity contrastPolarity, percent percent) AdaptorColor
		Analogous() analogous
		AsColorfulColor() colorful.Color
		AsGamutColor() gamut.Color
		AsGooKitColor() color.Color
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
		DistanceRgb(other AdaptorColor) float64
		DistancesRgb() Distance
		Similar(count int) adaptorPalette
		SplitComplementary() splitComplementary
		Triadic() triadic
		IsPastel() bool
		IsHappy() bool
		IsCool() bool
		IsWarm() bool
		PleasingPalette(colorCount uint) adaptorPalette
	}
	rgbAdaptorColor [3]uint8
	// FgBgTuple ...
	FgBgTuple struct {
		Fg string `yaml:"fg"`
		Bg string `yaml:"bg"`
	}
)

// DistanceRgb ...
func (rgb rgbAdaptorColor) DistanceRgb(other AdaptorColor) float64 {
	return rgb.AsColorfulColor().DistanceRgb(other.AsColorfulColor())
}

// DistancesRgb ...
func (rgb rgbAdaptorColor) DistancesRgb() Distance {
	return Distance{
		R: rgb.DistanceRgb(Red),
		G: rgb.DistanceRgb(Green),
		B: rgb.DistanceRgb(Blue),
	}
}

// IsCool ,,,
func (rgb rgbAdaptorColor) IsCool() bool {
	return gamut.Cool(rgb.AsColorfulColor())
}

// IsWarm ,,,
func (rgb rgbAdaptorColor) IsWarm() bool {
	return gamut.Warm(rgb.AsColorfulColor())
}

// IsHappy ...
func (rgb rgbAdaptorColor) IsHappy() bool {
	return gamut.HappyGenerator{}.Valid(rgb.AsColorfulColor())
}

// IsPastel ...
func (rgb rgbAdaptorColor) IsPastel() bool {
	return gamut.PastelGenerator{}.Valid(rgb.AsColorfulColor())
}

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

func (d Distance) min() float64 {
	var m float64
	for f, v := range []float64{d.R, d.G, d.B} {
		if f == 0 || v < m {
			m = v
		}
	}
	return m
}

func (d Distance) closestToRed() bool { return d.R == d.min() }

func (d Distance) closestToGreen() bool { return d.G == d.min() }

func (d Distance) closestToBlue() bool { return d.B == d.min() }
