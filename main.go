// Package main lauches the application.
// See README.md for more information.
package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pandich/couture/cmd"
	"github.com/pandich/couture/manual"
	"os"
)

// thin wrapper around cmd.
func main() {
	if len(os.Args) == 2 && os.Args[1] == "man" {
		program := tea.NewProgram(
			manual.NewProgram(),
			tea.WithAltScreen(),
			tea.WithANSICompressor(),
			tea.WithMouseAllMotion(),
			tea.WithMouseCellMotion(),
		)
		_, err := program.Run()
		if err != nil {
			panic(err)
		}
		return
	}

	cmd.Run()
}
