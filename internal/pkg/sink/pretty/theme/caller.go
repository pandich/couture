package theme

import (
	"couture/internal/pkg/tty"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

func (theme Theme) callerColors() (string, string, string) {
	var aContrast = 0.4
	var bContrast = 0.8
	var cContrast = 0.4
	if tty.IsDarkMode() {
		aContrast = 0.6
		bContrast = 0.0
		cContrast = 0.0
	}
	col := gamut.Hex(theme.BaseColor)
	q := gamut.Analogous(col)
	a, _ := colorful.MakeColor(gamut.Darker(col, aContrast))
	b, _ := colorful.MakeColor(gamut.Darker(q[0], bContrast))
	c, _ := colorful.MakeColor(gamut.Darker(q[1], cContrast))
	return a.Hex(), b.Hex(), c.Hex()
}

// ClassFg ...
func (theme Theme) ClassFg() string {
	_, v, _ := theme.callerColors()
	return v
}

// MethodDelimiterFg ...
func (theme Theme) MethodDelimiterFg() string {
	const contrast = 0.25
	return tty.Darker(theme.MethodFg(), contrast)
}

// MethodFg ...
func (theme Theme) MethodFg() string {
	v, _, _ := theme.callerColors()
	return v
}

// LineNumberDelimiterFg ...
func (theme Theme) LineNumberDelimiterFg() string {
	const contrast = 0.25
	if tty.IsDarkMode() {
		return tty.Darker(theme.LineNumberFg(), contrast)
	}
	return tty.Lighter(theme.LineNumberFg(), contrast)
}

// LineNumberFg ...
func (theme Theme) LineNumberFg() string {
	_, _, v := theme.callerColors()
	return v
}

// ThreadFg ...
func (theme Theme) ThreadFg() string {
	if tty.IsDarkMode() {
		return tty.SimilarBg(tty.Darker(theme.BaseColor, 0.5))
	}
	return tty.SimilarBg(tty.Lighter(theme.BaseColor, 0.5))
}

// CallerBg ...
func (theme Theme) CallerBg() string {
	if tty.IsDarkMode() {
		return "#202020"
	}
	return "#f0f0f0"
}
