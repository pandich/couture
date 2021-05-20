package source

import (
	"couture/internal/pkg/model"
	"reflect"
)

// Metadata ...
type Metadata struct {
	Name        string
	Type        reflect.Type
	CanHandle   func(url model.SourceURL) bool
	Creator     func(sourceURL model.SourceURL) (*interface{}, error)
	ExampleURLs []string
}
