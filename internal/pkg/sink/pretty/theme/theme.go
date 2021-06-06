package theme

import (
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/source"
	"math/rand"
)

var (
	// Default ...
	Default = prince

	// Registry is the registry of theme names to their structs.
	Registry = map[string]Theme{}

	// Names ...
	Names []string
)

type (
	// columnStyle ...
	columnStyle struct {
		Fg string
		Bg string
	}

	// Theme ...
	Theme struct {
		Legend          columnStyle
		Source          []columnStyle
		Timestamp       columnStyle
		Application     columnStyle
		Context         columnStyle
		Entity          columnStyle
		ActionDelimiter columnStyle
		Action          columnStyle
		LineDelimiter   columnStyle
		Line            columnStyle
		Level           map[level.Level]columnStyle
		Message         map[level.Level]columnStyle
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

func register(name string, theme Theme) {
	Names = append(Names, name)
	Registry[name] = theme
}

func style(fg string, bg string) columnStyle {
	return columnStyle{
		Fg: fg,
		Bg: bg,
	}
}
