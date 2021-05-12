package pretty

import (
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
	"github.com/olekukonko/ts"
	"reflect"
	"time"
)

const sourceColumnWidth = 40

type punctuation string

const (
	methodNameDelimiter punctuation = ":"
	lineNumberDelimiter punctuation = "#"
	newLine             punctuation = "\n"
)

const (
	applicationNameWidth = 20
	threadNameWidth      = 20
	callerWidth          = 48
)

func terminalWidth() int {
	const defaultTerminalWidth = 72
	if size, err := ts.GetSize(); err == nil {
		return size.Col()
	}
	return defaultTerminalWidth
}

//nolint:gomnd
var (
	messageWidth = terminalWidth()
	baseStyle    = func(color lipgloss.TerminalColor) lipgloss.Style {
		return lipgloss.NewStyle().Foreground(color)
	}
	columnStyle = func(color lipgloss.TerminalColor) lipgloss.Style {
		return baseStyle(color).MaxHeight(1).PaddingLeft(1)
	}
	levelColumnStyle = func(color lipgloss.TerminalColor) lipgloss.Style {
		return baseStyle(levelForegroundColor).
			MarginLeft(1).
			PaddingLeft(1).PaddingRight(1).
			Background(color)
	}

	messageStyle = func(color lipgloss.TerminalColor) lipgloss.Style {
		h, s, l := termenv.ConvertToRGB(colorProfile.Color(fmt.Sprint(color))).Hsl()
		borderColor := lipgloss.Color(colorful.Hsl(h, s*0.7, l*0.3).Hex())
		return columnStyle(color).
			BorderLeft(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(borderColor).
			UnsetMaxHeight().
			Bold(true).
			Width(messageWidth).
			MaxWidth(messageWidth)
	}

	globalStyles = map[interface{}]lipgloss.Style{
		model.LevelError: levelColumnStyle(errorColor),
		model.LevelWarn:  levelColumnStyle(warnColor),
		model.LevelInfo:  levelColumnStyle(infoColor),
		model.LevelDebug: levelColumnStyle(debugColor),
		model.LevelTrace: levelColumnStyle(traceColor),

		methodNameDelimiter: baseStyle(methodNameColor).MaxWidth(1).Bold(false),
		lineNumberDelimiter: baseStyle(lineNumberColor).MaxWidth(1).Bold(false),

		reflect.TypeOf(caller("")):                lipgloss.NewStyle().Bold(true).Align(lipgloss.Right).Width(callerWidth),
		reflect.TypeOf(model.Stamp("")):           columnStyle(timestampColor).Width(len(time.Stamp)),
		reflect.TypeOf(model.ApplicationName("")): columnStyle(applicationNameColor).Width(applicationNameWidth),
		reflect.TypeOf(model.ThreadName("")):      columnStyle(threadNameColor).Width(threadNameWidth),
		reflect.TypeOf(model.ClassName("")):       baseStyle(classNameColor).Bold(true),
		reflect.TypeOf(model.MethodName("")):      baseStyle(methodNameColor).Bold(true),
		reflect.TypeOf(model.LineNumber(0)):       baseStyle(lineNumberColor).Bold(true).Width(4),
		reflect.TypeOf(model.Message("")):         messageStyle(messageColor).Width(messageWidth),
		reflect.TypeOf(model.Unhighlighted("")):   baseStyle(messageColor).MarginLeft(1),
		reflect.TypeOf(model.Highlighted("")):     baseStyle(messageColor).MarginLeft(1).Reverse(true).PaddingLeft(1).PaddingRight(1),
		reflect.TypeOf(model.StackTrace("")):      messageStyle(errorColor),
	}
)

func (styler *styler) sourceStyle(src source.Source) lipgloss.Style {
	styler.sourceRegistryLock.Lock()
	defer styler.sourceRegistryLock.Unlock()

	if style, ok := styler.sourceRegistry[src]; ok {
		return style
	}

	sourceColor := <-styler.sourceColorCycle
	style := lipgloss.NewStyle().
		PaddingLeft(1).
		MaxHeight(1).
		Width(sourceColumnWidth).
		MaxWidth(sourceColumnWidth + 1 /* padding */).
		Foreground(sourceColor).
		BorderStyle(lipgloss.NormalBorder()).
		BorderLeft(true).
		BorderLeftForeground(sourceColor).
		BorderLeftBackground(sourceColor)

	styler.sourceRegistry[src] = style
	return style
}
