package theme

import (
	"fmt"
	"github.com/muesli/gamut"
	"github.com/muesli/gamut/palette"
	"github.com/pandich/couture/event/level"
	"github.com/pandich/couture/sink/color"
	"github.com/pandich/couture/source"
	"github.com/tidwall/pretty"
	"math/rand"
	"slices"
	"strings"
	"time"
)

var (
	random              = rand.New(rand.NewSource(time.Now().UnixNano()))
	consistentColorPool gamut.Colors
)

func init() {

	consistent := palette.Crayola.Colors()
	cl := make(gamut.Colors, len(consistent))
	copy(cl, consistent)
	slices.SortFunc(
		cl, func(a, b gamut.Color) int {
			return strings.Compare(a.Name, b.Name)
		},
	)
	consistentColorPool = cl
}

// Theme ...
type Theme struct {
	Source          []color.FgBgTuple               `yaml:"source"`
	Timestamp       color.FgBgTuple                 `yaml:"timestamp"`
	Application     color.FgBgTuple                 `yaml:"application"`
	Context         color.FgBgTuple                 `yaml:"context"`
	Entity          color.FgBgTuple                 `yaml:"entity"`
	ActionDelimiter color.FgBgTuple                 `yaml:"action_delimiter"`
	Action          color.FgBgTuple                 `yaml:"action"`
	LineDelimiter   color.FgBgTuple                 `yaml:"line_delimiter"`
	Line            color.FgBgTuple                 `yaml:"line"`
	Level           map[level.Level]color.FgBgTuple `yaml:"level"`
	Message         map[level.Level]color.FgBgTuple `yaml:"message"`
}

// AsHexPair returns a color for a source. When consistentColors is true, sources will get the same
// color across invocations of the application. Otherwise, the color selection randomized for each run.
func (theme *Theme) AsHexPair(consistentColors bool, src source.Source) color.FgBgTuple {
	if consistentColors {
		url := src.URL()
		bg := consistentColorPool[url.HashInt()%len(consistentColorPool)]
		rB, gB, bB, _ := bg.Color.RGBA()
		rF, gF, bF, _ := gamut.Contrast(bg.Color).RGBA()
		sc := func(n uint32) int { return int(float64(n) / float64(1<<16) * 255) }
		fgBg := color.FgBgTuple{
			Fg: fmt.Sprintf("#%02x%02x%02x", sc(rF), sc(gF), sc(bF)),
			Bg: fmt.Sprintf("#%02x%02x%02x", sc(rB), sc(gB), sc(bB)),
		}
		return fgBg
	}

	var index = random.Intn(len(theme.Source))
	return theme.Source[index]
}

// AsPrettyJSONStyle ...
func (theme *Theme) AsPrettyJSONStyle() *pretty.Style {
	valueColor := color.ByHex(theme.Action.Fg).AsPrettyJSONColor()
	keyColor := color.ByHex(theme.Timestamp.Fg).AsPrettyJSONColor()
	dimValueColor := color.ByHex(theme.Level[level.Trace].Bg).AsPrettyJSONColor()
	return &pretty.Style{
		Key:    keyColor,
		String: valueColor,
		Number: valueColor,
		True:   valueColor,
		False:  valueColor,
		Null:   dimValueColor,
		Escape: dimValueColor,
	}
}

func (theme *Theme) base() color.AdaptorColor {
	return color.ByHex(theme.Entity.Fg)
}
