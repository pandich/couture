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
				// TODO how to deal with un-parsable strings -- this applies to billing info etc.
				//		this will impact CloudWatch, too. Not all of its events are structured.
				if err := json.Unmarshal([]byte(scanner.Text()), &event); err != nil {
					callback(event)
				}
			}
		}
	}

	wg.Add(1)
	go f()
	return nil
}
