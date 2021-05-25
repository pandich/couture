package cli

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model"
	"github.com/muesli/termenv"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

// Start ...
func Start() *model.Manager {
	if runtime.GOOS == "windows" {
		parser.Fatalf("unsupported operating system: %s", runtime.GOOS)
	}

	// load config
	err := loadConfig()
	parser.FatalIfErrorf(err)

	// expand aliases, etc.
	args, err := expandAliases()
	parser.FatalIfErrorf(err)

	_, err = parser.Parse(args)
	parser.FatalIfErrorf(err)

	// get cli flags and args
	mgrOptions, err := flags()
	parser.FatalIfErrorf(err)
	sources, err := sourceArgs()
	parser.FatalIfErrorf(err)

	// create the manager and start it
	mgr, err := manager.New(append(mgrOptions, sources...)...)
	parser.FatalIfErrorf(err)

	err = (*mgr).Start()
	parser.FatalIfErrorf(err)

	trapInterrupt(mgr)

	return mgr
}

func trapInterrupt(mgr *model.Manager) {
	const stopGracePeriod = 2 * time.Second

	die := func() {
		println()
		termenv.Reset()
		os.Exit(0)
	}

	interruption := make(chan os.Signal)
	signal.Notify(interruption, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interruption
		(*mgr).Stop()
		go func() { time.Sleep(stopGracePeriod); die() }()
		(*mgr).Wait()
		die()
	}()
}
