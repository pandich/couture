package manual

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/olekukonko/ts"
	"slices"
	"strings"
	"sync"
)

const (
	modeList mode = iota
	modePage
	manualWidth = 72 - 1*2
)

var (
	listStyle            = lipgloss.NewStyle().Width(manualWidth)
	pageNotSelectedStyle = lipgloss.NewStyle().Copy().Foreground(lipgloss.Color("#854e0b"))
	pageSelectedStyle    = lipgloss.NewStyle().Copy().Foreground(lipgloss.Color("#FFD230"))
	titleStyle           = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder(), true).
				BorderForeground(lipgloss.Color("#3967FF")).
				Padding(0, 1).
				Width(manualWidth - 2)
)

func NewProgram() tea.Model {
	man := New()

	lines := make([]string, len(man.Pages()))
	titles := make([]string, len(man.Pages()))

	for i := range man.Pages() {
		title := man.Pages()[i]
		titles[i] = string(title)
	}
	slices.SortFunc(titles, strings.Compare)

	for i, title := range titles {
		content, err := man.Page(Page(title))
		if err != nil {
			panic(err)
		}
		lineCount := strings.Count(content, "\n")
		lines[i] = lipgloss.JoinHorizontal(
			lipgloss.Center,
			titles[i],
			" ",
			strings.Repeat(".", manualWidth-len([]rune(title))-len("000 lines")-5),
			" ",
			fmt.Sprintf("%4d lines", lineCount),
		)
	}

	w, h := terminalWidth()
	if w > 72 {
		w = 72
	}
	listView := viewport.New(w, h)
	pageView := viewport.New(w, h)
	hlp := help.New()
	hlp.ShowAll = true
	return &viewer{
		manual: man,
		titles: titles,

		listView:  listView,
		listKeys:  listKeys,
		listLines: lines,

		pageView: pageView,
		pageKeys: pageKeys,
		help:     hlp,
	}
}

func manualPageCmd(page Page) tea.Cmd {
	return func() tea.Msg {
		return manualPageMsg{page}
	}
}

type (
	mode int

	manualPageMsg struct {
		page Page
	}

	viewer struct {
		manual             Manual
		pageView, listView viewport.Model
		pageKeys, listKeys help.KeyMap
		pageIndex          int
		titles, listLines  []string
		mode               mode
		help               help.Model
		initOnce           sync.Once
	}

	modePageKeys struct {
		Help     key.Binding
		Close    key.Binding
		Home     key.Binding
		End      key.Binding
		Up       key.Binding
		Down     key.Binding
		PageUp   key.Binding
		PageDown key.Binding
	}

	modeListKeys struct {
		Help         key.Binding
		Quit         key.Binding
		Up           key.Binding
		Down         key.Binding
		Home         key.Binding
		End          key.Binding
		PageUp       key.Binding
		PageDown     key.Binding
		OpenSelected key.Binding
	}
)

func (p *viewer) Init() (cmd tea.Cmd) {
	p.initOnce.Do(
		func() {
			cmd = tea.Batch(
				p.listView.Init(),
				p.pageView.Init(),
			)
		},
	)

	return
}

func (p *viewer) View() string {
	view := ""
	switch p.mode {
	case modePage:
		view = p.viewPage()
	case modeList:
		view = p.viewList()
	default:
		view = ""
	}

	return lipgloss.JoinVertical(lipgloss.Center, p.help.View(p.keys()), view)
}

func (p *viewer) keys() help.KeyMap {
	switch p.mode {
	case modePage:
		return p.pageKeys
	case modeList:
		return p.listKeys
	default:
		return nil
	}
}

func (k modePageKeys) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k modePageKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down}, {k.Home, k.End}, {k.PageUp, k.PageDown},
	}
}

func (k modeListKeys) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k modeListKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.OpenSelected},
		{k.Up, k.Down},
		{k.Home, k.End},
		{k.PageUp, k.PageDown},
	}
}

