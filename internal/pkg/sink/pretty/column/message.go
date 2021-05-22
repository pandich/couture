package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"github.com/i582/cfmt/cmd/cfmt"
)

type messageColumn struct{}

// Name ...
func (m messageColumn) Name() string {
	return "message"
}

// Register ...
func (m messageColumn) Register(theme theme.Theme) {
	for _, lvl := range level.Levels {
		fgColor := theme.MessageColor()
		bgColor := theme.MessageBackgroundColor(lvl)
		cfmt.RegisterStyle(m.Name()+string(lvl), func(s string) string {
			return cfmt.Sprintf("{{%s}}::"+fgColor+"|bg"+bgColor, s)
		})
		cfmt.RegisterStyle("H"+m.Name()+string(lvl), func(s string) string {
			return cfmt.Sprintf("{{%s}}::bg"+fgColor+"|"+bgColor, s)
		})
	}
}

// Formatter ...
func (m messageColumn) Formatter(_ source.Source, _ model.Event) string {
	return "%s"
}

// Renderer ...
func (m messageColumn) Renderer(config config.Config, _ source.Source, event model.Event) []interface{} {
	var message = ""
	for _, chunk := range event.HighlightedMessage() {
		if message != "" {
			message += " "
		}
		switch chunk.(type) {
		case model.HighlightedMessage:
			message += cfmt.Sprintf("{{%s}}::H"+m.Name()+string(event.Level), chunk)
		case model.UnhighlightedMessage:
			message += cfmt.Sprintf("{{%s}}::"+m.Name()+string(event.Level), chunk)
		default:
			message += cfmt.Sprintf("{{%s}}::"+m.Name()+string(event.Level), chunk)
		}
	}
	var prefix = " "
	if config.MultiLine {
		prefix = "\n"
	}
	return []interface{}{prefix + message}
}
