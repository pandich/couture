package local

import (
	"github.com/pandich/couture/model"
	"github.com/pandich/couture/source"
	"github.com/pandich/couture/source/pipe"
	"os/exec"
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
		CanHandle:   func(url model.SourceURL) bool { return url.Scheme == "file" },
		Creator:     newSource,
		ExampleURLs: []string{"file://<path>"},
	}
}

type fileSource struct {
	source.BaseSource
	filename string
}

func newSource(_ *time.Time, sourceURL model.SourceURL) (*source.Source, error) {
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
	snkChan chan model.SinkEvent,
	errChan chan source.Error,
) error {
	// get the safe path to the file
	path, err := filepath.Abs(src.filename)
	if err != nil {
		return err
	}

	// create the command
	tail := exec.Command("tail", "-F", path)

	// capture its output
	in, err := tail.StdoutPipe()
	if err != nil {
		return err
	}

	// and start it
	if err = tail.Start(); err != nil {
		return err
	}
	return pipe.Start(wg, running, src, srcChan, snkChan, errChan, func() {}, in)
}
