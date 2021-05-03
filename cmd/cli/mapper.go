package cli

import (
	"couture/internal/pkg/sink"
	"errors"
	"github.com/alecthomas/kong"
	errors2 "github.com/pkg/errors"
	"reflect"
	"regexp"
)

var (
	errBadType  = errors.New("unknown type")
	errBadEvent = errors.New("could not decode event")
)
var (
	// coreMappers are kong.Mapper mappers exposed via kong.Option structs.
	coreMappers = []kong.Option{
		// regexp
		kong.TypeMapper(reflect.PtrTo(reflect.TypeOf(regexp.Regexp{})), regexpMapper{}),
		kong.TypeMapper(reflect.SliceOf(reflect.PtrTo(reflect.TypeOf(regexp.Regexp{}))), regexpMapper{}),
	}
)

type (
	// creator converts a string into a resource (e.g. source or sink).
	creator func(options sink.Options, config string) interface{}
	// sourceCreators maps reflect.Type to creator.
	creators map[reflect.Type]creator
	// creatorMapper implements the kong.Mapper interface.
	creatorMapper struct {
		creators creators
	}
	// regexpMapper uses regexp.Compile to compile the specified pattern.
	regexpMapper struct{}
)

// mapper creates a new kong.Option registering a kong.Mapper for a creator for required, optional, and slice types.
func mapper(i interface{}, creator creator) []kong.Option {
	t := reflect.TypeOf(i)
	return []kong.Option{
		kong.TypeMapper(t, creatorMapper{creators: creators{reflect.PtrTo(t): creator}}),
		kong.TypeMapper(reflect.PtrTo(t), creatorMapper{creators: creators{reflect.PtrTo(t): creator}}),
		kong.TypeMapper(reflect.SliceOf(t), creatorMapper{creators: creators{reflect.SliceOf(t): creator}}),
	}
}

// Decode ...
func (mapper creatorMapper) Decode(ctx *kong.DecodeContext, target reflect.Value) error {
	var arg string
	switch ctx.Scan.Peek().Type {
	case kong.PositionalArgumentToken:
	case kong.ShortFlagTailToken:
	case kong.FlagValueToken:
		arg = ctx.Scan.Pop().String()
	default:
		arg = ""
	}
	creator, ok := mapper.creators[target.Type()]
	if !ok {
		return errors2.WithMessagef(errBadType, "%v %T", target, target)
	}
	value := reflect.ValueOf(creator(cliSinkOptions, arg))
	switch target.Kind() {
	case reflect.Slice:
		target.Set(reflect.Append(target, value))
	case reflect.Ptr:
		target.Elem().Set(value)
	default:
		target.Set(value)
	}
	return nil
}

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
