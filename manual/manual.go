package manual

import (
	_ "embed"
	"github.com/pkg/errors"
)

func New() Manual {
	pages, names := mustLoadPages()
	return &manual{
		pages: pages,
		names: names,
	}
}

var errPageNotFound = errors.New("page not found")

type (
	Manual interface {
		Pages() []Page
		Page(name Page) (content string, err error)
	}

	manual struct {
		pages map[Page]content
		names []Page
	}
)

func (man *manual) Pages() []Page { return man.names }
func (man *manual) Page(name Page) (string, error) {
	body, found := man.pages[name]
	if !found {
		return "", errors.Wrapf(errPageNotFound, "page '%s'", name)
	}

	return string(body), nil
}
