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

type (
	columnStyle struct {
		Fg string `yaml:"fg"`
		Bg string `yaml:"bg"`
	}

	// Theme ...
	Theme struct {
		Legend          columnStyle                 `yaml:"legend"`
		Source          []columnStyle               `yaml:"source"`
		Timestamp       columnStyle                 `yaml:"timestamp"`
		Application     columnStyle                 `yaml:"application"`
		Context         columnStyle                 `yaml:"context"`
		Entity          columnStyle                 `yaml:"entity"`
		ActionDelimiter columnStyle                 `yaml:"action_delimiter"`
		Action          columnStyle                 `yaml:"action"`
		LineDelimiter   columnStyle                 `yaml:"line_delimiter"`
		Line            columnStyle                 `yaml:"line"`
		Level           map[level.Level]columnStyle `yaml:"level"`
		Message         map[level.Level]columnStyle `yaml:"message"`
	}
)

// ApplicationFg ...
func (theme Theme) ApplicationFg() string {
	return fgHex(theme.Application)
}

// ApplicationBg ...
func (theme Theme) ApplicationBg() string {
	return bgHex(theme.Application)
}

// TimestampFg ...
func (theme Theme) TimestampFg() string {
	return fgHex(theme.Timestamp)
}

// TimestampBg ...
func (theme Theme) TimestampBg() string {
	return bgHex(theme.Application)
}

// LevelColorFg ...
func (theme Theme) LevelColorFg(lvl level.Level) string {
	return theme.Level[lvl].Fg
}

// LevelColorBg ...
func (theme Theme) LevelColorBg(lvl level.Level) string {
	return theme.Level[lvl].Bg
}

// MessageFg ...
func (theme Theme) MessageFg() string {
	return fgHex(theme.Message[level.Info])
}

// MessageBg ...
func (theme Theme) MessageBg(lvl level.Level) string {
	return bgHex(theme.Message[lvl])
}

// HighlightFg ...
func (theme Theme) HighlightFg(lvl level.Level) string {
	return theme.MessageBg(lvl)
}

// HighlightBg ...
func (theme Theme) HighlightBg() string {
	return theme.MessageFg()
}

// StackTraceFg ...
func (theme Theme) StackTraceFg() string {
	return fgHex(theme.Message[level.Error])
}

// EntityFg ...
func (theme Theme) EntityFg() string {
	return fgHex(theme.Entity)
}

// ActionDelimiterFg ...
func (theme Theme) ActionDelimiterFg() string {
	return fgHex(theme.ActionDelimiter)
}

// ActionFg ...
func (theme Theme) ActionFg() string {
	return fgHex(theme.Action)
}

// LineNumberDelimiterFg ...
func (theme Theme) LineNumberDelimiterFg() string {
	return fgHex(theme.LineDelimiter)
}

// LineNumberFg ...
func (theme Theme) LineNumberFg() string {
	return fgHex(theme.Line)
}

// ContextFg ...
func (theme Theme) ContextFg() string {
	return fgHex(theme.Context)
}

// CallerBg ...
func (theme Theme) CallerBg() string {
	return bgHex(theme.Entity)
}

// SourceColor returns a color for a source. When consistentColors is true, sources will get the same
// color across invocations of the application. Otherwise, the color selection randomized for each run.
func (theme Theme) SourceColor(consistentColors bool, src source.Source) (string, string) {
	//nolint:gosec
	var index = rand.Intn(len(theme.Source))
	if consistentColors {
		index = src.URL().Hash() % len(theme.Source)
	}
	c := theme.Source[index]
	return c.Fg, c.Bg
}

func fgHex(style columnStyle) string {
	return style.Fg
}

func bgHex(style columnStyle) string {
	return style.Bg
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
