package source

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gagglepanda/couture/model"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"
) // Source ...

// Source ...
type (
	// Source of events. Responsible for ingest and conversion to the standard format.
	Source interface {
		// Sigil represents the type of source in a single character.
		Sigil() rune
		// URL is the URL from which the events come.
		URL() model.SourceURL
		// Start collecting events.
		Start(wg *sync.WaitGroup, running func() bool, srcChan chan Event, snkChan chan model.SinkEvent, errChan chan Error) error
	}

	// BaseSource ...
	BaseSource struct {
		id        string
		sigil     rune
		sourceURL model.SourceURL
	}

	// Event ...
	Event struct {
		Source Source
		Event  string
	}

	// Error ...
	Error struct {
		SourceURL model.SourceURL
		Error     error
	}

	// Metadata ...
	Metadata struct {
		Name        string
		Type        reflect.Type
		CanHandle   func(url model.SourceURL) bool
		Creator     func(since *time.Time, sourceURL model.SourceURL) (*Source, error)
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

// Sigil ...
func (b BaseSource) Sigil() rune {
	return b.sigil
}

// URL ...
func (b BaseSource) URL() model.SourceURL {
	return b.sourceURL
}

// Start ...
func (b BaseSource) Start(_ *sync.WaitGroup, _ func() bool, _ chan Event, _ chan model.SinkEvent, _ chan Error) error {
	panic("implement me")
}
