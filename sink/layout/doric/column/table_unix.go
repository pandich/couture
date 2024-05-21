//go:build !windows

package column

import (
	"os"
	"os/signal"
	"syscall"
)

func (table *Table) autoUpdateColumnWidths() {
	resize := make(chan os.Signal, 1)
	signal.Notify(resize, os.Interrupt, syscall.SIGWINCH)
	go func() {
		defer close(resize)
		for {
			<-resize
			table.updateColumnWidths()
		}
	}()
}
