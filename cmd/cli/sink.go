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
		return pretty.New(out), nil
	default:
		return nil, errors2.Errorf("unknown output format: %s", outputFormat)
	}
}

func getWriter(flags *pflag.FlagSet) (io.Writer, error) {
	//goland:noinspection SpellCheckingInspection
	const theLogNavigatorPaginator = "lnav"

	defaultOut := os.Stdout

	noPaginate, err := flags.GetBool(noPaginatorFlag)
	if err != nil {
		return nil, err
	}
	if noPaginate {
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
	if len(pagerArgs) == 0 {
		if pager == theLogNavigatorPaginator {
			pagerArgs = append(pagerArgs, "-t")
		}
	}
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
