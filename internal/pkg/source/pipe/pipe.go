package pipe

import (
	"bufio"
	"couture/internal/pkg/model"
	"encoding/json"
	"io"
	"sync"
)

// Start ...
func Start(
	wg *sync.WaitGroup,
	running func() bool,
	callback func(event model.Event),
	closer func(),
	out io.Reader,
) error {
	scanner := bufio.NewScanner(out)
	scanner.Split(bufio.ScanLines)
	f := func() {
		defer wg.Done()
		defer closer()
		for running() {
			for scanner.Scan() {
				var event model.Event
				// TODO we should be pushing raw strings to the callback and do event JSON parsing centrally
				_ = json.Unmarshal([]byte(scanner.Text()), &event)
				callback(event)
			}
		}
	}

	wg.Add(1)
	go f()
	return nil
}
