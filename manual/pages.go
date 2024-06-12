package manual

import (
	"embed"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"io/fs"
	"path"
	"slices"
	"strings"
)

const pagesDir = "pages"

var (
	//go:embed pages/*.md
	pageFS embed.FS
)

type (
	content string
	Page    string
)

func mustLoadPages() (map[Page]content, []Page) {
	dir, err := pageFS.ReadDir(pagesDir)
	if err != nil {
		panic(err)
	}
	pages := map[Page]content{}

	var names []Page
	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		if path.Ext(f.Name()) != ".md" {
			continue
		}

		name := path.Base(f.Name())
		name = name[:len(name)-len(".md")]
		name = strings.ReplaceAll(name, "__", " • ")
		name = strings.ReplaceAll(name, "_", " ")
		name = cases.Title(language.English).String(name)
		renderer := mustNewTermRenderer()

		var pageFile fs.File
		pageFile, err = pageFS.Open(path.Join(pagesDir, f.Name()))
		if err != nil {
			panic(err)
		}

		var md []byte
		md, err = io.ReadAll(pageFile)
		if err != nil {
			panic(err)
		}

		page := Page(name)
		names = append(names, page)
		var body string
		body, err = renderer.Render(string(md))
		if err != nil {
			panic(err)
		}
		pages[page] = content(body)
	}

	slices.SortFunc(
		names, func(a Page, b Page) int {
			return strings.Compare(string(a), string(b))
		},
	)

	return pages, names
}
