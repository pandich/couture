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
	out chan source.Event,
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
				// TODO we should be pushing raw strings to the out and do event JSON parsing centrally
				_ = json.Unmarshal([]byte(scanner.Text()), &event)
				out <- source.Event{Source: src, Event: event}
			}
		}
	}

	wg.Add(1)
	go f()
	return nil
}
