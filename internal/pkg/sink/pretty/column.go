package pretty

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"github.com/muesli/reflow/padding"
)

const highlighted = "Highlighted"

// ColumnName ...
type ColumnName string
type columnRegister func(Theme)
type columnFormat func(source.Source, model.Event) string
type columnValue func(source.Source, model.Event) string

type column struct {
	format   columnFormat
	value    columnValue
	register columnRegister
}

const (
	applicationColumn ColumnName = "application"
	callerColumn      ColumnName = "caller"
	levelColumn       ColumnName = "level"
	messageColumn     ColumnName = "message"
	sourceColumn      ColumnName = "source"
	stackTraceColumn  ColumnName = "stackTrace"
	threadColumn      ColumnName = "thread"
	timestampColumn   ColumnName = "timestamp"
)

var defaultColumnOrder = []ColumnName{
	timestampColumn,
	applicationColumn,
	threadColumn,
	callerColumn,
	levelColumn,
	messageColumn,
	stackTraceColumn,
}

var columns = map[ColumnName]column{
	applicationColumn: {
		format: func(src source.Source, evt model.Event) string {
			return "{{%s}}::" + string(applicationColumn)
		},
		value: func(src source.Source, event model.Event) string {
			return string(event.ApplicationNameOrBlank())
		},
		register: func(theme Theme) {
			cfmt.RegisterStyle(string(applicationColumn), func(s string) string {
				return cfmt.Sprintf("{{ %-20.20s }}::"+sink.WithFaintBg(theme.ApplicationColor), s)
			})
		},
	},

	callerColumn: {
		format: func(src source.Source, evt model.Event) string {
			return "%s"
		},
		value: func(src source.Source, event model.Event) string {
			const classNameWidth = 30
			const callerWidth = 55
			caller := padding.String(cfmt.Sprintf(
				"{{%s}}::Class{{/}}::MethodDelimiter{{%s}}::Method{{#}}::LineNumberDelimiter{{%d}}::LineNumber  ",
				event.ClassName.Abbreviate(classNameWidth),
				event.MethodName,
				event.LineNumber,
			), callerWidth)
			return caller
		},
		register: func(theme Theme) {
			const contrastPercent = 0.25
			callerPalette := func(center string) (string, string, string) {
				col := gamut.Hex(center)
				q := gamut.Analogous(col)
				a, _ := colorful.MakeColor(col)
				b, _ := colorful.MakeColor(q[0])
				c, _ := colorful.MakeColor(q[1])
				return a.Hex(), b.Hex(), c.Hex()
			}
			methodColor, classColor, lineNumberColor := callerPalette(theme.BaseColor)

			var methodDelimiterColor = sink.Lighter(methodColor, contrastPercent)
			var lineNumberDelimiterColor = sink.Lighter(lineNumberColor, contrastPercent)
			if sink.IsDark(methodColor) {
				methodDelimiterColor = sink.Darker(methodColor, contrastPercent)
				lineNumberDelimiterColor = sink.Darker(lineNumberColor, contrastPercent)
			}

			cfmt.RegisterStyle("Class", func(s string) string {
				return cfmt.Sprintf("{{%.30s}}::"+classColor, s)
			})
			cfmt.RegisterStyle("MethodDelimiter", func(s string) string {
				return cfmt.Sprintf("{{%s}}::"+methodDelimiterColor, s)
			})
			cfmt.RegisterStyle("Method", func(s string) string {
				return cfmt.Sprintf("{{%.30s}}::"+methodColor, s)
			})
			cfmt.RegisterStyle("LineNumberDelimiter", func(s string) string {
				return cfmt.Sprintf("{{%s}}::"+lineNumberDelimiterColor, s)
			})
			cfmt.RegisterStyle("LineNumber", func(s string) string {
				return cfmt.Sprintf("{{%s}}::"+lineNumberColor, s)
			})
		},
	},

	levelColumn: {
		format: func(src source.Source, evt model.Event) string {
			return "{{%s}}::" + string(levelColumn) + string(evt.Level)
		},
		value: func(src source.Source, event model.Event) string {
			return string(event.Level)
		},
		register: func(theme Theme) {
			reg := func(lvl level.Level, bgColor string) {
				fgColor := sink.Contrast(bgColor)
				cfmt.RegisterStyle(string(levelColumn)+string(lvl), func(s string) string {
					return cfmt.Sprintf("{{ %1.1s }}::bg"+bgColor+"|"+fgColor, s)
				})
			}
			reg(level.Trace, theme.TraceColor)
			reg(level.Debug, theme.DebugColor)
			reg(level.Info, theme.InfoColor)
			reg(level.Warn, theme.WarnColor)
			reg(level.Error, theme.ErrorColor)
		},
	},

	messageColumn: {
		format: func(src source.Source, evt model.Event) string { return "%s" },
		value: func(src source.Source, event model.Event) string {
			var message = ""
			for _, chunk := range event.HighlightedMessage() {
				message += " "
				switch chunk.(type) {
				case model.HighlightedMessage:
					message += cfmt.Sprintf("{{%s}}::"+highlighted+string(messageColumn)+string(event.Level), chunk)
				case model.UnhighlightedMessage:
					message += cfmt.Sprintf("{{%s}}::"+string(messageColumn)+string(event.Level), chunk)
				default:
					message += cfmt.Sprintf("{{%s}}::"+string(messageColumn)+string(event.Level), chunk)
				}
			}
			return message
		},
		register: func(theme Theme) {
			reg := func(lvl level.Level, bgColor string) {
				messageBgColor := sink.Fainter(bgColor, 0.90)
				cfmt.RegisterStyle(string(messageColumn)+string(lvl), func(s string) string {
					return cfmt.Sprintf("{{%s}}::"+theme.MessageColor+"|bg"+messageBgColor, s)
				})
				cfmt.RegisterStyle(highlighted+string(messageColumn)+string(lvl), func(s string) string {
					return cfmt.Sprintf("{{%s}}::bg"+theme.MessageColor+"|"+messageBgColor, s)
				})
			}
			reg(level.Trace, theme.TraceColor)
			reg(level.Debug, theme.DebugColor)
			reg(level.Info, theme.InfoColor)
			reg(level.Warn, theme.WarnColor)
			reg(level.Error, theme.ErrorColor)
		},
	},

	sourceColumn: {
		format: func(src source.Source, evt model.Event) string {
			styleName := "source" + src.ID()
			return "{{%s}}::" + styleName
		},
		value: func(src source.Source, event model.Event) string {
			return src.URL().ShortForm()
		},
		register: func(theme Theme) {},
	},

	stackTraceColumn: {
		format: func(src source.Source, evt model.Event) string { return "%s" },
		value: func(src source.Source, event model.Event) string {
			var stackTrace = ""
			for _, chunk := range event.HighlightedStackTrace() {
				if stackTrace == "" {
					stackTrace += "\n"
				} else {
					stackTrace += " "
				}
				switch chunk.(type) {
				case model.HighlightedStackTrace:
					stackTrace += cfmt.Sprintf("{{%s}}::"+highlighted+string(stackTraceColumn), chunk)
				case model.UnhighlightedStackTrace:
					stackTrace += cfmt.Sprintf("{{%s}}::"+string(stackTraceColumn), chunk)
				default:
					stackTrace += cfmt.Sprintf("{{%s}}::"+string(stackTraceColumn), chunk)
				}
			}
			return stackTrace
		},
		register: func(theme Theme) {
			cfmt.RegisterStyle(string(stackTraceColumn), func(s string) string {
				return cfmt.Sprintf("{{%s}}::"+sink.WithFaintBg(theme.StackTraceColor), s)
			})
			cfmt.RegisterStyle(highlighted+string(stackTraceColumn), func(s string) string {
				return cfmt.Sprintf("{{%s}}::bg"+sink.WithFaintBg(theme.ErrorColor), s)
			})
		},
	},

	threadColumn: {
		format: func(src source.Source, evt model.Event) string {
			return "{{%s}}::" + string(threadColumn)
		},
		value: func(src source.Source, event model.Event) string {
			return string(event.ThreadNameOrBlank())
		},
		register: func(theme Theme) {
			d, _ := colorful.MakeColor(gamut.Darker(gamut.Hex(theme.BaseColor), 0.5))
			cfmt.RegisterStyle(string(threadColumn), func(s string) string {
				return cfmt.Sprintf("{{ %-15.15s }}::"+sink.WithFaintBg(d.Hex()), s)
			})
		},
	},

	timestampColumn: {
		format: func(src source.Source, evt model.Event) string {
			return "{{%s}}::" + string(timestampColumn)
		},
		value: func(src source.Source, event model.Event) string {
			return event.Timestamp.Stamp()
		},
		register: func(theme Theme) {
			cfmt.RegisterStyle(string(timestampColumn), func(s string) string {
				return cfmt.Sprintf("{{ %s }}::"+sink.WithFaintBg(theme.TimestampColor), s)
			})
		},
	},
}
