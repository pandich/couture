package source

import (
	"couture/internal/pkg/model"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strings"
)

// Pushable ...
type (
	Pushable Source

	// Pushing Source.
	Pushing struct {
		Pushable
		id        string
		sigil     rune
		sourceURL model.SourceURL
	}
)

// New base Source.
func New(sigil rune, sourceURL model.SourceURL) *Pushing {
	u := url.URL(sourceURL)
	s := u.String()
	hasher := sha256.New()
	hasher.Write([]byte(s))
	return &Pushing{
		id:        "Source" + strings.ToUpper(hex.EncodeToString(hasher.Sum(nil))[0:15]),
		sigil:     sigil,
		sourceURL: sourceURL,
	}
}

// URL ...
func (source Pushing) URL() model.SourceURL {
	return source.sourceURL
}

// ID ...
func (source Pushing) ID() string {
	return source.id
}

// Sigil ...
func (source Pushing) Sigil() rune {
	return source.sigil
}
