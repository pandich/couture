package theme

import (
	"couture/internal/pkg/tty"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

func (theme Theme) callerColors() (string, string, string) {
	const darkness = 0.4
	col := gamut.Hex(theme.BaseColor)
	q := gamut.Analogous(col)
	a, _ := colorful.MakeColor(gamut.Darker(col, darkness))
	b, _ := colorful.MakeColor(q[0])
	c, _ := colorful.MakeColor(q[1])
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
	return tty.Darker(theme.LineNumberFg(), contrast)
}

// LineNumberFg ...
func (theme Theme) LineNumberFg() string {
	_, _, v := theme.callerColors()
	return v
}

// ThreadFg ...
func (theme Theme) ThreadFg() string {
	return tty.SimilarBg(tty.Darker(theme.BaseColor, 0.5))
}

// CallerBg ...
func (theme Theme) CallerBg() string {
	return "#202020"
}
