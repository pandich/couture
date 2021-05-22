package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"couture/internal/pkg/tty"
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
	"time"
)

// TODO column widths should adapt to the terminal

const (
	applicationColumn = "application"
	callerColumn      = "caller"
	levelColumn       = "level"
	messageColumn     = "message"
	sourceColumn      = "source"
	stackTraceColumn  = "stackTrace"
	threadColumn      = "thread"
	timestampColumn   = "timestamp"
)

var defaultColumns = []string{
	timestampColumn,
	applicationColumn,
	threadColumn,
	callerColumn,
	levelColumn,
	messageColumn,
	stackTraceColumn,
}

const highlight = "HL"

// ColumnName ...
type (
	registerer func(theme.Theme)
	formatter  func(source.Source, model.Event) string
	renderer   func(config.Config, source.Source, model.Event) []interface{}

	column struct {
		name      string
		Formatter formatter
		Renderer  renderer
		register  registerer
	}
)

// TODO handle config.ShowSigil
// TODO it is clumsy that the column name has to be specified three times per column

var columns = []column{
	{
		name: applicationColumn,
		Formatter: func(src source.Source, evt model.Event) string {
			return "{{ § %-20.20s }}::" + string(applicationColumn)
		},
		Renderer: func(_ config.Config, src source.Source, event model.Event) []interface{} {
			return []interface{}{string(event.ApplicationNameOrBlank())}
		},
		register: func(theme theme.Theme) {
			cfmt.RegisterStyle(applicationColumn, func(s string) string {
				return cfmt.Sprintf("{{%s}}::"+theme.ApplicationColor(), s)
			})
		},
	},

	{
		name: callerColumn,
		Formatter: func(src source.Source, evt model.Event) string {
			return "{{ ☎︎ %s}}::Class{{/}}::MethodDelimiter{{%s}}::Method{{#}}::LineNumberDelimiter{{%s }}::LineNumber "
		},
		Renderer: func(_ config.Config, src source.Source, event model.Event) []interface{} {
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
		},
		register: func(theme theme.Theme) {
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
		},
	},

	{
		name: levelColumn,
		Formatter: func(src source.Source, evt model.Event) string {
			return "{{ %1.1s }}::" + string(levelColumn) + string(evt.Level)
		},
		Renderer: func(_ config.Config, src source.Source, event model.Event) []interface{} {
			return []interface{}{string(event.Level)}
		},
		register: func(theme theme.Theme) {
			for _, lvl := range level.Levels {
				bgColor := theme.LevelColor(lvl)
				fgColor := tty.Contrast(bgColor)
				cfmt.RegisterStyle(levelColumn+string(lvl), func(s string) string {
					return cfmt.Sprintf("{{%s}}::bg"+bgColor+"|"+fgColor, s)
				})
			}
		},
	},

	{
		name:      messageColumn,
		Formatter: func(src source.Source, evt model.Event) string { return "%s" },
		Renderer: func(config config.Config, src source.Source, event model.Event) []interface{} {
			var message = ""
			for _, chunk := range event.HighlightedMessage() {
				if message != "" {
					message += " "
				}
				switch chunk.(type) {
				case model.HighlightedMessage:
					message += cfmt.Sprintf("{{%s}}::"+highlight+string(messageColumn)+string(event.Level), chunk)
				case model.UnhighlightedMessage:
					message += cfmt.Sprintf("{{%s}}::"+string(messageColumn)+string(event.Level), chunk)
				default:
					message += cfmt.Sprintf("{{%s}}::"+string(messageColumn)+string(event.Level), chunk)
				}
			}
			var prefix = " "
			if config.MultiLine {
				prefix = "\n"
			}
			return []interface{}{prefix + message}
		},
		register: func(theme theme.Theme) {
			for _, lvl := range level.Levels {
				fgColor := theme.MessageColor()
				bgColor := theme.MessageBackgroundColor(lvl)
				cfmt.RegisterStyle(messageColumn+string(lvl), func(s string) string {
					return cfmt.Sprintf("{{%s}}::"+fgColor+"|bg"+bgColor, s)
				})
				cfmt.RegisterStyle(highlight+string(messageColumn)+string(lvl), func(s string) string {
					return cfmt.Sprintf("{{%s}}::bg"+fgColor+"|"+bgColor, s)
				})
			}
		},
	},

	{
		name: sourceColumn,
		Formatter: func(src source.Source, evt model.Event) string {
			return "{{" + string(src.Sigil()) + " %-30.30s }}::" + src.ID()
		},
		Renderer: func(_ config.Config, src source.Source, event model.Event) []interface{} {
			return []interface{}{src.URL().ShortForm()}
		},
		register: func(theme theme.Theme) {},
	},

	{
		name:      stackTraceColumn,
		Formatter: func(src source.Source, evt model.Event) string { return "%s" },
		Renderer: func(_ config.Config, src source.Source, event model.Event) []interface{} {
			var stackTrace = ""
			for _, chunk := range event.HighlightedStackTrace() {
				if stackTrace == "" {
					stackTrace += "\n"
				} else {
					stackTrace += " "
				}
				switch chunk.(type) {
				case model.HighlightedStackTrace:
					stackTrace += cfmt.Sprintf("{{%s}}::"+highlight+string(stackTraceColumn), chunk)
				case model.UnhighlightedStackTrace:
					stackTrace += cfmt.Sprintf("{{%s}}::"+string(stackTraceColumn), chunk)
				default:
					stackTrace += cfmt.Sprintf("{{%s}}::"+string(stackTraceColumn), chunk)
				}
			}
			return []interface{}{stackTrace}
		},
		register: func(theme theme.Theme) {
			c := theme.StackTraceColor()
			cfmt.RegisterStyle(stackTraceColumn, func(s string) string {
				return cfmt.Sprintf("{{%s}}::"+c, s)
			})
			cfmt.RegisterStyle(highlight+string(stackTraceColumn), func(s string) string {
				return cfmt.Sprintf("{{%s}}::bg"+c, s)
			})
		},
	},

	{
		name: threadColumn,
		Formatter: func(src source.Source, evt model.Event) string {
			return "{{ ⇶ %-15.15s }}::" + string(threadColumn)
		},
		Renderer: func(_ config.Config, src source.Source, event model.Event) []interface{} {
			return []interface{}{string(event.ThreadNameOrBlank())}
		},
		register: func(theme theme.Theme) {
			cfmt.RegisterStyle(threadColumn, func(s string) string {
				return cfmt.Sprintf("{{%s}}::"+theme.ThreadColor(), s)
			})
		},
	},

	{
		name: timestampColumn,
		Formatter: func(src source.Source, evt model.Event) string {
			return "{{ ⌚ %s }}::" + string(timestampColumn)
		},
		Renderer: func(config config.Config, src source.Source, event model.Event) []interface{} {
			return []interface{}{time.Time(event.Timestamp).Format(config.TimeFormat)}
		},
		register: func(theme theme.Theme) {
			cfmt.RegisterStyle(timestampColumn, func(s string) string {
				return cfmt.Sprintf("{{%s}}::"+theme.TimestampColor(), s)
			})
		},
	},
}

// EffectiveColumns ...
func EffectiveColumns(cfg config.Config) []string {
	var effectiveColumns = defaultColumns
	if len(cfg.Columns) > 0 {
		effectiveColumns = cfg.Columns
	}
	effectiveColumns = append([]string{sourceColumn}, effectiveColumns...)
	return effectiveColumns
}
