package schema

import (
	"couture/internal/pkg/couture"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/oriser/regroup"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

type definition struct {
	Format     format            `yaml:"format,omitempty"`
	Priority   priority          `yaml:"priority,omitempty"`
	Predicates map[string]string `yaml:"predicates,omitempty"`
	Mapping    map[string]string `yaml:"mapping,omitempty"`
	Display    map[string]string `yaml:"display,omitempty"`
}

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

	var schemaByName map[string]definition
	err = yaml.Unmarshal(buf, &schemaByName)
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

func (d definition) predicateFields() []string {
	var predicateFields []string
	for fieldName := range d.Predicates {
		predicateFields = append(predicateFields, fieldName)
	}
	return predicateFields
}

func (d definition) predicatePatterns() map[string]*regexp.Regexp {
	predicatePatterns := map[string]*regexp.Regexp{}
	for fieldName, pattern := range d.Predicates {
		if pattern != "" {
			predicatePatterns[fieldName] = regexp.MustCompile(pattern)
		} else {
			predicatePatterns[fieldName] = nil
		}
	}
	return predicatePatterns
}

func (d definition) canHandlePredicate() predicate {
	predicateFields := d.predicateFields()
	predicatePatterns := d.predicatePatterns()

	var canHandle predicate
	switch d.Format {
	case JSON:
		canHandle = func(s string) bool {
			values := gjson.GetMany(s, predicateFields...)
			for i := range predicateFields {
				value := values[i]
				field := predicateFields[i]
				pattern := predicatePatterns[field]
				if pattern == nil {
					if !value.Exists() {
						return false
					}
				} else {
					stringValue := value.String()
					if !pattern.MatchString(stringValue) {
						return false
					}
				}
			}
			return true
		}
	case Text:
		const textRootPredicate = "_"
		pattern := predicatePatterns[textRootPredicate]
		canHandle = func(s string) bool {
			return pattern.MatchString(strings.TrimRight(s, "\n"))
		}
	default:
		return nil
	}

	return canHandle
}

func (d definition) textPattern() *regroup.ReGroup {
	const textRootPredicate = "_"

	var textPattern *regroup.ReGroup
	if d.Format == Text {
		pattern := d.predicatePatterns()[textRootPredicate]
		textPattern = regroup.MustCompile(pattern.String())
	}
	return textPattern
}

func (d definition) inputFields() []string {
	var inputFields []string
	for _, v := range d.Mapping {
		inputFields = append(inputFields, v)
	}
	return inputFields
}

func (d definition) inverseMapping() map[string]string {
	inverseMapping := map[string]string{}
	for k, v := range d.Mapping {
		inverseMapping[v] = k
	}
	return inverseMapping
}
