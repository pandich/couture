package pretty

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/source"
	"couture/internal/pkg/tty"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"github.com/muesli/reflow/padding"
	"time"
)

// TODO render's have a lot of potentially expensive (or at least highly redundant) operations
// TODO column widths should adapt to the terminal

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

const highlight = "HL"

// ColumnName ...
type (
	// ColumnName ...
	ColumnName string

	columnRegisterer func(Theme)
	columnFormatter  func(source.Source, model.Event) string
	columnRenderer   func(Config, source.Source, model.Event) string

	column struct {
		formatter columnFormatter
		renderer  columnRenderer
		register  columnRegisterer
	}

	columnRegistry map[ColumnName]column
)

func (r columnRegistry) init(theme Theme) {
	for _, col := range r {
		col.register(theme)
	}
}

var columns = columnRegistry{
	applicationColumn: {
		formatter: func(src source.Source, evt model.Event) string {
			return "{{%s}}::" + string(applicationColumn)
		},
		renderer: func(_ Config, src source.Source, event model.Event) string {
			return string(event.ApplicationNameOrBlank())
		},
		register: func(theme Theme) {
			cfmt.RegisterStyle(string(applicationColumn), func(s string) string {
				return cfmt.Sprintf("{{ %-20.20s }}::"+theme.applicationColor(), s)
			})
		},
	},

	callerColumn: {
		formatter: func(src source.Source, evt model.Event) string {
			return "%s"
		},
		renderer: func(_ Config, src source.Source, event model.Event) string {
			const classNameWidth = 30
			const callerWidth = 55
			caller := padding.String(cfmt.Sprintf(
				"{{ %s}}::Class{{/}}::MethodDelimiter{{%s}}::Method{{#}}::LineNumberDelimiter{{%d }}::LineNumber",
				event.ClassName.Abbreviate(classNameWidth),
				event.MethodName,
				event.LineNumber,
			), callerWidth)
			return caller
		},
		register: func(theme Theme) {
			cfmt.RegisterStyle("Class", func(s string) string {
				return cfmt.Sprintf("{{%.30s}}::"+theme.classColor(), s)
			})
			cfmt.RegisterStyle("MethodDelimiter", func(s string) string {
				return cfmt.Sprintf("{{%s}}::"+theme.methodDelimiterColor(), s)
			})
			cfmt.RegisterStyle("Method", func(s string) string {
				return cfmt.Sprintf("{{%.30s}}::"+theme.methodColor(), s)
			})
			cfmt.RegisterStyle("LineNumberDelimiter", func(s string) string {
				return cfmt.Sprintf("{{%s}}::"+theme.lineNumberDelimiterColor(), s)
			})
			cfmt.RegisterStyle("LineNumber", func(s string) string {
				return cfmt.Sprintf("{{%s}}::"+theme.lineNumberColor(), s)
			})
		},
	},

	levelColumn: {
		formatter: func(src source.Source, evt model.Event) string {
			return "{{%s}}::" + string(levelColumn) + string(evt.Level)
		},
		renderer: func(_ Config, src source.Source, event model.Event) string {
			return string(event.Level)
		},
		register: func(theme Theme) {
			for _, lvl := range level.Levels {
				bgColor := theme.levelColor(lvl)
				fgColor := tty.Contrast(bgColor)
				cfmt.RegisterStyle(string(levelColumn)+string(lvl), func(s string) string {
					return cfmt.Sprintf("{{ %1.1s }}::bg"+bgColor+"|"+fgColor, s)
				})
			}
		},
	},

	messageColumn: {
		formatter: func(src source.Source, evt model.Event) string { return "%s" },
		renderer: func(config Config, src source.Source, event model.Event) string {
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
			return prefix + message
		},
		register: func(theme Theme) {
			for _, lvl := range level.Levels {
				fgColor := theme.messageColor()
				bgColor := theme.messageBackgroundColor(lvl)
				cfmt.RegisterStyle(string(messageColumn)+string(lvl), func(s string) string {
					return cfmt.Sprintf("{{%s}}::"+fgColor+"|bg"+bgColor, s)
				})
				cfmt.RegisterStyle(highlight+string(messageColumn)+string(lvl), func(s string) string {
					return cfmt.Sprintf("{{%s}}::bg"+fgColor+"|"+bgColor, s)
				})
			}
		},
	},

	sourceColumn: {
		formatter: func(src source.Source, evt model.Event) string {
			return "{{%s}}::" + src.ID()
		},
		renderer: func(_ Config, src source.Source, event model.Event) string {
			return src.URL().ShortForm()
		},
		register: func(theme Theme) {},
	},

	stackTraceColumn: {
		formatter: func(src source.Source, evt model.Event) string { return "%s" },
		renderer: func(_ Config, src source.Source, event model.Event) string {
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
			return stackTrace
		},
		register: func(theme Theme) {
			c := theme.stackTraceColor()
			cfmt.RegisterStyle(string(stackTraceColumn), func(s string) string {
				return cfmt.Sprintf("{{%s}}::"+c, s)
			})
			cfmt.RegisterStyle(highlight+string(stackTraceColumn), func(s string) string {
				return cfmt.Sprintf("{{%s}}::bg"+c, s)
			})
		},
	},

	threadColumn: {
		formatter: func(src source.Source, evt model.Event) string {
			return "{{%s}}::" + string(threadColumn)
		},
		renderer: func(_ Config, src source.Source, event model.Event) string {
			return string(event.ThreadNameOrBlank())
		},
		register: func(theme Theme) {
			d := tty.SimilarBg(tty.Darker(theme.BaseColor, 0.5))
			cfmt.RegisterStyle(string(threadColumn), func(s string) string {
				return cfmt.Sprintf("{{ %-15.15s }}::"+d, s)
			})
		},
	},

	timestampColumn: {
		formatter: func(src source.Source, evt model.Event) string {
			return "{{%s}}::" + string(timestampColumn)
		},
		renderer: func(config Config, src source.Source, event model.Event) string {
			return time.Time(event.Timestamp).Format(config.TimeFormat)
		},
		register: func(theme Theme) {
			const degrees60 = 60 / 360.0
			var yellow = colorful.Hcl(degrees60, 1, 1)
			cf, _ := colorful.MakeColor(gamut.Tints(gamut.Complementary(gamut.Hex(theme.BaseColor)), 3)[1])
			timestampColor := gamut.Blends(yellow, cf, 16)[3]
			timestampCf, _ := colorful.MakeColor(timestampColor)
			cfmt.RegisterStyle(string(timestampColumn), func(s string) string {
				return cfmt.Sprintf("{{ %s }}::"+tty.SimilarBg(timestampCf.Hex()), s)
			})
		},
	},
}

func registerSourceStyle(src source.Source, styleColor string) {
	cfmt.RegisterStyle(src.ID(), func(s string) string { return cfmt.Sprintf("{{/%-30.30s }}::"+styleColor, s) })
}
