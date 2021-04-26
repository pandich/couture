package sink

import (
	"couture/pkg/couture/model"
	log "github.com/sirupsen/logrus"
)

var (
	Logrus = logrusSink{}
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
}

type (
	logrusSink struct {
	}
)

func (l logrusSink) ConsumeEvent(event *model.Event) {
	log.Printf("%#v\n", *event)
}
