package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/iancoleman/strcase"
	"github.com/pandich/couture/couture"
	"github.com/pandich/couture/mapping"
	"github.com/pandich/couture/sink/theme"
	"github.com/pkg/errors"
	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
	"reflect"
	"strings"
)

// completionsHook generates bash/zsh auti-completions suitable for evaluation
// in an init script.
func completionsHook(_ *kong.Kong) error {
	// generate a set of flag predictors by introspecting the cli struct.
	flagPredictors := map[string]complete.Predictor{}

	// introspect the value of clie
	cliVal := reflect.ValueOf(cli)

	// go through each cli field
	for i := 0; i < cliVal.NumField(); i++ {
		// introspect the field. use its value to get the field's type
		fieldValue := cliVal.Field(i)
		field := cliVal.Type().Field(i)

		// if the field is an argument (and not a flag), move on
		if field.Tag.Get("arg") == "true" {
			continue
		}

		flag := strcase.ToKebab(field.Name)
		enum := field.Tag.Get("enum")

		switch {
		// the field is a boolean toggle
		case fieldValue.Type().Kind() == reflect.Bool:
			// there is no argument to it
			flagPredictors[flag] = predict.Nothing

		// if the field has no acceptable value enim
		case enum == "":
			// it has no predictions
			flagPredictors[flag+"="] = predict.Nothing

		// otherwise setup autocompletions for it
		default:
			const start = "${"
			const end = "}"

			// if this is a variable declaration, the it is a dynamically generated enum
			// and must be expanded
			if strings.HasPrefix(enum, start) && strings.HasSuffix(enum, end) {
				// get its name
				varName := strings.TrimSuffix(strings.TrimPrefix(enum, start), end)

				// lookup the value
				ok := false
				if enum, ok = (kong.Vars{
					"timeFormatNames": strings.Join(timeFormatNames, ","),
					"columnNames":     strings.Join(mapping.Names(), ","),
					"specialThemes":   strings.Join(theme.Names(), ","),
				})[varName]; !ok {
					return errors.Errorf("could not find enum expansion for name %s", varName)
				}
			}

			// split the enum up into its values
			for _, s := range strings.Split(enum, ",") {
				flagPredictors[flag+"="+s] = predict.Nothing
			}
		}
	}

	// delete the placeholder time format argument
	// then create all valid time-formats names
	delete(flagPredictors, "time-format=")
	for _, n := range timeFormatNames {
		flagPredictors["time-format="+n] = predict.Nothing
	}

	// generate the completions
	cmd := &complete.Command{Flags: flagPredictors}
	cmd.Complete(couture.Name)
	return nil
}
