package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"github.com/i582/cfmt/cmd/cfmt"
)

type stackTraceColumn struct {
}

// Name ...
func (s stackTraceColumn) Name() string {
	return "error"
}

// Register ...
func (s stackTraceColumn) Register(theme theme.Theme) {
	c := theme.StackTraceColor()
	cfmt.RegisterStyle(s.Name(), func(s string) string {
		return cfmt.Sprintf("{{%s}}::"+c, s)
	})
	cfmt.RegisterStyle("H"+s.Name(), func(s string) string {
		return cfmt.Sprintf("{{%s}}::bg"+c, s)
	})
}

// Formatter ...
func (s stackTraceColumn) Formatter(_ source.Source, _ model.Event) string {
	return "%s"
}

// Renderer ...
func (s stackTraceColumn) Renderer(_ config.Config, _ source.Source, event model.Event) []interface{} {
	var stackTrace = ""
	for _, chunk := range event.HighlightedStackTrace() {
		if stackTrace == "" {
			stackTrace += "\n"
		} else {
			stackTrace += " "
		}
		switch chunk.(type) {
		case model.HighlightedStackTrace:
			stackTrace += cfmt.Sprintf("{{%s}}::H"+s.Name(), chunk)
		case model.UnhighlightedStackTrace:
			stackTrace += cfmt.Sprintf("{{%s}}::"+s.Name(), chunk)
		default:
			stackTrace += cfmt.Sprintf("{{%s}}::"+s.Name(), chunk)
		}
	}
	return []interface{}{stackTrace}
}
