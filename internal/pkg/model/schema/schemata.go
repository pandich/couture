package schema

import (
	"couture/internal/pkg/couture"
	"encoding/json"
	"github.com/gobuffalo/packr"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

const jsonExtension = ".json"

type schemata []Schema

// LoadSchemas ...
func LoadSchemas() ([]Schema, error) {
	schemata := schemata{}
	err := schemata.addAll()
	if err != nil {
		return nil, err
	}
	schemata.sort()
	return schemata, nil
}

func (s *schemata) add(schemaFilename string, schemaJSON string) error {
	name := path.Base(schemaFilename[0 : len(schemaFilename)-len(jsonExtension)])
	var schemaDefinition = definition{}
	err := json.Unmarshal([]byte(schemaJSON), &schemaDefinition)
	if err != nil {
		return err
	}
	*s = append(*s, newSchema(name, schemaDefinition))
	return nil
}

func (s *schemata) addAll() error {
	err := s.addBundled()
	if err != nil {
		return err
	}
	err = s.addUser()
	if err != nil {
		return err
	}
	return nil
}

func (s *schemata) addBundled() error {
	schemaBox := packr.NewBox("definitions")
	for _, schemaFilename := range schemaBox.List() {
		if path.Ext(schemaFilename) == jsonExtension {
			schemaJSON, err := schemaBox.FindString(schemaFilename)
			if err != nil {
				return err
			}
			err = s.add(schemaFilename, schemaJSON)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *schemata) addUser() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir := path.Join(home, ".config", couture.Name, "schemas")
	return filepath.Walk(configDir, func(schemaFilename string, info fs.FileInfo, err error) error {
		if path.Ext(schemaFilename) == jsonExtension {
			if info != nil && info.IsDir() {
				return nil
			}
			schemaJSON, err := ioutil.ReadFile(schemaFilename)
			if err != nil {
				return err
			}
			err = s.add(schemaFilename, string(schemaJSON))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *schemata) sort() {
	sort.Slice(*s, func(i, j int) bool { return strings.Compare((*s)[i].Name(), (*s)[j].Name()) <= 0 })
}
