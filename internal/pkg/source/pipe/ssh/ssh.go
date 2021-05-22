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
		Name:      "SSH",
		Type:      reflect.TypeOf(sshSource{}),
		CanHandle: func(url model.SourceURL) bool { return url.Scheme == "ssh" },
		Creator: func(sourceURL model.SourceURL) (*interface{}, error) {
			src, err := newSource(sourceURL)
			if err != nil {
				return nil, err
			}
			var i interface{} = src
			return &i, nil
		},
		ExampleURLs: []string{"ssh://user:passphrase@host:port/<path>"},
	}
}

// sshSource ...
type sshSource struct {
	*source.Pushing
	ssh      *goph.Client
	filename string
}

// newSource ...
func newSource(sourceURL model.SourceURL) (*source.Pushable, error) {
	client, err := getClient(sourceURL)
	if err != nil {
		return nil, err
	}
	src := sshSource{
		Pushing:  source.New(' ', sourceURL),
		ssh:      client,
		filename: sourceURL.Path,
	}
	var p source.Pushable = src
	return &p, nil
}

// Start ...
func (src sshSource) Start(wg *sync.WaitGroup, running func() bool, callback func(event model.Event)) error {
	// create the command
	cmd, err := src.ssh.Command("tail", "-F", src.filename)
	if err != nil {
		return err
	}

	// capture its output
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// and start it
	if err := cmd.Start(); err != nil {
		return err
	}
	return pipe.Start(wg, running, callback, func() { _ = cmd.Close() }, out)
}
