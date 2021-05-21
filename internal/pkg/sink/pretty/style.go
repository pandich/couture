package pretty

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink"
	"crypto/sha256"
	"encoding/hex"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"github.com/muesli/termenv"
	"net/url"
)

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

type palette struct {
	defaultColor string
	sourceColors chan string
}

func newPalette(theme Theme) palette {
	faintly := func(hex string) string { return hex + "|bg" + fainter(hex, 0.96) }

	methodColor, classColor, lineNumberColor, threadColor := caller(theme.BaseColor)
	const contrastPercent = 0.25
	var methodDelimiterColor = lighter(methodColor, contrastPercent)
	var lineNumberDelimiterColor = lighter(lineNumberColor, contrastPercent)
	if isDark(methodColor) {
		methodDelimiterColor = darker(methodColor, contrastPercent)
		lineNumberDelimiterColor = darker(lineNumberColor, contrastPercent)
	}

	reg := cfmt.RegisterStyle
	regLog := func(lvl level.Level, bgColor string) {
		messageBgColor := fainter(bgColor, 0.90)
		fgColor := contrast(bgColor)
		reg("Level"+string(lvl), func(s string) string { return cfmt.Sprintf("{{ %1.1s }}::bg"+bgColor+"|"+fgColor, s) })
		reg("Message"+string(lvl), func(s string) string { return cfmt.Sprintf("{{%s}}::"+theme.MessageColor+"|bg"+messageBgColor, s) })
		reg("HighlightedMessage"+string(lvl), func(s string) string {
			return cfmt.Sprintf("{{%s}}::bg"+theme.MessageColor+"|"+messageBgColor, s)
		})
	}

	reg("Default", func(s string) string { return s })
	reg("MethodDelimiter", func(s string) string { return cfmt.Sprintf("{{%s}}::"+methodDelimiterColor, s) })
	reg("LineNumberDelimiter", func(s string) string { return cfmt.Sprintf("{{%s}}::"+lineNumberDelimiterColor, s) })

	reg("Timestamp", func(s string) string { return cfmt.Sprintf("{{ %s }}::"+faintly(theme.TimestampColor), s) })
	reg("Application", func(s string) string {
		return cfmt.Sprintf("{{ %-20.20s }}::"+faintly(theme.ApplicationColor), s)
	})
	reg("Thread", func(s string) string { return cfmt.Sprintf("{{ %-15.15s }}::"+faintly(threadColor), s) })
	reg("Class", func(s string) string { return cfmt.Sprintf("{{%.30s}}::"+classColor, s) })
	reg("Method", func(s string) string { return cfmt.Sprintf("{{%.30s}}::"+methodColor, s) })
	reg("LineNumber", func(s string) string { return cfmt.Sprintf("{{%s}}::"+lineNumberColor, s) })

	regLog(level.Trace, theme.TraceColor)
	regLog(level.Debug, theme.DebugColor)
	regLog(level.Info, theme.InfoColor)
	regLog(level.Warn, theme.WarnColor)
	regLog(level.Error, theme.ErrorColor)

	reg("StackTrace", func(s string) string { return cfmt.Sprintf("{{%s}}::"+faintly(theme.StackTraceColor), s) })
	reg("HighlightedStackTrace", func(s string) string { return cfmt.Sprintf("{{%s}}::bg"+faintly(theme.ErrorColor), s) })

	return palette{
		defaultColor: theme.DefaultColor,
		sourceColors: sink.NewColorCycle(theme.SourceColors, theme.DefaultColor),
	}
}

func (p *palette) sourceStyle(sourceURL model.SourceURL) string {
	u := url.URL(sourceURL)
	if s := u.String(); s != "" {
		hasher := sha256.New()
		hasher.Write([]byte(s))
		return "Source" + hex.EncodeToString(hasher.Sum(nil))
	}
	return "Default"
}

func (p *palette) registerSource(sourceURL model.SourceURL) {
	styleName := p.sourceStyle(sourceURL)
	sourceColor := <-p.sourceColors
	cfmt.RegisterStyle(styleName, func(s string) string { return cfmt.Sprintf("{{/%-30.30s }}::"+sourceColor, s) })
}

func caller(center string) (string, string, string, string) {
	col := gamut.Hex(center)
	q := gamut.Analogous(col)
	a, _ := colorful.MakeColor(col)
	b, _ := colorful.MakeColor(q[0])
	c, _ := colorful.MakeColor(q[1])
	d, _ := colorful.MakeColor(gamut.Darker(col, 0.5))
	return a.Hex(), b.Hex(), c.Hex(), d.Hex()
}

func contrast(hex string) string {
	cf, _ := colorful.MakeColor(gamut.Contrast(gamut.Hex(hex)))
	return cf.Hex()
}

func lighter(hex string, percent float64) string {
	cf, _ := colorful.MakeColor(gamut.Lighter(gamut.Hex(hex), percent))
	return cf.Hex()
}

func darker(hex string, percent float64) string {
	cf, _ := colorful.MakeColor(gamut.Darker(gamut.Hex(hex), percent))
	return cf.Hex()
}
func fainter(hex string, percent float64) string {
	const count = 1000
	i := int(count * percent)
	bg := termenv.ConvertToRGB(termenv.BackgroundColor())
	fainter := gamut.Blends(bg, gamut.Hex(hex), count)[count-i]
	col, _ := colorful.MakeColor(fainter)
	return col.Hex()
}

func isDark(hex string) bool {
	const gray = 0.5
	col, _ := colorful.MakeColor(gamut.Hex(hex))
	_, _, l := col.Hcl()
	return l < gray
}
