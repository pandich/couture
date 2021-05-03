package pretty

import (
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/gamut"
	errors2 "github.com/pkg/errors"
	"reflect"
	"strings"
	"sync"
)

// styler ...
type styler struct {
	sourceRegistryLock sync.Mutex
	sourceRegistry     map[source.Source]lipgloss.Style
	sourceColorCycle   chan lipgloss.Color
}

// newStyler ...
func newStyler() *styler {
	const sourceColorCycleLength = 50
	return &styler{
		sourceRegistryLock: sync.Mutex{},
		sourceRegistry:     map[source.Source]lipgloss.Style{},
		sourceColorCycle:   newColorCycle(sourceColorCycleLength, gamut.HappyGenerator{}),
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
		case model.Level:
			sa = append(sa, globalStyles[v].Render(string(v[0])))
		case model.Timestamp:
			sa = append(sa, globalStyles[reflect.TypeOf(v)].Render(v.Stamp()))
		case model.Caller:
			sa = append(sa, globalStyles[reflect.TypeOf(v)].Render(fmt.Sprintf(
				"%s:%s#%d",
				abbreviateClassName(v.ClassName, 32),
				v.MethodName,
				v.LineNumber,
			)))
		default:
			if style, ok := globalStyles[reflect.TypeOf(i)]; ok {
				sa = append(sa, style.Render(fmt.Sprint(i)))
			} else {
				panic(errors2.Errorf("unknown type: %+v %T", i, i))
			}
		}
	}
	return strings.Join(sa, "")
}
