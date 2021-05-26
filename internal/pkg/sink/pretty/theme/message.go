package theme

import (
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/tty"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

// MessageFg ...
func (theme Theme) MessageFg() string {
	const shadeCount = 64
	var index = 3
	if tty.IsDarkMode() {
		index = shadeCount - index
	}
	var cf, _ = colorful.MakeColor(gamut.Tints(gamut.Hex(theme.BaseColor), shadeCount)[index])
	if !tty.IsDarkMode() {
		cf, _ = colorful.MakeColor(gamut.Darker(cf, 0.2))
	}
	return cf.Hex()
}

// MessageBg ...
func (theme Theme) MessageBg(lvl level.Level) string {
	return tty.Fainter(theme.LevelColor(lvl), 0.90)
}

// StackTraceFg ...
func (theme Theme) StackTraceFg() string {
	return tty.SimilarBg(theme.levelTint(theme.errorLevelColor()))
}

// HighlightBg ...
func (theme Theme) HighlightBg(lvl level.Level) string {
	if tty.IsDarkMode() {
		return tty.Fainter(theme.LevelColor(lvl), 0.60)
	}
	return tty.Lighter(theme.LevelColor(lvl), 0.10)
}
