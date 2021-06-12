package theme

// TODO every thing must have a light and dark mode

import (
	"embed"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/pandich/couture/model/level"
	"github.com/pandich/couture/source"
	"gopkg.in/yaml.v3"
	"math/rand"
	"path"
)

// Prince ...
const Prince = "prince"

//go:embed themes/*.yaml
var fs embed.FS

// Registry is the registry of theme names to their structs.
var Registry = map[string]Theme{
	Prince: mustLoad(Prince),
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

func load(name string) (*Theme, error) {
	b, err := fs.ReadFile(path.Join("themes", name+".yaml"))
	if err != nil {
		return nil, err
	}

	var theme Theme
	err = yaml.Unmarshal(b, &theme)
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

func mustLoad(name string) Theme {
	theme, err := load(name)
	if err != nil {
		panic(err)
	}
	return *theme
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
