package config

import (
	"github.com/alecthomas/kong"
)

func Load(_ *kong.Context) error { return nil }
func Sources() []interface{}     { return []interface{}{} }
func Sinks() []interface{}       { return []interface{}{} }
func Options() []interface{}     { return []interface{}{} }

func MustLoad(ctx *kong.Context) {
	if err := Load(ctx); err != nil {
		panic(err)
	}
}
