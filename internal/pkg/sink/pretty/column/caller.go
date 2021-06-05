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

func newCallerColumn() column {
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

	bgColor := theme.CallerBg()

	cfmt.RegisterStyle("Class", func(s string) string {
		return cfmt.Sprintf("{{"+prefix+"︎%s}}::bg"+bgColor+"|"+theme.ClassFg(), s)
	})
	cfmt.RegisterStyle("MethodDelimiter", func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+bgColor+"|"+theme.MethodDelimiterFg(), s)
	})
	cfmt.RegisterStyle("Method", func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+bgColor+"|"+theme.MethodFg(), s)
	})
	cfmt.RegisterStyle("LineNumberDelimiter", func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+bgColor+"|"+theme.LineNumberDelimiterFg(), s)
	})
	cfmt.RegisterStyle("Line", func(s string) string {
		return cfmt.Sprintf("{{%s }}::bg"+bgColor+"|"+theme.LineNumberFg(), s)
	})
}

// Format ...
func (col callerColumn) Format(_ uint, evt model.SinkEvent) string {
	var s = "{{%s}}::Class"
	if evt.Method != "" {
		s += "{{∕}}::MethodDelimiter"
	}
	s += "{{%s}}::Method"
	if evt.Line != 0 {
		s += "{{#}}::LineNumberDelimiter"
	}
	s += "{{%s}}::Line"
	return s
}

// Render ...
func (col callerColumn) Render(_ config.Config, event model.SinkEvent) []interface{} {
	const maxClassNameWidth = 30
	const maxWidth = 60

	var padding = ""
	var className = orNoValue(string(event.Class.Abbreviate(maxClassNameWidth)))
	var methodName = string(event.Method)
	var lineNumber = ""
	if event.Line != 0 {
		lineNumber = fmt.Sprintf("%4d", event.Line)
	}
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
