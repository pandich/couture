package theme

import (
	"couture/internal/pkg/source"
	"crypto/sha256"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/styles"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"github.com/muesli/gamut/palette"
	"sort"
	"strings"
)

// TODO change themes to be predefined colors
//		have the definitions use color names from gamut.Palette?

const (
	// White ...
	White = "#ffffff"
	// Black ....
	Black      = "#000000"
	purpleRain = "#ae99bf"
	merlot     = "#a01010"
	ocean      = "#5198eb"
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
const (
	// BlackAndWhite ...
	BlackAndWhite = "none"
	// Prince ...
	Prince = "prince"
	// Brougham ...
	Brougham = "brougham"
	// Ocean ...
	Ocean = "ocean"
)

// Registry ...
var Registry = map[string]Theme{
	BlackAndWhite: newTheme(White, palette.Wikipedia, styles.BlackWhite),
	Prince:        newTheme(purpleRain, palette.Crayola, styles.Monokai),
	Brougham:      newTheme(merlot, palette.RAL, styles.Autumn),
	Ocean:         newTheme(ocean, palette.CSS, styles.Tango),
}

func newTheme(baseColor string, sourceColors gamut.Palette, jsonColorScheme *chroma.Style) Theme {
	return Theme{
		BaseColor:      baseColor,
		sourceColors:   sourceColors,
		JSONColorTheme: jsonColorScheme,
	}
}

// Names ...
func Names() []string {
	var themeNames []string
	for name := range Registry {
		themeNames = append(themeNames, name)
	}
	return themeNames
}

// Theme ...
type Theme struct {
	BaseColor      string
	sourceColors   gamut.Palette
	JSONColorTheme *chroma.Style
}

// SourceColor ,,,
func (theme Theme) SourceColor(consistentColors bool, src source.Source) string {
	colors := theme.sourceColors.Colors()
	if consistentColors {
		sort.Slice(colors, func(i, j int) bool { return strings.Compare(colors[i].Name, colors[j].Name) <= 0 })
	}
	index := sourceHash(src, len(colors))
	cf, _ := colorful.MakeColor(colors[index].Color)
	return cf.Hex()
}

func sourceHash(src source.Source, max int) int {
	var sum int64
	for _, v := range sha256.Sum256([]byte(src.URL().String())) {
		sum += int64(v)
	}
	return int(sum % int64(max))
}
