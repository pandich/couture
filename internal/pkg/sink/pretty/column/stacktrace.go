package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"github.com/i582/cfmt/cmd/cfmt"
)

type stackTraceColumn struct{}

// Name ...
func (col stackTraceColumn) name() string { return "error" }

// weight ...
func (col stackTraceColumn) weight() weight { return 0 }

// weightType ...
func (col stackTraceColumn) weightType() weightType { return filling }

// RegisterStyles ...
func (col stackTraceColumn) RegisterStyles(theme theme.Theme) {
	c := theme.StackTraceColor()
	cfmt.RegisterStyle(col.name(), func(s string) string {
		return cfmt.Sprintf("{{%s}}::"+c, s)
	})
	cfmt.RegisterStyle("H"+col.name(), func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+c, s)
	})
}

// Format ...
func (col stackTraceColumn) Format(_ uint, _ source.Source, _ model.Event) string {
	return "%s"
}

// Render ...
func (col stackTraceColumn) Render(_ config.Config, _ source.Source, event model.Event) []interface{} {
	var stackTrace = ""
	for _, chunk := range event.HighlightedStackTrace() {
		if stackTrace == "" {
			stackTrace += "\n"
		} else {
			stackTrace += " "
		}
		switch chunk.(type) {
		case model.HighlightedStackTrace:
			stackTrace += cfmt.Sprintf("{{%s}}::H"+col.name(), chunk)
		case model.UnhighlightedStackTrace:
			stackTrace += cfmt.Sprintf("{{%s}}::"+col.name(), chunk)
		default:
			stackTrace += cfmt.Sprintf("{{%s}}::"+col.name(), chunk)
		}
	}
	return []interface{}{stackTrace}
}
