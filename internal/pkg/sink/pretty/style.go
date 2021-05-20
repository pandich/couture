package pretty

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"image/color"
	"net/url"
)

func init() {
	if !sink.IsTTY() {
		cfmt.DisableColors()
	}
}

type palette struct {
	defaultColor string
	sourceColors chan string
}

func newPalette(baseColor color.Color) palette {
	const defaultColor = "#ffffff"
	const timestampColor = "#877FD7"
	const errorColor = "#DD2A12"
	const traceColor = "#868686"
	const debugColor = "#F6F6F6"
	const infoColor = "#66A71E"
	const warnColor = "#FFE127"
	const messageColor = "#FBF0D7"
	const stackTraceColor = errorColor
	const punctuationColor = "#FEC8D8"
	sourceColorGenerator := gamut.PastelGenerator{}

	methodColor, classColor, lineNumberColor := sink.Triple(baseColor)
	threadRawColor, _ := colorful.MakeColor(gamut.Darker(gamut.Hex(classColor), 0.4))
	threadColor := threadRawColor.Hex()
	applicationRawColor, _ := colorful.MakeColor(gamut.Lighter(gamut.Hex(classColor), 0.4))
	applicationColor := applicationRawColor.Hex()

	reg := cfmt.RegisterStyle
	regLog := func(lvl level.Level, bgColor string) {
		fgColor := sink.HexContrast(bgColor)
		reg("Level"+string(lvl), func(s string) string { return cfmt.Sprintf("{{ %s }}::bg"+bgColor+"|"+fgColor, string(s[0])) })
	}

	reg("Default", func(s string) string { return s })
	reg("Punctuation", func(s string) string { return cfmt.Sprintf("{{%s}}::"+punctuationColor, s) })

	reg("Timestamp", func(s string) string { return cfmt.Sprintf("{{%s}}::"+timestampColor, s) })
	reg("ApplicationName", func(s string) string { return cfmt.Sprintf("{{%-20.20s}}::"+applicationColor, s) })
	reg("ThreadName", func(s string) string { return cfmt.Sprintf("{{%-15.15s}}::"+threadColor, s) })
	reg("ClassName", func(s string) string { return cfmt.Sprintf("{{%.30s}}::"+classColor, s) })
	reg("MethodName", func(s string) string { return cfmt.Sprintf("{{%.30s}}::"+methodColor, s) })
	reg("LineNumber", func(s string) string { return cfmt.Sprintf("{{%s}}::"+lineNumberColor, s) })

	regLog(level.Trace, traceColor)
	regLog(level.Debug, debugColor)
	regLog(level.Info, infoColor)
	regLog(level.Warn, warnColor)
	regLog(level.Error, errorColor)

	reg("Message", func(s string) string { return cfmt.Sprintf("{{%s}}::"+messageColor, s) })
	reg("HighlightedMessage", func(s string) string { return cfmt.Sprintf("{{%s}}::reverse|"+messageColor, s) })
	reg("StackTrace", func(s string) string { return cfmt.Sprintf("{{%s}}::"+stackTraceColor, s) })
	reg("HighlightedStackTrace", func(s string) string { return cfmt.Sprintf("{{%s}}::reverse|"+errorColor, s) })

	return palette{
		defaultColor: defaultColor,
		sourceColors: sink.NewColorCycle(sourceColorGenerator, defaultColor),
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
	name := p.sourceStyle(sourceURL)
	bgHex := <-p.sourceColors
	fgHex := sink.HexContrast(bgHex)
	cfmt.RegisterStyle(name, func(s string) string {
		return cfmt.Sprintf(fmt.Sprintf("{{%%-30.30s}}::%s|bg%s", fgHex, bgHex), s)
	})
}
