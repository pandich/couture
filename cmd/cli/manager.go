package cli

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model"
	"github.com/muesli/termenv"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run runs the manager using the CLI arguments.
func Run() {
	// load config
	err := loadAliasConfig()
	parser.FatalIfErrorf(err)

	// expand aliases, etc.
	args, err := expandAliases()
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
