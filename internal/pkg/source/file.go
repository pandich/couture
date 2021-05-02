package source

import (
	"couture/internal/pkg/model"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"reflect"
	"sync"
	"time"
)

// https://stackoverflow.com/questions/31120987/tail-f-like-generator

//init registers the type with the typeRegistry.
func init() {
	typeRegistry[reflect.TypeOf(File{})] = func(srcUrl url.URL) interface{} {
		return File{baseSource: baseSource{srcUrl: srcUrl}}
	}
	registry = append(registry, File{})
}

type File struct {
	io.ReadCloser
	baseSource
	running  bool
	callback PushingCallback
}

func (source File) CanHandle(url url.URL) bool {
	return url.Scheme == "file"
}

func (source File) String() string {
	return path.Base(source.srcUrl.Path)
}

func (source File) GoString() string {
	return "❯ " + source.String()
}

func (source File) Start(wg *sync.WaitGroup, callback PushingCallback) error {
	logFile, err := os.Open(source.srcUrl.Path)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(logFile)
	f := func() {
		defer func() { _ = logFile.Close() }()
		defer wg.Done()
		for source.running {
			var event = model.Event{}
			if err := decoder.Decode(&event); err != nil {
				if err != io.EOF {
					log.Println(errors.Wrap(err, source.String()+": could not decode event"))
					return
				}
			} else {
				callback(event)
			}
		}
	}

	source.running = true
	wg.Add(1)
	go f()
	return nil
}

func (source File) Stop() {
	source.running = false
}

func (source File) Read(b []byte) (int, error) {
	for {
		n, err := source.ReadCloser.Read(b)
		if n > 0 {
			return n, nil
		}
		if err != io.EOF {
			return 0, err
		}
		time.Sleep(10 * time.Millisecond)
	}
}
