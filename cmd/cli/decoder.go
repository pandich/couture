package cli

import (
	"github.com/alecthomas/kong"
	"github.com/araddon/dateparse"
	errors2 "github.com/pkg/errors"
	"reflect"
	"regexp"
	"time"
)

func regexpDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("pattern", &value); err != nil {
			return err
		}
		r, err := regexp.Compile(value)
		if err != nil {
			return err
		}
		target.Set(reflect.ValueOf(*r))
		return nil
	}
}

func timeLikeDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("duration", &value); err != nil {
			return err
		}
		var t time.Time
		d, err := time.ParseDuration(value)
		if err == nil {
			t = time.Now().Add(-d)
		} else {
			t, err = dateparse.ParseAny(value)
			if err != nil {
				return errors2.Errorf("expected duration but got %q: %s", value, err)
			}
		}
		target.Set(reflect.ValueOf(t))
		return nil
	}
}
