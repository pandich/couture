package color

import (
	"github.com/gookit/color"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut/palette"
	imgcolor "image/color"
	"regexp"
)

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

func byImageColors(in []imgcolor.Color) []AdaptorColor {
	var out []AdaptorColor
	for _, c := range in {
		out = append(out, ByImageColor(c))
	}
	return out
}
