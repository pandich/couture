package theme

import (
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/tty"
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

func (theme Theme) levelTint(hex string) string {
	return tty.Tinted(theme.BaseColor, hex)
}

func (theme Theme) traceLevelColor() string {
	const traceColor = "#868686"
	return theme.levelTint(traceColor)
}
func (theme Theme) debugLevelColor() string {
	const debugColor = "#f6f6f6"
	return theme.levelTint(debugColor)
}
func (theme Theme) infoLevelColor() string {
	const infoColor = "#00ff00"
	return theme.levelTint(infoColor)
}
func (theme Theme) warnLevelColor() string {
	const warnColor = "#ffff00"
	return theme.levelTint(warnColor)
}
func (theme Theme) errorLevelColor() string {
	const errorColor = "#ff0000"
	return theme.levelTint(errorColor)
}
