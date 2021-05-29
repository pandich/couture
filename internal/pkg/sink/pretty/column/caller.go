package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
)

type callerColumn struct {
	baseColumn
}

func newCallerColumn() callerColumn {
	const width = 50
	sigil := '☎'
	return callerColumn{baseColumn{
		columnName:  "caller",
		widthMode:   fixed,
		widthWeight: width,
		sigil:       &sigil,
	}}
}

// RegisterStyles ...
func (col callerColumn) RegisterStyles(theme theme.Theme) {
	var prefix = ""
	if col.sigil != nil {
		prefix = " " + string(*col.sigil) + " "
	}

	cfmt.RegisterStyle("Class", func(s string) string {
		return cfmt.Sprintf("{{"+prefix+"︎%s}}::bg"+theme.CallerBg()+"|"+theme.ClassFg(), s)
	})
	cfmt.RegisterStyle("MethodDelimiter", func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+theme.CallerBg()+"|"+theme.MethodDelimiterFg(), s)
	})
	cfmt.RegisterStyle("Method", func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+theme.CallerBg()+"|"+theme.MethodFg(), s)
	})
	cfmt.RegisterStyle("LineNumberDelimiter", func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+theme.CallerBg()+"|"+theme.LineNumberDelimiterFg(), s)
	})
	cfmt.RegisterStyle("Line", func(s string) string {
		return cfmt.Sprintf("{{%s }}::bg"+theme.CallerBg()+"|"+theme.LineNumberFg(), s)
	})
}

// Format ...
func (col callerColumn) Format(_ uint, _ model.SinkEvent) string {
	return "{{%s}}::Class" +
		"{{∕}}::MethodDelimiter" + "{{%s}}::Method" +
		"{{#}}::LineNumberDelimiter" + "{{%s}}::Line"
}

// Render ...
func (col callerColumn) Render(_ config.Config, event model.SinkEvent) []interface{} {
	const maxClassNameWidth = 30
	const maxWidth = 60

	var padding = ""
	className := string(event.Class.Abbreviate(maxClassNameWidth))
	var methodName = string(event.Method)
	lineNumber := fmt.Sprintf("%4d", event.Line)
	totalLength := len(className) + len(methodName) + len(lineNumber)
	for i := totalLength; i < maxWidth; i++ {
		padding += " "
	}
	extraChars := totalLength - maxWidth
	if extraChars > 0 {
		methodName = methodName[:len(methodName)-extraChars-1]
	}

	return []interface{}{
		padding + className,
		methodName,
		lineNumber,
	}
}
