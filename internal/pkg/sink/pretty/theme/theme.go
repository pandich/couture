package theme

import (
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/source"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"math/rand"
)

// Theme ...
type (
	Theme struct {
		Source          []lipgloss.Style
		Timestamp       lipgloss.Style
		Application     lipgloss.Style
		Thread          lipgloss.Style
		Class           lipgloss.Style
		MethodDelimiter lipgloss.Style
		Method          lipgloss.Style
		LineDelimiter   lipgloss.Style
		Line            lipgloss.Style
		Level           map[level.Level]lipgloss.Style
		Message         map[level.Level]lipgloss.Style
	}
)

// ApplicationFg ...
func (theme Theme) ApplicationFg() string {
	return FgHex(theme.Application)
}

// ApplicationBg ...
func (theme Theme) ApplicationBg() string {
	return BgHex(theme.Application)
}

// TimestampFg ...
func (theme Theme) TimestampFg() string {
	return FgHex(theme.Timestamp)
}

// TimestampBg ...
func (theme Theme) TimestampBg() string {
	return BgHex(theme.Application)
}

// LevelColor ...
func (theme Theme) LevelColor(lvl level.Level) lipgloss.Style {
	return theme.Level[lvl]
}

// MessageFg ...
func (theme Theme) MessageFg() string {
	return FgHex(theme.Message[level.Info])
}

// MessageBg ...
func (theme Theme) MessageBg(lvl level.Level) string {
	return BgHex(theme.Message[lvl])
}

// HighlightFg ...
func (theme Theme) HighlightFg() string {
	messageFg := gamut.Hex(theme.MessageFg())
	errorFg := FgHex(theme.LevelColor(level.Error))
	fg := gamut.Blends(messageFg, gamut.Hex(errorFg), 64)[40]
	cf, _ := colorful.MakeColor(fg)
	return cf.Hex()
}

// StackTraceFg ...
func (theme Theme) StackTraceFg() string {
	return FgHex(theme.Message[level.Error])
}

// ClassFg ...
func (theme Theme) ClassFg() string {
	return FgHex(theme.Class)
}

// MethodDelimiterFg ...
func (theme Theme) MethodDelimiterFg() string {
	return FgHex(theme.MethodDelimiter)
}

// MethodFg ...
func (theme Theme) MethodFg() string {
	return FgHex(theme.Method)
}

// LineNumberDelimiterFg ...
func (theme Theme) LineNumberDelimiterFg() string {
	return FgHex(theme.LineDelimiter)
}

// LineNumberFg ...
func (theme Theme) LineNumberFg() string {
	return FgHex(theme.Line)
}

// ThreadFg ...
func (theme Theme) ThreadFg() string {
	return FgHex(theme.Thread)
}

// CallerBg ...
func (theme Theme) CallerBg() string {
	return BgHex(theme.Class)
}

// SourceColor returns a color for a source. When consistentColors is true, sources will get the same
// color across invocations of the application. Otherwise, the color selection randomized for each run.
func (theme Theme) SourceColor(consistentColors bool, src source.Source) lipgloss.Style {
	//nolint:gosec
	var index = rand.Intn(len(theme.Source))
	if consistentColors {
		index = src.URL().Hash() % len(theme.Source)
	}
	return theme.Source[index]
}

// FgHex ...
func FgHex(style lipgloss.Style) string {
	c := style.GetForeground()
	cf, _ := colorful.MakeColor(c)
	return cf.Hex()
}

// BgHex ...
func BgHex(style lipgloss.Style) string {
	c := style.GetBackground()
	cf, _ := colorful.MakeColor(c)
	return cf.Hex()
}
