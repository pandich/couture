package cli

import (
	"couture/internal/pkg/sink/json"
	"couture/internal/pkg/sink/pretty"
	errors2 "github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"os"
	"os/exec"
	"strings"
)

func sinkOption(flags *pflag.FlagSet) (interface{}, error) {
	out, err := getWriter(flags)
	if err != nil {
		return nil, err
	}
	outputFormat, err := flags.GetString(outputFormatFlag)
	if err != nil {
		return nil, err
	}
	switch outputFormat {
	case "json":
		return json.New(out), nil
	case "pretty":
		var wrap, err = flags.GetInt(wrapFlag)
		if err != nil {
			return nil, err
		}
		noWrap, err := flags.GetBool(noWrapFlag)
		if err != nil {
			return nil, err
		}
		if noWrap {
			wrap = pretty.NoWrap
		}
		return pretty.New(out, wrap), nil
	default:
		return nil, errors2.Errorf("unknown output format: %s", outputFormat)
	}
}

func getWriter(flags *pflag.FlagSet) (io.Writer, error) {
	defaultOut := os.Stdout

	paginate, err := flags.GetBool(paginateFlag)
	if err != nil {
		return nil, err
	}
	if !paginate {
		return defaultOut, nil
	}

	var pager = viper.GetString(paginatorEnvKey)
	if pager == "" {
		pager = viper.GetString(paginatorConfigKey)
	}
	if pager == "" {
		var err error
		pager, err = flags.GetString(paginatorFlag)
		if err != nil {
			return nil, err
		}
	}
	if pager == "" {
		return defaultOut, nil
	}

	var pagerArgs = strings.Split(pager, " \t\n")
	pager, pagerArgs = pagerArgs[0], pagerArgs[1:]
	pagerCmd := exec.Command(pager, pagerArgs...)

	// I/O
	pagerCmd.Stdout, pagerCmd.Stderr = os.Stdout, os.Stderr
	writer, err := pagerCmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	if err = pagerCmd.Start(); err != nil {
		return nil, err
	}
	return writer, nil
}
