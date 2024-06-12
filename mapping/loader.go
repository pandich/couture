package mapping

import (
	"embed"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/pandich/couture/couture"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path"
	"sort"
)

//go:embed mappings.yaml
var fs embed.FS

// LoadMappings ...
func LoadMappings() ([]Mapping, error) {
	const mappingsFilename = "mappings.yaml"

	var mappings []Mapping

	bundledConfig, err := fs.Open(mappingsFilename)
	if err != nil {
		return nil, err
	}
	bundledConfigFile, err := loadMappingFile(bundledConfig)
	if err != nil {
		return nil, err
	}
	mappings = append(mappings, bundledConfigFile...)

	userDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	userConfigFilename := path.Join(userDir, ".config", couture.Name, mappingsFilename)
	if fileutil.Exist(userConfigFilename) {
		userConfigFile, err := os.Open(userConfigFilename)
		if err != nil {
			return nil, err
		}
		defer userConfigFile.Close()
		userMappings, err := loadMappingFile(userConfigFile)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, userMappings...)
	}

	return mappings, nil
}

func loadMappingFile(mappingsFile io.ReadCloser) ([]Mapping, error) {
	var mappings []Mapping
	defer mappingsFile.Close()
	buf, err := io.ReadAll(mappingsFile)
	if err != nil {
		return nil, err
	}

	var definitionsByName map[string]Mapping
	err = yaml.Unmarshal(buf, &definitionsByName)
	if err != nil {
		return nil, err
	}

	for name, mapping := range definitionsByName {
		mapping.init(name)
		mappings = append(mappings, mapping)
	}
	sort.Slice(
		mappings, func(i, j int) bool {
			a, b := mappings[i], mappings[j]
			if a.Priority == b.Priority {
				return i < j
			}
			return a.Priority > b.Priority
		},
	)
	return mappings, nil
}
