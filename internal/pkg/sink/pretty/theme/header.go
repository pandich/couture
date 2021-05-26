package theme

import (
	"couture/internal/pkg/tty"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

// ApplicationFg ...
func (theme Theme) ApplicationFg() string {
	cf, _ := colorful.MakeColor(gamut.Quadratic(gamut.Hex(theme.BaseColor))[2])
	h, _, _ := cf.Hcl()
	applicationColor := colorful.Hcl(h, 0.9, 0.8).Hex()
	return tty.SimilarBg(applicationColor)
}

// ApplicationBg ...
func (theme Theme) ApplicationBg() string {
	return tty.Contrast(theme.ApplicationFg())
}

// TimestampFg ...
func (theme Theme) TimestampFg() string {
	const degrees60 = 60 / 360.0
	var yellow = colorful.Hcl(degrees60, 1, 1)
	cf, _ := colorful.MakeColor(gamut.Tints(gamut.Complementary(gamut.Hex(theme.BaseColor)), 3)[0])
	timestampColor := gamut.Blends(yellow, cf, 16)[3]
	timestampCf, _ := colorful.MakeColor(timestampColor)
	return tty.SimilarBg(timestampCf.Hex())
}

// TimestampBg ...
func (theme Theme) TimestampBg() string {
	return tty.Contrast(theme.TimestampFg())
}
