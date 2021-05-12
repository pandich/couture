package pretty

import (
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
	"reflect"
)

const sourceColumnWidth = 40

type punctuation string

const (
	methodNameDelimiter punctuation = ":"
	lineNumberDelimiter punctuation = "#"
)

const (
	classNameColumnWidth  = 40
	methodNameColumnWidth = 30
)

//nolint:gomnd
var (
	columnStyle = func(color lipgloss.Color) lipgloss.Style {
		h, s, l := termenv.ConvertToRGB(colorProfile.Color(fmt.Sprint(color))).Hsl()
		borderColor := colorful.Hsl(h, s, l*0.7)
		return lipgloss.NewStyle().
			MaxHeight(1).
			BorderLeft(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(borderColor.Hex())).
			PaddingLeft(1).
			Foreground(color)
	}

	levelColumnStyle = func(color lipgloss.Color) lipgloss.Style {
		return lipgloss.NewStyle().
			MarginLeft(1).
			PaddingLeft(1).
			PaddingRight(1).
			Background(color).
			Foreground(levelForegroundColor)
	}

	globalStyles = map[interface{}]lipgloss.Style{
		methodNameDelimiter: lipgloss.NewStyle().Foreground(methodNameColor).MaxWidth(1).Bold(false),
		lineNumberDelimiter: lipgloss.NewStyle().Foreground(lineNumberColor).MaxWidth(1).Bold(false),

		model.LevelError: levelColumnStyle(errorColor),
		model.LevelWarn:  levelColumnStyle(warnColor),
		model.LevelInfo:  levelColumnStyle(infoColor),
		model.LevelDebug: levelColumnStyle(debugColor),
		model.LevelTrace: levelColumnStyle(traceColor),

		reflect.TypeOf(model.Timestamp{}):         columnStyle(timestampColor).Width(16),
		reflect.TypeOf(model.ApplicationName("")): columnStyle(applicationNameColor).Width(13).MaxWidth(13),
		reflect.TypeOf(model.ThreadName("")):      columnStyle(threadNameColor).Width(20).MaxWidth(20),
		reflect.TypeOf(model.ClassName("")):       columnStyle(classNameColor).Bold(true).Align(lipgloss.Right).Width(classNameColumnWidth + 2),
		reflect.TypeOf(model.MethodName("")):      lipgloss.NewStyle().Foreground(methodNameColor).Bold(true).Width(methodNameColumnWidth),
		reflect.TypeOf(model.LineNumber(0)):       lipgloss.NewStyle().Foreground(lineNumberColor).Bold(true).Width(3).MaxWidth(4),
		reflect.TypeOf(model.Message("")):         columnStyle(messageColor).UnsetMaxHeight().Bold(true).Width(132),
		reflect.TypeOf(model.StackTrace("")):      columnStyle(errorColor).UnsetMaxHeight(),
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
		MaxWidth(sourceColumnWidth + 1 /* padding */ + 1 /* bordr */).
		Width(sourceColumnWidth + +1 /* padding */ + 1 /* border */).
		PaddingLeft(1).
		MaxHeight(1).
		Foreground(sourceColor).
		BorderStyle(lipgloss.NormalBorder()).
		BorderLeft(true).
		BorderLeftForeground(sourceColor).
		BorderLeftBackground(sourceColor)

	styler.sourceRegistry[src] = style
	return style
}
