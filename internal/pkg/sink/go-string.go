package sink

import (
	"couture/internal/pkg/model"
	"fmt"
)

func NewGoString(_ string) interface{} {
	return GoString{}
}

//GoString uses the GoStringer interface to display.
type GoString struct {
}

func (s GoString) Accept(event *model.Event) {
	fmt.Printf("%+v\n", *event)
}
