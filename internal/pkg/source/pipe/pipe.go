package pipe

import (
	"bufio"
	"couture/internal/pkg/model"
	"couture/internal/pkg/source"
	"encoding/json"
	"io"
	"sync"
)

// Start ...
func Start(
	wg *sync.WaitGroup,
	running func() bool,
	src source.Source,
	srcChan chan source.Event,
	errChan chan source.Error,
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
				var event model.Event
				line := scanner.Text()
				err := json.Unmarshal([]byte(line), &event)
				if err != nil {
					errChan <- source.Error{Source: src, Error: err}
				} else {
					srcChan <- source.Event{Source: src, Event: event}
				}
			}
		}
	}

	wg.Add(1)
	go f()
	return nil
}
