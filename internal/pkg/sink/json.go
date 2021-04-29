package sink

import (
	"couture/internal/pkg/model"
	"encoding/json"
	"fmt"
	"log"
)

func NewJson(config string) interface{} {
	return Json{pretty: config == "pretty"}
}

//Json uses json.Marshal to display the value.
type Json struct {
	//pretty determines whether or not to pretty-print the JSON.
	pretty bool
}

func (s Json) Accept(event *model.Event) {
	var txt []byte
	var err error
	if s.pretty {
		txt, err = json.MarshalIndent(&event, "", "  ")
	} else {
		txt, err = json.Marshal(&event)

	}
	if err != nil {
		log.Println(fmt.Errorf("failed to marshal event: %v", err))
	}
	fmt.Println(string(txt))
}