var (
	Help     = newK("ctrl+k", "^k", "help")
	Quit     = newK("ctrl+b", "esc", "quit")
	Open     = newK("enter", "⏎", "open")
	Close    = newK("esc", "esc", "list")
	Up       = newK("up", "↑", "up")
	Down     = newK("down", "↓", "down")
	PageUp   = newK("pgup", "⇞", "pgup")
	PageDown = newK("pgdn", "⇟", "pgdn")
	Home     = newK("home", "⇱", "home")
	End      = newK("end", "⇲", "end")

	pageKeys = modePageKeys{
		Close:    Close,
		Home:     Home,
		End:      End,
		Up:       Up,
		Down:     Down,
		PageUp:   PageUp,
		PageDown: PageDown,
		Help:     Help,
	}

	listKeys = modeListKeys{
		Quit:         Quit,
		Help:         Help,
		Up:           Up,
		Down:         Down,
		Home:         Home,
		End:          End,
		PageUp:       PageUp,
		PageDown:     PageDown,
		OpenSelected: Open,
	}
)

func (p *viewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case manualPageMsg:
		p.pageIndex = slices.Index(p.titles, string(m.page))
		cnt, err := p.manual.Page(m.page)
		if err != nil {
			p.pageView.SetContent(err.Error())
			return p, nil
		}
		p.mode = modePage
		p.pageView.SetContent(cnt)
		p.pageView.GotoTop()
		return p, nil
	}
	switch p.mode {
	case modePage:
		return p, p.updatePage(msg)
	case modeList:
		return p, p.updateList(msg)
	default:
		return p, nil
	}

}

func (p *viewer) updatePage(msg tea.Msg) tea.Cmd {
	switch t := msg.(type) {
	case tea.KeyMsg:
		switch t.String() {
		case "esc":
			p.mode = modeList
		case "pgup":
			p.pageView.ViewUp()
			return nil
		case "pgdn":
			p.pageView.ViewDown()
			return nil
		case "down":
			p.pageView.LineDown(1)
			return nil
		case "up":
			p.pageView.LineUp(1)
		case "home":
			p.pageView.GotoTop()
			return nil
		case "end":
			p.pageView.GotoBottom()
		default:
			return nil
		}
	}

	var cmd tea.Cmd
	p.pageView, cmd = p.pageView.Update(msg)

	return cmd
}

func (p *viewer) updateList(msg tea.Msg) tea.Cmd {
	switch t := msg.(type) {
	case tea.KeyMsg:
		switch t.String() {
		case "up":
			if p.pageIndex > 0 {
				p.pageIndex--
			}
			return nil
		case "down":
			if p.pageIndex < len(p.manual.Pages())-1 {
				p.pageIndex++
			}
			return nil
		case "home":
			p.pageIndex = 0
			return nil
		case "end":
			p.pageIndex = len(p.manual.Pages()) - 1
			return nil
		case "enter":
			pg := p.manual.Pages()[p.pageIndex]
			return manualPageCmd(pg)
		case "esc":
			return tea.Quit
		}
	}

	var cmd tea.Cmd
	p.listView, cmd = p.listView.Update(msg)
	return cmd
}

func (p *viewer) viewList() string {
	lines := make([]string, len(p.listLines))
	for i, title := range p.listLines {
		st := pageNotSelectedStyle
		if i == p.pageIndex {
			st = pageSelectedStyle
		}

		lines[i] = st.Render(title)
	}

	p.listView.SetContent(lipgloss.JoinVertical(lipgloss.Left, lines...))
	return listStyle.Render(p.listView.View())
}

func (p *viewer) viewPage() string {
	header := titleStyle.Render(p.titles[p.pageIndex])
	return lipgloss.JoinVertical(lipgloss.Left, header, p.pageView.View())
}

func (p *viewer) Background() bool { return false }

func newK(k, name, description string) key.Binding {
	return key.NewBinding(key.WithKeys(k), key.WithHelp(name, description))
}
func terminalWidth() (int, int) {
	var tw = 0
	var th = 0
	if size, err := ts.GetSize(); err == nil {
		tw = size.Col()
		th = size.Row()
	}
	return tw, th
}
