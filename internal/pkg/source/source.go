package source

import (
	"couture/pkg/couture/model"
)

type (

	// Source of events. Responsible for ingest and conversion to the standard format.
	Source interface {
		// ProvideEvent if no events are available, nil is returned with no error.
		ProvideEvent() (*model.Event, error)
	}
)

/*

Implementations go here. Each implementation struct should be unexported and exposed with a var.

Example:

	var (
		Something Source = somethingSource{}
	)

	type (
		somethingSource struct {
		}
	)

	func (s somethingSource) ProvideEvent() (*model.Event, error) {
		return &model.Event{
			// values
		}, nil
	}

*/
