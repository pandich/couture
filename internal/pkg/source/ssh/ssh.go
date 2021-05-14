package ssh

import (
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/pushing"
	"couture/pkg/model"
	"reflect"
	"sync"
)

// Metadata ...
func Metadata() source.Metadata {
	return source.Metadata{
		Type:        reflect.TypeOf(sshSource{}),
		CanHandle:   func(url model.SourceURL) bool { return url.Scheme == "ssh" },
		Creator:     create,
		ExampleURLs: []string{"ssh://host/<path>"},
	}
}

// sshSource ...
type sshSource struct {
	source.Base
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
func newSource(_ model.SourceURL) (*pushing.Source, error) {
	panic("implement me")
}

// Start ...
func (src sshSource) Start(wg *sync.WaitGroup, running func() bool, _ func(event model.Event)) error {
	f := func() {
		defer wg.Done()
		for running() {
			// TODO
		}
	}
	wg.Add(1)
	go f()
	return nil
}
