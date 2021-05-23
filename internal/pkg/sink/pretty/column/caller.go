package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
)

type callerColumn struct{}

// Name ...
func (col callerColumn) name() string { return "caller" }

// weight ...
func (col callerColumn) weight() weight {
	const columnWidth = 65
	return columnWidth
}

// weightType ...
func (col callerColumn) weightType() weightType { return fixed }

// RegisterStyles ...
func (col callerColumn) RegisterStyles(theme theme.Theme) {
	cfmt.RegisterStyle("Class", func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+theme.CallerBgColor()+"|"+theme.ClassColor(), s)
	})
	cfmt.RegisterStyle("MethodDelimiter", func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+theme.CallerBgColor()+"|"+theme.MethodDelimiterColor(), s)
	})
	cfmt.RegisterStyle("Method", func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+theme.CallerBgColor()+"|"+theme.MethodColor(), s)
	})
	cfmt.RegisterStyle("LineNumberDelimiter", func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+theme.CallerBgColor()+"|"+theme.LineNumberDelimiterColor(), s)
	})
	cfmt.RegisterStyle("LineNumber", func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+theme.CallerBgColor()+"|"+theme.LineNumberColor(), s)
	})
}

// Format ...
func (col callerColumn) Format(_ uint, _ source.Source, _ model.Event) string {
	return "{{ ☎︎ %s}}::Class{{∕}}::MethodDelimiter{{%s}}::Method{{#}}::LineNumberDelimiter{{%s }}::LineNumber "
}

// Render ...
func (col callerColumn) Render(_ config.Config, _ source.Source, event model.Event) []interface{} {
	const maxClassNameWidth = 30
	const maxWidth = 60

	var padding = ""
	className := string(event.ClassName.Abbreviate(maxClassNameWidth))
	var methodName = string(event.MethodName)
	lineNumber := fmt.Sprintf("%4d", event.LineNumber)
	totalLength := len(className) + len(methodName) + len(lineNumber)
	for i := totalLength; i < maxWidth; i++ {
		padding += " "
	}
	extraChars := totalLength - maxWidth
	if extraChars > 0 {
		methodName = methodName[:len(methodName)-extraChars-1]
	}

	ia := []interface{}{
		padding + className,
		methodName,
		lineNumber,
	}
	return ia
}
