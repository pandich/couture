package theme

import (
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/source"
	"math/rand"
)

type (
	// columnStyle ...
	columnStyle struct {
		Fg string
		Bg string
	}

	// Theme ...
	Theme struct {
		Legend           columnStyle
		Source           []columnStyle
		Timestamp        columnStyle
		Application      columnStyle
		Thread           columnStyle
		Class            columnStyle
		MethodDelimiter  columnStyle
		Method           columnStyle
		LineDelimiter    columnStyle
		Line             columnStyle
		Level            map[level.Level]columnStyle
		Message          map[level.Level]columnStyle
		MessageHighlight string
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
func (theme Theme) HighlightFg() string {
	return theme.MessageHighlight
}

// StackTraceFg ...
func (theme Theme) StackTraceFg() string {
	return fgHex(theme.Message[level.Error])
}

// ClassFg ...
func (theme Theme) ClassFg() string {
	return fgHex(theme.Class)
}

// MethodDelimiterFg ...
func (theme Theme) MethodDelimiterFg() string {
	return fgHex(theme.MethodDelimiter)
}

// MethodFg ...
func (theme Theme) MethodFg() string {
	return fgHex(theme.Method)
}

// LineNumberDelimiterFg ...
func (theme Theme) LineNumberDelimiterFg() string {
	return fgHex(theme.LineDelimiter)
}

// LineNumberFg ...
func (theme Theme) LineNumberFg() string {
	return fgHex(theme.Line)
}

// ThreadFg ...
func (theme Theme) ThreadFg() string {
	return fgHex(theme.Thread)
}

// CallerBg ...
func (theme Theme) CallerBg() string {
	return bgHex(theme.Class)
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
