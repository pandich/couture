package cli

import (
	"couture/internal/pkg/couture"
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/schema"
	"encoding/json"
	"github.com/gobuffalo/packr"
	"github.com/muesli/termenv"
	"io/fs"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"
)

// Run runs the manager using the CLI arguments.
func Run() {
	schemas := loadSchemas()
	var args = os.Args[1:]

	// load config
	err := loadAliasConfig()
	parser.FatalIfErrorf(err)

	// expand aliases, etc.
	args, err = expandAliases(args)
	parser.FatalIfErrorf(err)

	// parse CLI args
	_, err = parser.Parse(args)
	parser.FatalIfErrorf(err)

	// get manager config
	mgrConfig := manager.Config{
		Level:          cli.Level,
		Since:          &cli.Since,
		IncludeFilters: cli.Include,
		ExcludeFilters: cli.Exclude,
		Schemas:        schemas,
	}

	// get sources and sinks
	mgrOptions, err := getSourceAndSinkOptions()
	parser.FatalIfErrorf(err)

	// create the manager
	mgr, err := manager.New(mgrConfig, mgrOptions...)
	parser.FatalIfErrorf(err)
	// start it
	trapInterrupt(mgr)
	err = (*mgr).Start()
	parser.FatalIfErrorf(err)
	// wait for it to die
	(*mgr).Wait()
	os.Exit(0)
}

func loadSchemas() []schema.Schema {
	const jsonExtension = ".json"

	schemaBox := packr.NewBox("./schemas")

	var schemas []schema.Schema
	addSchema := func(schemaFilename string, schemaJSON string) {
		name := path.Base(schemaFilename[0 : len(schemaFilename)-len(jsonExtension)])
		var schemaDefinition = schema.Definition{}
		err := json.Unmarshal([]byte(schemaJSON), &schemaDefinition)
		parser.FatalIfErrorf(err)
		schemas = append(schemas, schema.NewSchema(name, schemaDefinition))
	}

	for _, schemaFilename := range schemaBox.List() {
		schemaJSON, err := schemaBox.FindString(schemaFilename)
		parser.FatalIfErrorf(err)
		addSchema(schemaFilename, schemaJSON)
	}

	home, err := os.UserHomeDir()
	parser.FatalIfErrorf(err)
	schemasDir := path.Join(home, ".config", couture.Name, "schemas")
	err = filepath.Walk(schemasDir, func(schemaFilename string, info fs.FileInfo, err error) error {
		if path.Ext(schemaFilename) == jsonExtension {
			if info != nil && info.IsDir() {
				return nil
			}
			schemaJSON, err := ioutil.ReadFile(schemaFilename)
			if err != nil {
				return err
			}
			addSchema(schemaFilename, string(schemaJSON))
		}
		return nil
	})
	parser.FatalIfErrorf(err)

	sort.Slice(schemas, func(i, j int) bool { return strings.Compare(schemas[i].Name(), schemas[j].Name()) <= 0 })

	return schemas
}

func trapInterrupt(mgr *model.Manager) {
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		const stopGracePeriod = 250 * time.Millisecond
		defer close(interrupt)

		cleanup := func() { termenv.Reset(); os.Exit(0) }

		<-interrupt
		(*mgr).Stop()

		go func() { time.Sleep(stopGracePeriod); cleanup() }()
		(*mgr).Wait()
		cleanup()
	}()
}
