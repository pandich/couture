package source

import (
	"couture/internal/pkg/model"
	"reflect"
)

// Metadata ...
type Metadata struct {
	Type        reflect.Type
	CanHandle   func(url model.SourceURL) bool
	Creator     func(sourceURL model.SourceURL) (*interface{}, error)
	ExampleURLs []string
}

// MetadataGroup ...
type MetadataGroup []Metadata

// ExampleURLs ...
func (grep MetadataGroup) ExampleURLs() []string {
	var exampleURLs []string
	for _, src := range grep {
		exampleURLs = append(exampleURLs, src.ExampleURLs...)
	}
	return exampleURLs
}
