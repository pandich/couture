package cli

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model"
	errors2 "github.com/pkg/errors"
	"gopkg.in/multierror.v1"
	"os"
)

func getArgs() ([]interface{}, error) {
	var sources []interface{}
	var violations []error
	for _, u := range cli.Log.Source {
		sourceURL := model.SourceURL(u)
		var handled bool
		for _, metadata := range manager.SourceMetadata {
			if !metadata.CanHandle(sourceURL) {
				continue
			}
			handled = true
			configuredSource, err := metadata.Creator(sourceURL)
			if err != nil {
				violations = append(violations, err)
			} else {
				sources = append(sources, *configuredSource)
			}
			break
		}
		if !handled {
			violations = append(violations, errors2.Errorf("invalid source URL: %+v\n", sourceURL))
		}
	}
	if len(violations) > 0 {
		return nil, multierror.New(violations)
	}
	return sources, nil
}

func evaluateArgs() []string {
	if len(os.Args) >= 2 && os.Args[1] == installCompletionsCommand {
		return os.Args[1:]
	}
	// use the log command by default
	args := append([]string{logCommand}, os.Args[1:]...)
	expandArgAliases(args[1:])
	return args
}

func expandArgAliases(args []string) {
	aliases := aliasConfig()
	for i := range args {
		if alias, ok := aliases[args[i]]; ok {
			args[i] = alias
		}
	}
}
