package tty

import (
	"github.com/mattn/go-isatty"
	"github.com/muesli/reflow/wordwrap"
	"github.com/olekukonko/ts"
	"os"
)

// ResetSequence ...
const ResetSequence = "\x1b[2m"

// ClearScreenSequence ...
const ClearScreenSequence = "\x1b[2J"

// HomeCursorSequence ...
const HomeCursorSequence = "\x1b[0;0H"

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
func Wrap(s string, width uint) string {
	if width <= 0 {
		return s
	}
	wrapper := wordwrap.NewWriter(int(width))
	wrapper.Breakpoints = []rune(" ")
	wrapper.KeepNewlines = true
	if _, err := wrapper.Write([]byte(s)); err != nil {
		return s
	}
	return wrapper.String()
}
