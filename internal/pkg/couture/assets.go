package couture

import (
	_ "couture/internal/pkg/assets" // for statik file system initialization
	"github.com/rakyll/statik/fs"
	"net/http"
)

var fileSystem http.FileSystem

func init() {
	var err error
	fileSystem, err = fs.New()
	if err != nil {
		panic(err)
	}
}

// Open ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func Open(name string) (http.File, error) {
	return fileSystem.Open(name)
}

// MustOpen ...
func MustOpen(name string) http.File {
	file, err := Open(name)
	if err != nil {
		panic(err)
	}
	return file
}
