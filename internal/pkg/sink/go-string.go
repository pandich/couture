package sink

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/source"
	"fmt"
)

//NewGoString provides a configured GoString sink.
func NewGoString(_ string) interface{} {
	return GoString{}
}

//GoString uses the GoStringer interface to display.
type GoString struct {
}

func (s GoString) Accept(src source.Source, event model.Event) {
	fmt.Printf("%s %+v\n", src.Name(), event)
}
