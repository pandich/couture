// Package local provides a source that tails a local file.
// TODO replace with github.com/nxadm/tail
package local

import (
	"bufio"
	"github.com/gagglepanda/couture/event"
	"github.com/gagglepanda/couture/source"
	"github.com/gagglepanda/couture/source/pipe"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"
)

// Metadata ...
func Metadata() source.Metadata {
	return source.Metadata{
		Name:        "Local File",
		Type:        reflect.TypeOf(fileSource{}),
		CanHandle:   func(url event.SourceURL) bool { return url.Scheme == "file" },
		Creator:     newSource,
		ExampleURLs: []string{"file://<path>"},
	}
}

type fileSource struct {
	source.BaseSource
	filename string
}

func newSource(_ *time.Time, sourceURL event.SourceURL) (*source.Source, error) {
	sourceURL.Normalize()
	var src source.Source = fileSource{
		BaseSource: source.New('⫽', sourceURL),
		filename:   sourceURL.Path,
	}
	return &src, nil
}

// Start ...
func (src fileSource) Start(
	wg *sync.WaitGroup,
	running func() bool,
	srcChan chan source.Event,
	snkChan chan event.SinkEvent,
	errChan chan source.Error,
) error {
	// get the safe path to the file
	path, err := filepath.Abs(src.filename)
	if err != nil {
		return err
	}

	out := make(chan string)
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		writer := bufio.NewWriter(pw)
		for line := range out {
			_, err := writer.WriteString(line)
			if err != nil {
				errChan <- source.Error{Error: err}
			} else {
				writer.Flush()
			}
		}
	}()

	go func() {
		var err error
		if src.filename == "/-" {
			err = tailStdin(out)
		} else {
			err = tailFile(path, out)
		}
		if err != nil {
			errChan <- source.Error{Error: err}
			close(out)
		}
	}()

	in := bufio.NewReader(pr)
	return pipe.Start(wg, running, src, srcChan, snkChan, errChan, func() {}, in)
}

func tailStdin(out chan<- string) error {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		out <- line
	}
	close(out)
	return nil
}

func tailFile(filePath string, out chan<- string) error {
	for {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}

		// TODO lookback
		_, err = file.Seek(0, io.SeekEnd)
		if err != nil {
			return err
		}

		reader := bufio.NewReader(file)

		for {
			var line string
			line, err = reader.ReadString('\n')
			if err == nil {
				out <- line
				continue
			}

			if err != io.EOF {
				return err
			}

			if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
				file.Close()
				break
			}

			// wait for the file to be written to
			time.Sleep(100 * time.Millisecond)
		}

		// wait for the file to exist
		time.Sleep(100 * time.Millisecond)
	}
}
