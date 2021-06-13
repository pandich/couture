package theme

import (
	"github.com/pandich/couture/model/level"
	"github.com/pandich/couture/source"
	"github.com/pandich/couture/theme/color"
	"github.com/tidwall/pretty"
	"math/rand"
)

// Default ...
const Default = "prince"

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
	SourceStyle     string                          `yaml:"-"`
}

// AsHexPair returns a color for a source. When consistentColors is true, sources will get the same
// color across invocations of the application. Otherwise, the color selection randomized for each run.
func (theme Theme) AsHexPair(consistentColors bool, src source.Source) color.FgBgTuple {
	//nolint: gosec
	var index = rand.Intn(len(theme.Source))
	if consistentColors {
		index = src.URL().Hash() % len(theme.Source)
	}
	return theme.Source[index]
}

// AsPrettyJSONStyle ...
func (theme Theme) AsPrettyJSONStyle() *pretty.Style {
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
