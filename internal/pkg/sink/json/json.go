package json

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"encoding/json"
	"fmt"
	errors2 "github.com/pkg/errors"
	"log"
	"strings"
)

// Sink uses json.Marshal to display the value.
type Sink struct {
	sink.Base
	// pretty determines whether or not to pretty-print the Sink.
	pretty bool
}

// New provides a configured Sink sink.
func New(options sink.Options, config string) interface{} {
	return Sink{Base: sink.New(options), pretty: config == "pretty"}
}

// Accept an event.
func (sink Sink) Accept(src source.Source, event model.Event) {
	var ba []byte
	var err error
	if sink.pretty {
		ba, err = json.MarshalIndent(&event, "", "  ")
	} else {
		ba, err = json.Marshal(&event)
	}
	if err != nil {
		log.Println(errors2.Wrapf(err, "failed to marshal event: %+v", event))
	}
	var txt = string(ba)
	txt = strings.TrimRight(txt, "\n\t ")
	fmt.Printf("%s|%s\n", src, txt)
}
