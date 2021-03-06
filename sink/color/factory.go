package color

import (
	"github.com/gookit/color"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"github.com/muesli/gamut/palette"
	imgcolor "image/color"
	"regexp"
	"strings"
)

var sourcePalettes = palette.AllPalettes()

// ByTuple ...

// ByHex ...
func ByHex(hex string) AdaptorColor {
	values := color.Hex(hex).Values()
	r, g, b := uint8(values[0]), uint8(values[1]), uint8(values[2])
	return rgbAdaptorColor{r, g, b}
}

// ByName ...
func ByName(name string) (AdaptorColor, bool) {
	if hex, ok := tryHex(name); ok {
		return ByHex(hex), true
	}

	if c, ok := sourcePalettes.Color(name); ok {
		return byImageColor(c), true
	}

	if c, ok := sourcePalettes.Color(normalizeColorName(name)); ok {
		return byImageColor(c), true
	}

	return nil, false
}

func normalizeColorName(name string) string {
	var wordBreaks = regexp.MustCompile(`[ \t_./\-]+`)
	words := wordBreaks.Split(name, -1)
	for i, word := range words {
		if len(word) > 1 {
			words[i] = strings.ToUpper(word[0:1]) + word[1:]
		} else {
			words[i] = strings.ToUpper(word)
		}
	}
	name = strings.Join(words, " ")
	return name
}

// MustByName ...
func MustByName(name string) AdaptorColor {
	if c, ok := ByName(name); ok {
		return c
	}
	panic(name)
}

func byGamutColor(c gamut.Color) AdaptorColor {
	return byImageColor(c.Color)
}

func byColorfulColor(colorfulColor colorful.Color) AdaptorColor {
	return ByHex(colorfulColor.Hex())
}

func byImageColor(imgColor imgcolor.Color) AdaptorColor {
	cf, _ := colorful.MakeColor(imgColor)
	return byColorfulColor(cf)
}

func byImageColors(in []imgcolor.Color) adaptorPalette {
	var out adaptorPalette
	for _, c := range in {
		out = append(out, byImageColor(c))
	}
	return out
}

func tryHex(hex string) (string, bool) {
	hexPattern := regexp.MustCompile("^#?([0-9A-Fa-f]{3}|[0-9A-Fa-f]{6})$")
	if hexPattern.MatchString(hex) {
		if hex[0] != '#' {
			hex = "#" + hex
		}
		return hex, true
	}
	return hex, false
}
