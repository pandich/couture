package cli

import (
	"couture/internal/pkg/manager"
	"github.com/alecthomas/kong"
)

//coreOptions contain the core kong cli options.
var coreOptions = []kong.Option{
	kong.Name("couture"),
	kong.Description("Tail multiple log sources."),
	kong.ShortUsageOnError(),
	kong.ConfigureHelp(kong.HelpOptions{
		Tree:     true,
		Indenter: kong.TreeIndenter,
	}),
	kong.ExplicitGroups([]kong.Group{
		{Key: "aws", Title: "AWS Options"},
		{Key: "display", Title: "Display Options"},
		{Key: "behavior", Title: "Behavioral Options"},
		{Key: "source", Title: "Input Options"},
		{Key: "sink", Title: "Output Options"},
		{Key: "filter", Title: "Filtering Options"},
	}),
}

//Options from the parsed CLI.
func Options() []interface{} {
	var opts = []interface{}{
		manager.VerboseDisplayOption(coreCli.Verbosity),
	}
	if coreCli.ShowPrefix {
		opts = append(opts, manager.ShowPrefixDisplayOption())
	}
	if coreCli.ShortNames {
		opts = append(opts, manager.ShortNamesDisplayOption())
	}
	return opts
}
