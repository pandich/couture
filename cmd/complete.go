package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/gagglepanda/couture/couture"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
	"reflect"
	"strings"
)

func completionsHook(_ *kong.Kong) error {
	cliVal := reflect.ValueOf(cli)
	flagPredictors := map[string]complete.Predictor{}
	for i := 0; i < cliVal.NumField(); i++ {
		fieldValue := cliVal.Field(i)
		field := cliVal.Type().Field(i)
		if field.Tag.Get("arg") == "true" {
			continue
		}
		fieldName := field.Name
		flagName := strcase.ToKebab(fieldName)
		var enum = field.Tag.Get("enum")

		switch {
		case fieldValue.Type().Kind() == reflect.Bool:
			flagPredictors[flagName] = predict.Nothing
		case enum == "":
			flagPredictors[flagName+"="] = predict.Nothing
		default:
			const start = "${"
			const end = "}"
			if strings.HasPrefix(enum, start) && strings.HasSuffix(enum, end) {
				varName := strings.TrimSuffix(strings.TrimPrefix(enum, start), end)
				if varValue, ok := parserVars[varName]; ok {
					enum = varValue
				} else {
					return errors.Errorf("could not parse alias: %s", varValue)
				}
			}
			for _, s := range strings.Split(enum, ",") {
				flagPredictors[flagName+"="+s] = predict.Nothing
			}
		}
	}
	delete(flagPredictors, "time-format=")
	for _, n := range timeFormatNames {
		flagPredictors["time-format="+n] = predict.Nothing
	}
	cmd := &complete.Command{Flags: flagPredictors}
	cmd.Complete(couture.Name)
	return nil
}
