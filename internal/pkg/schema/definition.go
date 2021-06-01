package schema

import (
	"couture/internal/pkg/couture"
	"github.com/BurntSushi/toml"
	"github.com/coreos/etcd/pkg/fileutil"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
)

type definition struct {
	Format     format            `json:"format"`
	Priority   priority          `json:"priority"`
	Predicates map[string]string `json:"predicates"`
	Mapping    map[string]string `json:"mapping"`
}

// LoadSchemas ...
func LoadSchemas() ([]Schema, error) {
	const schemataFilename = "schemata.toml"

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

	var schemaByName map[string]definition
	err = toml.Unmarshal(buf, &schemaByName)
	if err != nil {
		return nil, err
	}

	for name, schema := range schemaByName {
		s, err := newSchema(name, schema)
		if err != nil {
			return nil, err
		}
		schemas = append(schemas, *s)
	}
	sort.Slice(schemas, func(i, j int) bool {
		a, b := schemas[i], schemas[j]
		return a.Priority() > b.Priority()
	})
	return schemas, nil
}
