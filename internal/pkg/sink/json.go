package sink

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/source"
	"encoding/json"
	"fmt"
	"log"
)

//NewJson provides a configured Json sink.
func NewJson(options Options, config string) interface{} {
	return Json{baseSink: baseSink{options: options}, pretty: config == "pretty"}
}

//Json uses json.Marshal to display the value.
type Json struct {
	baseSink
	//pretty determines whether or not to pretty-print the JSON.
	pretty bool
}

func (sink Json) Accept(src source.Source, event model.Event) {
	var txt []byte
	var err error
	if sink.pretty {
		txt, err = json.MarshalIndent(&event, "", "  ")
	} else {
		txt, err = json.Marshal(&event)

	}
	if err != nil {
		log.Println(fmt.Errorf("failed to marshal event: %s\n%v", txt, err))
	}
	fmt.Println(fmt.Sprintf(src.String() + "|" + string(txt)))
}
