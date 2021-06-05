package schema

import (
	"couture/internal/pkg/couture"
	"github.com/coreos/etcd/pkg/fileutil"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
)

// LoadSchemas ...
func LoadSchemas() ([]Schema, error) {
	const schemataFilename = "schemata.yaml"

	var schemas []Schema

	bundledConfig := couture.MustOpen("/" + schemataFilename)
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
	buf, err := ioutil.ReadAll(schemataFile)
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
