package theme

import (
	"couture/internal/pkg/couture"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/source"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"math/rand"
)

// Prince ...
const Prince = "prince"

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

// Style ...
type (
	// Style ...
	Style struct {
		Fg string `yaml:"fg"`
		Bg string `yaml:"bg"`
	}

	// Theme ...
	Theme struct {
		Legend          Style                 `yaml:"legend"`
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

// SourceColor returns a color for a source. When consistentColors is true, sources will get the same
// color across invocations of the application. Otherwise, the color selection randomized for each run.
func (theme Theme) SourceColor(consistentColors bool, src source.Source) Style {
	//nolint:gosec
	var index = rand.Intn(len(theme.Source))
	if consistentColors {
		index = src.URL().Hash() % len(theme.Source)
	}
	return theme.Source[index]
}

func load(name string) (*Theme, error) {
	f, err := couture.Open("/themes/" + name + ".yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
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
