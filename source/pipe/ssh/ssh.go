// Package ssh provides a source for reading from a remote file over SSH.
// TODO replace with https://pkg.go.dev/golang.org/x/crypto/ssh
package ssh

import (
	"github.com/pandich/couture/event"
	"github.com/pandich/couture/source"
	"github.com/pandich/couture/source/pipe"
	"github.com/melbahja/goph"
	"reflect"
	"sync"
	"time"
)

// Metadata ...
func Metadata() source.Metadata {
	return source.Metadata{
		Name:        "SSH",
		Type:        reflect.TypeOf(sshSource{}),
		CanHandle:   func(url event.SourceURL) bool { return url.Scheme == "ssh" },
		Creator:     source.Single(newSource),
		ExampleURLs: []string{"ssh://user:passphrase@host:port/<path>"},
	}
}

type sshSource struct {
	source.BaseSource
	ssh      *goph.Client
	filename string
}

func newSource(_ *time.Time, sourceURL event.SourceURL) (*source.Source, error) {
	client, err := sshURL(sourceURL).getClient()
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
	snkChan chan event.SinkEvent,
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
