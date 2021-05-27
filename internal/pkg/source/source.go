package source

import (
	"couture/internal/pkg/model"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"reflect"
	"strings"
	"sync"
) // Source ...

// Source ...
type (
	// Source of events. Responsible for ingest and conversion to the standard format.
	Source interface {
		// ID is the unique id for this source.
		ID() string
		// Sigil represents the type of source in a single character.
		Sigil() rune
		// URL is the URL from which the events come.
		URL() model.SourceURL
		// Start collecting events.
		Start(wg *sync.WaitGroup, running func() bool, srcChan chan Event, errChan chan Error) error
	}
	// BaseSource ...
	BaseSource struct {
		id        string
		sigil     rune
		sourceURL model.SourceURL
	}

	// Event ...
	Event struct {
		model.Event
		Source Source
	}

	// Error ...
	Error struct {
		Source Source
		Error  error
	}

	// Metadata ...
	Metadata struct {
		Name        string
		Type        reflect.Type
		CanHandle   func(url model.SourceURL) bool
		Creator     func(sourceURL model.SourceURL) (*interface{}, error)
		ExampleURLs []string
	}
)

// New base Source.
func New(sigil rune, sourceURL model.SourceURL) BaseSource {
	u := url.URL(sourceURL)
	s := u.String()
	hasher := sha256.New()
	hasher.Write([]byte(s))
	return BaseSource{
		id:        "Source" + strings.ToUpper(hex.EncodeToString(hasher.Sum(nil))[0:15]),
		sigil:     sigil,
		sourceURL: sourceURL,
	}
}

// ID ...
func (b BaseSource) ID() string {
	return b.id
}

// Sigil ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (b BaseSource) Sigil() rune {
	return b.sigil
}

// URL ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (b BaseSource) URL() model.SourceURL {
	return b.sourceURL
}
