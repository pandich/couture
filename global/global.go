package global

import "github.com/charmbracelet/lipgloss"

type Resource struct {
	Kind, Name string
}

var DiscoveryBus = make(chan Resource, 1000)

func init() {
	go func() {
		kindStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a0a0"))
		nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#336E4E"))
		for {
			select {
			case msg := <-DiscoveryBus:
				println(lipgloss.JoinHorizontal(lipgloss.Left, kindStyle.Render(msg.Kind), nameStyle.Render(msg.Name)))
			}
		}
	}()
}
