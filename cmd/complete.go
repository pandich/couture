package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/iancoleman/strcase"
	"github.com/pandich/couture/couture"
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
			if enum[0] == '$' {
				enum = enum[2 : len(enum)-1]
				if s, ok := parserVars[enum]; ok {
					enum = s
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
