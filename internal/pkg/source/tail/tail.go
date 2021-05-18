package tail

import (
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/nxadm/tail"
	"os"
	"reflect"
	"sync"
)

// TODO file lookback lines
// TODO file retry on not found like tail -F

// Metadata ...
func Metadata() source.Metadata {
	return source.Metadata{
		Type:      reflect.TypeOf(fileSource{}),
		CanHandle: func(url model.SourceURL) bool { return url.Scheme == "file" },
		Creator: func(sourceURL model.SourceURL) (*interface{}, error) {
			src, err := newSource(sourceURL)
			if err != nil {
				return nil, err
			}
			var i interface{} = src
			return &i, nil
		},
		ExampleURLs: []string{"file://<path>"},
	}
}

// fileSource ...
type fileSource struct {
	*source.Pushing
	tailer *tail.Tail
	file   *os.File
}

// newSource ...
func newSource(sourceURL model.SourceURL) (*source.Pushable, error) {
	sourceURL.Normalize()
	file, err := os.Open(sourceURL.Path)
	if err != nil {
		return nil, err
	}

	tailer, err := tail.TailFile(file.Name(), tail.Config{Follow: true})
	if err != nil {
		return nil, err
	}

	var src source.Pushable = fileSource{
		Pushing: source.New(sourceURL),
		tailer:  tailer,
		file:    file,
	}
	return &src, nil
}

// Start ...
func (src fileSource) Start(wg *sync.WaitGroup, running func() bool, callback func(event model.Event)) error {
	pusher := func() {
		defer func() { _ = src.file.Close() }()
		defer wg.Done()
		for running() {
			for line := range src.tailer.Lines {
				var event model.Event
				err := json.Unmarshal([]byte(line.Text), &event)
				if err != nil {
					// TODO how to deal with un-parsable strings -- this applies to billing info etc.
					//		this will impact CloudWatch, too. Not all of its events are structured.
					fmt.Println(err)
				}
				if err == nil {
					callback(event)
				}
			}
		}
	}
	wg.Add(1)
	go pusher()
	return nil
}
