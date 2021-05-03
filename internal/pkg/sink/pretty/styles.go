package pretty

import (
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"github.com/charmbracelet/lipgloss"
	"reflect"
)

const sourceColumnWidth = 40

//nolint:gomnd
var (
	columnStyle = func(color lipgloss.Color) lipgloss.Style {
		return lipgloss.NewStyle().
			BorderLeft(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(color).
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
		model.LevelError: levelColumnStyle(errorColor),
		model.LevelWarn:  levelColumnStyle(warnColor),
		model.LevelInfo:  levelColumnStyle(infoColor),
		model.LevelDebug: levelColumnStyle(debugColor),
		model.LevelTrace: levelColumnStyle(traceColor),

		reflect.TypeOf(model.Timestamp{}):         columnStyle(timestampColor).Width(16),
		reflect.TypeOf(model.ApplicationName("")): columnStyle(applicationNameColor).Width(13),
		reflect.TypeOf(model.ThreadName("")):      columnStyle(threadNameColor).Width(19),
		reflect.TypeOf(model.Caller{}):            columnStyle(callerColor).Align(lipgloss.Right).Width(32),
		// TODO wordwrap.WrapString(line, sink.Options().Wrap())
		reflect.TypeOf(model.Message("")): columnStyle(messageColor).
			Bold(true).
			Width(72),
		reflect.TypeOf(model.StackTrace("")): columnStyle(errorColor),
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
		Width(sourceColumnWidth).
		Foreground(sourceColor).
		BorderStyle(lipgloss.NormalBorder()).
		BorderLeft(true).
		BorderLeftForeground(sourceColor).
		BorderLeftBackground(sourceColor)

	styler.sourceRegistry[src] = style
	return style
}
