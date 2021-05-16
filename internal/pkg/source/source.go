package source

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

	Metadata struct {
		Type        reflect.Type
		CanHandle   func(url model.SourceURL) bool
		Creator     func(sourceURL model.SourceURL) (*interface{}, error)
		ExampleURLs []string
	}
)

// New base source.
func New(sourceURL model.SourceURL) Base {
	return Base{
		sourceURL: sourceURL,
	}
}

// URL ...
func (source Base) URL() model.SourceURL {
	return source.sourceURL
}

// MetadataGroup ...
type MetadataGroup []Metadata

// ExampleURLs ...
func (grep MetadataGroup) ExampleURLs() []string {
	var exampleURLs []string
	for _, src := range grep {
		exampleURLs = append(exampleURLs, src.ExampleURLs...)
		exampleURLs = append(exampleURLs, "\n")
	}
	return exampleURLs
}
