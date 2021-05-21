package pretty

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"github.com/muesli/termenv"
)

func caller(center string) (string, string, string, string) {
	col := gamut.Hex(center)
	q := gamut.Analogous(col)
	a, _ := colorful.MakeColor(col)
	b, _ := colorful.MakeColor(q[0])
	c, _ := colorful.MakeColor(q[1])
	d, _ := colorful.MakeColor(gamut.Darker(col, 0.5))
	return a.Hex(), b.Hex(), c.Hex(), d.Hex()
}

func contrast(hex string) string {
	cf, _ := colorful.MakeColor(gamut.Contrast(gamut.Hex(hex)))
	return cf.Hex()
}

func lighter(hex string, percent float64) string {
	cf, _ := colorful.MakeColor(gamut.Lighter(gamut.Hex(hex), percent))
	return cf.Hex()
}

func darker(hex string, percent float64) string {
	cf, _ := colorful.MakeColor(gamut.Darker(gamut.Hex(hex), percent))
	return cf.Hex()
}
func fainter(hex string, percent float64) string {
	const count = 1000
	i := int(count * percent)
	bg := termenv.ConvertToRGB(termenv.BackgroundColor())
	fainter := gamut.Blends(bg, gamut.Hex(hex), count)[count-i]
	col, _ := colorful.MakeColor(fainter)
	return col.Hex()
}

func isDark(hex string) bool {
	const gray = 0.5
	col, _ := colorful.MakeColor(gamut.Hex(hex))
	_, _, l := col.Hcl()
	return l < gray
}
