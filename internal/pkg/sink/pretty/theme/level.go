package theme

import (
	"couture/internal/pkg/model/level"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

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

func (theme Theme) levelColor(hex string) string {
	complementary, _ := colorful.MakeColor(gamut.Complementary(gamut.Hex(hex)))
	c := gamut.Hex(hex)
	cf, _ := colorful.MakeColor(c)
	return complementary.BlendHcl(cf, 0.85).Hex()
}

func (theme Theme) traceLevelColor() string {
	const traceColor = "#868686"
	return theme.levelColor(traceColor)
}
func (theme Theme) debugLevelColor() string {
	const debugColor = "#f6f6f6"
	return theme.levelColor(debugColor)
}
func (theme Theme) infoLevelColor() string {
	const infoColor = "#00ff00"
	return theme.levelColor(infoColor)
}
func (theme Theme) warnLevelColor() string {
	const warnColor = "#ffff00"
	return theme.levelColor(warnColor)
}
func (theme Theme) errorLevelColor() string {
	const errorColor = "#ff0000"
	return theme.levelColor(errorColor)
}
