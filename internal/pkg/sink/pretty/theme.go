package pretty

import (
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

var reg = cfmt.RegisterStyle

// Theme ...
type Theme struct {
	BaseColor        string
	ApplicationColor string
	DefaultColor     string
	TimestampColor   string
	ErrorColor       string
	TraceColor       string
	DebugColor       string
	InfoColor        string
	WarnColor        string
	MessageColor     string
	StackTraceColor  string
	SourceColors     gamut.ColorGenerator
}

func (theme Theme) init() {
	if !sink.IsTTY() || theme.BaseColor == "" {
		cfmt.DisableColors()
	}
	reg("Default", func(s string) string { return s })
	reg("Timestamp", func(s string) string { return cfmt.Sprintf("{{ %s }}::"+sink.WithFaintBg(theme.TimestampColor), s) })
	reg("Application", func(s string) string {
		return cfmt.Sprintf("{{ %-20.20s }}::"+sink.WithFaintBg(theme.ApplicationColor), s)
	})
	theme.initCaller()
	theme.initMessage()
}

func (theme Theme) initCaller() {
	const contrastPercent = 0.25
	callerPalette := func(center string) (string, string, string, string) {
		col := gamut.Hex(center)
		q := gamut.Analogous(col)
		a, _ := colorful.MakeColor(col)
		b, _ := colorful.MakeColor(q[0])
		c, _ := colorful.MakeColor(q[1])
		d, _ := colorful.MakeColor(gamut.Darker(col, 0.5))
		return a.Hex(), b.Hex(), c.Hex(), d.Hex()
	}
	methodColor, classColor, lineNumberColor, threadColor := callerPalette(theme.BaseColor)

	var methodDelimiterColor = sink.Lighter(methodColor, contrastPercent)
	var lineNumberDelimiterColor = sink.Lighter(lineNumberColor, contrastPercent)
	if sink.IsDark(methodColor) {
		methodDelimiterColor = sink.Darker(methodColor, contrastPercent)
		lineNumberDelimiterColor = sink.Darker(lineNumberColor, contrastPercent)
	}

	reg("Thread", func(s string) string { return cfmt.Sprintf("{{ %-15.15s }}::"+sink.WithFaintBg(threadColor), s) })
	reg("Class", func(s string) string { return cfmt.Sprintf("{{%.30s}}::"+classColor, s) })
	reg("MethodDelimiter", func(s string) string { return cfmt.Sprintf("{{%s}}::"+methodDelimiterColor, s) })
	reg("Method", func(s string) string { return cfmt.Sprintf("{{%.30s}}::"+methodColor, s) })
	reg("LineNumberDelimiter", func(s string) string { return cfmt.Sprintf("{{%s}}::"+lineNumberDelimiterColor, s) })
	reg("LineNumber", func(s string) string { return cfmt.Sprintf("{{%s}}::"+lineNumberColor, s) })
	reg("StackTrace", func(s string) string { return cfmt.Sprintf("{{%s}}::"+sink.WithFaintBg(theme.StackTraceColor), s) })
	reg("HighlightedStackTrace", func(s string) string { return cfmt.Sprintf("{{%s}}::bg"+sink.WithFaintBg(theme.ErrorColor), s) })
}

func (theme Theme) initMessage() {
	regLog := func(lvl level.Level, bgColor string) {
		messageBgColor := sink.Fainter(bgColor, 0.90)
		fgColor := sink.Contrast(bgColor)
		reg("Level"+string(lvl), func(s string) string { return cfmt.Sprintf("{{ %1.1s }}::bg"+bgColor+"|"+fgColor, s) })
		reg("Message"+string(lvl), func(s string) string { return cfmt.Sprintf("{{%s}}::"+theme.MessageColor+"|bg"+messageBgColor, s) })
		reg("HighlightedMessage"+string(lvl), func(s string) string {
			return cfmt.Sprintf("{{%s}}::bg"+theme.MessageColor+"|"+messageBgColor, s)
		})
	}
	regLog(level.Trace, theme.TraceColor)
	regLog(level.Debug, theme.DebugColor)
	regLog(level.Info, theme.InfoColor)
	regLog(level.Warn, theme.WarnColor)
	regLog(level.Error, theme.ErrorColor)
}
