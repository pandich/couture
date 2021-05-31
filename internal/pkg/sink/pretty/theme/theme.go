package theme

import (
	"couture/internal/pkg/source"
	"crypto/sha256"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"github.com/muesli/gamut/palette"
	"github.com/muesli/termenv"
	"image/color"
	"math/rand"
	"sort"
	"strings"
)

// Default is the default theme.
const Default = prince

const (
	prince        = "prince"
	blackAndWhite = "none"
	brougham      = "brougham"
	ocean         = "ocean"
)

// Names returns all theme names.
func Names() []string {
	var themeNames []string
	for name := range Registry {
		themeNames = append(themeNames, name)
	}
	return themeNames
}

// Registry is the registry of theme names to their structs.
var Registry = map[string]Theme{
	blackAndWhite: newTheme("White", palette.Crayola),
	prince:        newTheme("#ae99bf", palette.Crayola),
	brougham:      newTheme("Red", palette.Crayola),
	ocean:         newTheme("Navy Blue", palette.Crayola),
}

func newTheme(baseColor string, sourceColorGamut gamut.Palette) Theme {
	baseColor = determineBaseColor(baseColor)
	return Theme{
		BaseColor:    baseColor,
		sourceColors: buildSourceColors(baseColor, sourceColorGamut),
	}
}

// Theme contains the minimal information to programatically generate the output palette.
type Theme struct {
	// BaseColor drives most color generation options.
	BaseColor string
	// sourceColors is a list of available colors for displaying the source column.
	sourceColors []color.Color
}

// SourceColor returns a color for a source. When consistentColors is true, sources will get the same
// color across invocations of the application. Otherwise, the color selection randomized for each run.
func (theme Theme) SourceColor(consistentColors bool, src source.Source) string {
	//nolint:gosec
	var index = rand.Intn(len(theme.sourceColors))
	if consistentColors {
		index = sourceHash(src, len(theme.sourceColors))
	}
	cf, _ := colorful.MakeColor(theme.sourceColors[index])
	return cf.Hex()
}

func sourceHash(src source.Source, max int) int {
	var sum int64
	for _, v := range sha256.Sum256([]byte(src.URL().String())) {
		sum += int64(v)
	}
	return int(sum % int64(max))
}

//
// Helpers
//

func buildSourceColors(baseColor string, sourceColorGamut gamut.Palette) []color.Color {
	sourceColorGamutColors := sourceColorGamut.Colors()
	sort.Slice(sourceColorGamutColors, func(i, j int) bool {
		return strings.Compare(sourceColorGamutColors[i].Name, sourceColorGamutColors[j].Name) <= 0
	})
	baseColorHex := determineBaseColor(baseColor)
	var sourceColors []color.Color
	for i := range sourceColorGamutColors {
		c1 := sourceColorGamutColors[i].Color
		c2 := gamut.Hex(baseColorHex)
		c := gamut.Blends(c1, c2, 4)[1]
		sourceColors = append(sourceColors, c)
	}
	return sourceColors
}

func determineBaseColor(baseColor string) string {
	if baseColor[0] == '#' {
		return baseColor
	}
	if c, ok := palette.AllPalettes().Color(baseColor); ok {
		cf, _ := colorful.MakeColor(c)
		baseColor = cf.Hex()
	}
	if baseColor[0] != '#' {
		panic(baseColor)
	}
	return baseColor
}

func similarBg(hex string) string {
	return hex + "|bg" + fainter(hex, 0.96)
}

func fainter(hex string, percent float64) string {
	const count = 1000
	i := int(count * percent)
	bg := termenv.ConvertToRGB(termenv.BackgroundColor())
	fainter := gamut.Blends(bg, gamut.Hex(hex), count)[count-i]
	col, _ := colorful.MakeColor(fainter)
	return col.Hex()
}

func lighter(hex string, percent float64) string {
	cf, _ := colorful.MakeColor(gamut.Lighter(gamut.Hex(hex), percent))
	return cf.Hex()
}

func darker(hex string, percent float64) string {
	cf, _ := colorful.MakeColor(gamut.Darker(gamut.Hex(hex), percent))
	return cf.Hex()
}

func contrast(hex string) string {
	cf, _ := colorful.MakeColor(gamut.Contrast(gamut.Hex(hex)))
	return cf.Hex()
}
