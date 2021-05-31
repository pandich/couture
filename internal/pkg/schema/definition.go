package schema

import (
	"couture/internal/pkg/couture"
	"github.com/BurntSushi/toml"
	"io"
	"io/ioutil"
	"sort"
)

var schemataFile = couture.MustOpen("/schemata.toml")

type definition struct {
	Priority   priority          `json:"priority"`
	Predicates map[string]string `json:"predicates"`
	Mapping    map[string]string `json:"mapping"`
}

// LoadSchemas ...
func LoadSchemas() ([]Schema, error) {
	bundledSchemas, err := loadSchemaFile(schemataFile)
	if err != nil {
		return nil, err
	}
	var schemas []Schema
	schemas = append(schemas, bundledSchemas...)
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
		schemas = append(schemas, newSchema(name, schema))
	}
	sort.Slice(schemas, func(i, j int) bool {
		a, b := schemas[i], schemas[j]
		return a.Priority() > b.Priority()
	})
	return schemas, nil
}
