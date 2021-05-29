package ssh

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/pipe"
	"github.com/melbahja/goph"
	"reflect"
	"sync"
)

// Metadata ...
func Metadata() source.Metadata {
	return source.Metadata{
		Name:        "SSH",
		Type:        reflect.TypeOf(sshSource{}),
		CanHandle:   func(url model.SourceURL) bool { return url.Scheme == "ssh" },
		Creator:     newSource,
		ExampleURLs: []string{"ssh://user:passphrase@host:port/<path>"},
	}
}

type sshSource struct {
	source.BaseSource
	ssh      *goph.Client
	filename string
}

func newSource(sourceURL model.SourceURL) (*source.Source, error) {
	client, err := getClient(sourceURL)
	if err != nil {
		return nil, err
	}
	src := sshSource{
		BaseSource: source.New(' ', sourceURL),
		ssh:        client,
		filename:   sourceURL.Path,
	}
	var p source.Source = src
	return &p, nil
}

// Start ...
func (src sshSource) Start(
	wg *sync.WaitGroup,
	running func() bool,
	srcChan chan source.Event,
	snkChan chan model.SinkEvent,
	errChan chan source.Error,
) error {
	// create the command
	cmd, err := src.ssh.Command("tail", "-F", src.filename)
	if err != nil {
		return err
	}

	// capture its output
	in, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// and start it
	if err := cmd.Start(); err != nil {
		return err
	}
	return pipe.Start(wg, running, src, srcChan, snkChan, errChan, func() { _ = cmd.Close() }, in)
}
