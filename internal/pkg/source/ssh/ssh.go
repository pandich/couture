package ssh

// TODO ssh source implementation

import (
	"couture/internal/pkg/source"
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
	source.Pushing
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
func newSource(_ model.SourceURL) (*source.Pushable, error) {
	panic("implement me")
}

// Start ...
func (src sshSource) Start(wg *sync.WaitGroup, running func() bool, _ func(event model.Event)) error {
	f := func() {
		defer wg.Done()
		for running() {
		}
	}
	wg.Add(1)
	go f()
	return nil
}
