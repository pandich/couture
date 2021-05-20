package cli

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/sink/json"
	"couture/internal/pkg/sink/pretty"
	errors2 "github.com/pkg/errors"
	"os"
)

func getFlags() ([]interface{}, error) {
	snk, err := sinkFlag()
	if err != nil {
		return nil, err
	}
	return []interface{}{
		manager.LogLevelOption(cli.Log.Level),
		manager.FilterOption(cli.Log.Include, cli.Log.Exclude),
		manager.SinceOption(cli.Log.Since),
		snk,
	}, nil
}

func sinkFlag() (interface{}, error) {
	wrap := wrapFlag()
	switch cli.Log.OutputFormat {
	case "json":
		return json.New(os.Stdout), nil
	case "pretty":
		return pretty.New(os.Stdout, wrap), nil
	default:
		return nil, errors2.Errorf("unknown output format: %s\n", cli.Log.OutputFormat)
	}
}

func wrapFlag() int {
	var wrap = pretty.AutoWrap
	if cli.Log.NoWrap {
		wrap = pretty.NoWrap
	} else if cli.Log.Wrap != pretty.NoWrap {
		wrap = int(cli.Log.Wrap)
	}
	return wrap
}
