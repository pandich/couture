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
func (col messageColumn) name() string { return "message" }

// weight ...
func (col messageColumn) weight() weight { return 0 }

// weightType ...
func (col messageColumn) weightType() weightType { return filling }

// RegisterStyles ...
func (col messageColumn) RegisterStyles(theme theme.Theme) {
	for _, lvl := range level.Levels {
		fgColor := theme.MessageColor()
		bgColor := theme.MessageBackgroundColor(lvl)
		cfmt.RegisterStyle(col.name()+string(lvl), func(s string) string {
			return cfmt.Sprintf("{{%s}}::"+fgColor+"|bg"+bgColor, s)
		})
		cfmt.RegisterStyle("H"+col.name()+string(lvl), func(s string) string {
			return cfmt.Sprintf("{{%s}}::bg"+fgColor+"|"+bgColor, s)
		})
	}
}

// Format ...
func (col messageColumn) Format(_ uint, _ source.Source, _ model.Event) string {
	return "%s"
}

// Render ...
func (col messageColumn) Render(config config.Config, _ source.Source, event model.Event) []interface{} {
	var message = ""
	for _, chunk := range event.HighlightedMessage() {
		if message != "" {
			message += " "
		}
		switch chunk.(type) {
		case model.HighlightedMessage:
			message += cfmt.Sprintf("{{%s}}::H"+col.name()+string(event.Level), chunk)
		case model.UnhighlightedMessage:
			message += cfmt.Sprintf("{{%s}}::"+col.name()+string(event.Level), chunk)
		default:
			message += cfmt.Sprintf("{{%s}}::"+col.name()+string(event.Level), chunk)
		}
	}
	var prefix = " "
	if config.MultiLine {
		prefix = "\n"
	}
	return []interface{}{prefix + message}
}
