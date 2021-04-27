package source

import (
	"couture/pkg/couture/model"
)

type (
	// Source of events. Responsible for ingest and conversion to the standard format.
	// Implementations go in this package. Each implementation struct should be unexported and exposed with a var.
	Source interface {
		// ProvideEvent if no events are available, nil is returned with no error.
		ProvideEvent() (*model.Event, error)
	}
)
