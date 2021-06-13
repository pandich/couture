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
	"regexp"
)

// White ...
var White = ByHex("#ffffff")

// Black ...
var Black = ByHex("#000000")

var specialNames = map[string]string{
	"prince":    "Logan",
	"halloween": "Burnt Orange",
}

// ByHex ...
func ByHex(hex string) AdaptorColor {
	values := color.Hex(hex).Values()
	r, g, b := uint8(values[0]), uint8(values[1]), uint8(values[2])
	return rgbColor{r, g, b}
}

// ByName ...
func ByName(name string) (AdaptorColor, bool) {
	hexPattern := regexp.MustCompile("^#?[0-9A-Fa-f]{6}$")
	if hexPattern.MatchString(name) {
		if name[0] != '#' {
			name = "#" + name
		}
		return ByHex(name), true
	}
	if s, ok := specialNames[name]; ok {
		name = s
	}
	if c, ok := palette.AllPalettes().Color(name); ok {
		cf, _ := colorful.MakeColor(c)
		return ByHex(cf.Hex()), true
	}
	return nil, false
}

// MustByName ...
func MustByName(name string) AdaptorColor {
	if c, ok := ByName(name); ok {
		return c
	}
	panic(name)
}

// ByImageColor ...
func ByImageColor(imgColor imgcolor.Color) AdaptorColor {
	cf, _ := colorful.MakeColor(imgColor)
	return ByHex(cf.Hex())
}

type shades [256]AdaptorColor
type splitComplementary [2]AdaptorColor
type analogous [2]AdaptorColor
type triadic [2]AdaptorColor

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
	AdjustConstrast(polarity contrastPolarity, amount float64) AdaptorColor
	Blend(other AdaptorColor, blendPercent int) AdaptorColor
	Contrast() AdaptorColor
	HueOffset(degrees int) AdaptorColor
	SplitComplementary() splitComplementary
	Complementary() AdaptorColor
	Analogous() analogous
	Monochromatic() shades
	Triadic() triadic
	Lighter(percent float64) AdaptorColor
	Darker(percent float64) AdaptorColor
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

// HexPair ...
type HexPair struct {
	Fg string `yaml:"fg"`
	Bg string `yaml:"bg"`
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
