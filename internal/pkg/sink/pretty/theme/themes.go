package theme

import (
	"github.com/charmbracelet/lipgloss"
)

// Default ...
var Default = prince

// Registry is the registry of theme names to their structs.
var Registry = map[string]Theme{}

// Names ...
var Names []string

func register(name string, theme Theme) {
	Names = append(Names, name)
	Registry[name] = theme
}

func style(fg string, bg string) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(fg)).
		Background(lipgloss.Color(bg))
}
