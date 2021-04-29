package cli

import (
	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
	"reflect"
)

type (
	//creator converts a string into a resource (e.g. source or sink)
	creator func(config string) interface{}
	//creators maps reflect.Type to creator.
	creators map[reflect.Type]creator
	//creatorMapper implements the kong.Mapper interface.
	creatorMapper struct {
		creators creators
	}
)

//mapper creates a new kong.Option registering a kong.Mapper for a creator.
func mapper(t reflect.Type, creator creator) kong.Option {
	return kong.TypeMapper(t, creatorMapper{creators: creators{
		t: creator,
	}})
}

//one specified that a single one instance of the type is being mapped.
func one(i interface{}) reflect.Type {
	return reflect.PtrTo(reflect.TypeOf(i))
}

//many specified that a slice of instances of the type is being mapped.
func many(i interface{}) reflect.Type {
	return reflect.SliceOf(reflect.TypeOf(i))
}

//Decode decodes expecting a string argument to a type know it.
func (m creatorMapper) Decode(ctx *kong.DecodeContext, target reflect.Value) error {
	if ctx.Scan.Peek().Type == kong.FlagValueToken {
		token := ctx.Scan.Pop()
		switch config := token.Value.(type) {
		case string:
			var t = target.Type()
			if t.Kind() != reflect.Ptr {
				t = t.Elem()
			}
			creator, ok := m.creators[target.Type()]
			if !ok {
				return errors.Errorf("unknown type (%T) with config %s", token.Value, config)
			}
			c := creator(config)
			switch target.Kind() {
			case reflect.Slice:
				target.Set(reflect.Append(target, reflect.ValueOf(c)))
			case reflect.Ptr:
				target.Elem().Set(reflect.ValueOf(c))
			default:
				target.Set(reflect.ValueOf(c))
			}

		default:
			return errors.Errorf("expected string but got %q (%T)", token.Value, token.Value)
		}
	}
	return nil
}
