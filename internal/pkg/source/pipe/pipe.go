package pipe

import (
	"bufio"
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
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
				line := scanner.Text()
				err := json.Unmarshal([]byte(line), &event)
				if err != nil {
					out <- source.Event{Source: src, Event: model.Event{
						Timestamp: model.Timestamp{},
						Level:     level.Info,
						Message:   model.Message(line),
					}}
				} else {
					out <- source.Event{Source: src, Event: event}
				}
			}
		}
	}

	wg.Add(1)
	go f()
	return nil
}
