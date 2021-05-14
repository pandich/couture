package tail

import (
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/pushing"
	"couture/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/nxadm/tail"
	"os"
	"reflect"
	"sync"
)

// Metadata ...
func Metadata() source.Metadata {
	return source.Metadata{
		Type:        reflect.TypeOf(fileSource{}),
		CanHandle:   func(url model.SourceURL) bool { return url.Scheme == "file" },
		Creator:     create,
		ExampleURLs: []string{"file://<path>"},
	}
}

// fileSource ...
type fileSource struct {
	source.Base
	tailer *tail.Tail
	file   *os.File
}

// create CloudFormation source casted to an *interface{}.
func create(sourceURL model.SourceURL) (*interface{}, error) {
	src, err := newSource(sourceURL)
	if err != nil {
		return nil, err
	}
	var i interface{} = src
	return &i, nil
}

// newSource ...
func newSource(sourceURL model.SourceURL) (*pushing.Source, error) {
	sourceURL.Normalize()
	file, err := os.Open(sourceURL.Path)
	if err != nil {
		return nil, err
	}

	tailer, err := tail.TailFile(file.Name(), tail.Config{Follow: true})
	if err != nil {
		return nil, err
	}

	var src pushing.Source = fileSource{
		Base:   source.New(sourceURL),
		tailer: tailer,
		file:   file,
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
					// TODO do something?
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
