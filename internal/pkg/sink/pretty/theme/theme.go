package theme

import (
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/tty"
	"github.com/lucasb-eyer/go-colorful"
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
	BlackAndWhite: {BaseColor: White, SourceColors: gamut.PastelGenerator{}},
	Prince:        {BaseColor: purpleRain, SourceColors: gamut.PastelGenerator{}},
	Brougham:      {BaseColor: merlot, SourceColors: gamut.WarmGenerator{}},
	Ocean:         {BaseColor: ocean, SourceColors: gamut.HappyGenerator{}},
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
	SourceColors gamut.ColorGenerator
}

func (theme Theme) tinted(hex string) string {
	return tty.Tinted(theme.BaseColor, hex)
}

func (theme Theme) callerColors() (string, string, string) {
	const darkness = 0.4
	col := gamut.Hex(theme.BaseColor)
	q := gamut.Analogous(col)
	a, _ := colorful.MakeColor(gamut.Darker(col, darkness))
	b, _ := colorful.MakeColor(q[0])
	c, _ := colorful.MakeColor(q[1])
	return a.Hex(), b.Hex(), c.Hex()
}

// ClassColor ...
func (theme Theme) ClassColor() string {
	_, v, _ := theme.callerColors()
	return v
}

// MethodDelimiterColor ...
func (theme Theme) MethodDelimiterColor() string {
	const contrast = 0.25
	return tty.Darker(theme.MethodColor(), contrast)
}

// MethodColor ...
func (theme Theme) MethodColor() string {
	v, _, _ := theme.callerColors()
	return v
}

// LineNumberDelimiterColor ...
func (theme Theme) LineNumberDelimiterColor() string {
	const contrast = 0.25
	return tty.Darker(theme.LineNumberColor(), contrast)
}

// LineNumberColor ...
func (theme Theme) LineNumberColor() string {
	_, _, v := theme.callerColors()
	return v
}

//
// Level
//

func (theme Theme) traceLevelColor() string {
	const traceColor = "#868686"
	return theme.tinted(traceColor)
}
func (theme Theme) debugLevelColor() string {
	const debugColor = "#f6f6f6"
	return theme.tinted(debugColor)
}
func (theme Theme) infoLevelColor() string {
	const infoColor = "#00ff00"
	return theme.tinted(infoColor)
}
func (theme Theme) warnLevelColor() string {
	const warnColor = "#ffff00"
	return theme.tinted(warnColor)
}
func (theme Theme) errorLevelColor() string {
	const errorColor = "#ff0000"
	return theme.tinted(errorColor)
}

// LevelColor ...
func (theme Theme) LevelColor(lvl level.Level) string {
	switch lvl {
	case level.Debug:
		return theme.debugLevelColor()
	case level.Info:
		return theme.infoLevelColor()
	case level.Warn:
		return theme.warnLevelColor()
	case level.Error:
		return theme.errorLevelColor()
	case level.Trace:
		fallthrough
	default:
		return theme.traceLevelColor()
	}
}

// MessageColor ...
func (theme Theme) MessageColor() string {
	cf, _ := colorful.MakeColor(gamut.Tints(gamut.Hex(theme.BaseColor), 64)[60])
	return cf.Hex()
}

// MessageBackgroundColor ...
func (theme Theme) MessageBackgroundColor(lvl level.Level) string {
	return tty.Fainter(theme.LevelColor(lvl), 0.90)
}

// StackTraceColor ...
func (theme Theme) StackTraceColor() string {
	return tty.SimilarBg(theme.tinted(theme.errorLevelColor()))
}

//
// Other
//

// TimestampColor ...
func (theme Theme) TimestampColor() string {
	const degrees60 = 60 / 360.0
	var yellow = colorful.Hcl(degrees60, 1, 1)
	cf, _ := colorful.MakeColor(gamut.Tints(gamut.Complementary(gamut.Hex(theme.BaseColor)), 3)[1])
	timestampColor := gamut.Blends(yellow, cf, 16)[3]
	timestampCf, _ := colorful.MakeColor(timestampColor)
	return tty.SimilarBg(timestampCf.Hex())
}

// ApplicationColor ...
func (theme Theme) ApplicationColor() string {
	cf, _ := colorful.MakeColor(gamut.Quadratic(gamut.Hex(theme.BaseColor))[2])
	h, _, _ := cf.Hcl()
	applicationColor := colorful.Hcl(h, 0.9, 0.8).Hex()
	return tty.SimilarBg(applicationColor)
}

// ThreadColor ...
func (theme Theme) ThreadColor() string {
	return tty.SimilarBg(tty.Darker(theme.BaseColor, 0.5))
}

// CallerBgColor ...
func (theme Theme) CallerBgColor() string {
	return "#202020"
}
