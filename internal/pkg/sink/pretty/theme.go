package pretty

import (
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/tty"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

const (
	blackAndWhite = ""
	purpleRain    = "#ae99bf"
	merlot        = "#a01010"
	ocean         = "#5198eb"
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

// ThemeByName TODO theme color tweaks
var ThemeByName = map[string]Theme{
	BlackAndWhite: {BaseColor: blackAndWhite, SourceColors: gamut.PastelGenerator{}},
	Prince:        {BaseColor: purpleRain, SourceColors: gamut.PastelGenerator{}},
	Brougham:      {BaseColor: merlot, SourceColors: gamut.WarmGenerator{}},
	Ocean:         {BaseColor: ocean, SourceColors: gamut.HappyGenerator{}},
}

const (
	errorColor = "#ff0000"
	traceColor = "#868686"
	debugColor = "#f6f6f6"
	infoColor  = "#00ff00"
	warnColor  = "#ffff00"
)

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

func (theme Theme) applicationColor() string {
	cf, _ := colorful.MakeColor(gamut.Quadratic(gamut.Hex(theme.BaseColor))[2])
	h, _, _ := cf.Hcl()
	applicationColor := colorful.Hcl(h, 0.9, 0.8).Hex()
	return tty.SimilarBg(applicationColor)
}

func (theme Theme) classColor() string {
	_, v, _ := theme.callerColors()
	return v
}

func (theme Theme) methodColor() string {
	v, _, _ := theme.callerColors()
	return v
}

func (theme Theme) lineNumberColor() string {
	_, _, v := theme.callerColors()
	return v
}

func (theme Theme) methodDelimiterColor() string {
	const contrast = 0.25
	return tty.Darker(theme.methodColor(), contrast)
}

func (theme Theme) lineNumberDelimiterColor() string {
	const contrast = 0.25
	return tty.Darker(theme.lineNumberColor(), contrast)
}

func (theme Theme) traceLevelColor() string {
	return theme.tinted(traceColor)
}
func (theme Theme) debugLevelColor() string {
	return theme.tinted(debugColor)
}
func (theme Theme) infoLevelColor() string {
	return theme.tinted(infoColor)
}
func (theme Theme) warnLevelColor() string {
	return theme.tinted(warnColor)
}
func (theme Theme) errorLevelColor() string {
	return theme.tinted(errorColor)
}

func (theme Theme) messageColor() string {
	cf, _ := colorful.MakeColor(gamut.Tints(gamut.Hex(theme.BaseColor), 64)[60])
	return cf.Hex()
}

func (theme Theme) messageBackgroundColor(lvl level.Level) string {
	return tty.Fainter(theme.levelColor(lvl), 0.90)
}

func (theme Theme) levelColor(lvl level.Level) string {
	switch lvl {
	case level.Trace:
		return theme.traceLevelColor()
	case level.Debug:
		return theme.debugLevelColor()
	case level.Info:
		return theme.infoLevelColor()
	case level.Warn:
		return theme.warnLevelColor()
	case level.Error:
		return theme.errorLevelColor()
	default:
		return "#fffff"
	}
}

func (theme Theme) stackTraceColor() string {
	return tty.SimilarBg(theme.tinted(theme.errorLevelColor()))
}
