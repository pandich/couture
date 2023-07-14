package source

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gagglepanda/couture/event"
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
		URL() event.SourceURL
		// Start collecting events.
		Start(wg *sync.WaitGroup, running func() bool, srcChan chan Event, snkChan chan event.SinkEvent, errChan chan Error) error
	}

	// BaseSource ...
	BaseSource struct {
		id        string
		sigil     rune
		sourceURL event.SourceURL
	}

	// Event ...
	Event struct {
		Source Source
		Event  string
	}

	// Error ...
	Error struct {
		SourceURL event.SourceURL
		Error     error
	}

	// Metadata ...
	Metadata struct {
		Name        string
		Type        reflect.Type
		CanHandle   func(url event.SourceURL) bool
		Creator     func(since *time.Time, sourceURL event.SourceURL) (*Source, error)
		ExampleURLs []string
	}
)

// New base Source.
func New(sigil rune, sourceURL event.SourceURL) BaseSource {
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
func (b BaseSource) URL() event.SourceURL {
	return b.sourceURL
}

// Start ...
func (b BaseSource) Start(_ *sync.WaitGroup, _ func() bool, _ chan Event, _ chan event.SinkEvent, _ chan Error) error {
	panic("implement me")
}
