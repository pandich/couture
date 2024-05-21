package manual

import (
	_ "embed"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

//go:embed theme.json
var theme []byte

func init() {
	lipgloss.SetColorProfile(termenv.TrueColor)
}
func mustNewTermRenderer() *glamour.TermRenderer {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithPreservedNewLines(),
		glamour.WithEmoji(),
		glamour.WithStylesFromJSONBytes(theme),
	)
	if err != nil {
		panic(err)
	}
	return renderer
}
