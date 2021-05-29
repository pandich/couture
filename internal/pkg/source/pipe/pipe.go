package pipe

import (
	"bufio"
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/schema"
	"couture/internal/pkg/source"
	"io"
	"sync"
)

// Start ...
func Start(
	wg *sync.WaitGroup,
	running func() bool,
	src source.Source,
	srcChan chan source.Event,
	_ chan model.SinkEvent,
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
				srcChan <- source.Event{
					Source: src,
					Event:  scanner.Text(),
					Schema: schema.Logstash,
				}
			}
		}
	}
	wg.Add(1)
	go f()
	return nil
}
