package source

import (
	"couture/internal/pkg/model"
)

// Pushable ...
type (
	Pushable Source

	// Pushing Source.
	Pushing struct {
		Pushable
		sourceURL model.SourceURL
	}
)

// New base Source.
func New(sourceURL model.SourceURL) *Pushing {
	return &Pushing{
		sourceURL: sourceURL,
	}
}

// URL ...
func (source Pushing) URL() model.SourceURL {
	return source.sourceURL
}
