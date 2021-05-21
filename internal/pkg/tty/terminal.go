package tty

import (
	"github.com/mattn/go-isatty"
	"github.com/muesli/reflow/wordwrap"
	"github.com/olekukonko/ts"
	"os"
)

// IsTTY ...
func IsTTY() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}

// TerminalWidth ...
func TerminalWidth() int {
	var terminalWidth = 0
	if size, err := ts.GetSize(); err == nil {
		terminalWidth = size.Col()
	}
	return terminalWidth
}

// Wrap ...
func Wrap(s string, width int) string {
	wrapper := wordwrap.NewWriter(width)
	wrapper.Breakpoints = []rune(" \t")
	wrapper.KeepNewlines = true
	if _, err := wrapper.Write([]byte(s)); err != nil {
		return s
	}
	return wrapper.String()
}
