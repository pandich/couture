package mapping

import (
	"embed"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/gagglepanda/couture/couture"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path"
	"sort"
)

//go:embed schemata.yaml
var fs embed.FS

// LoadSchemas ...
func LoadSchemas() ([]Schema, error) {
	const schemataFilename = "schemata.yaml"

	var schemas []Schema

	bundledConfig, err := fs.Open(schemataFilename)
	if err != nil {
		return nil, err
	}
	bundledConfigFile, err := loadSchemaFile(bundledConfig)
	if err != nil {
		return nil, err
	}
	schemas = append(schemas, bundledConfigFile...)

	userDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	userConfigFilename := path.Join(userDir, ".config", couture.Name, schemataFilename)
	if fileutil.Exist(userConfigFilename) {
		userConfigFile, err := os.Open(userConfigFilename)
		if err != nil {
			return nil, err
		}
		defer userConfigFile.Close()
		userSchemas, err := loadSchemaFile(userConfigFile)
		if err != nil {
			return nil, err
		}
		schemas = append(schemas, userSchemas...)
	}

	return schemas, nil
}

func loadSchemaFile(schemataFile io.ReadCloser) ([]Schema, error) {
	var schemas []Schema
	defer schemataFile.Close()
	buf, err := io.ReadAll(schemataFile)
	if err != nil {
		return nil, err
	}

	var definitionsByName map[string]Schema
	err = yaml.Unmarshal(buf, &definitionsByName)
	if err != nil {
		return nil, err
	}

	for name, schema := range definitionsByName {
		schema.init(name)
		schemas = append(schemas, schema)
	}
	sort.Slice(schemas, func(i, j int) bool {
		a, b := schemas[i], schemas[j]
		if a.Priority == b.Priority {
			return i < j
		}
		return a.Priority > b.Priority
	})
	return schemas, nil
}
