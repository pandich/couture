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

const degrees60 = 60 / 360.0

var yellow = colorful.Hcl(degrees60, 1, 1)

const highlighted = "H"

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
			cf, _ := colorful.MakeColor(gamut.Quadratic(gamut.Hex(theme.BaseColor))[2])
			h, _, _ := cf.Hcl()
			applicationColor := colorful.Hcl(h, 0.9, 0.8).Hex()
			cfmt.RegisterStyle(string(applicationColumn), func(s string) string {
				return cfmt.Sprintf("{{ %-20.20s }}::"+tty.SimilarBg(applicationColor), s)
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
			const contrastPercent = 0.25
			col := gamut.Hex(theme.BaseColor)
			q := gamut.Analogous(col)
			a, _ := colorful.MakeColor(gamut.Darker(col, 0.4))
			b, _ := colorful.MakeColor(q[0])
			c, _ := colorful.MakeColor(q[1])
			methodColor, classColor, lineNumberColor := a.Hex(), b.Hex(), c.Hex()

			var methodDelimiterColor = tty.Lighter(methodColor, contrastPercent)
			var lineNumberDelimiterColor = tty.Lighter(lineNumberColor, contrastPercent)
			if tty.IsDark(methodColor) {
				methodDelimiterColor = tty.Darker(methodColor, contrastPercent)
				lineNumberDelimiterColor = tty.Darker(lineNumberColor, contrastPercent)
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
		formatter: func(src source.Source, evt model.Event) string {
			return "{{%s}}::" + string(levelColumn) + string(evt.Level)
		},
		renderer: func(_ Config, src source.Source, event model.Event) string {
			return string(event.Level)
		},
		register: func(theme Theme) {
			reg := func(lvl level.Level, bgColor string) {
				fgColor := tty.Contrast(bgColor)
				cfmt.RegisterStyle(string(levelColumn)+string(lvl), func(s string) string {
					return cfmt.Sprintf("{{ %1.1s }}::bg"+bgColor+"|"+fgColor, s)
				})
			}
			reg(level.Trace, theme.tinted(traceColor))
			reg(level.Debug, theme.tinted(debugColor))
			reg(level.Info, theme.tinted(infoColor))
			reg(level.Warn, theme.tinted(warnColor))
			reg(level.Error, theme.tinted(errorColor))
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
					message += cfmt.Sprintf("{{%s}}::"+highlighted+string(messageColumn)+string(event.Level), chunk)
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
			reg := func(lvl level.Level, bgColor string) {
				cf, _ := colorful.MakeColor(gamut.Tints(gamut.Hex(theme.BaseColor), 64)[60])
				messageColor := cf.Hex()
				messageBgColor := tty.Fainter(bgColor, 0.90)
				cfmt.RegisterStyle(string(messageColumn)+string(lvl), func(s string) string {
					return cfmt.Sprintf("{{%s}}::"+messageColor+"|bg"+messageBgColor, s)
				})
				cfmt.RegisterStyle(highlighted+string(messageColumn)+string(lvl), func(s string) string {
					return cfmt.Sprintf("{{%s}}::bg"+messageColor+"|"+messageBgColor, s)
				})
			}
			reg(level.Trace, theme.tinted(traceColor))
			reg(level.Debug, theme.tinted(debugColor))
			reg(level.Info, theme.tinted(infoColor))
			reg(level.Warn, theme.tinted(warnColor))
			reg(level.Error, theme.tinted(errorColor))
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
				return cfmt.Sprintf("{{%s}}::"+tty.SimilarBg(theme.tinted(errorColor)), s)
			})
			cfmt.RegisterStyle(highlighted+string(stackTraceColumn), func(s string) string {
				return cfmt.Sprintf("{{%s}}::bg"+tty.SimilarBg(theme.tinted(errorColor)), s)
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
