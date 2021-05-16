package pretty

import (
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"couture/pkg/model/level"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
	"github.com/olekukonko/ts"
	errors2 "github.com/pkg/errors"
	"reflect"
	"strings"
	"sync"
	"time"
)

// styler ...
type styler struct {
	sourceRegistryLock sync.Mutex
	sourceRegistry     map[source.Source]lipgloss.Style
	sourceColorCycle   chan lipgloss.TerminalColor
}

// newStyler ...
func newStyler() *styler {
	return &styler{
		sourceRegistryLock: sync.Mutex{},
		sourceRegistry:     map[source.Source]lipgloss.Style{},
		sourceColorCycle:   pastels(),
	}
}

// render ...
func (styler *styler) render(ia ...interface{}) string {
	var sa []string
	for _, i := range ia {
		switch v := i.(type) {
		case string:
			sa = append(sa, v)
		case source.Source:
			sa = append(sa, styler.sourceStyle(v).Render(v.URL().ShortForm()))
		case level.Level:
			sa = append(sa, globalStyles[v].Render(string(v[0])))
		case model.Stamp:
			sa = append(sa, globalStyles[reflect.TypeOf(v)].Render(string(v)))
		case punctuation:
			sa = append(sa, globalStyles[v].Render(string(v)))
		default:
			if style, ok := globalStyles[reflect.TypeOf(i)]; ok {
				sa = append(sa, style.Render(fmt.Sprint(i)))
			} else {
				panic(errors2.Errorf("unknown type: %+v %T\n", i, i))
			}
		}
	}
	return strings.Join(sa, "")
}

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
		level.Error: levelColumnStyle(errorColor),
		level.Warn:  levelColumnStyle(warnColor),
		level.Info:  levelColumnStyle(infoColor),
		level.Debug: levelColumnStyle(debugColor),
		level.Trace: levelColumnStyle(traceColor),

		reflect.TypeOf(model.Stamp("")):           columnStyle(timestampColor).Width(len(time.Stamp)),
		reflect.TypeOf(model.ApplicationName("")): columnStyle(applicationNameColor).Width(applicationNameWidth),
		reflect.TypeOf(model.ThreadName("")):      columnStyle(threadNameColor).Width(threadNameWidth),

		reflect.TypeOf(caller("")):           lipgloss.NewStyle().Bold(true).Align(lipgloss.Right).Width(callerWidth),
		reflect.TypeOf(model.ClassName("")):  baseStyle(classNameColor).Bold(true),
		methodNameDelimiter:                  baseStyle(methodNameColor).MaxWidth(1).Bold(false),
		reflect.TypeOf(model.MethodName("")): baseStyle(methodNameColor).Bold(true),
		lineNumberDelimiter:                  baseStyle(lineNumberColor).MaxWidth(1).Bold(false),
		reflect.TypeOf(model.LineNumber(0)):  baseStyle(lineNumberColor).Bold(true).Width(4),

		// TODO this approach currently messes up the line breaks in the message and exception
		//		rework to ensure
		reflect.TypeOf(model.Message("")):              messageStyle(messageColor).Width(messageWidth),
		reflect.TypeOf(model.UnhighlightedMessage("")): baseStyle(messageColor).MarginLeft(1),
		reflect.TypeOf(model.HighlightedMessage("")): baseStyle(messageColor).MarginLeft(1).
			Reverse(true).PaddingLeft(1).PaddingRight(1),
		reflect.TypeOf(model.StackTrace("")):              messageStyle(errorColor).Width(messageWidth),
		reflect.TypeOf(model.UnhighlightedStackTrace("")): baseStyle(errorColor).MarginLeft(1),
		reflect.TypeOf(model.HighlightedStackTrace("")): baseStyle(errorColor).MarginLeft(1).
			Reverse(true).PaddingLeft(1).PaddingRight(1),
	}
)

func (styler *styler) sourceStyle(src source.Source) lipgloss.Style {
	const sourceColumnWidth = 40

	styler.sourceRegistryLock.Lock()
	defer styler.sourceRegistryLock.Unlock()

	if style, ok := styler.sourceRegistry[src]; ok {
		return style
	}
	sourceColor := <-styler.sourceColorCycle
	style := lipgloss.NewStyle().
		MarginLeft(1).
		PaddingLeft(1).
		MaxHeight(1).
		Width(sourceColumnWidth).
		MaxWidth(sourceColumnWidth).
		Foreground(sourceColor).
		MarginBackground(sourceColor)

	styler.sourceRegistry[src] = style
	return style
}
