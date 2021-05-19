package pretty

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"couture/pkg/model/level"
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
	"math/rand"
)

func init() {
	if !sink.IsTTY() {
		cfmt.DisableColors()
	}
}

func init() {
	reg := cfmt.RegisterStyle
	regLog := func(lvl level.Level, color string) {
		reg("Log"+string(lvl), func(s string) string { return cfmt.Sprintf("{{ %s }}::"+color+"|reverse", string(s[0])) })
	}

	reg("Punctuation", func(s string) string { return cfmt.Sprintf("{{%s}}::#FEC8D8", s) })

	reg("Timestamp", func(s string) string { return cfmt.Sprintf("{{%s}}::#877FD7", s) })
	reg("ApplicationName", func(s string) string { return cfmt.Sprintf("{{%-20.20s}}::#957DAD", s) })
	reg("ThreadName", func(s string) string { return cfmt.Sprintf("{{%-15.15s}}::#808080", s) })
	reg("ClassName", func(s string) string { return cfmt.Sprintf("{{%.30s}}::#D291BC", s) })
	reg("MethodName", func(s string) string { return cfmt.Sprintf("{{%.30s}}::#E0BBE4", s) })
	reg("LineNumber", func(s string) string { return cfmt.Sprintf("{{%s}}::#FFDFD3", s) })

	regLog(level.Trace, "#868686")
	regLog(level.Debug, "#F6F6F6")
	regLog(level.Info, "#66A71E")
	regLog(level.Warn, "#FFE127")
	regLog(level.Error, "#DD2A12")

	const messageColor = "#FBF0D7"
	const stackTraceColor = "#DD2A12"
	reg("Message", func(s string) string { return cfmt.Sprintf("{{%s}}::"+messageColor, s) })
	reg("HighlightedMessage", func(s string) string { return cfmt.Sprintf("{{%s}}::reverse|"+messageColor, s) })
	reg("StackTrace", func(s string) string { return cfmt.Sprintf("{{%s}}::"+stackTraceColor, s) })
	reg("HighlightedStackTrace", func(s string) string { return cfmt.Sprintf("{{%s}}::reverse|"+stackTraceColor, s) })
}

func (snk *prettySink) getSourceStyleName(src source.Pushable) string {
	{
		snk.sourceStyleMutex.RLock()
		name, ok := snk.sourceStyle[src.URL()]
		snk.sourceStyleMutex.RUnlock()
		if ok {
			return name
		}
	}
	snk.sourceStyleMutex.Lock()
	defer snk.sourceStyleMutex.Unlock()
	if name, ok := snk.sourceStyle[src.URL()]; ok {
		return name
	}
	//nolint:gosec
	name := fmt.Sprintf("%d", rand.Uint32())
	cfmt.RegisterStyle(name, func(s string) string { return cfmt.Sprintf("{{%-30.30s}}::reverse|"+<-snk.sourceColors, s) })
	snk.sourceStyle[src.URL()] = name
	return name
}
