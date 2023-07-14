package pipe

import (
	"bufio"
	"github.com/gagglepanda/couture/event"
	"github.com/gagglepanda/couture/source"
	"io"
	"sync"
)

// Start ...
func Start(
	wg *sync.WaitGroup,
	running func() bool,
	src source.Source,
	srcChan chan source.Event,
	_ chan event.SinkEvent,
	_ chan source.Error,
	closer func(),
	in io.Reader,
) error {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)
	f := func() {
		defer wg.Done()
		defer closer()
		for running() {
			for scanner.Scan() {
				srcChan <- source.Event{Source: src, Event: scanner.Text()}
			}
		}
	}
	wg.Add(1)
	go f()
	return nil
}
