package theme

import (
	"couture/internal/pkg/tty"
	"github.com/muesli/gamut"
)

const (
	// White ...
	White      = "#ffffff"
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
	BlackAndWhite: {BaseColor: White, SourceColors: tty.NewColorCycle(gamut.PastelGenerator{})},
	Prince:        {BaseColor: purpleRain, SourceColors: tty.NewColorCycle(gamut.PastelGenerator{})},
	Brougham:      {BaseColor: merlot, SourceColors: tty.NewColorCycle(gamut.WarmGenerator{})},
	Ocean:         {BaseColor: ocean, SourceColors: tty.NewColorCycle(gamut.HappyGenerator{})},
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
	BaseColor    string
	SourceColors chan string
}
