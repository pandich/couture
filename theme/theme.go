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
	Source          []color.HexPair               `yaml:"source"`
	Timestamp       color.HexPair                 `yaml:"timestamp"`
	Application     color.HexPair                 `yaml:"application"`
	Context         color.HexPair                 `yaml:"context"`
	Entity          color.HexPair                 `yaml:"entity"`
	ActionDelimiter color.HexPair                 `yaml:"action_delimiter"`
	Action          color.HexPair                 `yaml:"action"`
	LineDelimiter   color.HexPair                 `yaml:"line_delimiter"`
	Line            color.HexPair                 `yaml:"line"`
	Level           map[level.Level]color.HexPair `yaml:"level"`
	Message         map[level.Level]color.HexPair `yaml:"message"`
}

// SourceStyle returns a color for a source. When consistentColors is true, sources will get the same
// color across invocations of the application. Otherwise, the color selection randomized for each run.
func (theme Theme) SourceStyle(consistentColors bool, src source.Source) color.HexPair {
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
