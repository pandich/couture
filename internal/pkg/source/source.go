package source

import (
	"couture/internal/pkg/model"
	"fmt"
	"net/url"
	"reflect"
	"sync"
)

var (
	typeRegistry = map[reflect.Type]Creator{}
	registry     []Source
)

func Available() []Source {
	return registry
}

type (
	//Implementations go in this package. Each implementation struct should be unexported and exposed with a var.
	//For each implementation, update cmd/couture/cli/source.

	//Source of events. Responsible for ingest and conversion to the standard format.
	Source interface {
		fmt.Stringer
		fmt.GoStringer
		CanHandle(url url.URL) bool
	}

	//Creator is a function which uses a URL to create a Source.
	Creator func(srcUrl url.URL) interface{}

	//PushingCallback is called by a PushingSource for each model.Event.
	PushingCallback func(event model.Event)

	//PushingSource calls a callback for each event.
	PushingSource interface {
		Source
		//Start collecting events.
		Start(wg *sync.WaitGroup, callback PushingCallback) error
		//Stop collecting events.
		Stop()
	}

	//PollableSource of events. Responsible for ingest and conversion to the standard format.
	PollableSource interface {
		Source
		//Poll performs a non-blocking poll for an event. Nil is returned if no event is available.
		Poll() (model.Event, error)
	}

	//baseSource for all Source implementations.
	baseSource struct {
		srcUrl url.URL
	}
)

//CreatorFor returns a Creator fo the specified interface.
func CreatorFor(i interface{}) (Creator, error) {
	creator, ok := typeRegistry[reflect.TypeOf(i)]
	if !ok {
		return nil, fmt.Errorf("no source handler for %v %T", i, i)
	}
	return creator, nil
}
