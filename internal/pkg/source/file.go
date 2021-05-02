package source

import (
	"couture/internal/pkg/model"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

// https://stackoverflow.com/questions/31120987/tail-f-like-generator

func NewFile(filename string) interface{} {
	return File{filename: filename}
}

type File struct {
	io.ReadCloser
	running  bool
	filename string
	callback PushingCallback
}

func (s File) String() string {
	return path.Base(s.filename)
}

func (s File) GoString() string {
	return "📁" + s.String()
}

func (s File) Name() string {
	return s.filename
}

func (s File) Start(wg *sync.WaitGroup, callback PushingCallback) error {
	logFile, err := os.Open(s.filename)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(logFile)
	f := func() {
		defer func() { _ = logFile.Close() }()
		defer wg.Done()
		for s.running {
			var evt = model.Event{}
			if err := decoder.Decode(&evt); err != nil {
				if err != io.EOF {
					log.Println(errors.Wrap(err, s.Name()+": could not decode event"))
					return
				}
			} else {
				callback(evt)
			}
		}
	}

	s.running = true
	wg.Add(1)
	go f()
	return nil
}

func (s File) Stop() {
	s.running = false
}

func (s File) Read(b []byte) (int, error) {
	for {
		n, err := s.ReadCloser.Read(b)
		if n > 0 {
			return n, nil
		}
		if err != io.EOF {
			return 0, err
		}
		time.Sleep(10 * time.Millisecond)
	}
}
