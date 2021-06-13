package theme

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/pandich/couture/model/level"
	"github.com/pandich/couture/source"
	"github.com/pandich/couture/theme/color"
	"github.com/tidwall/pretty"
	"math/rand"
)

const (
	// Default ...
	Default = "prince"

	purpleRain baseColor = "#9b99bf"
)

// Registry is the registry of theme names to their structs.
var Registry = map[string]Theme{}

func init() {
	Registry["prince"] = splitComplementary(purpleRain)
}

// Names all available theme names.
func Names() []string {
	var names []string
	for k := range Registry {
		names = append(names, k)
	}
	return names
}

type (
	// Style ...
	Style struct {
		Fg string `yaml:"fg"`
		Bg string `yaml:"bg"`
	}

	// Theme ...
	Theme struct {
		Source          []Style               `yaml:"source"`
		Timestamp       Style                 `yaml:"timestamp"`
		Application     Style                 `yaml:"application"`
		Context         Style                 `yaml:"context"`
		Entity          Style                 `yaml:"entity"`
		ActionDelimiter Style                 `yaml:"action_delimiter"`
		Action          Style                 `yaml:"action"`
		LineDelimiter   Style                 `yaml:"line_delimiter"`
		Line            Style                 `yaml:"line"`
		Level           map[level.Level]Style `yaml:"level"`
		Message         map[level.Level]Style `yaml:"message"`
	}
)

// SourceStyle returns a color for a source. When consistentColors is true, sources will get the same
// color across invocations of the application. Otherwise, the color selection randomized for each run.
func (theme Theme) SourceStyle(consistentColors bool, src source.Source) Style {
	//nolint:gosec
	var index = rand.Intn(len(theme.Source))
	if consistentColors {
		index = src.URL().Hash() % len(theme.Source)
	}
	return theme.Source[index]
}

// AsPrettyJSONStyle ...
func (theme Theme) AsPrettyJSONStyle() *pretty.Style {
	valueColor := color.Hex(theme.Action.Fg).AsPrettyJSONColor()
	keyColor := color.Hex(theme.Timestamp.Fg).AsPrettyJSONColor()
	dimValueColor := color.Hex(theme.Level[level.Trace].Bg).AsPrettyJSONColor()
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

// Reverse ...
func (s Style) Reverse() Style {
	return Style{
		Fg: s.Bg,
		Bg: s.Fg,
	}
}

// Format ...
func (s Style) Format() func(value string) string {
	return func(value string) string {
		return cfmt.Sprintf("{{%s}}::"+s.Fg+"|bg"+s.Bg, value)
	}
}
