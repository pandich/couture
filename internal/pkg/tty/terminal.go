package tty

import (
	"bufio"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
	"github.com/olekukonko/ts"
	"io"
	"os"
)

// IsTTY ...
func IsTTY() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}

// IsDarkMode ...
func IsDarkMode() bool {
	return termenv.HasDarkBackground()
}

// TerminalWidth ...
func TerminalWidth() int {
	var terminalWidth = 0
	if size, err := ts.GetSize(); err == nil {
		terminalWidth = size.Col()
	}
	return terminalWidth
}

// NewTTYWriter ...
func NewTTYWriter(target io.Writer) chan string {
	delegate := make(chan string)
	go func() {
		defer close(delegate)
		writer := bufio.NewWriter(target)
		for {
			message := <-delegate
			_, err := writer.WriteString(message + "\n")
			if err != nil {
				panic(err)
			}
			err = writer.Flush()
			if err != nil {
				panic(err)
			}
		}
	}()
	return delegate
}
