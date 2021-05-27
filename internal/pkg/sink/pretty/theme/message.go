package theme

import (
	"couture/internal/pkg/model/level"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"github.com/muesli/termenv"
)

// MessageFg ...
func (theme Theme) MessageFg() string {
	const shadeCount = 64
	var index = 3
	if termenv.HasDarkBackground() {
		index = shadeCount - index
	}
	var cf, _ = colorful.MakeColor(gamut.Tints(gamut.Hex(theme.BaseColor), shadeCount)[index])
	if !termenv.HasDarkBackground() {
		cf, _ = colorful.MakeColor(gamut.Darker(cf, 0.2))
	}
	return cf.Hex()
}

// MessageBg ...
func (theme Theme) MessageBg(lvl level.Level) string {
	return fainter(theme.LevelColor(lvl), 0.90)
}

// HighlightFg ...
func (theme Theme) HighlightFg() string {
	messageFg := gamut.Hex(theme.MessageFg())
	errorFg := gamut.Hex(theme.LevelColor(level.Error))
	fg := gamut.Blends(messageFg, errorFg, 64)[40]
	cf, _ := colorful.MakeColor(fg)
	return cf.Hex()
}

// StackTraceFg ...
func (theme Theme) StackTraceFg() string {
	return similarBg(theme.levelColor(theme.errorLevelColor()))
}
