package theme

import (
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/tty"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

// MessageFg ...
func (theme Theme) MessageFg() string {
	cf, _ := colorful.MakeColor(gamut.Tints(gamut.Hex(theme.BaseColor), 64)[60])
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
	return tty.Fainter(theme.LevelColor(lvl), 0.60)
}
