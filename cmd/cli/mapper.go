package cli

import (
	"errors"
	"github.com/alecthomas/kong"
	"reflect"
	"regexp"
)

var (
	errBadType  = errors.New("unknown type")
	errBadEvent = errors.New("could not decode event")
)

// coreMappers are kong.Mapper mappers exposed via kong.Option structs.
var coreMappers = []kong.Option{
	// regexp
	kong.TypeMapper(reflect.PtrTo(reflect.TypeOf(regexp.Regexp{})), regexpMapper{}),
	kong.TypeMapper(reflect.SliceOf(reflect.PtrTo(reflect.TypeOf(regexp.Regexp{}))), regexpMapper{}),
}

// regexpMapper uses regexp.Compile to compile the specified pattern.
type regexpMapper struct{}

// Decode ...
func (mapper regexpMapper) Decode(ctx *kong.DecodeContext, target reflect.Value) error {
	token := ctx.Scan.Pop()
	switch pattern := token.Value.(type) {
	case string:
		filter, err := regexp.Compile(pattern)
		if err != nil {
			return errBadEvent
		}
		target.Set(reflect.Append(target, reflect.ValueOf(filter)))
	default:
		return errBadType
	}
	return nil
}
