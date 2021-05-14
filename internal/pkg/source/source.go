package source

// TODO ssh+file:// and file:// – see https://github.com/nxadm/tail
// TODO unified way to do lookback date: take it out of the URLs and put it into CLI args (i.e. --since)

import (
	"couture/pkg/model"
	"reflect"
)

// Source ...
type (
	// Source of events. Responsible for ingest and conversion to the standard format.
	Source interface {
		URL() model.SourceURL
	}

	// Base for all Source implementations.
	Base struct {
		Source
		sourceURL model.SourceURL
	}

	creator       func(sourceURL model.SourceURL) (*interface{}, error)
	canHandleTest func(url model.SourceURL) bool
	Metadata      struct {
		Type        reflect.Type
		CanHandle   canHandleTest
		Creator     creator
		ExampleURLs []string
	}
)

// URL ...
func (source Base) URL() model.SourceURL {
	return source.sourceURL
}

// New base source.
func New(sourceURL model.SourceURL) Base {
	return Base{
		sourceURL: sourceURL,
	}
}
